#!/bin/bash
set -e

echo "🚀 Deploying Public Reports Feature"
echo "===================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Step 1: Database Migration
echo -e "${YELLOW}Step 1: Running database migration...${NC}"
cd "$(dirname "$0")/.."

if command -v supabase &> /dev/null; then
    echo "Applying Supabase migration..."
    supabase db push || {
        echo -e "${RED}Migration failed. Please check the error above.${NC}"
        echo "You can also apply the migration manually via Supabase Dashboard SQL Editor."
        exit 1
    }
    echo -e "${GREEN}✓ Migration applied successfully${NC}"
else
    echo -e "${YELLOW}⚠ Supabase CLI not found. Please apply migration manually:${NC}"
    echo "  1. Go to Supabase Dashboard → SQL Editor"
    echo "  2. Copy contents of: supabase/migrations/20251125123406_add_public_reports.sql"
    echo "  3. Paste and run"
    read -p "Press Enter after migration is applied..."
fi

echo ""

# Step 2: Deploy Backend
echo -e "${YELLOW}Step 2: Deploying backend to Cloud Run...${NC}"

# Check if GCP_PROJECT_ID is set
if [ -z "$GCP_PROJECT_ID" ]; then
    GCP_PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
    if [ -z "$GCP_PROJECT_ID" ]; then
        echo -e "${RED}Error: GCP_PROJECT_ID not set${NC}"
        echo "Set it with: export GCP_PROJECT_ID=your-project-id"
        exit 1
    fi
fi

export GCP_PROJECT_ID
export GCP_REGION=${GCP_REGION:-us-central1}

echo "Project: $GCP_PROJECT_ID"
echo "Region: $GCP_REGION"
echo ""

# Build and deploy
echo "Building Docker image..."
make docker-build || {
    echo -e "${RED}Docker build failed${NC}"
    exit 1
}

echo "Pushing to Artifact Registry..."
make docker-push || {
    echo -e "${RED}Docker push failed${NC}"
    exit 1
}

echo "Deploying to Cloud Run..."
make deploy-image || {
    echo -e "${RED}Deployment failed${NC}"
    exit 1
}

echo -e "${GREEN}✓ Backend deployed successfully${NC}"
echo ""

# Step 3: Verify APP_URL is set
echo -e "${YELLOW}Step 3: Checking environment variables...${NC}"
APP_URL=$(gcloud run services describe barracuda-api \
    --platform managed \
    --region $GCP_REGION \
    --format="value(spec.template.spec.containers[0].env)" 2>/dev/null | grep -o 'APP_URL=[^,]*' | cut -d'=' -f2 || echo "")

if [ -z "$APP_URL" ]; then
    echo -e "${YELLOW}⚠ APP_URL not set. Setting it now...${NC}"
    read -p "Enter your frontend URL (e.g., https://app.barracudaseo.com): " FRONTEND_URL
    if [ -n "$FRONTEND_URL" ]; then
        gcloud run services update barracuda-api \
            --platform managed \
            --region $GCP_REGION \
            --update-env-vars="APP_URL=$FRONTEND_URL" \
            --quiet
        echo -e "${GREEN}✓ APP_URL set to $FRONTEND_URL${NC}"
    fi
else
    echo -e "${GREEN}✓ APP_URL is set: $APP_URL${NC}"
fi

echo ""

# Step 4: Deploy Frontend
echo -e "${YELLOW}Step 4: Deploying frontend to Vercel...${NC}"
read -p "Deploy frontend now? (y/n): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    cd web
    echo "Deploying to Vercel..."
    vercel --prod || {
        echo -e "${RED}Vercel deployment failed${NC}"
        exit 1
    }
    echo -e "${GREEN}✓ Frontend deployed successfully${NC}"
    cd ..
else
    echo -e "${YELLOW}⚠ Skipping frontend deployment. Deploy manually with:${NC}"
    echo "  cd web && vercel --prod"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Deployment Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Next steps:"
echo "1. Test creating a public report in the app"
echo "2. Verify the public URL works without login"
echo "3. Check that project name/URL and issue URLs display correctly"

