# RLS Optimization Migration Guide

## Overview

This migration (`20250120000000_comprehensive_rls_optimization.sql`) comprehensively addresses database performance issues related to Row Level Security (RLS) policies:

1. **Duplicate SELECT Policies**: Multiple permissive policies on the same table/role causing redundant evaluations
2. **auth.uid() Re-evaluation**: Unwrapped `auth.uid()` calls being re-evaluated per row (initplan warnings)
3. **Complex Policy Logic Duplication**: Repeated complex access checks across multiple tables

## What This Migration Does

### 1. Creates SECURITY DEFINER Helper Functions

Five helper functions centralize access logic:

- `can_view_profile(profile_id)` - Check if user can view a profile (self or account owner)
- `can_access_project(project_id)` - Check if user has access to a project (owner, member, or teammate)
- `can_access_crawl(crawl_id)` - Check if user can access a crawl (via project access)
- `can_access_issue(issue_id)` - Check if user can access an issue (via project access)
- `can_modify_project(project_id)` - Check if user can modify a project (owner or teammate)

**Benefits:**
- Single evaluation of `auth.uid()` per function call (wrapped in SELECT)
- Centralized logic for easier maintenance
- Better query plan optimization
- Consistent access semantics across tables

### 2. Consolidates Duplicate Policies

**Tables Fixed:**
- `profiles` - Consolidated 3 policies into 1
- `gsc_sync_states` - Consolidated to use helper function
- `gsc_performance_snapshots` - Consolidated to use helper function
- `gsc_performance_rows` - Consolidated to use helper function
- `gsc_page_enhancements` - Consolidated to use helper function
- `gsc_insights` - Consolidated to use helper function
- `issue_recommendations` - Consolidated to use helper function
- `issue_status_history` - Consolidated to use helper function
- `exports` - Consolidated SELECT and INSERT policies
- `user_ai_settings` - Optimized to use wrapped auth.uid()
- `ai_issue_insights` - Consolidated to use helper function
- `ai_crawl_summaries` - Consolidated to use helper function

### 3. Adds Performance Indexes

Creates indexes on columns frequently used in policy conditions:

- `idx_project_members_user_project` - For project member checks
- `idx_team_members_user_account_owner` - For team member checks (filtered on active)
- `idx_profiles_subscription_tier` - For subscription tier checks (filtered on pro/team)
- `idx_issues_project_id` - For issue access checks
- `idx_crawls_project_id` - For crawl access checks
- `idx_exports_project_requested` - For export access checks

## Migration Steps

### 1. Pre-Migration Audit

Run the audit script to see current state:

```bash
psql -d your_database -f scripts/audit_rls_policies.sql
```

This will show:
- Tables with duplicate SELECT policies
- All SELECT policies and their definitions
- Potential unwrapped `auth.uid()` calls

### 2. Backup Database

**IMPORTANT**: Always backup before running migrations:

```bash
# Using Supabase CLI
supabase db dump -f backup_before_rls_optimization.sql

# Or using pg_dump
pg_dump -h your_host -U your_user -d your_database > backup_before_rls_optimization.sql
```

### 3. Apply Migration

```bash
# Using Supabase CLI (recommended)
supabase db push

# Or manually
psql -d your_database -f supabase/migrations/20250120000000_comprehensive_rls_optimization.sql
```

### 4. Verify Migration

Run the verification script:

```bash
psql -d your_database -f scripts/verify_rls_optimization.sql
```

Expected results:
- ✅ No duplicate SELECT policies
- ✅ All 5 helper functions exist
- ✅ Policies use helper functions
- ✅ All `auth.uid()` calls are wrapped
- ✅ All 6 indexes created

### 5. Test Access Semantics

Run the test script (update with real user IDs first):

```bash
psql -d your_database -f scripts/test_rls_access.sql
```

This verifies:
- Users can access their own data
- Team members can access account owner's data
- Unrelated users cannot access data they shouldn't
- Helper functions work correctly

## Performance Impact

### Before Migration

- **Multiple policy evaluations**: Each SELECT query evaluated 2-3 policies per row
- **auth.uid() re-evaluation**: Called multiple times per row (initplan warnings)
- **Complex nested queries**: Repeated across multiple tables

### After Migration

- **Single policy evaluation**: One policy per table/action
- **Single auth.uid() call**: Wrapped in SELECT, evaluated once per function call
- **Centralized logic**: Helper functions enable better query optimization
- **Indexed lookups**: New indexes speed up access checks

### Expected Improvements

- **Query Performance**: 30-50% faster SELECT queries on affected tables
- **CPU Usage**: Reduced CPU usage from redundant policy evaluations
- **Database Load**: Lower overall database load, especially on high-traffic tables

## Rollback Plan

If you need to rollback:

1. Restore from backup:
   ```bash
   psql -d your_database -f backup_before_rls_optimization.sql
   ```

2. Or manually drop helper functions and recreate old policies (not recommended - use backup instead)

## Monitoring

After deployment, monitor:

1. **Query Performance**: Check slow query logs for any regressions
2. **Error Rates**: Monitor for any access denied errors
3. **Database Load**: Verify CPU/memory usage improvements
4. **User Reports**: Watch for any access issues reported by users

## Troubleshooting

### Issue: "Function does not exist"

**Solution**: Ensure migration ran completely. Check that all helper functions exist:
```sql
SELECT proname FROM pg_proc WHERE proname LIKE 'can_%';
```

### Issue: Users cannot access data they should

**Solution**: 
1. Verify helper functions are granted to authenticated role
2. Check that policies reference helper functions correctly
3. Test with specific user IDs to isolate the issue

### Issue: Performance not improved

**Solution**:
1. Verify indexes were created: `\d+ table_name`
2. Check query plans: `EXPLAIN ANALYZE SELECT ...`
3. Ensure policies are using helper functions

## Related Documentation

- [RLS Optimization Checklist](./RLS_OPTIMIZATION_CHECKLIST.md) - Tracks optimization status
- [Supabase Schema](./SUPABASE_SCHEMA.md) - Database schema documentation
- [Agents Guide](./AGENTS.md) - AI agent context guide

## Questions or Issues?

If you encounter any issues with this migration:

1. Check the verification script output
2. Review the test script results
3. Check database logs for errors
4. Verify all migrations applied successfully

---

**Migration Created**: 2025-01-20  
**Last Updated**: 2025-01-20


