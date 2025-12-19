drop extension if exists "pg_net";

drop policy if exists "Project members and teammates can view integrations" on "public"."api_integrations";

drop policy if exists "Project owners and teammates can manage integrations" on "public"."api_integrations";

drop policy if exists "Team members can view account owner profile" on "public"."profiles";

drop policy if exists "Users can view their own profile" on "public"."profiles";

drop policy if exists "Users can create AI crawl summaries for accessible crawls" on "public"."ai_crawl_summaries";

drop policy if exists "Users can view AI crawl summaries for accessible crawls" on "public"."ai_crawl_summaries";

drop policy if exists "Users can create their own AI issue insights" on "public"."ai_issue_insights";

drop policy if exists "Users can view their own AI issue insights" on "public"."ai_issue_insights";

drop policy if exists "Project members and teammates can create crawls" on "public"."crawls";

drop policy if exists "Project members and teammates can delete crawls" on "public"."crawls";

drop policy if exists "Project members and teammates can update crawls" on "public"."crawls";

drop policy if exists "Project members and teammates can view crawls" on "public"."crawls";

drop policy if exists "Project members and teammates can create exports" on "public"."exports";

drop policy if exists "Project members and teammates can view exports" on "public"."exports";

drop policy if exists "Project members can view gsc insights" on "public"."gsc_insights";

drop policy if exists "Project members can view gsc page enhancements" on "public"."gsc_page_enhancements";

drop policy if exists "Project members can view gsc performance rows" on "public"."gsc_performance_rows";

drop policy if exists "Project members can view gsc snapshots" on "public"."gsc_performance_snapshots";

drop policy if exists "Project members can view gsc sync state" on "public"."gsc_sync_states";

drop policy if exists "Project members and teammates can view gsc insights" on "public"."gsc_insights";

drop policy if exists "Project members and teammates can view gsc page enhancements" on "public"."gsc_page_enhancements";

drop policy if exists "Project members and teammates can view gsc performance rows" on "public"."gsc_performance_rows";

drop policy if exists "Project members and teammates can view gsc snapshots" on "public"."gsc_performance_snapshots";

drop policy if exists "Project members and teammates can view gsc sync state" on "public"."gsc_sync_states";

drop policy if exists "Project members and teammates can create recommendations" on "public"."issue_recommendations";

drop policy if exists "Project members and teammates can view recommendations" on "public"."issue_recommendations";

drop policy if exists "Project members and teammates can create status history" on "public"."issue_status_history";

drop policy if exists "Project members and teammates can view status history" on "public"."issue_status_history";

drop policy if exists "Project members and teammates can insert issues" on "public"."issues";

drop policy if exists "Project members and teammates can update issue status" on "public"."issues";

drop policy if exists "Project members and teammates can view issues" on "public"."issues";

drop policy if exists "Project members can insert keyword snapshots" on "public"."keyword_rank_snapshots";

drop policy if exists "Project members can view keyword snapshots" on "public"."keyword_rank_snapshots";

drop policy if exists "Project members can insert keyword tasks" on "public"."keyword_tasks";

drop policy if exists "Project members can update keyword tasks" on "public"."keyword_tasks";

drop policy if exists "Project members can view keyword tasks" on "public"."keyword_tasks";

drop policy if exists "Project members can delete keywords" on "public"."keywords";

drop policy if exists "Project members can insert keywords" on "public"."keywords";

drop policy if exists "Project members can update keywords" on "public"."keywords";

drop policy if exists "Project members can view keywords" on "public"."keywords";

drop policy if exists "Project members and teammates can insert pages" on "public"."pages";

drop policy if exists "Project members and teammates can view pages" on "public"."pages";

drop policy if exists "Project members and teammates can view projects" on "public"."projects";

drop policy if exists "Project owners and teammates can delete projects" on "public"."projects";

drop policy if exists "Project owners and teammates can update projects" on "public"."projects";

drop policy if exists "Users can create public reports for their projects" on "public"."public_reports";

drop policy if exists "Users can delete their own public reports" on "public"."public_reports";

drop policy if exists "Users can update their own public reports" on "public"."public_reports";

drop policy if exists "Users can view their own public reports" on "public"."public_reports";

drop policy if exists "Account owners can delete team members" on "public"."team_members";

drop policy if exists "Account owners can insert team members" on "public"."team_members";

drop policy if exists "Account owners can update team members" on "public"."team_members";

drop policy if exists "Users can view team members" on "public"."team_members";

drop policy if exists "Users can insert their own AI settings" on "public"."user_ai_settings";

drop policy if exists "Users can update their own AI settings" on "public"."user_ai_settings";

drop policy if exists "Users can view their own AI settings" on "public"."user_ai_settings";

drop function if exists "public"."can_access_crawl"(crawl_id uuid);

drop function if exists "public"."can_access_issue"(issue_id bigint);

