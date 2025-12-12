# PKCE Magic Link Setup - Complete Guide

## Overview

You've successfully configured the **PKCE (Proof Key for Code Exchange) flow** for magic links, which is more secure than the implicit flow. This is the recommended approach per [Supabase documentation](https://supabase.com/docs/guides/auth/auth-email-passwordless).

## How It Works

### Traditional Implicit Flow (Less Secure)
```
Email → Click Link → Redirect with tokens in URL hash
→ https://app.barracudaseo.com#access_token=XXX&refresh_token=YYY
```
**Issue**: Tokens exposed in URL (visible in browser history, logs, etc.)

### PKCE Flow (More Secure) ✅
```
Email → Click Link → Redirect with token_hash
→ https://app.barracudaseo.com/#/auth/confirm?token_hash=XXX&type=email
→ App exchanges token_hash for session server-side
→ Redirect to dashboard with session established
```
**Benefit**: No sensitive tokens exposed in URL

---

## Configuration Checklist

### ✅ Step 1: Supabase Email Template

You've already done this! Your magic link email template now uses:

```html
<a href="{{ .SiteURL }}/auth/confirm?token_hash={{ .TokenHash }}&type=email">Log In</a>
```

This sends the token hash instead of full tokens.

### ✅ Step 2: Auth Confirm Route

The `/auth/confirm` route has been created at:
- `web/src/routes/AuthConfirm.svelte`
- Added to App.svelte routes
- Marked as public page (accessible without login)

This route:
1. Extracts `token_hash` and `type` from URL
2. Calls `supabase.auth.verifyOtp()` to exchange for session
3. Redirects to dashboard on success

### ✅ Step 3: Redirect URL Configuration

**In Supabase Dashboard:**

Go to: **Authentication → URL Configuration**

**Site URL:**
```
https://app.barracudaseo.com
```

**Redirect URLs (add all):**
```
https://app.barracudaseo.com/auth/confirm
https://app.barracudaseo.com/#/auth/confirm
http://localhost:3000/auth/confirm
http://localhost:3000/#/auth/confirm
http://127.0.0.1:3000/auth/confirm
```

**Why both formats?**
- Without hash: `/auth/confirm` - Supabase redirect target
- With hash: `/#/auth/confirm` - SPA routing format

---

## Testing the PKCE Flow

### Local Testing

1. **Start Supabase**:
   ```bash
   supabase start
   ```

2. **Update local Supabase config** (`supabase/config.toml`):
   ```toml
   [auth]
   site_url = "http://127.0.0.1:3000"
   additional_redirect_urls = [
     "http://127.0.0.1:3000/auth/confirm",
     "http://127.0.0.1:3000/#/auth/confirm"
   ]
   ```

3. **Restart Supabase** after config change:
   ```bash
   supabase stop
   supabase start
   ```

4. **Update email template in Supabase Studio**:
   - Go to: http://localhost:54323
   - Authentication → Email Templates → Magic Link
   - Update to use `token_hash` format

5. **Start web app**:
   ```bash
   cd web
   npm run dev
   ```

6. **Test magic link**:
   - Go to http://localhost:3000
   - Request magic link
   - Check Inbucket: http://localhost:54324
   - Click the link
   - Should redirect to `/#/auth/confirm?token_hash=...`
   - Should show "Signing you in..." → "Successfully signed in!"
   - Should redirect to dashboard

### Production Testing

1. **Deploy the code**:
   ```bash
   git add .
   git commit -m "feat: implement PKCE flow for magic links"
   git push origin main
   ```

2. **Configure Supabase Production**:
   - Site URL: `https://app.barracudaseo.com`
   - Add redirect URLs (see Step 3 above)

3. **Test with real email**:
   - Request magic link
   - Click link in email
   - Check browser console for logs
   - Should see successful verification

---

## What the Console Should Show

