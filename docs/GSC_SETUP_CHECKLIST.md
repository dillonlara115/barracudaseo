# Google Search Console API Setup Checklist

## Step 1: Enable the Search Console API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project: **barracuda-477122**
3. Navigate to **APIs & Services** → **Library**
4. Search for **"Google Search Console API"**
5. Click on it and click **Enable**

## Step 2: Configure OAuth Consent Screen

1. Go to **APIs & Services** → **OAuth consent screen**
2. Fill out the required fields:
   - **User Type**: Choose "External" (unless you're using Google Workspace)
   - **App Name**: Barracuda SEO Crawler
   - **User support email**: Your email
   - **Developer contact**: Your email
3. Click **Save and Continue**

4. On the **Scopes** page:
   - Click **Add or Remove Scopes**
   - Search for: `https://www.googleapis.com/auth/webmasters.readonly`
   - **Check the box** to add it
   - Click **Update** then **Save and Continue**

5. On the **Test users** page (if your app is in Testing mode):
   - **IMPORTANT**: Click **+ ADD USERS**
   - Add your Google account email address (e.g., `dillonlara115@gmail.com`)
   - You can add multiple test users if needed
   - Click **Save and Continue**

6. Review and **Back to Dashboard**

## Common Issues

### Issue: "Access blocked" Error

If you see "Access blocked: Barracuda has not completed the Google verification process":

**Solution**: Add yourself as a test user:
1. Go to **APIs & Services** → **OAuth consent screen**
2. Scroll to **Test users** section
3. Click **+ ADD USERS**
4. Enter your Google account email: `dillonlara115@gmail.com`
5. Click **Add**
6. Try connecting again

**Note**: Test users can only access the app while it's in Testing mode. For production use, you'll need to publish the app (requires Google verification).

### Issue: "Google hasn't verified this app" Warning (Even After Verification)

If you've completed Google's verification process but users still see "Google hasn't verified this app":

**This is normal if your app is still in "Testing" mode.** Even verified apps show this warning in Testing mode.

**To remove the warning:**

1. Go to **APIs & Services** → **OAuth consent screen** in Google Cloud Console
2. Check the **Publishing status** at the top of the page
3. If it says "Testing", click **PUBLISH APP** button
4. Confirm the publishing action
5. Wait a few minutes for the change to propagate

**Important Notes:**
- **Testing Mode**: Only test users (added manually) can use the app. Shows verification warning.
- **Published Mode**: Anyone can use the app. No verification warning (if app is verified).
- **Verification vs Publishing**: These are separate steps:
  - **Verification** = Google reviews your app for compliance (can take days/weeks)
  - **Publishing** = Makes your app available to all users (instant, but requires verification for sensitive scopes)

**If your app is already published but still shows the warning:**
- There may be a delay (wait 15-30 minutes)
- Check that verification status shows "Verified" in the OAuth consent screen
- Ensure all required scopes are approved
- Try clearing browser cache and cookies

## Step 3: Update Authorized Redirect URI

1. Go to **APIs & Services** → **Credentials**
2. Click on your OAuth 2.0 Client ID
3. Under **Authorized redirect URIs**, make sure you have:
   - `http://localhost:8080/api/gsc/callback`
4. Click **Save**

## Step 4: Verify Setup

The scope you need is:
- **Scope**: `https://www.googleapis.com/auth/webmasters.readonly`
- **Display Name**: "View Search Console data for your verified sites"

This is already configured in the code - you just need to enable it in Google Cloud Console!

## Testing

After setup:
1. Run `barracuda serve --results results.json`
2. Go to Recommendations tab
3. Click "Connect Google Search Console"
4. You should see the consent screen asking for permission to "View Search Console data"
5. Authorize and you're done!

## Step 5: Schedule Automatic Syncs (Optional but Recommended)

1. Set the shared cron secret in both your API environment and Supabase Edge functions:
   - `GSC_SYNC_SECRET`
   - `CLOUD_RUN_API_URL` / `BARRACUDA_API_URL`
   - Optional `GSC_SYNC_LOOKBACK_DAYS` to tweak the data window (default 30).
2. Deploy the `gsc-sync` Edge function (`supabase/functions/gsc-sync`) and enable the cron job defined in `supabase/config.toml` (`gsc_daily_sync` runs daily at 06:00 UTC).
3. Once scheduled, the cron job will invoke `/api/internal/gsc/sync` with the shared secret and refresh cached Search Console data for every connected project.

This keeps the dashboard up-to-date without requiring manual refreshes.
