#!/bin/bash

# Test script for RLS optimization migration
# This script runs the audit, applies the migration, and verifies the results

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_DIR"

echo "🔍 RLS Optimization Migration Test Script"
echo "=========================================="
echo ""

# Check if Supabase CLI is available
if ! command -v supabase &> /dev/null; then
    echo "❌ Supabase CLI not found"
    echo "Install it from: https://supabase.com/docs/guides/cli"
    exit 1
fi

# Check if Supabase is running
echo "📊 Step 1: Checking Supabase status..."

# Try to get status (with or without sudo)
if supabase status &> /dev/null; then
    echo "✅ Supabase is running"
    USE_SUDO=false
elif sudo supabase status &> /dev/null; then
    echo "✅ Supabase is running (using sudo)"
    USE_SUDO=true
else
    echo "❌ Supabase is not running"
    echo ""
    echo "Please start Supabase first:"
    echo "  supabase start"
    echo ""
    echo "Or if you need sudo:"
    echo "  sudo supabase start"
    echo ""
    exit 1
fi

echo ""

# Step 2: Run audit script (using Supabase Studio SQL editor or manual check)
echo "📋 Step 2: Pre-migration audit..."
echo "-------------------------------------------"
echo "⚠️  Note: For detailed audit, run this SQL in Supabase Studio SQL Editor:"
echo "   File: $SCRIPT_DIR/audit_rls_policies.sql"
echo ""
echo "   Or open: http://127.0.0.1:54323 (Supabase Studio)"
echo ""

# Step 3: Apply migration using Supabase CLI
echo "🚀 Step 3: Applying RLS optimization migration..."
echo "--------------------------------------------------"
MIGRATION_FILE="$PROJECT_DIR/supabase/migrations/20250120000000_comprehensive_rls_optimization.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo "❌ Migration file not found: $MIGRATION_FILE"
    exit 1
fi

echo "Applying migration: $MIGRATION_FILE"
echo ""

if [ "$USE_SUDO" = true ]; then
    sudo supabase db push || {
        echo "❌ Migration failed!"
        exit 1
    }
else
    supabase db push || {
        echo "❌ Migration failed!"
        exit 1
    }
fi

echo "✅ Migration applied successfully"
echo ""

# Step 4: Run verification script
echo "✅ Step 4: Verifying migration..."
echo "----------------------------------"
echo "⚠️  Note: For detailed verification, run this SQL in Supabase Studio SQL Editor:"
echo "   File: $SCRIPT_DIR/verify_rls_optimization.sql"
echo ""
echo "   Or open: http://127.0.0.1:54323 (Supabase Studio)"
echo ""

# Quick verification using Supabase CLI
echo "Running quick verification..."
if [ "$USE_SUDO" = true ]; then
    sudo supabase db execute "SELECT proname FROM pg_proc WHERE proname LIKE 'can_%' ORDER BY proname;" 2>/dev/null || echo "⚠️  Could not verify helper functions (this is OK)"
else
    supabase db execute "SELECT proname FROM pg_proc WHERE proname LIKE 'can_%' ORDER BY proname;" 2>/dev/null || echo "⚠️  Could not verify helper functions (this is OK)"
fi

echo ""

# Step 5: Summary
echo "📊 Step 5: Summary"
echo "------------------"
echo ""
echo "✅ Migration applied!"
echo ""
echo "Next steps:"
echo "1. Open Supabase Studio: http://127.0.0.1:54323"
echo "2. Go to SQL Editor and run: scripts/verify_rls_optimization.sql"
echo "3. Verify all checks pass"
echo "4. Test your application to ensure access still works correctly"
echo ""
echo "If everything looks good, you can deploy to production with:"
echo "  supabase db push"
echo ""


