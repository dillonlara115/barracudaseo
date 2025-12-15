# Resend SMTP Setup - Quick Guide

## Problem
Getting rate limited: `email rate limit exceeded` (429 error)

## Solution
Configure Resend SMTP in Supabase to bypass default email limits.

---

## Step-by-Step Setup

### Step 1: Get Resend API Key

1. **Sign up/Login**: https://resend.com
2. **Go to**: https://resend.com/api-keys
3. **Create API Key**:
   - Name: `Supabase Auth`
   - Permission: `Sending access`
   - Copy the API key (starts with `re_...`)

---

### Step 2: Configure Domain (Optional but Recommended)

For production, you should use your own domain:

1. **Go to**: https://resend.com/domains
2. **Add Domain**: `mail.barracudaseo.com` (or your preferred subdomain)
3. **Add DNS Records**:
   - **SPF** (TXT):
     ```
     v=spf1 include:resend.com ~all
     ```
   - **DKIM** (TXT): Provided by Resend
   - **DMARC** (TXT):
     ```
     v=DMARC1; p=none; rua=mailto:admin@barracudaseo.com
     ```
4. **Wait for verification** (usually 5-10 minutes)

**Note**: You can use Resend's default domain (`onboarding@resend.dev`) for testing, but it's not recommended for production.

---

### Step 3: Configure Supabase Production Dashboard

1. **Go to**: [Supabase Dashboard](https://app.supabase.com) → Your Project
2. **Navigate to**: **Project Settings** → **Auth** → **Email**
3. **Scroll to**: **SMTP Settings**
4. **Enable SMTP**: Toggle **"Enable Custom SMTP"** ON
5. **Fill in**:
   ```
   SMTP Host: smtp.resend.com
   SMTP Port: 465
   SMTP User: resend
   SMTP Password: re_YOUR_API_KEY_HERE
   Sender Email: noreply@mail.barracudaseo.com
   Sender Name: Barracuda SEO
   ```
6. **Click**: "Save"

---

### Step 4: Test Email Delivery

1. **Go to**: **Authentication** → **Email Templates**
2. **Click**: "Send test email"
3. **Enter**: Your email address
4. **Click**: "Send"
5. **Check inbox** - Should arrive within seconds!

---

### Step 5: Update Rate Limits (Optional)

With Resend configured, you can increase rate limits:

1. **Go to**: **Authentication** → **Rate Limits**
2. **Email Sent**: Increase to `30 per hour` (or higher)
3. **Save**

**Note**: Resend has generous limits (100 emails/day free, 50k/month on paid plans), so you can safely increase Supabase rate limits.

---

## Verification Checklist

After setup:

- [ ] Resend API key created
- [ ] Domain added to Resend (optional)
- [ ] DNS records added (if using custom domain)
- [ ] SMTP configured in Supabase dashboard
- [ ] Test email sent successfully
- [ ] Rate limits updated (optional)

---

## Troubleshooting

### Issue: "SMTP authentication failed"

**Causes**:
- Wrong API key
- API key doesn't have sending permissions
- Wrong SMTP host/port

**Fix**:
1. Double-check API key in Resend dashboard
2. Verify SMTP settings match exactly:
   - Host: `smtp.resend.com`
   - Port: `465`
   - User: `resend`
   - Password: Your API key (starts with `re_`)

### Issue: "Domain not verified"

**Cause**: Using custom domain without DNS setup

**Fix**:
1. Use Resend's default domain for testing: `onboarding@resend.dev`
2. Or complete DNS setup for your domain

### Issue: Emails still rate limited

**Causes**:
- SMTP not enabled
- Rate limits not updated
- Using default Supabase email (not Resend)

**Fix**:
1. Verify SMTP is enabled (green toggle)
2. Check rate limits in Authentication → Rate Limits
3. Send test email to confirm Resend is being used

---

## Resend Limits

### Free Plan
- **100 emails/day**
- **3,000 emails/month**
- Perfect for development/testing

### Paid Plans
- **50,000 emails/month** ($20/month)
- **Unlimited** on higher tiers

**For production**: Free plan should be fine initially, upgrade when needed.

---

## Quick Reference

### Supabase SMTP Settings
```
Host: smtp.resend.com
Port: 465
User: resend
Password: re_YOUR_API_KEY
Sender: noreply@mail.barracudaseo.com
```

### Resend Dashboard
- **API Keys**: https://resend.com/api-keys
- **Domains**: https://resend.com/domains
- **Logs**: https://resend.com/emails (see delivery status)

---

## Next Steps

After configuring Resend:

1. ✅ **Test magic link** - Should work without rate limits
2. ✅ **Monitor Resend dashboard** - Check email delivery logs
3. ✅ **Update rate limits** - Increase Supabase limits
4. ✅ **Customize email template** - Brand your magic link emails

---

**Status**: Ready to configure! Follow steps above to set up Resend SMTP.

