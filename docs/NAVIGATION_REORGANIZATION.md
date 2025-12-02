# Navigation Reorganization Proposal

## Current Issues

1. **Inconsistent Layouts**: Some pages (Rank Tracker, GSC Dashboard, GSC Keywords) don't use the Dashboard sidebar layout
2. **Nested Structure**: Rank Tracker appears after GSC section, but should be independent
3. **Missing Navigation**: Users lose navigation context when navigating to standalone pages

## Proposed Solution: Shared Layout Component

### Option 1: Unified Sidebar Layout (Recommended)

Create a shared `ProjectLayout.svelte` component that wraps all project-related pages with consistent sidebar navigation.

#### Navigation Structure:

```
📊 Core SEO
  ├─ Dashboard
  ├─ Results  
  ├─ Issues
  ├─ Recommendations
  └─ Link Graph

📈 Rank Tracking
  ├─ Rank Tracker
  ├─ Discover Keywords
  └─ Impact-First View

🔍 Google Search Console (only if connected)
  ├─ GSC Dashboard
  └─ GSC Keywords

⚙️ Project
  └─ Settings
```

#### Implementation Steps:

1. **Create `ProjectLayout.svelte`** component:
   - Contains the sidebar navigation
   - Wraps page content in `<main>` section
   - Handles active route highlighting
   - Shows/hides GSC section based on connection status

2. **Update Routes**:
   - `RankTracker.svelte` → Wrap with `ProjectLayout`
   - `GSCDashboard.svelte` → Wrap with `ProjectLayout`
   - `GSCKeywords.svelte` → Wrap with `ProjectLayout`
   - `ImpactFirstView.svelte` → Wrap with `ProjectLayout`
   - `ProjectView.svelte` → Already uses Dashboard component (can refactor)

3. **Navigation Groups**:
   - Use visual separators/dividers between groups
   - Group headers (optional): "Core SEO", "Rank Tracking", etc.
   - Consistent icon styling

### Option 2: Top Navigation Bar (Alternative)

If sidebar becomes too crowded, consider a horizontal top navigation with dropdowns:

```
[Logo] [Dashboard] [Results] [Issues] [Recommendations] [Rank Tracking ▼] [GSC ▼] [Settings]
```

Where dropdowns contain:
- **Rank Tracking**: Rank Tracker, Impact-First View
- **GSC**: GSC Dashboard, GSC Keywords (only if connected)

## Benefits

1. **Consistent UX**: All pages have the same navigation structure
2. **Better Organization**: Logical grouping of related features
3. **Always Accessible**: Rank Tracker visible regardless of GSC connection
4. **Scalable**: Easy to add new features to appropriate groups
5. **Clear Hierarchy**: Users understand feature relationships

## Discover Keywords Feature

Currently, "Discover Keywords" is implemented as a modal component (`KeywordDiscovery.svelte`) that opens from Rank Tracker. 

**Options for navigation:**

### Option A: Keep as Modal + Add to Navigation (Recommended)
- Keep the modal functionality within Rank Tracker (quick access)
- Also add as a standalone page/route: `/project/:id/discover-keywords`
- Add to sidebar navigation under "Rank Tracking" section
- Users can bookmark, share links, and access directly
- Icon: `Binoculars` (already implemented)

### Option B: Modal Only
- Keep current implementation (modal only)
- Accessible via button in Rank Tracker
- Not in sidebar navigation

### Option C: Page Only
- Convert to full page route only
- Remove modal functionality
- Always accessible from sidebar

**Recommendation: Option A** - Best of both worlds: quick access from Rank Tracker AND standalone page for power users.

## Implementation Priority

1. ✅ Create `ProjectLayout.svelte` component
2. ✅ Move sidebar navigation from `Dashboard.svelte` to `ProjectLayout.svelte`
3. ✅ Create `/project/:id/discover-keywords` route (standalone page)
4. ✅ Update `RankTracker.svelte` to use `ProjectLayout`
5. ✅ Update `GSCDashboard.svelte` to use `ProjectLayout`
6. ✅ Update `GSCKeywords.svelte` to use `ProjectLayout`
7. ✅ Update `ImpactFirstView.svelte` to use `ProjectLayout`
8. ✅ Refactor `ProjectView.svelte` to use `ProjectLayout` instead of `Dashboard.svelte`
9. ✅ Add "Discover Keywords" to sidebar navigation
10. ✅ Update route highlighting logic

## Visual Design Notes

- Use Lucide icons consistently:
  - 📊 Dashboard: `LayoutDashboard`
  - 📈 Rank Tracker: `TrendingUp` or `BarChart`
  - 🔍 GSC: `Search` or `Google`
  - ⚙️ Settings: `Settings`
- Add subtle background colors or borders to group sections
- Highlight active route with accent color
- Responsive: Collapse to hamburger menu on mobile