### When requesting magic link:
```
🔍 Requesting magic link for: user@example.com
🔍 Redirect URL: https://app.barracudaseo.com
🔍 Current origin: https://app.barracudaseo.com
✅ Magic link sent successfully
```

### When clicking magic link:
```
🔍 Auth confirm page loaded
🔍 Full URL: https://app.barracudaseo.com/#/auth/confirm?token_hash=...&type=email
🔍 Hash: #/auth/confirm?token_hash=...&type=email
🔍 Token hash: present
🔍 Type: email
🔄 Verifying OTP with token hash...
✅ OTP verified successfully: { user: {...}, session: {...} }
🔄 Redirecting to dashboard...
```

---

## Troubleshooting

### Issue: "Missing token_hash or type parameter"

**Cause**: URL doesn't contain the required parameters

**Check**:
1. Email template uses `{{ .TokenHash }}` (not `{{ .Token }}`)
2. Email link goes to `/auth/confirm` endpoint
3. URL in email looks like: `https://app.barracudaseo.com/auth/confirm?token_hash=...`

### Issue: "Invalid redirect URL"

**Cause**: `/auth/confirm` not in Supabase allowlist

**Fix**:
1. Add to Supabase Dashboard → URL Configuration → Redirect URLs
2. Include both `/auth/confirm` and `/#/auth/confirm`
3. Save and test again

### Issue: OTP verification fails

**Causes**:
- Link already used (one-time use only)
- Link expired (1 hour default)
- Token hash invalid

**Solutions**:
1. Request a new magic link
2. Don't click link multiple times
3. Click link within 1 hour

### Issue: Blank page at /auth/confirm

**Cause**: Route not properly configured

**Check**:
1. `AuthConfirm.svelte` exists in `web/src/routes/`
2. Route added to `App.svelte` routes object
3. Route marked as public page

---

## Security Benefits of PKCE

1. **No tokens in URL**: Token hash is not sensitive, can't be used directly
2. **Server-side exchange**: Tokens only exist in memory after verification
3. **One-time use**: Token hash can't be reused
4. **Short-lived**: Expires after 1 hour
5. **Browser history safe**: No sensitive data in URL history

---

## Comparison: Implicit vs PKCE

### Implicit Flow
```javascript
// Email link:
https://app.barracudaseo.com#access_token=eyJ...&refresh_token=abc123

// Tokens directly in URL ❌
// Visible in:
// - Browser history
// - Server logs
// - Referrer headers
// - Browser extensions
```

### PKCE Flow ✅
```javascript
// Email link:
https://app.barracudaseo.com/auth/confirm?token_hash=xyz789&type=email

// Token hash is NOT a usable token ✅
// Must be exchanged server-side:
const { data } = await supabase.auth.verifyOtp({
  token_hash: 'xyz789',
  type: 'email'
});
// Now have session with tokens
```

---

## Migration Notes

### Existing Users
- Old magic links (implicit flow) will stop working
- Users need to request new magic links
- Session cookies still valid (no forced logout)

### Email Template
- **Before**: `{{ .ConfirmationURL }}`
- **After**: `{{ .SiteURL }}/auth/confirm?token_hash={{ .TokenHash }}&type=email`

### Redirect URLs
- **Before**: `https://app.barracudaseo.com` (root)
- **After**: `https://app.barracudaseo.com/auth/confirm` (specific endpoint)

---

## Next Steps

1. ✅ Code deployed
2. ✅ Email template updated
3. ⏳ **Update Supabase redirect URLs** (see Step 3 above)
4. ⏳ **Test with real email**
5. ⏳ **Monitor console logs**

---

## Support

If you encounter issues:

1. **Check console logs** - All debug info logged
2. **Check Supabase Auth Logs** - Dashboard → Logs
3. **Verify email template** - Should use `TokenHash` not `Token`
4. **Verify redirect URLs** - Must include `/auth/confirm`

---

**Last Updated**: December 11, 2024
**PKCE Documentation**: https://supabase.com/docs/guides/auth/auth-email-passwordless
