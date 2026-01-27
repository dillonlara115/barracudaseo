#!/bin/bash

# Quick diagnostic script to check local Supabase setup for magic links

echo "🔍 Checking Local Supabase Setup..."
echo ""

# Check if Supabase CLI is installed
if ! command -v supabase &> /dev/null; then
    echo "❌ Supabase CLI not found. Install it first:"
    echo "   https://supabase.com/docs/guides/cli"
    exit 1
fi

echo "✅ Supabase CLI found"
echo ""

# Check Supabase status
echo "📊 Checking Supabase status..."
if supabase status &> /dev/null; then
    echo "✅ Supabase is running"
    
    # Extract key info
    API_URL=$(supabase status 2>/dev/null | grep "API URL" | awk '{print $3}')
    INBUCKET_URL=$(supabase status 2>/dev/null | grep "Inbucket URL" | awk '{print $3}')
    ANON_KEY=$(supabase status 2>/dev/null | grep "anon key" | awk '{print $3}')
    
    echo "   API URL: $API_URL"
    echo "   Inbucket URL: $INBUCKET_URL"
    echo ""
    
    # Check if Inbucket is accessible
    if curl -s "$INBUCKET_URL" &> /dev/null; then
        echo "✅ Inbucket is accessible at $INBUCKET_URL"
        echo "   📧 Open this URL to view magic link emails"
    else
        echo "⚠️  Inbucket might not be fully started yet"
        echo "   Try: supabase stop && supabase start"
    fi
else
    echo "❌ Supabase is not running"
    echo ""
    echo "   Start it with:"
    echo "   cd /home/dillon/Sites/barracuda"
    echo "   supabase start"
    exit 1
fi

echo ""
echo "📝 Checking environment variables..."

# Check web/.env.local
if [ -f "web/.env.local" ]; then
    echo "✅ Found web/.env.local"
    
    ENV_URL=$(grep "PUBLIC_SUPABASE_URL" web/.env.local | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    ENV_KEY=$(grep "PUBLIC_SUPABASE_ANON_KEY" web/.env.local | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -n "$ENV_URL" ]; then
        echo "   PUBLIC_SUPABASE_URL: $ENV_URL"
        
        if [[ "$ENV_URL" == *"127.0.0.1:54321"* ]] || [[ "$ENV_URL" == *"localhost:54321"* ]]; then
            echo "   ✅ Points to local Supabase"
        else
            echo "   ⚠️  Points to production/remote Supabase"
            echo "   💡 For local magic links, should be: http://127.0.0.1:54321"
        fi
    else
        echo "   ⚠️  PUBLIC_SUPABASE_URL not found in .env.local"
    fi
    
    if [ -n "$ENV_KEY" ]; then
        KEY_PREFIX="${ENV_KEY:0:20}..."
        echo "   PUBLIC_SUPABASE_ANON_KEY: $KEY_PREFIX"
        
        # Compare with running Supabase
        if [ -n "$ANON_KEY" ] && [ "$ENV_KEY" != "$ANON_KEY" ]; then
            echo "   ⚠️  Anon key doesn't match running Supabase instance"
            echo "   💡 Update .env.local with: supabase status | grep 'anon key'"
        else
            echo "   ✅ Anon key matches running Supabase"
        fi
    else
        echo "   ⚠️  PUBLIC_SUPABASE_ANON_KEY not found in .env.local"
    fi
else
    echo "❌ web/.env.local not found"
    echo ""
    echo "   Create it with:"
    echo "   cd web"
    echo "   echo 'PUBLIC_SUPABASE_URL=http://127.0.0.1:54321' > .env.local"
    echo "   echo 'PUBLIC_SUPABASE_ANON_KEY=<your-anon-key>' >> .env.local"
    echo ""
    echo "   Get anon key with: supabase status | grep 'anon key'"
fi

echo ""
echo "🧪 Quick Test:"
echo ""
echo "1. Make sure Supabase is running:"
echo "   supabase status"
echo ""
echo "2. Start your web app:"
echo "   cd web && npm run dev"
echo ""
echo "3. Request a magic link in the app"
echo ""
echo "4. Check Inbucket for the email:"
if [ -n "$INBUCKET_URL" ]; then
    echo "   $INBUCKET_URL"
else
    echo "   http://localhost:54324"
fi
echo ""
echo "📖 Full guide: docs/LOCAL_MAGIC_LINK_SETUP.md"
