package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// EmailService handles sending emails via different providers
type EmailService interface {
	SendTeamInvite(email, inviteURL string, isNewUser bool) error
}

// EmailProvider represents the email provider type
type EmailProvider string

const (
	EmailProviderSupabase   EmailProvider = "supabase"
	EmailProviderResend     EmailProvider = "resend"
	EmailProviderElastic    EmailProvider = "elastic"
	EmailProviderNone       EmailProvider = "none" // For testing/development
)

// SupabaseEmailService sends emails via Supabase Admin API
type SupabaseEmailService struct {
	supabaseURL string
	serviceKey  string
	logger      *zap.Logger
}

// ResendEmailService sends emails via Resend API
type ResendEmailService struct {
	apiKey string
	logger *zap.Logger
}

// ElasticEmailService sends emails via Elastic Email API
type ElasticEmailService struct {
	apiKey string
	logger *zap.Logger
}

// NewEmailService creates an email service based on configuration
func NewEmailService(supabaseURL, supabaseServiceKey string, logger *zap.Logger) EmailService {
	// Check which email provider to use
	emailProvider := EmailProvider(os.Getenv("EMAIL_PROVIDER"))
	if emailProvider == "" {
		// Default to Resend if RESEND_API_KEY is set, otherwise Supabase
		if os.Getenv("RESEND_API_KEY") != "" {
			emailProvider = EmailProviderResend
		} else if os.Getenv("ELASTIC_EMAIL_API_KEY") != "" {
			emailProvider = EmailProviderElastic
		} else {
			emailProvider = EmailProviderSupabase
		}
	}

	switch emailProvider {
	case EmailProviderResend:
		apiKey := os.Getenv("RESEND_API_KEY")
		if apiKey == "" {
			logger.Warn("EMAIL_PROVIDER=resend but RESEND_API_KEY not set, falling back to Supabase")
			return &SupabaseEmailService{
				supabaseURL: supabaseURL,
				serviceKey:  supabaseServiceKey,
				logger:      logger,
			}
		}
		logger.Info("Using Resend for sending emails")
		return &ResendEmailService{
			apiKey: apiKey,
			logger: logger,
		}
	case EmailProviderElastic:
		apiKey := os.Getenv("ELASTIC_EMAIL_API_KEY")
		if apiKey == "" {
			logger.Warn("EMAIL_PROVIDER=elastic but ELASTIC_EMAIL_API_KEY not set, falling back to Supabase")
			return &SupabaseEmailService{
				supabaseURL: supabaseURL,
				serviceKey:  supabaseServiceKey,
				logger:      logger,
			}
		}
		logger.Info("Using Elastic Email for sending emails")
		return &ElasticEmailService{
			apiKey: apiKey,
			logger: logger,
		}
	case EmailProviderNone:
		logger.Info("Email service disabled (EMAIL_PROVIDER=none)")
		return &NoOpEmailService{logger: logger}
	default:
		logger.Info("Using Supabase for sending emails (configure Resend SMTP in Supabase Dashboard for production)")
		return &SupabaseEmailService{
			supabaseURL: supabaseURL,
			serviceKey:  supabaseServiceKey,
			logger:      logger,
		}
	}
}

// NoOpEmailService is a no-op email service for testing/development
type NoOpEmailService struct {
	logger *zap.Logger
}

func (s *NoOpEmailService) SendTeamInvite(email, inviteURL string, isNewUser bool) error {
	s.logger.Info("Email service disabled - invite URL", zap.String("email", email), zap.String("invite_url", inviteURL))
	return nil
}

