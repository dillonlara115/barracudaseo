# 🚀 Production Deployment Guide

This guide covers deploying changes to both the **App** (app.barracudaseo.com) and **Marketing Site** (barracudaseo.com) to production.

---

## 📋 Pre-Deployment Checklist

Before deploying, ensure you have:

- [ ] All code changes committed to git
- [ ] Environment variables updated (if needed)
- [ ] Database migrations applied (if any)
- [ ] Tested changes locally
- [ ] Reviewed deployment checklist for breaking changes

---

## 🎯 Deployment Overview

| Component | Platform | Auto-Deploy | Manual Deploy |
|-----------|----------|-------------|---------------|
| **App Frontend** | Vercel | ✅ Git push to `main` | Vercel Dashboard |
| **Marketing Site** | Vercel | ✅ Git push to `main` | Vercel Dashboard |
| **Backend API** | Google Cloud Run | ❌ Manual | `make deploy-backend` |
| **Database** | Supabase | ❌ Manual | Supabase Dashboard/CLI |

---

## 📱 Part 1: Deploy App Frontend (app.barracudaseo.com)

The app frontend is located in `/web` and deploys automatically via Vercel when you push to the main branch.

### Option A: Automatic Deployment (Recommended)

```bash
# Navigate to project root
cd /home/dillon/Sites/barracuda

# Stage all changes
git add .

# Commit with descriptive message
git commit -m "feat: Description of changes"

# Push to main branch (triggers Vercel auto-deploy)
git push origin main
```

**What happens:**
1. Vercel detects the push to `main`
2. Builds the SvelteKit app from `/web` directory
3. Deploys to `app.barracudaseo.com`
4. Sends deployment notification (if configured)

### Option B: Manual Deployment via Vercel CLI

```bash
cd /home/dillon/Sites/barracuda/web

# Deploy to production
npx vercel --prod

# Or if Vercel CLI is installed globally
vercel --prod
```

### Verify App Deployment

1. **Check Vercel Dashboard**
   - Visit https://vercel.com/dashboard
   - Find your `barracuda` project
   - Verify latest deployment succeeded

2. **Test Production Site**
   - Visit https://app.barracudaseo.com
   - Check browser console for errors
   - Test key functionality (login, crawls, etc.)

3. **Check Build Logs**
   - In Vercel dashboard, click on deployment
   - Review build logs for warnings/errors

---

## 🌐 Part 2: Deploy Marketing Site (barracudaseo.com)

The marketing site is located in `/marketing` and also deploys automatically via Vercel.

### Option A: Automatic Deployment (Recommended)

```bash
# Navigate to project root
cd /home/dillon/Sites/barracuda

# Stage marketing site changes
git add marketing/

# Commit with descriptive message
git commit -m "feat(marketing): Fix canonical URLs and add www redirect"

# Push to main branch (triggers Vercel auto-deploy)
git push origin main
```

**What happens:**
1. Vercel detects the push to `main`
2. Builds the SvelteKit marketing site from `/marketing` directory
3. Deploys to `barracudaseo.com`
4. Applies redirects from `vercel.json` (www → non-www)

### Option B: Manual Deployment via Vercel CLI

```bash
cd /home/dillon/Sites/barracuda/marketing

# Deploy to production
npx vercel --prod

# Or if Vercel CLI is installed globally
vercel --prod
```

### Verify Marketing Site Deployment

1. **Test www Redirect**
   ```bash
   curl -I https://www.barracudaseo.com
   # Should return 301 redirect to https://barracudaseo.com
   ```

2. **Check Canonical Tags**
   - Visit any page on https://barracudaseo.com
   - View page source
   - Verify canonical URLs don't have trailing slashes (except root)

3. **Check Sitemap**
   - Visit https://barracudaseo.com/sitemap.xml
   - Verify URLs are normalized (no trailing slashes)

