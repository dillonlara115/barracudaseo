-- Lock billing columns on public.profiles.
--
-- The "Users can update their own profile" RLS policy only checks row ownership
-- ((select auth.uid()) = id). RLS cannot restrict which columns change, so a
-- signed-in user could set subscription_tier / team_size on their own row
-- straight through the REST API. Column privileges close that gap: browser
-- roles may only write display fields. Billing columns are written exclusively
-- by the API server using the service role (Stripe webhooks), which is
-- unaffected by these grants.

revoke insert, update on public.profiles from anon, authenticated;

-- Signup creates the row with id + display_name (web/src/lib/auth.js).
grant insert (id, display_name, avatar_url) on public.profiles to authenticated;

-- Users may edit only their display fields.
grant update (display_name, avatar_url) on public.profiles to authenticated;