// SendTeamInvite sends a team invite email via Supabase Admin API
func (s *SupabaseEmailService) SendTeamInvite(email, inviteURL string, isNewUser bool) error {
	// For Supabase, we need to get the user ID first
	// Query for user by email
	checkURL := fmt.Sprintf("%s/auth/v1/admin/users?email=%s", s.supabaseURL, email)
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("apikey", s.serviceKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.serviceKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to check user: status %d, body: %s", resp.StatusCode, string(body))
	}

	var usersResponse struct {
		Users []struct {
			ID string `json:"id"`
		} `json:"users"`
	}

	if err := json.Unmarshal(body, &usersResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(usersResponse.Users) == 0 {
		return fmt.Errorf("user not found: %s", email)
	}

	userID := usersResponse.Users[0].ID

	if isNewUser {
		// Send magic link for new users
		magicLinkURL := fmt.Sprintf("%s/auth/v1/admin/users/%s/generate_link", s.supabaseURL, userID)
		magicLinkPayload := map[string]interface{}{
			"type":         "magiclink",
			"redirect_to":  inviteURL,
		}

		jsonData, err := json.Marshal(magicLinkPayload)
		if err != nil {
			return fmt.Errorf("failed to marshal magic link payload: %w", err)
		}

		req, err = http.NewRequest("POST", magicLinkURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create magic link request: %w", err)
		}
		req.Header.Set("apikey", s.serviceKey)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.serviceKey))
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to generate magic link: %w", err)
		}
		defer resp.Body.Close()

		s.logger.Info("Magic link generated for new user", zap.String("email", email))
		return nil
	}

	// For existing users, send invite email
	inviteURLAPI := fmt.Sprintf("%s/auth/v1/admin/users/%s/invite", s.supabaseURL, userID)
	invitePayload := map[string]interface{}{
		"data": map[string]interface{}{
			"invite_url": inviteURL,
			"team_invite": true,
		},
		"redirect_to": inviteURL,
	}

	jsonData, err := json.Marshal(invitePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal invite payload: %w", err)
	}

	req, err = http.NewRequest("POST", inviteURLAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create invite request: %w", err)
	}
	req.Header.Set("apikey", s.serviceKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.serviceKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send invite: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read invite response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("invite endpoint returned error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	s.logger.Info("Team invite email sent via Supabase", zap.String("email", email))
	return nil
}

// SendTeamInvite sends a team invite email via Resend API
func (s *ResendEmailService) SendTeamInvite(email, inviteURL string, isNewUser bool) error {
	fromEmail := os.Getenv("EMAIL_FROM_ADDRESS")
	if fromEmail == "" {
		fromEmail = "noreply@mail.barracudaseo.com" // Default to verified domain
	}

	// Determine subject and content based on whether user is new
	subject := "You've been invited to join a team"
	if isNewUser {
		subject = "Welcome! You've been invited to join a team"
	}

	// HTML email content
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Team Invitation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<div style="background-color: #f8f9fa; padding: 30px; border-radius: 8px;">
		<h1 style="color: #2563eb; margin-top: 0;">%s</h1>
		<p>You've been invited to join a team on Barracuda SEO.</p>
		%s
		<div style="margin: 30px 0;">
			<a href="%s" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; font-weight: bold;">Accept Invitation</a>
		</div>
		<p style="color: #666; font-size: 14px; margin-top: 30px;">
			Or copy and paste this link into your browser:<br>
			<code style="background-color: #f1f3f5; padding: 4px 8px; border-radius: 4px; word-break: break-all;">%s</code>
		</p>
		<p style="color: #666; font-size: 12px; margin-top: 30px; border-top: 1px solid #e5e7eb; padding-top: 20px;">
			If you didn't expect this invitation, you can safely ignore this email.
		</p>
	</div>
</body>
</html>
`, subject, 
		func() string {
			if isNewUser {
				return "<p>Click the button below to create your account and join the team.</p>"
			}
			return "<p>Click the button below to accept the invitation and join the team.</p>"
		}(),
		inviteURL, inviteURL)

	// Plain text content
	textContent := fmt.Sprintf(`%s

You've been invited to join a team on Barracuda SEO.

%s

Accept your invitation by clicking this link:
%s

If you didn't expect this invitation, you can safely ignore this email.
`, subject,
		func() string {
			if isNewUser {
				return "Click the link below to create your account and join the team."
			}
			return "Click the link below to accept the invitation and join the team."
		}(),
		inviteURL)

	// Resend API endpoint
	apiURL := "https://api.resend.com/emails"

	// Build request payload
	payload := map[string]interface{}{
		"from":    fromEmail,
		"to":      []string{email},
		"subject": subject,
		"html":    htmlContent,
		"text":    textContent,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend API returned error: status %d, body: %s", resp.StatusCode, string(body))
	}

	s.logger.Info("Team invite email sent via Resend",
		zap.String("email", email),
		zap.Bool("is_new_user", isNewUser))

	return nil
}

// SendTeamInvite sends a team invite email via Elastic Email API
func (s *ElasticEmailService) SendTeamInvite(email, inviteURL string, isNewUser bool) error {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		// Default to localhost for local development, production URL otherwise
		// Check if running locally (common indicators)
		if os.Getenv("PORT") == "" || os.Getenv("PORT") == "8080" {
			appURL = "http://localhost:5173"
		} else {
			appURL = "https://app.barracudaseo.com"
		}
	}

	// Determine subject and content based on whether user is new
	subject := "You've been invited to join a team"
	if isNewUser {
		subject = "Welcome! You've been invited to join a team"
	}

	// HTML email content
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Team Invitation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<div style="background-color: #f8f9fa; padding: 30px; border-radius: 8px;">
		<h1 style="color: #2563eb; margin-top: 0;">%s</h1>
		<p>You've been invited to join a team on Barracuda SEO.</p>
		%s
		<div style="margin: 30px 0;">
			<a href="%s" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; font-weight: bold;">Accept Invitation</a>
		</div>
		<p style="color: #666; font-size: 14px; margin-top: 30px;">
			Or copy and paste this link into your browser:<br>
			<code style="background-color: #f1f3f5; padding: 4px 8px; border-radius: 4px; word-break: break-all;">%s</code>
		</p>
		<p style="color: #666; font-size: 12px; margin-top: 30px; border-top: 1px solid #e5e7eb; padding-top: 20px;">
			If you didn't expect this invitation, you can safely ignore this email.
		</p>
	</div>
</body>
</html>
`, subject, 
		func() string {
			if isNewUser {
				return "<p>Click the button below to create your account and join the team.</p>"
			}
			return "<p>Click the button below to accept the invitation and join the team.</p>"
		}(),
		inviteURL, inviteURL)

	// Plain text content
	textContent := fmt.Sprintf(`%s

You've been invited to join a team on Barracuda SEO.

%s

Accept your invitation by clicking this link:
%s

If you didn't expect this invitation, you can safely ignore this email.
`, subject,
		func() string {
			if isNewUser {
				return "Click the link below to create your account and join the team."
			}
			return "Click the link below to accept the invitation and join the team."
		}(),
		inviteURL)

	// Elastic Email API endpoint
	apiURL := "https://api.elasticemail.com/v4/emails/transactional"

	// Build request payload
	payload := map[string]interface{}{
		"Recipients": map[string]interface{}{
			"To": []string{email},
		},
		"Content": map[string]interface{}{
			"Body": []map[string]interface{}{
				{
					"ContentType": "HTML",
					"Content":     htmlContent,
				},
				{
					"ContentType": "PlainText",
					"Content":     textContent,
				},
			},
			"Subject": subject,
			"From":    os.Getenv("EMAIL_FROM_ADDRESS"), // e.g., "noreply@mail.barracudaseo.com"
		},
		"Options": map[string]interface{}{
			"ChannelName": "Team Invites",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-ElasticEmail-ApiKey", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("elastic email API returned error: status %d, body: %s", resp.StatusCode, string(body))
	}

	s.logger.Info("Team invite email sent via Elastic Email",
		zap.String("email", email),
		zap.Bool("is_new_user", isNewUser))

	return nil
}

