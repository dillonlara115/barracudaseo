# Deployment Checklist

## Summary of Changes

This deployment includes:
- ✅ New "Crawls" page with dedicated crawl management
- ✅ Navigation reorganization (crawl section removed from most pages)
- ✅ Page descriptions added for better UX
- ✅ Impact-First View bug fix (missing loadGSCStatus function)

## Pre-Deployment Checklist

### 1. ✅ Vercel (Frontend) - Already Done
- [x] Environment variables added to Vercel
- [ ] **Next:** Push code and deploy (Vercel will auto-deploy on push)

### 2. Google Cloud Run (Backend API)

#### Required Environment Variables

Add these to your `.env` file (or update Cloud Run directly):

```bash
# DataForSEO Integration (REQUIRED for rank tracking)
DATAFORSEO_LOGIN=your-dataforseo-email@example.com
DATAFORSEO_PASSWORD=your-dataforseo-api-password

# Existing variables (verify these are set)
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key  # Use Secret Manager

# Stripe (if using billing)
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_PRICE_ID_PRO=price_...
STRIPE_PRICE_ID_PRO_ANNUAL=price_...
STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true

# Email (if using Resend)
EMAIL_PROVIDER=resend
RESEND_API_KEY=re_...
EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com
APP_URL=https://app.barracudaseo.com

# GSC OAuth (if using GSC integration)
GSC_CLIENT_ID=your-client-id
GSC_CLIENT_SECRET=your-client-secret
GSC_REDIRECT_URL=https://app.barracudaseo.com/integrations
```

#### Steps to Deploy Backend

**Option 1: Update Environment Variables Only**
```bash
# Add DataForSEO variables to .env file
echo "DATAFORSEO_LOGIN=your-email@example.com" >> .env
echo "DATAFORSEO_PASSWORD=your-password" >> .env

# Update Cloud Run environment variables
./scripts/update-cloud-run-env.sh --production
```

**Option 2: Full Deployment (if code changes)**
```bash
# Build and push Docker image
make docker-build
make docker-push

# Deploy to Cloud Run
make deploy-backend
```

**Option 3: Manual gcloud Command**
```bash
gcloud run services update barracuda-api \
  --region us-central1 \
  --update-env-vars="DATAFORSEO_LOGIN=your-email@example.com,DATAFORSEO_PASSWORD=your-password" \
  --quiet
```

### 3. Supabase (Database)

#### Check Migrations

No new migrations are required for this deployment. The changes are primarily frontend UI improvements.

However, verify existing migrations are applied:
```bash
# Check migration status in Supabase dashboard
# Or use Supabase CLI:
supabase migration list
```

#### Existing Migrations (should already be applied)
- ✅ `20250120_add_dataforseo_rank_tracking.sql` - Rank tracking tables
- ✅ `20250121_add_keyword_scheduling_and_usage.sql` - Keyword scheduling

## Deployment Steps

### Step 1: Update Cloud Run Environment Variables

```bash
# Make sure .env has all required variables
# Then run:
./scripts/update-cloud-run-env.sh --production
```

**Note:** The script doesn't currently include DataForSEO variables. You'll need to add them manually:

```bash
gcloud run services update barracuda-api \
  --region us-central1 \
  --update-env-vars="DATAFORSEO_LOGIN=your-email,DATAFORSEO_PASSWORD=your-password" \
  --quiet
```

### Step 2: Deploy Backend Code (if needed)

If you made backend changes:
```bash
make docker-build
make docker-push
make deploy-backend
```

### Step 3: Deploy Frontend

Vercel will auto-deploy when you push to your main branch:

```bash
git add .
git commit -m "feat: Add Crawls page, reorganize navigation, improve UX"
git push origin main
```

Or trigger manual deployment in Vercel dashboard.

### Step 4: Verify Deployment

1. **Frontend (Vercel)**
   - Visit https://app.barracudaseo.com
   - Check browser console for errors
   - Verify "Crawls" link appears in sidebar
   - Test navigation to new Crawls page

2. **Backend (Cloud Run)**
   - Check Cloud Run logs for errors
   - Test API health endpoint: `curl https://your-api-url.run.app/health`
   - Verify DataForSEO integration works (test rank tracking)

3. **Database (Supabase)**
   - Verify no migration errors
   - Check that existing data is intact

## Important Notes

1. **DataForSEO Variables**: These MUST be set in Cloud Run for rank tracking to work
2. **No Database Changes**: This deployment doesn't require new migrations
3. **Backward Compatible**: All changes are additive - existing functionality remains intact
4. **Environment Variables**: Make sure Vercel has `PUBLIC_SUPABASE_URL` and `PUBLIC_SUPABASE_ANON_KEY` set

## Troubleshooting

### If rank tracking doesn't work:
- Check Cloud Run logs for DataForSEO authentication errors
- Verify `DATAFORSEO_LOGIN` and `DATAFORSEO_PASSWORD` are set correctly
- Test DataForSEO API credentials manually

### If navigation doesn't work:
- Clear browser cache
- Check Vercel build logs for errors
- Verify environment variables are set in Vercel

### If crawls page doesn't load:
- Check browser console for errors
- Verify route is registered in `App.svelte`
- Check that `ProjectPageLayout` is working correctly

