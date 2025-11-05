#!/bin/bash
# Quick Vercel deployment helper

set -e

cd "$(dirname "$0")/../web"

echo "🚀 Deploying Barracuda Frontend to Vercel"
echo "=========================================="
echo ""

# Check if vercel CLI is installed
if ! command -v vercel &> /dev/null; then
    echo "Vercel CLI not found. Installing..."
    npm install -g vercel
fi

# Check if .env.local exists
if [ ! -f .env.local ]; then
    echo "⚠️  No .env.local file found."
    echo ""
    echo "Create .env.local with:"
    echo "  PUBLIC_SUPABASE_URL=https://your-project.supabase.co"
    echo "  PUBLIC_SUPABASE_ANON_KEY=your-anon-key"
    echo "  VITE_CLOUD_RUN_API_URL=https://barracuda-api-7paxg34svq-uc.a.run.app"
    echo ""
    read -p "Continue anyway? (y/n): " CONTINUE
    if [ "$CONTINUE" != "y" ] && [ "$CONTINUE" != "Y" ]; then
        exit 1
    fi
fi

# Build first to check for errors
echo "Building frontend..."
npm run build

echo ""
echo "✓ Build successful!"
echo ""

# Deploy
echo "Deploying to Vercel..."
vercel

echo ""
echo "✓ Deployment complete!"
echo ""
echo "Next steps:"
echo "1. Add environment variables in Vercel dashboard"
echo "2. Redeploy to apply environment variables"
echo "3. Visit your deployment URL"