4. **Test Key Pages**
   - Homepage loads correctly
   - Blog posts accessible
   - Navigation works
   - Forms/submissions work (if any)

---

## ⚙️ Part 3: Deploy Backend API (Google Cloud Run)

The backend API is deployed to Google Cloud Run and requires manual deployment.

### Step 1: Update Environment Variables (if needed)

If you added new environment variables:

```bash
# Update Cloud Run environment variables
./scripts/update-cloud-run-env.sh --production
```

**Or manually via gcloud:**

```bash
gcloud run services update barracuda-api \
  --region us-central1 \
  --update-env-vars="KEY1=value1,KEY2=value2" \
  --quiet
```

**Common environment variables:**
- `DATAFORSEO_LOGIN` - DataForSEO API email
- `DATAFORSEO_PASSWORD` - DataForSEO API password
- `GA4_CLIENT_ID` - Google Analytics OAuth client ID
- `GA4_CLIENT_SECRET` - Google Analytics OAuth secret
- `GSC_CLIENT_ID` - Google Search Console OAuth client ID
- `GSC_CLIENT_SECRET` - Google Search Console OAuth secret

### Step 2: Build and Deploy Backend Code

**If you made backend code changes:**

```bash
# Navigate to project root
cd /home/dillon/Sites/barracuda

# Build Docker image
make docker-build

# Push to Google Container Registry
make docker-push

# Deploy to Cloud Run
make deploy-backend
```

**Or manually:**

```bash
# Build Docker image
docker build -t gcr.io/YOUR_PROJECT_ID/barracuda-api:latest .

# Push to registry
docker push gcr.io/YOUR_PROJECT_ID/barracuda-api:latest

# Deploy to Cloud Run
gcloud run deploy barracuda-api \
  --image gcr.io/YOUR_PROJECT_ID/barracuda-api:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated
```

### Verify Backend Deployment

1. **Check Cloud Run Logs**
   ```bash
   gcloud run services logs read barracuda-api --region us-central1 --limit 50
   ```

2. **Test Health Endpoint**
   ```bash
   curl https://YOUR_API_URL.run.app/health
   # Should return: {"status":"healthy"}
   ```

3. **Test API Endpoints**
   - Test authentication
   - Test key API endpoints
   - Verify integrations (DataForSEO, GA4, GSC) work

---

## 🗄️ Part 4: Database Migrations (Supabase)

If you have new database migrations:

### Option A: Via Supabase Dashboard

1. Go to https://supabase.com/dashboard
2. Select your project
3. Navigate to **SQL Editor**
4. Run migration SQL files from `/supabase/migrations/`

### Option B: Via Supabase CLI

```bash
# Apply all pending migrations
supabase migration up

# Or apply specific migration
supabase migration up --version 20250130000000
```

### Verify Migrations

```bash
# List applied migrations
supabase migration list

# Check database schema
supabase db diff
```

---

## 🔍 Post-Deployment Verification

### App Frontend (app.barracudaseo.com)

- [ ] Site loads without errors
- [ ] Authentication works
- [ ] Projects load correctly
- [ ] Crawls can be created/viewed
- [ ] Keyword tracking works (if applicable)
- [ ] No console errors
- [ ] Environment variables loaded correctly

### Marketing Site (barracudaseo.com)

- [ ] Homepage loads
- [ ] All pages accessible
- [ ] Blog posts load correctly
- [ ] Canonical tags correct (no trailing slashes)
- [ ] www redirects work (301 to non-www)
- [ ] Sitemap accessible and correct
- [ ] No console errors
- [ ] Analytics tracking works (if applicable)

### Backend API

- [ ] Health endpoint responds
- [ ] Authentication endpoints work
- [ ] API endpoints return correct data
- [ ] Integrations work (DataForSEO, GA4, GSC)
- [ ] No errors in Cloud Run logs
- [ ] Response times acceptable

### Database

