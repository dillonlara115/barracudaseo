# Magic Link Authentication - Implementation Guide

## Overview

Barracuda now uses **Magic Link** authentication as the primary login method, following modern SaaS best practices (like Slack, Notion, Linear). This provides:

- ✅ **Better Security** - No passwords to forget, leak, or reuse
- ✅ **Simpler UX** - One-click login from email
- ✅ **Lower Support Costs** - No "forgot password" flows needed
- ✅ **Mobile-Friendly** - Easier on mobile devices
- ✅ **7-Day Sessions** - Users stay logged in for a week

## User Flow

### Sign Up (New Users)
1. User enters: Email, First Name, Last Name
2. Clicks "Send magic link"
3. Receives email with magic link
4. Clicks link → Auto-logged in → Redirected to dashboard

### Sign In (Existing Users)
1. User enters: Email
2. Clicks "Send magic link"
3. Receives email with magic link
4. Clicks link → Auto-logged in → Redirected to dashboard

### Optional: Password Login
- Users can still use password login if they have one
- Click "Use password instead" on login page
- Password reset is effectively obsolete (just send another magic link)

## Configuration Changes

### Local Development (supabase/config.toml)

Key changes made:
```toml
[auth]
# Extended session duration to 7 days (604800 seconds)
jwt_expiry = 604800

# Added hash-based routes to redirect URLs
additional_redirect_urls = [
  "https://127.0.0.1:3000",
  "http://127.0.0.1:3000/#/",
  "http://localhost:3000",
  "http://localhost:3000/#/"
]

[auth.email]
# Increased throttling for magic link requests (60 seconds between requests)
max_frequency = "60s"

# Magic link expires after 1 hour
otp_expiry = 3600
```

### Production Supabase Dashboard Setup

⚠️ **IMPORTANT**: You must configure these settings in your production Supabase dashboard:

1. **Navigate to**: Supabase Dashboard → Authentication → URL Configuration

2. **Site URL**: Set to your production domain
   ```
   https://your-production-domain.com
   ```

3. **Redirect URLs**: Add these URLs to the allowlist:
   ```
   https://your-production-domain.com
   https://your-production-domain.com/#/
   https://www.your-production-domain.com
   https://www.your-production-domain.com/#/
   ```

4. **Navigate to**: Authentication → Settings

5. **JWT Expiry**: Set to `604800` (7 days in seconds)

6. **Enable Email Auth**: Make sure "Enable Email Signup" is ON

7. **Email Templates** (optional but recommended):
   - Customize the "Magic Link" email template
   - Use your brand colors and messaging
   - Keep the `{{ .ConfirmationURL }}` variable intact

8. **Rate Limiting**: Configure in Authentication → Rate Limits
   ```
   Email Sent: 3 per hour (prevents abuse)
   Token Verifications: 30 per 5 minutes
   ```

## Testing Instructions

### Local Testing

1. **Start Supabase locally**:
   ```bash
   cd /home/dillon/Sites/barracuda
   supabase start
   ```

2. **Start the web app**:
   ```bash
   cd web
   npm run dev
   ```

3. **Test Sign Up Flow**:
   - Navigate to `http://localhost:3000`
   - Enter test email (e.g., `test@example.com`)
   - Enter first/last name
   - Click "Send magic link"
   - Check Inbucket (Supabase email testing): `http://localhost:54324`
   - Click the magic link in the email
   - Should redirect to dashboard, logged in

4. **Test Sign In Flow**:
   - Sign out
   - Navigate back to login
   - Enter same email
   - Click "Send magic link"
   - Check Inbucket for new magic link
   - Click link → Should auto-login

5. **Test Password Login (Optional)**:
   - Sign out
   - Click "Use password instead"
   - Enter email + password (if user has one set)
   - Should sign in with password

### Production Testing

1. **Deploy changes** to your production environment

2. **Test with real email**:
   - Use a real email address you have access to
   - Complete sign-up flow
   - Verify magic link email arrives (check spam folder)
   - Click link and verify redirect works

3. **Test redirect URLs**:
   - Ensure magic link redirects to `https://yourdomain.com/#/`
   - User should be logged in automatically
   - Session should persist for 7 days

4. **Test rate limiting**:
   - Request multiple magic links quickly
   - Should be rate-limited after 60 seconds (per config)

## Troubleshooting

### Issue: Magic link email not received

**Causes:**
- Email not configured in Supabase production
- SMTP settings incorrect
- Rate limiting hit (too many requests)

**Solutions:**
1. Check Supabase → Project Settings → Auth → SMTP settings
2. Verify email provider is working (send test email from dashboard)
3. Check user's spam folder
4. Wait 60 seconds between magic link requests

### Issue: Magic link redirects to wrong URL

**Causes:**
- Redirect URLs not configured in Supabase dashboard
- Hash routing not included in allowlist

**Solutions:**
1. Add `https://yourdomain.com/#/` to Supabase redirect URLs
2. Make sure to include the `/#/` hash for SPA routing
3. Clear browser cache and test again

