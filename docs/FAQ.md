# 🐟 Barracuda SEO — Frequently Asked Questions

## 1. What makes Barracuda different from tools like Screaming Frog?

Screaming Frog is a powerful desktop crawler—but it’s single-device, manual, and lacks team collaboration unless you export CSVs and spreadsheets.

Barracuda offers a hybrid model:

- **Fast local crawls (CLI — coming soon)**
- **Cloud-based dashboard for teams**
- **AI-powered issue explanations (coming soon)**
- **Built-in integrations** with Google Search Console, GA4, Microsoft Clarity, Slack, and Google Drive

Think of Barracuda as *Screaming Frog + team dashboard + automation layer*.

---

## 2. How is Barracuda different from SEMrush and Ahrefs?

SEMrush and Ahrefs run crawls on their servers, which means:

- You have limited control over crawl behavior  
- Crawl speeds are throttled  
- You don’t get raw crawl data  
- Pricing is tied to credits  
- They may miss pages blocked to outside IPs

Barracuda:

- Runs **local crawls via CLI (coming soon)** for full control  
- Allows uploading data to the cloud dashboard  
- Stores complete crawl datasets  
- Has **no credits or crawl caps**  
- Mirrors what your site serves to *real browsers*  

Barracuda is **your own crawler, your own infrastructure**.

---

## 3. Why choose Barracuda over Looker Studio dashboards?

Looker Studio is excellent for reporting—but it can’t crawl your site or diagnose SEO issues.

Barracuda does:

- Website crawling (CLI coming soon)  
- SEO issue detection  
- Data collection from GSC, GA4, and Clarity  
- Dashboard generation automatically  
- Full issue prioritization and insights

Looker Studio = reporting  
**Barracuda = crawling + diagnosing + reporting + recommendations**

---

## 4. Does Barracuda replace my current SEO tools?

Barracuda enhances your SEO stack without replacing everything.

**Replaces:**
- Screaming Frog for most audits (CLI coming soon)
- Manual CSV exports
- Custom Looker Studio builds
- GA4/GSC spreadsheet workflows

**Complements:**
- SEMrush/Ahrefs for keyword research
- GSC for search insights
- GA4 for user behavior
- Clarity for UX signals

Barracuda becomes your **technical SEO command center**.

---

## 5. How fast is the Barracuda crawler?

Once released, the CLI will be extremely fast due to:

- A Go-based worker-pool engine  
- Concurrent crawling  
- Intelligent robots.txt + sitemap handling  
- Lightweight architecture

Expected speed: **100–500 pages/min**, depending on target site performance.

---

## 6. What makes Barracuda’s issue detection unique?

Most crawlers simply report raw issues.

Barracuda:

- Detects issues  
- Prioritizes them by severity + impact  
- Groups them by URL structure  
- Connects crawl data with GSC/GA4  
- (Coming soon) Uses AI to explain root causes and exact fixes

This intelligence layer goes beyond traditional crawlers.

---

## 7. What can I do in the Barracuda Cloud Workspace?

The dashboard lets you:

- Save historical crawls  
- Compare issues over time  
- Invite team members  
- View page-level insights  
- Send issues to Slack  
- Push exports to Google Sheets  
- Connect GSC/GA4/Clarity for blended insights  
- View recommendations tied to real metrics

This makes Barracuda collaborative and data-driven.

---

## 8. How are crawls stored and secured?

Cloud-stored crawls are:

- Saved in Supabase  
- Protected by Row-Level Security  
- Accessible only to the project team  
- Supported by Cloud Run secrets for API keys  
- Fully isolated per project

Your data stays private and accessible only to your team.

---

## 9. How does billing work?

No credits. No crawl caps.

**Free tier:**
- 100-page crawl limit  
- Issues only (no recommendations)  
- No integrations  

**Pro tier:**
- Unlimited crawls (via CLI—coming soon)  
- Full recommendations  
- All integrations  
- Priority scoring  

**Team add-ons:**  
- Multi-user support for agencies and in-house teams

Billing uses secure Stripe subscriptions.

---

## 10. Who is Barracuda designed for?

Barracuda is ideal for:

- SEO agencies  
- Freelancers  
- In-house SEO teams  
- Developers  
- Technical SEOs  
- Anyone tired of managing SEO audits manually  

The dashboard is non-technical.  
The CLI (coming soon) is perfect for developers and power users.

---

## 11. Does Barracuda support team accounts?

Yes. You can:

- Invite collaborators  
- Assign roles  
- Share projects  
- Track crawl history together  
- Maintain secure access controls

Team collaboration is built into the platform.

---

## 12. What features are coming soon?

Our roadmap includes:

- **Full CLI release** for local crawls  
- **AI-powered issue insights & crawl summaries** ✓ **NOW AVAILABLE**  
- **Google Search Console integration** (traffic-based prioritization)  
- **AI-enhanced priority scoring** (combining crawl + GSC data)  
- **Slack + email audit summaries**  
- **Google Drive auto-sync**  
- **GA4/GSC blended insights**  
- **Template-level issue clustering**  
- **Automated priority scoring using traffic + crawl data**

Barracuda is rapidly evolving to give you an end-to-end technical SEO workflow.

---
