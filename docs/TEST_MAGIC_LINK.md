# Magic Link Testing Guide

## Quick Test Checklist

Use this checklist to verify the magic link implementation works correctly.

---

## Local Testing

### Prerequisites
- [ ] Supabase running locally (`supabase start`)
- [ ] Web app running (`npm run dev` in `web/` directory)
- [ ] Inbucket accessible at http://localhost:54324

### Test 1: Sign Up with Magic Link

1. **Open app**: http://localhost:3000
2. **Fill out form**:
   - First Name: `Test`
   - Last Name: `User`
   - Email: `test@example.com`
3. **Click**: "Send magic link"
4. **Verify**: Success message shows "Check your email! We sent you a magic link..."
5. **Open Inbucket**: http://localhost:54324
6. **Find email**: Should see email to `test@example.com`
7. **Click magic link** in email
8. **Verify**: 
   - Redirected to dashboard
   - User is logged in
   - Can see user email in top right

**Expected Result**: ✅ User successfully signed up and logged in

---

### Test 2: Sign In with Magic Link

1. **Sign out**: Click user menu → "Sign Out"
2. **Return to login**: Should see login page
3. **Enter email**: `test@example.com` (from Test 1)
4. **Click**: "Send magic link"
5. **Verify**: Success message shows
6. **Open Inbucket**: http://localhost:54324
7. **Find new email**: Should see second email
8. **Click magic link**
9. **Verify**: Logged in and redirected to dashboard

**Expected Result**: ✅ Existing user successfully signed in

---

### Test 3: Password Login (Optional)

1. **Sign out** 
2. **Click**: "Use password instead"
3. **Verify**: Password field appears
4. **Note**: This only works if user has a password set
   - New users created with magic link don't have passwords
   - Test with a user created via old password flow

**Expected Result**: ✅ Password option is available but not required

---

### Test 4: Magic Link Sent State

1. **Go to login page**
2. **Enter email**: `newuser@example.com`
3. **Click**: "Send magic link"
4. **Verify UI shows**:
   - ✅ Email icon/illustration
   - ✅ "Check your email" message
   - ✅ Shows the email address entered
   - ✅ "Try a different email" button

**Expected Result**: ✅ Clear feedback that email was sent

---

### Test 5: Rate Limiting

1. **Request magic link** for `test@example.com`
2. **Wait 10 seconds**
3. **Request again** for same email
4. **Immediately request again** (within 60 seconds)
5. **Verify**: Error or warning about rate limiting

**Expected Result**: ✅ Can't spam magic link requests

---

### Test 6: Expired Magic Link

1. **Request magic link**
2. **In Supabase Studio**: 
   - Go to http://localhost:54323
   - Navigate to Authentication → Users
   - Find the test user
   - Click "..." → "Sign Out All Sessions"
3. **Try clicking old magic link**
4. **Verify**: Should show error or redirect to login

**Expected Result**: ✅ Old links don't work after session invalidation

---

### Test 7: Session Persistence

1. **Sign in with magic link**
2. **Close browser tab**
3. **Open new tab**: http://localhost:3000
4. **Verify**: Still logged in (no redirect to login)
5. **Check browser DevTools**:
   - Application → Local Storage
   - Should see Supabase session data

**Expected Result**: ✅ Session persists across tabs/browser restarts

---

### Test 8: Invite Token Flow

1. **Get an invite token** (create one or use existing)
2. **Open**: `http://localhost:3000/#/auth?invite_token=YOUR_TOKEN`
3. **Verify**: 
   - Shows "Team Invitation" banner
   - User is prompted to sign up
4. **Complete magic link signup**
5. **Verify**: After login, redirected to `/team/accept?token=...`

**Expected Result**: ✅ Invite flow works with magic links

---

## Production Testing

### Prerequisites
- [ ] Production Supabase configured (see `PRODUCTION_MAGIC_LINK_SETUP.md`)
- [ ] App deployed to production
- [ ] Real email account for testing

### Test 1: Production Magic Link Delivery

1. **Open production app**: https://your-domain.com
2. **Enter real email**: your-email@gmail.com
3. **Click**: "Send magic link"
4. **Check inbox** (and spam folder)
5. **Verify email**:
   - From: Barracuda SEO
   - Subject: Your magic link to Barracuda
   - Contains branded button/link
6. **Click link**
7. **Verify**: 
   - Redirects to `https://your-domain.com/#/`
   - User is logged in
   - Session persists

**Expected Result**: ✅ Email arrives within 1-2 minutes, link works

---

### Test 2: Production Session Duration

1. **Sign in with magic link**
2. **Wait 1 hour**
3. **Refresh page**
4. **Verify**: Still logged in
5. **Check DevTools**: Token should auto-refresh

**Expected Result**: ✅ User stays logged in for 7 days

---

### Test 3: Mobile Testing

1. **Open production app on mobile** (iPhone/Android)
2. **Request magic link**
3. **Verify**: 
   - Can switch to email app
   - Can click link from email app
   - Redirects back to browser
   - User is logged in

**Expected Result**: ✅ Mobile flow works smoothly

---

### Test 4: Cross-Device Login

1. **Request magic link on Desktop**
2. **Open email on Phone**
3. **Click magic link on Phone**
4. **Verify**: User is logged in on Phone (not Desktop)

**Expected Result**: ✅ Magic link works on any device

---

## Automated Testing (Future)

Consider adding these Playwright/Cypress tests:

```javascript
// Test: Magic Link Sign Up
test('user can sign up with magic link', async ({ page }) => {
  await page.goto('http://localhost:3000');
  await page.fill('input[type="email"]', 'test@example.com');
  await page.fill('input[placeholder="First name"]', 'Test');
  await page.fill('input[placeholder="Last name"]', 'User');
  await page.click('button:has-text("Send magic link")');
  
  // Wait for success message
  await expect(page.locator('text=Check your email')).toBeVisible();
  
  // TODO: Add Inbucket API call to retrieve magic link
  // TODO: Navigate to magic link URL
  // TODO: Assert user is logged in
});
```

---

## Monitoring Checklist

After launch, monitor these metrics:

### Week 1
- [ ] Magic link delivery rate (target: >99%)
- [ ] Click-through rate (target: >80%)
- [ ] Failed auth attempts
- [ ] Support tickets related to login

### Ongoing
- [ ] Email bounce rate (< 2%)
- [ ] Session duration (should average ~3-5 days)
- [ ] Password vs magic link usage ratio
- [ ] User retention after auth change

---

## Rollback Triggers

Consider reverting if:
- 🚨 Magic link delivery rate < 90%
- 🚨 >10% of users can't log in
- 🚨 >50% of support tickets about auth
- 🚨 Email provider issues persist >6 hours

---

## Success Criteria

Magic link implementation is successful if:
- ✅ >95% magic link delivery rate
- ✅ <5% support tickets related to auth
- ✅ >70% of users prefer magic link over password
- ✅ Average session duration 3-5 days
- ✅ Zero security incidents

---

**Testing completed on**: _____________
**Tested by**: _____________
**Production deployment date**: _____________

