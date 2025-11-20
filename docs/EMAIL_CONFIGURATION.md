# Email Configuration Guide

This project supports multiple email providers for sending team invite emails. **Resend** is recommended for both development and production.

## Quick Start

### Option 1: Use Resend via Supabase SMTP (Recommended)

Configure Resend SMTP in Supabase Dashboard - works for both local and production.

**Setup Steps:**

1. **Get Resend API Key:**
   - Sign up at https://resend.com
   - Go to **Settings** → **API Keys** → **Create**
   - Copy your API key

2. **Verify your domain** in Resend dashboard (required for production)

3. **Configure Supabase to use Resend SMTP:**
   - Go to your Supabase project → **Project Settings** → **Authentication** tab
   - Find the **SMTP** section and toggle **Enable Custom SMTP**
   - Enter your sender email and name (e.g., `noreply@mail.barracudaseo.com` and `Barracuda SEO`)
   - Enter Resend SMTP credentials:
     - **SMTP Host**: `smtp.resend.com`
     - **SMTP Port**: `465`
     - **SMTP User**: `resend`
     - **SMTP Password**: Your Resend API key
   - Click **Save**

4. **For Local Development:**
   - Emails will be captured by Supabase Inbucket (view at http://localhost:54324)
   - Or configure Resend SMTP in Supabase for real emails in dev too

**No additional environment variables needed** - Supabase handles email sending automatically via Resend SMTP.

### Option 2: Use Resend API Directly

For more control over email templates and sending, use Resend API directly:

1. **Get your Resend API Key:**
   - Sign up at https://resend.com
   - Go to **Settings** → **API Keys** → **Create**
   - Copy your API key

2. **Set Environment Variables:**
   ```bash
   EMAIL_PROVIDER=resend
   RESEND_API_KEY=re_your-api-key-here
   EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com
   APP_URL=https://app.barracudaseo.com
   ```

3. **Verify your domain** in Resend dashboard (required for sending emails)

### Option 3: Use Elastic Email

To use Elastic Email instead of Supabase:

1. **Get your Elastic Email API Key:**
   - Sign up at https://elasticemail.com
   - Go to **Settings** → **Manage API Keys** → **Create**
   - Copy your API key

2. **Set Environment Variables:**
   ```bash
   # Set email provider
   EMAIL_PROVIDER=elastic
   
   # Set Elastic Email API key
   ELASTIC_EMAIL_API_KEY=your-api-key-here
   
   # Set sender email address (must be verified in Elastic Email)
   EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com
   ```

3. **Verify your domain** in Elastic Email dashboard (required for sending emails)

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `EMAIL_PROVIDER` | No | `supabase` | Email provider: `supabase`, `resend`, `elastic`, or `none` |
| `RESEND_API_KEY` | Yes (if using Resend API) | - | Your Resend API key (starts with `re_`) |
| `ELASTIC_EMAIL_API_KEY` | Yes (if using Elastic) | - | Your Elastic Email API key |
| `EMAIL_FROM_ADDRESS` | Yes (if using Resend/Elastic API) | - | Sender email address (must be verified) |
| `APP_URL` | No | Auto-detected | Base URL for invite links. Defaults to `http://localhost:5173` for local dev, `https://app.barracudaseo.com` for production |

## Configuration Examples

### Local Development with Resend (via Supabase SMTP)

Configure Resend SMTP in Supabase Dashboard:
- Emails will be sent via Resend (or captured by Inbucket if SMTP not configured)
- View emails at: http://localhost:54324 (if using Inbucket)
- **APP_URL**: Automatically defaults to `http://localhost:5173` for local development (no need to set)

### Production with Resend (via Supabase SMTP) - Recommended

```bash
# .env (production)
APP_URL=https://app.barracudaseo.com
# Configure Resend SMTP in Supabase Dashboard:
# SMTP Host: smtp.resend.com
# SMTP Port: 465
# SMTP User: resend
# SMTP Password: your-resend-api-key
```

### Production with Resend API Directly

```bash
# .env (production)
EMAIL_PROVIDER=resend
RESEND_API_KEY=re_your-api-key-here
EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com
APP_URL=https://app.barracudaseo.com
```

### Production with Elastic Email

```bash
# .env (production)
EMAIL_PROVIDER=elastic
ELASTIC_EMAIL_API_KEY=your-elastic-api-key
EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com
APP_URL=https://app.barracudaseo.com
```

## How It Works

### Resend via Supabase SMTP Flow (Recommended)

1. **New Users:**
   - User is created in Supabase Auth
   - Magic link is generated via Supabase Admin API
   - Supabase sends email with magic link via Resend SMTP
   - User clicks link → signs up → redirected to accept invite

2. **Existing Users:**
   - Invite email sent via Supabase Admin API `/invite` endpoint
   - Supabase sends email via Resend SMTP
   - User clicks link → redirected to accept invite

### Resend API Direct Flow

1. **New Users:**
   - User is created in Supabase Auth
   - Custom HTML email sent via Resend API
   - Email contains invite link
   - User clicks link → signs up → redirected to accept invite

2. **Existing Users:**
   - Custom HTML email sent via Resend API
   - Email contains invite link
   - User clicks link → redirected to accept invite

### Supabase Email Flow

1. **New Users:**
   - User is created in Supabase Auth
   - Magic link is generated via Supabase Admin API
   - Supabase sends email with magic link (via configured SMTP)
   - User clicks link → signs up → redirected to accept invite

2. **Existing Users:**
   - Invite email sent via Supabase Admin API `/invite` endpoint
   - Supabase sends email (via configured SMTP)
   - User clicks link → redirected to accept invite

### Elastic Email Flow

1. **New Users:**
   - User is created in Supabase Auth
   - Custom HTML email sent via Elastic Email API
   - Email contains invite link
   - User clicks link → signs up → redirected to accept invite

2. **Existing Users:**
   - Custom HTML email sent via Elastic Email API
   - Email contains invite link
   - User clicks link → redirected to accept invite

## Testing

### Disable Email Sending (Development)

Set `EMAIL_PROVIDER=none` to disable email sending (useful for testing):
```bash
EMAIL_PROVIDER=none
```

Emails will be logged but not sent.

### View Supabase Emails (Local)

When using Supabase locally, emails are captured by Inbucket:
- Open: http://localhost:54324
- View all sent emails

## Troubleshooting

### Resend Emails Not Sending

1. **Check SMTP Configuration (if using Supabase SMTP):**
   - Verify SMTP is enabled in Supabase Dashboard
   - Check Resend SMTP credentials are correct:
     - Host: `smtp.resend.com`
     - Port: `465`
     - User: `resend`
     - Password: Your Resend API key
   - Test SMTP connection in Supabase

2. **Check API Key (if using Resend API):**
   - Verify `RESEND_API_KEY` is set correctly (starts with `re_`)
   - Check API key has send permissions

3. **Verify Domain:**
   - Domain must be verified in Resend dashboard (e.g., `mail.barracudaseo.com`)
   - `EMAIL_FROM_ADDRESS` must use the verified domain (e.g., `noreply@mail.barracudaseo.com`)

4. **Check Logs:**
   - Look for "Team invite email sent via Resend" messages
   - Check Resend dashboard for delivery status

### Supabase Emails Not Sending

1. **Check SMTP Configuration:**
   - Verify SMTP is enabled in Supabase Dashboard
   - Check SMTP credentials are correct
   - Test SMTP connection

2. **Check Email Templates:**
   - Go to **Authentication** → **Email Templates**
   - Verify templates are configured correctly

3. **Check Logs:**
   - Look for "Magic link generated" or "Team invite email sent" messages
   - Check for SMTP errors

### Elastic Email Not Sending

1. **Verify API Key:**
   - Check `ELASTIC_EMAIL_API_KEY` is set correctly
   - Verify API key has send permissions

2. **Verify Domain:**
   - Domain must be verified in Elastic Email dashboard
   - `EMAIL_FROM_ADDRESS` must match verified domain

3. **Check API Response:**
   - Look for errors in server logs
   - Check Elastic Email dashboard for delivery status

## References

- [Resend SMTP with Supabase](https://resend.com/docs/send-with-supabase-smtp) - Recommended setup
- [Resend API Documentation](https://resend.com/docs/api-reference/emails/send-email)
- [Supabase SMTP Configuration](https://supabase.com/docs/guides/auth/auth-smtp)
- [Elastic Email API Documentation](https://elasticemail.com/developers/api-documentation)
- [Elastic Email Go Library](https://elasticemail.com/developers/api-libraries/go)

