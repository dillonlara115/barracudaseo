# Magic Link Debug - Empty Hash Issue

## Problem
The URL hash is empty after clicking magic link, meaning Supabase isn't redirecting with auth tokens.

## Root Cause Analysis

This happens when:
1. **Redirect URL not in Supabase allowlist** (most common)
2. **Magic link email template is wrong**
3. **Supabase email settings misconfigured**

## Immediate Debug Steps

### Step 1: Check the Magic Link URL in Email

1. **Open the magic link email**
2. **Right-click the button/link** → "Copy Link Address"
3. **Paste into text editor**

**The URL should look like:**
```
https://YOUR-PROJECT.supabase.co/auth/v1/verify?token=...&type=magiclink&redirect_to=https://app.barracudaseo.com
```

**Check:**
- ✅ `redirect_to=https://app.barracudaseo.com` (correct)
- ❌ `redirect_to=` is missing or wrong

If `redirect_to` is missing, the issue is in how we're sending the magic link.

### Step 2: Check Supabase Dashboard Settings

Go to: **Supabase Dashboard → Authentication → URL Configuration**

**Verify these EXACT settings:**

#### Site URL
```
https://app.barracudaseo.com
```
**No trailing slash!**

#### Redirect URLs
Add these one by one:
```
https://app.barracudaseo.com
https://app.barracudaseo.com/
http://localhost:3000
http://127.0.0.1:3000
```

### Step 3: Check Supabase Auth Logs

1. Go to: **Supabase Dashboard → Logs → Auth Logs**
2. Look for recent events after clicking magic link
3. Check for errors like:
   - "Invalid redirect URL"
   - "Redirect URL not allowed"

## Quick Fix Options

### Option A: Test with Localhost First

If production isn't working, test locally:

1. **Request magic link on localhost**:
   ```bash
   cd web
   npm run dev
   # Open http://localhost:3000
   ```

2. **Check Inbucket**: http://localhost:54324

3. **Click magic link** - should work locally with local Supabase

4. **If local works but production doesn't** → It's definitely a Supabase production config issue

### Option B: Verify Redirect URL is Being Sent

Add this temporary debug code to `auth.js`:

```javascript
export async function signInWithMagicLink(email) {
  try {
    const redirectTo = typeof window !== 'undefined' 
      ? window.location.origin
      : undefined;
    
    console.log('🔍 Magic link redirect URL:', redirectTo);
    console.log('🔍 Current origin:', window.location.origin);

    const { data, error } = await supabase.auth.signInWithOtp({
      email,
      options: {
        emailRedirectTo: redirectTo,
        shouldCreateUser: false
      }
    });

    if (error) {
      console.error('🔴 Magic link error:', error);
      throw error;
    }
    
    console.log('✅ Magic link sent successfully:', data);
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}
```

Check console when requesting magic link - should show:
```
🔍 Magic link redirect URL: https://app.barracudaseo.com
🔍 Current origin: https://app.barracudaseo.com
✅ Magic link sent successfully: {...}
```

### Option C: Test Direct Supabase API

Test if Supabase is working at all:

```javascript
// Add to browser console
const { data, error } = await supabase.auth.signInWithOtp({
  email: 'your-email@example.com',
  options: {
    emailRedirectTo: 'https://app.barracudaseo.com'
  }
});
console.log('Result:', { data, error });
```

Check the email - does redirect_to appear in the link?

## Common Issues & Solutions

### Issue: "Site URL does not match redirect URL"

**Cause**: Site URL in Supabase doesn't match your domain

**Fix**:
1. Site URL: `https://app.barracudaseo.com` (exact match, no trailing slash)
2. Redirect URLs: Include `https://app.barracudaseo.com` in the list

### Issue: Email link goes to wrong domain

**Cause**: `redirectTo` is being set to wrong value

**Fix**: 
- Check `window.location.origin` in console - should be `https://app.barracudaseo.com`
- Make sure no environment variables are overriding it

### Issue: Tokens in URL but then disappear

**Cause**: App is redirecting before processing tokens

**Fix**: Already handled in latest code - tokens are processed before redirect

## Nuclear Option: Use Token-Based Flow

If magic links continue to fail, switch to OTP code flow:

```javascript
// User enters email
const { data, error } = await supabase.auth.signInWithOtp({
  email,
  options: {
    shouldCreateUser: false
  }
  // No redirect_to needed!
});

// User receives code in email
// User enters code in app
const { data, error } = await supabase.auth.verifyOtp({
  email,
  token: userEnteredCode,
  type: 'email'
});
```

This eliminates redirect URL issues entirely.

## What to Check Right Now

1. **Copy the magic link URL from email** - paste it here
2. **Check Supabase Site URL** - screenshot it
3. **Check Supabase Redirect URLs list** - screenshot it
4. **Check browser console** when requesting magic link - any errors?

Let me know what you find!
