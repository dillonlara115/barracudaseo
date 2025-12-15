# Production Magic Link Setup - Quick Reference

## ⚠️ CRITICAL: Production Supabase Configuration

Before deploying magic link authentication to production, you **MUST** configure these settings in your Supabase dashboard.

---

## 1. Authentication URLs

**Dashboard Path**: `Project Settings → Authentication → URL Configuration`

### Site URL
Set to your production domain (without trailing slash):
```
https://barracudaseo.com
```

### Redirect URLs (Allowlist)
Add ALL of these URLs:
```
https://barracudaseo.com
https://barracudaseo.com/#/
https://www.barracudaseo.com
https://www.barracudaseo.com/#/
```

> **Why?**: Your app uses hash-based routing (`/#/`), and magic links need to redirect to these exact URLs.

---

## 2. JWT Settings

**Dashboard Path**: `Authentication → Settings → JWT Settings`

### JWT Expiry
```
604800
```
(7 days in seconds)

### Enable Refresh Token Rotation
```
✓ Enabled
```

### Refresh Token Reuse Interval
```
10
```
(seconds)

---

## 3. Email Provider

**Dashboard Path**: `Project Settings → Auth → Email`

### Choose ONE:

#### Option A: Resend (Recommended)
```
SMTP Host: smtp.resend.com
Port: 465
Username: resend
Password: YOUR_RESEND_API_KEY
Sender Email: noreply@mail.barracudaseo.com
Sender Name: Barracuda SEO
```

Get API key: https://resend.com/api-keys

#### Option B: SendGrid
```
SMTP Host: smtp.sendgrid.net
Port: 587
Username: apikey
Password: YOUR_SENDGRID_API_KEY
Sender Email: noreply@mail.barracudaseo.com
Sender Name: Barracuda SEO
```

Get API key: https://app.sendgrid.com/settings/api_keys

---

## 4. Rate Limiting

**Dashboard Path**: `Authentication → Rate Limits`

### Email Sent
```
3 per hour
```

### Token Verifications
```
30 per 5 minutes
```

### Sign-in / Sign-ups
```
30 per 5 minutes
```

---

## 5. Email Templates

**Dashboard Path**: `Authentication → Email Templates → Magic Link`

### Subject
```
Your magic link to Barracuda
```

### Email Body (HTML)
```html
<h2>Sign in to Barracuda</h2>
<p>Click the button below to sign in to your Barracuda account:</p>

<a href="{{ .ConfirmationURL }}" 
   style="background: #8ec07c; 
          color: #3c3836; 
          padding: 12px 24px; 
          border-radius: 8px; 
          text-decoration: none; 
          display: inline-block; 
          font-weight: 600;">
  Sign in to Barracuda
</a>

<p>Or copy and paste this URL into your browser:</p>
<p style="word-break: break-all; color: #666;">{{ .ConfirmationURL }}</p>

<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;">

<p style="color: #666; font-size: 12px;">
  This link expires in 1 hour. If you didn't request this, you can safely ignore this email.
</p>
```

> **Important**: Keep the `{{ .ConfirmationURL }}` variable intact!

---

## 6. Test Email Delivery

After configuration, test that emails are being sent:

1. **Send Test Email** from Supabase dashboard:
   - Go to Authentication → Email Templates
   - Click "Send test email"
   - Enter your email
   - Verify receipt

2. **Test Magic Link Flow**:
   - Go to your production app
   - Enter email and request magic link
   - Check inbox (and spam folder)
   - Click link and verify redirect

---

## 7. DNS Records (For Custom Email Domain)

If using a custom domain like `mail.barracudaseo.com`:

### For Resend:
Add these DNS records in your domain registrar:

**SPF Record** (TXT):
```
Name: @
Value: v=spf1 include:resend.com ~all
```

**DKIM Record** (TXT):
```
Name: resend._domainkey
Value: [Provided by Resend]
```

**DMARC Record** (TXT):
```
Name: _dmarc
Value: v=DMARC1; p=none; rua=mailto:admin@barracudaseo.com
```

### For SendGrid:
Follow SendGrid's domain authentication wizard:
https://app.sendgrid.com/settings/sender_auth

---

## 8. Environment Variables

Make sure these are set in your **production environment** (Vercel/Netlify/etc.):

```bash
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
VITE_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
VITE_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

---

## Verification Checklist

Before going live, verify:

- [ ] Site URL set to production domain
- [ ] All redirect URLs added (with and without `/#/`)
- [ ] JWT expiry set to 604800 (7 days)
- [ ] Email provider configured (Resend or SendGrid)
- [ ] Test email successfully sent and received
- [ ] Magic link successfully redirects to production app
- [ ] User stays logged in after redirect
- [ ] Rate limiting configured
- [ ] Email template customized with branding
- [ ] DNS records added (if using custom email domain)
- [ ] Environment variables set in production

---

## Quick Test Commands

### Test from local to production Supabase:
```bash
# Update .env.local with production Supabase credentials
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-production-anon-key

# Run local dev server
cd web
npm run dev

# Test magic link flow
# Open http://localhost:3000 and request magic link
```

### Test production deployment:
```bash
# After deploying to Vercel/Netlify
# Visit your production URL
curl https://barracudaseo.com

# Test auth flow in browser
# 1. Enter email
# 2. Check inbox for magic link
# 3. Click link
# 4. Verify redirect to dashboard
```

---

## Troubleshooting

### Magic link email not received
1. Check Supabase logs: `Authentication → Logs`
2. Verify SMTP credentials are correct
3. Check spam/junk folder
4. Test with different email provider (Gmail, Outlook)

### "Invalid redirect URL" error
1. Verify redirect URLs include `/#/` for hash routing
2. Check for typos in URL configuration
3. Make sure URLs match exactly (https vs http)

### User not staying logged in
1. Verify JWT expiry is 604800 in production
2. Check browser isn't blocking cookies
3. Test in incognito mode

---

## Support

If issues persist:
1. Check Supabase logs: `Project → Logs`
2. Contact Supabase support: https://supabase.com/support
3. Check email provider status page:
   - Resend: https://status.resend.com
   - SendGrid: https://status.sendgrid.com

---

**Last Updated**: December 11, 2024


