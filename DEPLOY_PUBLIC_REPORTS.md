# Deployment Checklist: Public Reports Feature

This guide will help you deploy the new public reports feature.

## Prerequisites

- ✅ Database migration created: `20251125123406_add_public_reports.sql`
- ✅ Backend handlers implemented: `internal/api/report_handlers.go`
- ✅ Frontend components created: `PublicReportGenerator.svelte` and `PublicReportView.svelte`
- ✅ Routes registered in `App.svelte`

## Step 1: Run Database Migration

The new `public_reports` table needs to be created in Supabase.

### Option A: Using Supabase CLI (Recommended)

```bash
# Navigate to project root
cd /home/dillon/Sites/cli-scanner

# Apply migration
supabase db push

# Or if using local Supabase:
supabase migration up
```

### Option B: Using Supabase Dashboard

1. Go to your Supabase project dashboard
2. Navigate to SQL Editor
3. Copy the contents of `supabase/migrations/20251125123406_add_public_reports.sql`
4. Paste and run the SQL script

### Verify Migration

```sql
-- Check if table exists
SELECT * FROM public.public_reports LIMIT 1;
```

## Step 2: Deploy Backend (Cloud Run)

### Quick Deploy (Using Makefile)

```bash
# Set environment variables (if not already set)
export GCP_PROJECT_ID=your-project-id
export GCP_REGION=us-central1

# Build and deploy (preserves existing env vars)
make deploy-image

# Or full deploy (rebuilds everything)
make deploy
```

### Manual Deploy

```bash
# Build Docker image
make docker-build

# Push to Artifact Registry
make docker-push

# Deploy to Cloud Run (preserves env vars)
gcloud run services update barracuda-api \
  --image $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/barracuda/barracuda-api:latest \
  --platform managed \
  --region $GCP_REGION \
  --quiet
```

### Verify Backend Deployment

```bash
# Get service URL
export API_URL=$(gcloud run services describe barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --format="value(status.url)")

# Test health endpoint
curl $API_URL/health

# Check logs for any errors
gcloud run services logs read barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --limit 50
```

## Step 3: Deploy Frontend (Vercel)

### Option A: Automatic (Git Push)

If you have Vercel connected to your Git repository:

```bash
# Commit changes
git add .
git commit -m "Add public reports feature"
git push

# Vercel will automatically deploy
```

### Option B: Manual Deploy (Vercel CLI)

```bash
# Navigate to web directory
cd web

# Deploy to production
vercel --prod

# Or deploy to preview
vercel
```

### Verify Frontend Deployment

1. Visit your Vercel deployment URL
2. Log in and navigate to a crawl
3. Check the Dashboard tab for "Public Client Reports" section
4. Try creating a public report

## Step 4: Set Environment Variables

### Backend (Cloud Run)

The backend should already have these variables. Verify they're set:

```bash
# Check current env vars
gcloud run services describe barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --format="value(spec.template.spec.containers[0].env)"

# If APP_URL is not set, add it:
gcloud run services update barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --update-env-vars="APP_URL=https://app.barracudaseo.com"
```

### Frontend (Vercel)

Ensure these are set in Vercel Dashboard → Settings → Environment Variables:

- `PUBLIC_SUPABASE_URL` - Your Supabase URL
- `PUBLIC_SUPABASE_ANON_KEY` - Your Supabase anon key
- `VITE_CLOUD_RUN_API_URL` - Your Cloud Run API URL

## Step 5: Test the Feature

### Test Creating a Public Report

1. Log into the app
2. Navigate to a project with a crawl
3. Go to Dashboard tab
4. Scroll to "Public Client Reports" section
5. Click "Create Public Report"
6. Fill in the form (optional password, expiry, etc.)
7. Click "Create Report"

### Test Viewing a Public Report

1. Copy the generated public URL
2. Open in an incognito/private window (to test without login)
3. Verify the report loads correctly
4. Check that project name and URL are displayed
5. Verify issues show their URLs
6. Test password protection (if enabled)

## Troubleshooting

### Database Migration Issues

```bash
# Check migration status
supabase migration list

# If migration failed, check Supabase logs
# Or run migration manually in SQL Editor
```

### Backend Deployment Issues

```bash
# Check Cloud Run logs
gcloud run services logs read barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --limit 100

# Check service status
gcloud run services describe barracuda-api \
  --platform managed \
  --region $GCP_REGION
```

### Frontend Issues

- Check browser console for errors
- Verify environment variables are set correctly
- Check Vercel deployment logs
- Ensure `APP_URL` is set correctly in backend

### Common Issues

1. **Public report URL redirects to login**
   - Verify route is added to `App.svelte`
   - Check `isPublicPage` includes `/reports/`

2. **404 on public report route**
   - Verify route is registered: `/reports/:token`
   - Check hash routing is working

3. **Issues don't show URLs**
   - Verify backend is enriching issues with page URLs
   - Check that issues have `page_id` set

## Rollback Plan

If something goes wrong:

### Rollback Backend

```bash
# Deploy previous image version
gcloud run services update barracuda-api \
  --image $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/barracuda/barracuda-api:previous-tag \
  --platform managed \
  --region $GCP_REGION
```

### Rollback Frontend

```bash
# Revert to previous deployment in Vercel dashboard
# Or redeploy previous commit
```

### Rollback Database

```sql
-- Drop the table (only if needed)
DROP TABLE IF EXISTS public.public_reports CASCADE;
```

## Post-Deployment Checklist

- [ ] Database migration applied successfully
- [ ] Backend deployed and health check passes
- [ ] Frontend deployed and accessible
- [ ] Can create public reports
- [ ] Public reports accessible without login
- [ ] Project name and URL display correctly
- [ ] Issue URLs display correctly
- [ ] Password protection works (if tested)
- [ ] Expiry dates work (if tested)

## Next Steps

After successful deployment:

1. Monitor error logs for any issues
2. Test with real clients/users
3. Gather feedback on the feature
4. Consider adding analytics/tracking

