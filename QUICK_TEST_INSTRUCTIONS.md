# Quick Test Instructions for RLS Migration

Since Supabase is running (you just started it), here's how to test the migration:

## Option 1: Using Supabase CLI (Recommended)

```bash
cd /home/dillon/Sites/barracuda

# Apply the migration (this will apply all pending migrations including the new RLS optimization)
sudo supabase db push
```

This will:
- ✅ Apply the new migration: `20250120000000_comprehensive_rls_optimization.sql`
- ✅ Create helper functions
- ✅ Consolidate duplicate policies
- ✅ Add performance indexes

## Option 2: Using Supabase Studio SQL Editor

1. **Open Supabase Studio**: http://127.0.0.1:54323
2. **Go to SQL Editor** (left sidebar)
3. **Copy and paste** the migration file contents:
   ```bash
   cat supabase/migrations/20250120000000_comprehensive_rls_optimization.sql
   ```
4. **Click "Run"** to execute

## Verify the Migration

After applying, verify it worked:

### Quick Check (in Supabase Studio SQL Editor):

```sql
-- Check helper functions exist
SELECT proname FROM pg_proc WHERE proname LIKE 'can_%' ORDER BY proname;
-- Should return: can_access_crawl, can_access_issue, can_access_project, can_modify_project, can_view_profile

-- Check for duplicate SELECT policies (should return 0 rows)
SELECT 
  tablename,
  COUNT(*) as policy_count
FROM pg_policies
WHERE schemaname = 'public'
  AND cmd = 'SELECT'
  AND tablename IN ('profiles', 'gsc_sync_states', 'issue_recommendations', 'exports')
GROUP BY tablename
HAVING COUNT(*) > 1;
```

### Full Verification

Run the verification script in Supabase Studio SQL Editor:
- Open: `scripts/verify_rls_optimization.sql`
- Copy contents and run in SQL Editor

## Expected Results

After migration, you should see:

✅ **5 helper functions created**:
- `can_view_profile`
- `can_access_project`
- `can_access_crawl`
- `can_access_issue`
- `can_modify_project`

✅ **No duplicate SELECT policies** (each table has exactly 1 SELECT policy)

✅ **6 performance indexes created**

✅ **All policies use helper functions**

## Test Your Application

After migration, test that your app still works:

1. ✅ Login/logout
2. ✅ View projects
3. ✅ View crawl data
4. ✅ View issues and recommendations
5. ✅ View GSC data (if configured)
6. ✅ View AI insights (if configured)

## If Something Goes Wrong

**Rollback** (if needed):
```bash
# Reset database (WARNING: Deletes all data!)
sudo supabase db reset
```

Or restore from backup if you created one.

## Next Steps

Once local testing is successful:
1. ✅ Review verification output
2. ✅ Test application functionality  
3. ✅ Deploy to production: `supabase db push` (to remote)

---

**Migration File**: `supabase/migrations/20250120000000_comprehensive_rls_optimization.sql`
**Created**: 2025-01-20