drop function if exists "public"."can_access_project"(project_id uuid);

drop function if exists "public"."can_modify_project"(project_id uuid);

drop index if exists "public"."idx_crawls_project_id";

drop index if exists "public"."idx_exports_project_requested";

drop index if exists "public"."idx_issues_project_id";

drop index if exists "public"."idx_team_members_user_account_owner";

set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.can_view_profile(profile_id uuid)
 RETURNS boolean
 LANGUAGE sql
 STABLE SECURITY DEFINER
 SET search_path TO 'public'
AS $function$
  SELECT (SELECT auth.uid()) = profile_id
    OR EXISTS (
      SELECT 1
      FROM public.team_members tm
      WHERE tm.user_id = (SELECT auth.uid())
        AND tm.account_owner_id = profile_id
        AND tm.status = 'active'
    );
$function$
;


  create policy "Project members can create exports"
  on "public"."exports"
  as permissive
  for insert
  to public
with check ((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = exports.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))));



  create policy "Project members can view exports"
  on "public"."exports"
  as permissive
  for select
  to public
using ((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = exports.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))));



  create policy "Users can create AI crawl summaries for accessible crawls"
  on "public"."ai_crawl_summaries"
  as permissive
  for insert
  to public
with check (((auth.uid() = user_id) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = ai_crawl_summaries.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Users can view AI crawl summaries for accessible crawls"
  on "public"."ai_crawl_summaries"
  as permissive
  for select
  to public
using (((auth.uid() = user_id) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = ai_crawl_summaries.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Users can create their own AI issue insights"
  on "public"."ai_issue_insights"
  as permissive
  for insert
  to public
with check (((auth.uid() = user_id) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = ai_issue_insights.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Users can view their own AI issue insights"
  on "public"."ai_issue_insights"
  as permissive
  for select
  to public
using (((auth.uid() = user_id) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = ai_issue_insights.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Project members and teammates can create crawls"
  on "public"."crawls"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = crawls.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = crawls.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = p.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can delete crawls"
  on "public"."crawls"
  as permissive
  for delete
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = crawls.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = crawls.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = p.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can update crawls"
  on "public"."crawls"
  as permissive
  for update
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = crawls.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = crawls.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = p.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))))
with check (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = crawls.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = crawls.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = p.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view crawls"
  on "public"."crawls"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = crawls.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = crawls.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = p.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = crawls.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can create exports"
  on "public"."exports"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = exports.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = auth.uid())))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = exports.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = exports.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = auth.uid()) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view exports"
  on "public"."exports"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = exports.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = auth.uid())))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = exports.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = exports.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = auth.uid()) AND (tm.status = 'active'::text))))));



  create policy "Project members can view gsc insights"
  on "public"."gsc_insights"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = gsc_insights.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = gsc_insights.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view gsc page enhancements"
  on "public"."gsc_page_enhancements"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = gsc_page_enhancements.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = gsc_page_enhancements.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view gsc performance rows"
  on "public"."gsc_performance_rows"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = gsc_performance_rows.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = gsc_performance_rows.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view gsc snapshots"
  on "public"."gsc_performance_snapshots"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = gsc_performance_snapshots.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = gsc_performance_snapshots.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view gsc sync state"
  on "public"."gsc_sync_states"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = gsc_sync_states.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = gsc_sync_states.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members and teammates can create recommendations"
  on "public"."issue_recommendations"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM (public.issues i
     JOIN public.project_members pm ON ((pm.project_id = i.project_id)))
  WHERE ((i.id = issue_recommendations.issue_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.team_members tm1 ON ((tm1.user_id = auth.uid())))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((i.id = issue_recommendations.issue_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((i.id = issue_recommendations.issue_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = auth.uid()) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view recommendations"
  on "public"."issue_recommendations"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM (public.issues i
     JOIN public.project_members pm ON ((pm.project_id = i.project_id)))
  WHERE ((i.id = issue_recommendations.issue_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.team_members tm1 ON ((tm1.user_id = auth.uid())))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((i.id = issue_recommendations.issue_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((i.id = issue_recommendations.issue_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = auth.uid()) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can create status history"
  on "public"."issue_status_history"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM (public.issues i
     JOIN public.project_members pm ON ((pm.project_id = i.project_id)))
  WHERE ((i.id = issue_status_history.issue_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.team_members tm1 ON ((tm1.user_id = auth.uid())))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((i.id = issue_status_history.issue_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((i.id = issue_status_history.issue_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = auth.uid()) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view status history"
  on "public"."issue_status_history"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM (public.issues i
     JOIN public.project_members pm ON ((pm.project_id = i.project_id)))
  WHERE ((i.id = issue_status_history.issue_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.team_members tm1 ON ((tm1.user_id = auth.uid())))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((i.id = issue_status_history.issue_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (((public.issues i
     JOIN public.projects p ON ((p.id = i.project_id)))
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((i.id = issue_status_history.issue_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = auth.uid()) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can insert issues"
  on "public"."issues"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = issues.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = issues.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = issues.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can update issue status"
  on "public"."issues"
  as permissive
  for update
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = issues.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = issues.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = issues.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))))
with check (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = issues.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = issues.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = issues.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view issues"
  on "public"."issues"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = issues.project_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((p.id = issues.project_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM ((public.projects p
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((p.id = issues.project_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members can insert keyword snapshots"
  on "public"."keyword_rank_snapshots"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.project_members pm ON ((pm.project_id = k.project_id)))
  WHERE ((k.id = keyword_rank_snapshots.keyword_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.projects p ON ((p.id = k.project_id)))
  WHERE ((k.id = keyword_rank_snapshots.keyword_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view keyword snapshots"
  on "public"."keyword_rank_snapshots"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.project_members pm ON ((pm.project_id = k.project_id)))
  WHERE ((k.id = keyword_rank_snapshots.keyword_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.projects p ON ((p.id = k.project_id)))
  WHERE ((k.id = keyword_rank_snapshots.keyword_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can insert keyword tasks"
  on "public"."keyword_tasks"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.project_members pm ON ((pm.project_id = k.project_id)))
  WHERE ((k.id = keyword_tasks.keyword_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.projects p ON ((p.id = k.project_id)))
  WHERE ((k.id = keyword_tasks.keyword_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can update keyword tasks"
  on "public"."keyword_tasks"
  as permissive
  for update
  to public
using (((EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.project_members pm ON ((pm.project_id = k.project_id)))
  WHERE ((k.id = keyword_tasks.keyword_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.projects p ON ((p.id = k.project_id)))
  WHERE ((k.id = keyword_tasks.keyword_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view keyword tasks"
  on "public"."keyword_tasks"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.project_members pm ON ((pm.project_id = k.project_id)))
  WHERE ((k.id = keyword_tasks.keyword_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM (public.keywords k
     JOIN public.projects p ON ((p.id = k.project_id)))
  WHERE ((k.id = keyword_tasks.keyword_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can delete keywords"
  on "public"."keywords"
  as permissive
  for delete
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = keywords.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = keywords.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can insert keywords"
  on "public"."keywords"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = keywords.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = keywords.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can update keywords"
  on "public"."keywords"
  as permissive
  for update
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = keywords.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = keywords.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members can view keywords"
  on "public"."keywords"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = keywords.project_id) AND (pm.user_id = auth.uid())))) OR (EXISTS ( SELECT 1
   FROM public.projects p
  WHERE ((p.id = keywords.project_id) AND (p.owner_id = auth.uid()))))));



  create policy "Project members and teammates can insert pages"
  on "public"."pages"
  as permissive
  for insert
  to public
with check (((EXISTS ( SELECT 1
   FROM (public.crawls c
     JOIN public.project_members pm ON ((pm.project_id = c.project_id)))
  WHERE ((c.id = pages.crawl_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM (((public.crawls c
     JOIN public.projects p ON ((p.id = c.project_id)))
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((c.id = pages.crawl_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (((public.crawls c
     JOIN public.projects p ON ((p.id = c.project_id)))
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((c.id = pages.crawl_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view pages"
  on "public"."pages"
  as permissive
  for select
  to public
using (((EXISTS ( SELECT 1
   FROM (public.crawls c
     JOIN public.project_members pm ON ((pm.project_id = c.project_id)))
  WHERE ((c.id = pages.crawl_id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM (((public.crawls c
     JOIN public.projects p ON ((p.id = c.project_id)))
     JOIN public.team_members tm1 ON ((tm1.user_id = ( SELECT auth.uid() AS uid))))
     JOIN public.team_members tm2 ON ((tm2.account_owner_id = tm1.account_owner_id)))
  WHERE ((c.id = pages.crawl_id) AND (tm2.user_id = p.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (((public.crawls c
     JOIN public.projects p ON ((p.id = c.project_id)))
     JOIN public.profiles prof ON ((prof.id = p.owner_id)))
     JOIN public.team_members tm ON ((tm.account_owner_id = p.owner_id)))
  WHERE ((c.id = pages.crawl_id) AND ((prof.subscription_tier = 'pro'::text) OR (prof.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project members and teammates can view projects"
  on "public"."projects"
  as permissive
  for select
  to public
using (((owner_id = ( SELECT auth.uid() AS uid)) OR (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = projects.id) AND (pm.user_id = ( SELECT auth.uid() AS uid))))) OR (EXISTS ( SELECT 1
   FROM (public.team_members tm1
     JOIN public.team_members tm2 ON ((tm1.account_owner_id = tm2.account_owner_id)))
  WHERE ((tm1.user_id = ( SELECT auth.uid() AS uid)) AND (tm2.user_id = projects.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (public.profiles p
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = ( SELECT auth.uid() AS uid)) AND ((p.subscription_tier = 'pro'::text) OR (p.subscription_tier = 'team'::text)) AND (tm.user_id = projects.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (public.profiles p
     JOIN public.team_members tm ON ((tm.account_owner_id = projects.owner_id)))
  WHERE ((p.id = projects.owner_id) AND ((p.subscription_tier = 'pro'::text) OR (p.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project owners and teammates can delete projects"
  on "public"."projects"
  as permissive
  for delete
  to public
using (((owner_id = ( SELECT auth.uid() AS uid)) OR (EXISTS ( SELECT 1
   FROM (public.team_members tm1
     JOIN public.team_members tm2 ON ((tm1.account_owner_id = tm2.account_owner_id)))
  WHERE ((tm1.user_id = ( SELECT auth.uid() AS uid)) AND (tm2.user_id = projects.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (public.profiles p
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = ( SELECT auth.uid() AS uid)) AND ((p.subscription_tier = 'pro'::text) OR (p.subscription_tier = 'team'::text)) AND (tm.user_id = projects.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (public.profiles p
     JOIN public.team_members tm ON ((tm.account_owner_id = projects.owner_id)))
  WHERE ((p.id = projects.owner_id) AND ((p.subscription_tier = 'pro'::text) OR (p.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Project owners and teammates can update projects"
  on "public"."projects"
  as permissive
  for update
  to public
using (((owner_id = ( SELECT auth.uid() AS uid)) OR (EXISTS ( SELECT 1
   FROM (public.team_members tm1
     JOIN public.team_members tm2 ON ((tm1.account_owner_id = tm2.account_owner_id)))
  WHERE ((tm1.user_id = ( SELECT auth.uid() AS uid)) AND (tm2.user_id = projects.owner_id) AND (tm1.status = 'active'::text) AND (tm2.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (public.profiles p
     JOIN public.team_members tm ON ((tm.account_owner_id = ( SELECT auth.uid() AS uid))))
  WHERE ((p.id = ( SELECT auth.uid() AS uid)) AND ((p.subscription_tier = 'pro'::text) OR (p.subscription_tier = 'team'::text)) AND (tm.user_id = projects.owner_id) AND (tm.status = 'active'::text)))) OR (EXISTS ( SELECT 1
   FROM (public.profiles p
     JOIN public.team_members tm ON ((tm.account_owner_id = projects.owner_id)))
  WHERE ((p.id = projects.owner_id) AND ((p.subscription_tier = 'pro'::text) OR (p.subscription_tier = 'team'::text)) AND (tm.user_id = ( SELECT auth.uid() AS uid)) AND (tm.status = 'active'::text))))));



  create policy "Users can create public reports for their projects"
  on "public"."public_reports"
  as permissive
  for insert
  to public
with check (((auth.uid() = created_by) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = public_reports.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Users can delete their own public reports"
  on "public"."public_reports"
  as permissive
  for delete
  to public
using (((auth.uid() = created_by) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = public_reports.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Users can update their own public reports"
  on "public"."public_reports"
  as permissive
  for update
  to public
using (((auth.uid() = created_by) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = public_reports.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Users can view their own public reports"
  on "public"."public_reports"
  as permissive
  for select
  to public
using (((auth.uid() = created_by) AND (EXISTS ( SELECT 1
   FROM public.project_members pm
  WHERE ((pm.project_id = public_reports.project_id) AND (pm.user_id = auth.uid()))))));



  create policy "Account owners can delete team members"
  on "public"."team_members"
  as permissive
  for delete
  to authenticated
using ((account_owner_id = ( SELECT auth.uid() AS uid)));



  create policy "Account owners can insert team members"
  on "public"."team_members"
  as permissive
  for insert
  to authenticated
with check ((account_owner_id = ( SELECT auth.uid() AS uid)));



  create policy "Account owners can update team members"
  on "public"."team_members"
  as permissive
  for update
  to authenticated
using ((account_owner_id = ( SELECT auth.uid() AS uid)));



  create policy "Users can view team members"
  on "public"."team_members"
  as permissive
  for select
  to authenticated
using (((account_owner_id = ( SELECT auth.uid() AS uid)) OR (user_id = ( SELECT auth.uid() AS uid))));



  create policy "Users can insert their own AI settings"
  on "public"."user_ai_settings"
  as permissive
  for insert
  to public
with check ((auth.uid() = user_id));



  create policy "Users can update their own AI settings"
  on "public"."user_ai_settings"
  as permissive
  for update
  to public
using ((auth.uid() = user_id));



  create policy "Users can view their own AI settings"
  on "public"."user_ai_settings"
  as permissive
  for select
  to public
using ((auth.uid() = user_id));