- [ ] Migrations applied successfully
- [ ] No migration errors
- [ ] Existing data intact
- [ ] New tables/columns exist (if applicable)

---

## 🚨 Troubleshooting

### App Frontend Issues

**Build fails in Vercel:**
- Check build logs in Vercel dashboard
- Verify environment variables are set
- Check for TypeScript/ESLint errors locally first

**Environment variables not loading:**
- Verify variables are set in Vercel dashboard
- Check variable names match code (case-sensitive)
- Redeploy after adding variables

**Routes not working:**
- Check `App.svelte` route configuration
- Verify route files exist in `/web/src/routes/`
- Check browser console for errors

### Marketing Site Issues

**www redirect not working:**
- Verify `vercel.json` exists in `/marketing/`
- Check redirect syntax is correct
- Ensure Vercel project is connected to correct domain

**Canonical URLs incorrect:**
- Check `MetaTags.svelte` normalization logic
- Verify `SITE_URL` constant is correct
- Test locally with `npm run preview`

**Build fails:**
- Run `npm run build` locally to see errors
- Check for missing dependencies
- Verify SvelteKit adapter is configured correctly

### Backend API Issues

**Deployment fails:**
- Check Docker build logs
- Verify Google Cloud authentication (`gcloud auth login`)
- Check Cloud Run service limits/quota

**Environment variables not loading:**
- Verify variables are set in Cloud Run
- Check variable names match code
- Restart Cloud Run service after updating variables

**API errors:**
- Check Cloud Run logs: `gcloud run services logs read barracuda-api`
- Verify database connection
- Check external API credentials (DataForSEO, Google APIs)

---

## 📝 Quick Reference Commands

### App Frontend
```bash
# Build locally
cd web && npm run build

# Preview production build
cd web && npm run preview

# Deploy via git (auto)
git push origin main
```

### Marketing Site
```bash
# Build locally
cd marketing && npm run build

# Preview production build
cd marketing && npm run preview

# Deploy via git (auto)
git push origin main
```

### Backend API
```bash
# Build and deploy
make docker-build && make docker-push && make deploy-backend

# Update environment variables
./scripts/update-cloud-run-env.sh --production

# View logs
gcloud run services logs read barracuda-api --region us-central1 --limit 50
```

### Database
```bash
# Apply migrations
supabase migration up

# List migrations
supabase migration list
```

---

## 🔐 Environment Variables Reference

### Vercel (App Frontend)
- `PUBLIC_SUPABASE_URL` - Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Supabase anonymous key
- `VITE_CLOUD_RUN_API_URL` - Backend API URL

### Vercel (Marketing Site)
- Optional: `PUBLIC_SUPABASE_URL` (if using Supabase features)
- Optional: `PUBLIC_SUPABASE_ANON_KEY`

### Cloud Run (Backend API)
- `DATAFORSEO_LOGIN` - DataForSEO API email
- `DATAFORSEO_PASSWORD` - DataForSEO API password
- `GA4_CLIENT_ID` - Google Analytics OAuth client ID
- `GA4_CLIENT_SECRET` - Google Analytics OAuth secret
- `GSC_CLIENT_ID` - Google Search Console OAuth client ID
- `GSC_CLIENT_SECRET` - Google Search Console OAuth secret
- `SUPABASE_SERVICE_ROLE_KEY` - Supabase service role key (use Secret Manager)
- `STRIPE_SECRET_KEY` - Stripe secret key (use Secret Manager)
- `APP_URL` - Frontend app URL (https://app.barracudaseo.com)

---

## 📚 Additional Resources

- [Vercel Deployment Docs](https://vercel.com/docs)
- [Google Cloud Run Docs](https://cloud.google.com/run/docs)
- [Supabase Migrations](https://supabase.com/docs/guides/cli/local-development#database-migrations)
- [DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md) - Detailed checklist for specific deployments

---

**Last Updated:** January 2025
**Maintainer:** @dillonlara
