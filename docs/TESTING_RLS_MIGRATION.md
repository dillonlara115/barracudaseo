# Testing RLS Optimization Migration

This guide walks you through testing the RLS optimization migration locally.

## Prerequisites

1. **Docker** must be running
2. **Supabase CLI** installed (`which supabase` should work)
3. **Local Supabase** instance running

## Quick Start

### Option 1: Automated Test Script (Recommended)

```bash
cd /home/dillon/Sites/barracuda

# Make sure Supabase is running
supabase start

# Run the automated test script
./scripts/test-rls-migration.sh
```

This script will:
1. ✅ Check Supabase status
2. ✅ Run pre-migration audit
3. ✅ Apply the migration
4. ✅ Verify the migration
5. ✅ Show summary

### Option 2: Manual Step-by-Step

#### Step 1: Start Supabase

```bash
cd /home/dillon/Sites/barracuda
supabase start
```

Wait for Supabase to fully start, then note the **DB URL** from the output.

#### Step 2: Run Pre-Migration Audit

```bash
# Get the DB URL from supabase status
DB_URL=$(supabase status | grep "DB URL" | awk '{print $3}')

# Run audit
psql "$DB_URL" -f scripts/audit_rls_policies.sql
```

**What to look for:**
- Tables with multiple SELECT policies (these will be fixed)
- Unwrapped `auth.uid()` calls (these will be optimized)

#### Step 3: Apply Migration

```bash
# Apply the migration
psql "$DB_URL" -f supabase/migrations/20250120000000_comprehensive_rls_optimization.sql
```

**Expected output:**
- `CREATE FUNCTION` messages for helper functions
- `DROP POLICY` messages for old policies
- `CREATE POLICY` messages for new consolidated policies
- `CREATE INDEX` messages for new indexes

#### Step 4: Verify Migration

```bash
# Run verification
psql "$DB_URL" -f scripts/verify_rls_optimization.sql
```

**Expected results:**
- ✅ No duplicate SELECT policies
- ✅ All 5 helper functions exist
- ✅ Policies use helper functions
- ✅ All `auth.uid()` calls are wrapped
- ✅ All 6 indexes created

#### Step 5: Test Access Semantics

```bash
# Run access tests (update user IDs first!)
psql "$DB_URL" -f scripts/test_rls_access.sql
```

**Note:** Update the test script with real user IDs from your database:
```sql
-- Find user IDs:
SELECT id, email FROM auth.users LIMIT 5;
```

## Troubleshooting

### Issue: "Supabase is not running"

**Solution:**
```bash
# Start Supabase
supabase start

# If that fails, check Docker
docker ps

# If Docker isn't running, start it first
```

### Issue: "Permission denied" for Docker

**Solution:**
```bash
# Add your user to docker group (requires logout/login)
sudo usermod -aG docker $USER

# Or use sudo (not recommended for development)
sudo supabase start
```

### Issue: Migration fails with "function already exists"

**Solution:** This is OK! The migration uses `CREATE OR REPLACE FUNCTION`, so it's safe to run multiple times. If you see this, the migration is idempotent.

### Issue: "Policy does not exist" warnings

**Solution:** These are expected if policies were already dropped or don't exist. The migration uses `DROP POLICY IF EXISTS`, so warnings are harmless.

### Issue: Verification shows duplicate policies

**Solution:** 
1. Check if migration ran completely
2. Verify no errors occurred during migration
3. Check if there are other migrations that create policies
4. Run migration again (it's idempotent)

## What to Test

After migration, test these scenarios:

### 1. Profile Access
- ✅ User can view their own profile
- ✅ Team member can view account owner's profile
- ✅ Unrelated user cannot view other profiles

### 2. Project Access
- ✅ Project owner can view their projects
- ✅ Project member can view projects they're members of
- ✅ Team members can view projects via team access

### 3. Data Access
- ✅ Users can query GSC tables for accessible projects
- ✅ Users can query issue tables for accessible projects
- ✅ Users can query exports for accessible projects
- ✅ Users can query AI tables for their own data

### 4. Application Testing
- ✅ Login and view dashboard
- ✅ View projects
- ✅ View crawl data
- ✅ View issues and recommendations
- ✅ View GSC data (if configured)
- ✅ View AI insights (if configured)

## Performance Verification

After migration, you should see:

1. **Fewer policy evaluations**: Check with `EXPLAIN ANALYZE` on SELECT queries
2. **Faster queries**: Compare query times before/after
3. **Lower CPU usage**: Monitor database CPU usage

Example performance check:
```sql
-- Before migration: Multiple policy evaluations
EXPLAIN ANALYZE SELECT * FROM profiles WHERE id = auth.uid();

-- After migration: Single policy evaluation with helper function
EXPLAIN ANALYZE SELECT * FROM profiles WHERE id = auth.uid();
```

## Rollback (If Needed)

If you need to rollback:

```bash
# Option 1: Reset database (WARNING: Deletes all data!)
supabase db reset

# Option 2: Restore from backup
psql "$DB_URL" -f backup_before_rls_optimization.sql
```

## Next Steps

Once local testing is successful:

1. ✅ Review all verification outputs
2. ✅ Test application functionality
3. ✅ Check for any access issues
4. ✅ Deploy to production: `supabase db push`

## Questions?

- Check `docs/RLS_OPTIMIZATION_MIGRATION.md` for detailed migration info
- Review `scripts/verify_rls_optimization.sql` output for issues
- Check Supabase logs: `supabase logs`


