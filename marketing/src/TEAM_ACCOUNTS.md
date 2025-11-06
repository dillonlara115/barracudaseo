# Barracuda Team Accounts & Pricing Add-on Proposal

*Last updated: {{today}}*

## 1. Overview

Barracuda’s **Team Accounts** feature will enable collaboration within the web dashboard by allowing multiple users to access and manage the same projects. This feature is primarily aimed at agencies, in-house SEO teams, and technical departments where several contributors work on shared sites or clients.

The team functionality should integrate seamlessly with **Supabase Auth** and **project_members** tables (as outlined in `SUPABASE_SCHEMA.md`), while tying into a **Pro-tier subscription model** with optional paid team add-ons.

---

## 2. Feature Summary

### Core Capabilities

* **Shared Project Access**: Multiple users can view and collaborate on the same projects.
* **Role-Based Permissions**:

  * **Owner**: Full permissions, including billing, invitations, and project deletion.
  * **Editor**: Can trigger crawls, manage issues, and view recommendations.
  * **Viewer**: Read-only access to dashboard and reports.
* **Project Sharing**: Invite users by email via Supabase Auth.
* **Audit Logs (Future)**: Track who triggered crawls, added notes, or changed statuses.

### Collaboration UX

* In the dashboard sidebar, a new **Team** tab lists all members with role and status.
* Owners can invite or remove members directly from the UI.
* Editors and viewers receive email invites and authenticate via Supabase.

---

## 3. Pricing Model Options

### Option A: Pro + Team Add-on (Recommended)

**Base plan**: Pro includes 1 user.
**Add-on**: Each additional user is billed monthly.

| Plan | Base Price | Users Included | Additional Users | Crawl Limit   | Integrations     | Notes                    |
| ---- | ---------- | -------------- | ---------------- | ------------- | ---------------- | ------------------------ |
| Free | $0         | 1              | —                | 100 pages     | None             | No team access           |
| Pro  | $29/mo     | 1              | +$5/user         | 10,000 pages  | All integrations | Includes recommendations |
| Team | Custom     | 5+             | $5/user          | 25,000+ pages | All integrations | For agencies and orgs    |

**Advantages:**

* Flexible scaling for small teams.
* Simple Stripe integration (base plan + quantity add-on).
* Encourages upgrades from solo users.

**Technical Implementation:**

* Store `team_size` and `subscription_tier` in Supabase.
* Use Stripe `quantity` pricing for seat-based billing.
* Limit active users per project based on plan.

### Option B: Separate Team Plan

Create a standalone “Team” plan that includes a set number of seats (e.g., 5 users).

**Pros:** Clear positioning for agencies and larger orgs.
**Cons:** Less flexible, more pricing maintenance.

### Option C: Usage-Based (Seats + Crawls)

Tie billing to both crawl volume and team size.
E.g., $0.10 per 1,000 pages crawled + $5 per user/month.

**Pros:** Aligns cost to usage.
**Cons:** Harder to predict and manage costs for users.

---

## 4. Account & Role Management

**Supabase Tables**:

* `projects`: Defines ownership and settings.
* `project_members`: Stores member roles and invites.
* `profiles`: Stores user metadata (name, avatar, plan, etc.).

**Endpoints (Cloud Run API):**

* `POST /api/v1/projects/:id/invite` → Invite user to project.
* `PATCH /api/v1/project_members/:id` → Update role.
* `DELETE /api/v1/project_members/:id` → Remove member.

**UI Components:**

* `TeamPanel.svelte` → Display team members, roles, invites.
* `InviteMemberModal.svelte` → Send invites.
* `RoleBadge.svelte` → Display permissions visually.

---

## 5. Billing Integration

**Implementation:**

* Stripe integration with seat-based pricing (via `quantity` parameter).
* Billing portal accessible only to project owners.
* Supabase webhook updates team size on successful payment events.

**Environment Variables:**

* `STRIPE_PRICE_PRO` → Base price.
* `STRIPE_PRICE_TEAM_ADDON` → Per-seat price.
* `STRIPE_WEBHOOK_SECRET` → Verifies billing events.

---

## 6. User Experience Flow

1. **Pro User** signs up and starts a project.
2. User navigates to **Team** tab → clicks “Invite teammate.”
3. Enters email → teammate receives invite via Supabase email.
4. On acceptance, teammate gains dashboard access with defined role.
5. If user exceeds seat limit → modal prompts to upgrade or add seats.

---

## 7. Future Enhancements

* **Audit trail**: Record team activity (who changed what and when).
* **Notes & comments**: Allow commenting on issues.
* **Slack integration**: Notify team channels when new crawls or issues are created.
* **Team analytics**: Track crawl performance across projects.

---

## 8. Recommendation

Start with **Option A (Pro + Team Add-on)**.
This approach keeps pricing flexible and implementation simple within the existing subscription system, while offering a clear upgrade path for agencies and teams.

---

**Next Steps:**

1. Implement Supabase role management logic and invitations.
2. Add Stripe quantity-based billing for user seats.
3. Update dashboard UI with Team tab and management modals.
4. Add marketing copy to Pricing and FAQ pages explaining seat-based billing.
5. Gather feedback from early Pro users on pricing and team usability.
