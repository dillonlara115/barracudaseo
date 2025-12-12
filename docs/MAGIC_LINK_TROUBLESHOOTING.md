# Magic Link Not Logging In - Troubleshooting Guide

## Issue: Magic Link Goes to app.barracudaseo.com But Doesn't Log In

### Quick Diagnostic Steps

1. **Open Browser DevTools Console** (F12) when clicking magic link
2. **Look for these console messages**:
   - "Full URL hash: ..." - Should show access_token
   - "Found auth token in URL" 
   - "Magic link session set successfully"
   
3. **Check the URL after clicking magic link**:
   - Should contain: `#access_token=...&refresh_token=...&type=magiclink`
   - Example: `https://app.barracudaseo.com/#access_token=eyJ...&refresh_token=...&type=magiclink`

---

## Common Causes & Solutions

### Cause 1: Redirect URL Not in Supabase Allowlist

**Symptoms**: 
- Error message about "invalid redirect URL"
- Redirects to wrong domain
- Tokens not present in URL

**Solution**:
1. Go to: [Supabase Dashboard](https://app.supabase.com) → Your Project → Authentication → URL Configuration
2. Check **Site URL** is set to: `https://app.barracudaseo.com`
3. Check **Redirect URLs** includes:
   ```
   https://app.barracudaseo.com
   https://app.barracudaseo.com/#/
   https://app.barracudaseo.com/
   ```

---

### Cause 2: Hash Routing Interference

**Symptoms**:
- URL shows tokens but doesn't log in
- Console shows "Auth callback error"
- Page refreshes but no login

**Solution**: Already fixed in latest code update!
- Updated `App.svelte` to properly extract tokens from hash
- Handles multiple `#` symbols in URL
- Logs diagnostic info to console

**Test**: Refresh the page after clicking magic link

---

### Cause 3: PKCE Flow Not Enabled

**Symptoms**:
- Tokens appear but session doesn't persist
- "Invalid grant" errors in console

**Solution**: Already fixed!
- Updated `supabase.js` to use `flowType: 'pkce'`
- This is the modern, secure flow for magic links

---

### Cause 4: Email Link Format Issue

**Symptoms**:
- Email link goes to wrong URL
- No tokens in URL at all

**Solution**: Check the magic link configuration in your auth code

**Current config** (in `auth.js`):
```javascript
const redirectTo = typeof window !== 'undefined' 
  ? `${window.location.origin}/#/`
  : undefined;
```

This generates: `https://app.barracudaseo.com/#/`

**If this doesn't work**, try without the hash:
```javascript
const redirectTo = typeof window !== 'undefined' 
  ? window.location.origin
  : undefined;
```

---

## Step-by-Step Debugging

### Step 1: Check Magic Link Email

1. Open the magic link email
2. **Right-click** the "Sign in" button
3. **Copy link address**
4. **Paste into text editor**
5. **Verify it looks like**:
   ```
   https://your-project.supabase.co/auth/v1/verify?token=...&type=magiclink&redirect_to=https://app.barracudaseo.com/%23/
   ```

**Expected**:
- Should have `redirect_to=https://app.barracudaseo.com/...`
- The `%23` is URL-encoded `#` (this is correct)

---

### Step 2: Test Magic Link Manually

1. Click magic link in email
2. **Immediately open DevTools Console** (F12)
3. **Look for console logs**:
   ```
   Full URL hash: #access_token=...
   Found auth token in URL
   Auth callback detected: { type: 'magiclink', hasAccessToken: true, hasRefreshToken: true }
   Magic link session set successfully: { ... }
   ```

**If you see an error**, note the exact error message.

---

### Step 3: Check Supabase Dashboard Logs

1. Go to: Supabase Dashboard → Your Project → **Logs**
2. Filter by: **Authentication**
3. Look for recent events:
   - `user.otp.sent` - Magic link sent
   - `user.signin` - User signed in (should appear after clicking link)
   
**If `user.signin` is missing**, the token wasn't processed.

---

### Step 4: Verify Environment Variables

Check your production deployment (Vercel/Netlify) has:

```bash
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

**To verify in browser**:
1. Open DevTools Console
2. Type: `import.meta.env`
3. Check PUBLIC_SUPABASE_URL is correct

---

## Quick Fixes to Try

### Fix 1: Update Supabase Redirect URLs

Add these **exact** URLs to Supabase Dashboard → Authentication → URL Configuration → Redirect URLs:

```
https://app.barracudaseo.com
https://app.barracudaseo.com/
https://app.barracudaseo.com/#/
https://app.barracudaseo.com/#
```

**Why**: Supabase is very strict about exact URL matches.

---

### Fix 2: Change Magic Link Redirect in Code

Update `web/src/lib/auth.js`:

**Current**:
```javascript
const redirectTo = typeof window !== 'undefined' 
  ? `${window.location.origin}/#/`
  : undefined;
```

**Try This Instead**:
```javascript
const redirectTo = typeof window !== 'undefined' 
  ? window.location.origin // No hash
  : undefined;
```

Then in Supabase Dashboard, make sure redirect URLs include:
- `https://app.barracudaseo.com` (no trailing slash or hash)

---

### Fix 3: Force Supabase to Detect Session

If tokens are in URL but not logging in, add this to `App.svelte` right after `initAuth()`:

```javascript
// Force check URL for auth tokens
const params = new URLSearchParams(window.location.hash.substring(1));
if (params.get('access_token')) {
  console.log('Forcing session exchange...');
  await supabase.auth.exchangeCodeForSession(window.location.hash.substring(1));
}
```

---

## Testing Checklist

After making changes:

- [ ] Clear browser cache (Ctrl+Shift+Delete)
- [ ] Test in incognito/private window
- [ ] Request new magic link (old ones won't work)
- [ ] Check console for error messages
- [ ] Verify URL contains `access_token` after clicking link
- [ ] Check Supabase logs for `user.signin` event

---

## Nuclear Option: Simplify Redirect

If nothing works, temporarily simplify:

1. **Remove hash routing from redirect**:
   ```javascript
   const redirectTo = 'https://app.barracudaseo.com'
   ```

2. **In Supabase Dashboard**:
   - Site URL: `https://app.barracudaseo.com`
   - Redirect URLs: Only `https://app.barracudaseo.com`

3. **Test with simple URL first**, then add hash routing back.

---

## Get More Debug Info

Add this to `App.svelte` right at the start of `onMount`:

```javascript
console.log('=== AUTH DEBUG INFO ===');
console.log('Current URL:', window.location.href);
console.log('Hash:', window.location.hash);
console.log('Supabase URL:', import.meta.env.PUBLIC_SUPABASE_URL);
console.log('Session storage:', localStorage.getItem('supabase.auth.token'));

// Listen for all auth events
supabase.auth.onAuthStateChange((event, session) => {
  console.log('Auth state change:', event, session);
});
```

This will give you full visibility into what's happening.

---

## Report Back

When asking for help, provide:

1. **Console logs** (screenshot or copy/paste)
2. **Full magic link URL** (with tokens removed for security)
3. **Supabase redirect URLs** (screenshot of your config)
4. **Any error messages**

---

## Next Steps

Based on the console output:
- **If you see "Magic link session set successfully"** → Session is working, might be UI issue
- **If you see "Auth callback error"** → Check the error details
- **If you see nothing** → Tokens aren't in URL, check Supabase redirect config

Let me know what you see in the console!
