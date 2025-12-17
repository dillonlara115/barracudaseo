drop extension if exists "pg_net";

drop policy "Project members and teammates can view integrations" on "public"."api_integrations";

drop policy "Project owners and teammates can manage integrations" on "public"."api_integrations";

drop policy "Team members can view account owner profile" on "public"."profiles";

drop policy "Project members and teammates can create crawls" on "public"."crawls";

drop policy "Project members and teammates can delete crawls" on "public"."crawls";

drop policy "Project members and teammates can update crawls" on "public"."crawls";

drop policy "Project members and teammates can view crawls" on "public"."crawls";

drop policy "Project members and teammates can insert issues" on "public"."issues";

drop policy "Project members and teammates can update issue status" on "public"."issues";

drop policy "Project members and teammates can view issues" on "public"."issues";

drop policy "Project members and teammates can insert pages" on "public"."pages";

drop policy "Project members and teammates can view pages" on "public"."pages";

drop policy "Project members and teammates can view projects" on "public"."projects";

drop policy "Project owners and teammates can delete projects" on "public"."projects";

drop policy "Project owners and teammates can update projects" on "public"."projects";

drop policy "Account owners can delete team members" on "public"."team_members";

drop policy "Account owners can insert team members" on "public"."team_members";

drop policy "Account owners can update team members" on "public"."team_members";

drop policy "Users can view team members" on "public"."team_members";


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