### Issue: User not staying logged in

**Causes:**
- JWT expiry not set to 7 days in production
- Browser blocking cookies
- User clearing cache frequently

**Solutions:**
1. Verify `jwt_expiry = 604800` in production Supabase settings
2. Check browser cookie settings (should allow first-party cookies)
3. Test in incognito mode to rule out extensions

### Issue: "Invalid redirect URL" error

**Causes:**
- Production URL not in Supabase allowlist
- Typo in redirect URL configuration

**Solutions:**
1. Double-check Supabase → Authentication → URL Configuration
2. Make sure Site URL matches exactly
3. Include both with and without hash: `/` and `/#/`
4. Include www and non-www versions if applicable

## Migration Notes

### Existing Users with Passwords

- Users with existing passwords can still use them
- They can click "Use password instead" on login
- Gradually encourage migration to magic links
- Consider deprecating password auth in future

### Data Impact

- No database schema changes required
- Existing user accounts work unchanged
- `auth.users` table remains the same
- Only authentication method changes

## Email Templates

### Recommended Magic Link Email Content

**Subject**: Your magic link to Barracuda

**Body**:
```html
<h2>Sign in to Barracuda</h2>
<p>Click the button below to sign in to your Barracuda account:</p>

<a href="{{ .ConfirmationURL }}" 
   style="background: #8ec07c; color: #3c3836; padding: 12px 24px; 
          border-radius: 8px; text-decoration: none; display: inline-block;">
  Sign in to Barracuda
</a>

<p>Or copy and paste this URL into your browser:</p>
<p>{{ .ConfirmationURL }}</p>

<p><small>This link expires in 1 hour. If you didn't request this, you can safely ignore this email.</small></p>
```

## Security Considerations

### Magic Link Security

- **Expires after 1 hour** - `otp_expiry = 3600`
- **Rate limited** - Max 1 request per 60 seconds per email
- **Single use** - Link becomes invalid after first use
- **HTTPS only** - Links only work over secure connections
- **Token in URL** - Consider moving to PKCE flow for even better security (future enhancement)

### Session Security

- **7-day expiry** - Balances UX and security
- **Auto-refresh** - Tokens refresh automatically before expiry
- **Revocable** - Admin can revoke sessions from Supabase dashboard
- **Device tracking** - Supabase tracks sessions per device

## Analytics & Monitoring

### Key Metrics to Track

1. **Magic Link Success Rate**:
   - % of magic links clicked vs sent
   - Time between send and click

2. **Authentication Method Usage**:
   - Magic link vs password login ratio
   - Track in your analytics platform

3. **Email Deliverability**:
   - Monitor bounce rates
   - Check spam reports
   - Use SendGrid/Resend delivery metrics

4. **Session Duration**:
   - How long users stay logged in
   - Re-authentication frequency

### Supabase Dashboard Metrics

Check: **Authentication → Logs** for:
- `user.otp.sent` - Magic link emails sent
- `user.signin` - Successful sign-ins
- `token.refreshed` - Session refreshes

## Future Enhancements

### Potential Improvements

1. **PKCE Flow** - Move from URL-based magic links to OTP codes for even better security

2. **Biometric Auth** - Add WebAuthn/Passkeys for password-less, device-based auth

3. **SMS Magic Links** - Offer phone number login for mobile users

4. **Social OAuth** - Add Google/GitHub login options

5. **Remember Device** - Extend session to 30 days for trusted devices

6. **Email Verification** - Add badge to verified email addresses

## Support & Documentation

### User-Facing Documentation

Update your help docs to explain:
- How magic links work
- Where to find the email (check spam)
- What to do if link expires (request new one)
- How to set optional password (Settings page)

### Support Ticket Categories

Common issues to prepare for:
- "I didn't receive the magic link email"
- "The magic link expired"
- "I want to use a password instead"
- "How do I stay logged in longer?"

## Rollback Plan

If magic links cause issues:

1. **Immediate rollback** (code-level):
   ```bash
   git revert <commit-hash>
   ```

2. **Keep both methods** (hybrid approach):
   - Make password the default again
   - Keep magic link as "Use magic link instead" option

3. **Database rollback**: Not needed (no schema changes)

4. **User communication**: Email users about reverting to password auth

---

## Quick Reference

### Local Development URLs
- App: `http://localhost:3000`
- Supabase Studio: `http://localhost:54323`
- Email Testing (Inbucket): `http://localhost:54324`

### Production URLs (Update with your domain)
- App: `https://your-domain.com`
- Supabase Dashboard: `https://app.supabase.com/project/YOUR_PROJECT_ID`

### Key Files Modified
- `web/src/lib/auth.js` - Added magic link functions
- `web/src/components/Auth.svelte` - New UI with magic link primary
- `supabase/config.toml` - Updated JWT expiry and rate limits
- `web/src/routes/ResetPassword.svelte` - Updated for magic link flow

---

**Last Updated**: December 11, 2024
**Contact**: Update this with your team's contact info


