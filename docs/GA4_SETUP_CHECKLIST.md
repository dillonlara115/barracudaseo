# Google Analytics 4 (GA4) API Setup Checklist

## Step 1: Enable the Required APIs

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project (or create a new one for GA4)
3. Navigate to **APIs & Services** → **Library**
4. Enable both APIs:
   - Search for **"Google Analytics Data API"** and click **Enable**
   - Search for **"Google Analytics Admin API"** and click **Enable**

## Step 2: Configure OAuth Consent Screen

1. Go to **APIs & Services** → **OAuth consent screen**
2. Fill out the required fields:
   - **User Type**: Choose "External" (unless you're using Google Workspace)
   - **App Name**: Barracuda SEO Crawler (or your app name)
   - **User support email**: Your email
   - **Developer contact**: Your email
3. Click **Save and Continue**

4. On the **Scopes** page:
   - Click **Add or Remove Scopes**
   - Search for: `https://www.googleapis.com/auth/analytics.readonly`
   - **Check the box** to add it
   - Click **Update** then **Save and Continue**
   
   **Note**: The same scope (`analytics.readonly`) works for both the Data API (reading reports) and Admin API (listing properties). You only need to add it once.

5. On the **Test users** page (if your app is in Testing mode):
   - **IMPORTANT**: Click **+ ADD USERS**
   - Add your Google account email address (e.g., `your-email@gmail.com`)
   - You can add multiple test users if needed
   - Click **Save and Continue**

6. Review and **Back to Dashboard**

## Common Issue: "Access blocked" Error

If you see "Access blocked: Barracuda has not completed the Google verification process":

**Solution**: Add yourself as a test user:
1. Go to **APIs & Services** → **OAuth consent screen**
2. Scroll to **Test users** section
3. Click **+ ADD USERS**
4. Enter your Google account email
5. Click **Add**
6. Try connecting again

**Note**: Test users can only access the app while it's in Testing mode. For production use, you'll need to publish the app (requires Google verification).

## Step 3: Create OAuth 2.0 Client ID

1. Go to **APIs & Services** → **Credentials**
2. Click **+ CREATE CREDENTIALS** → **OAuth client ID**
3. Choose **Web application** as the application type
4. Give it a name (e.g., "Barracuda GA4 Integration")
5. Under **Authorized redirect URIs**, add:
   - For local development: `http://localhost:8080/api/ga4/callback`
   - For production: `https://your-api-url.com/api/ga4/callback`
6. Click **Create**
7. **Copy your Client ID and Client Secret** - you'll need these for environment variables

## Step 4: Set Environment Variables

Set the following environment variables:

```bash
export GA4_CLIENT_ID='your-client-id.apps.googleusercontent.com'
export GA4_CLIENT_SECRET='your-client-secret'
```

Or use JSON credentials:

```bash
export GA4_CREDENTIALS_JSON='{"web":{"client_id":"...","client_secret":"...","redirect_uris":["http://localhost:8080/api/ga4/callback"]}}'
```

## Step 5: Verify Setup

The scope you need is:
- **Scope**: `https://www.googleapis.com/auth/analytics.readonly`
- **Display Name**: "View your Google Analytics data"
- **Used for**: Both Data API (reading reports) and Admin API (listing properties)

This scope is already configured in the code - you just need to enable it in Google Cloud Console!

## Testing

After setup:
1. Start your API server with GA4 credentials set
2. Go to your project's integrations page
3. Click "Connect Google Analytics 4"
4. You should see the consent screen asking for permission to "View your Google Analytics data"
5. Authorize and select your GA4 property
6. Done! 🎉

## Important Notes

- **Separate from GSC**: GA4 uses its own OAuth credentials, allowing users to connect different Google accounts for GSC vs GA4
- **Same Scope for Both APIs**: The `analytics.readonly` scope covers both the Data API and Admin API for read-only operations
- **Property Listing**: The Admin API is used to list available GA4 properties, while the Data API is used to fetch performance metrics
