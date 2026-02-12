# Launch Metrics Tracking

## Overview

This document defines metrics to track during public launch and ongoing growth. Metrics align with PRD Part VIII (Business Model) success criteria for Month 12: 1,500 WAU, $10K MRR, 200 paying customers.

## Primary Metrics

### GitHub Engagement

**Stars:**
- Week 1 target: 500
- Month 1 target: 1000
- Month 3 target: 2000
- Month 6 target: 3500

**Clones (unique):**
- Proxy for actual usage (downloads)
- Week 1 target: 200
- Month 1 target: 1000
- Month 3 target: 3000

**Issues:**
- Week 1 target: 10+ (shows engagement)
- Quality over quantity (substantive issues)
- Track bug vs feature request ratio
- Response time target: <24 hours

**Pull Requests:**
- Month 1 target: First community PR
- Month 3 target: 5+ merged community PRs
- Track contributor count (unique contributors)

**Discussions:**
- Track question volume
- Track answer quality (resolved ratio)
- Target: 90% resolved within 48h

### Download Metrics

**Binary Downloads (from Releases):**
- Track per-platform (Linux, macOS, Windows)
- Week 1 target: 200
- Month 1 target: 1000
- Month 3 target: 5000

**Installation Success Rate:**
- Track via telemetry (opt-in)
- Target: >95% successful installs
- Monitor by platform

### Active Usage

**Weekly Active Users (WAU):**
- Track via opt-in telemetry
- Month 1 target: 100 WAU
- Month 3 target: 300 WAU
- Month 6 target: 750 WAU
- Month 12 target: 1,500 WAU (PRD target)

**Daily Active Users (DAU):**
- Track for engagement metrics
- Target DAU/WAU ratio: 0.4 (4 days per week)

**Commands Used:**
- Track most popular commands
- Track command success rate
- Identify pain points (failing commands)

**Generation Count:**
- Track `percepta generate` usage
- Track success rate (compilation)
- Track pattern storage rate (validated patterns)

**Observation Count:**
- Track `percepta observe` usage
- Track vision confidence scores
- Track validation success rate

### Community Feedback

**Sentiment:**
- Manual tracking from comments (HN, Reddit, GitHub)
- Categorize: Positive / Neutral / Negative
- Target: >70% positive sentiment

**Net Promoter Score (NPS):**
- Survey after 1 week of usage
- Question: "How likely are you to recommend Percepta (0-10)?"
- Month 12 target: NPS 60+ (PRD target)

**Feature Requests:**
- Track most requested features
- Prioritize roadmap based on requests
- Communicate roadmap publicly

## Secondary Metrics

### Documentation Engagement

**Docs Page Views:**
- Track most viewed docs pages
- Identify gaps (low views = poor discoverability)
- Track bounce rate (incomplete docs?)

**Time on Page:**
- Target: >2 minutes on getting-started
- Indicates thorough reading vs quick scan

**Search Queries:**
- If docs have search, track queries
- Identify missing documentation

### Content Performance

**Blog Post Views:**
- Announcing Percepta: Target 2000 views week 1
- Technical deep-dives: Target 500 views each
- Track referral sources (HN, Reddit, organic)

**Video Views:**
- Demo video: Target 1000 views week 1
- Track average watch time (target: >3 minutes)
- Track CTR to GitHub (target: >10%)

**Social Media:**
- Twitter followers (track growth)
- Tweet engagement (likes, retweets, replies)
- LinkedIn post views and engagement

### Launch Campaign Metrics

**Hacker News:**
- Peak rank (target: top 10)
- Total points (target: >100)
- Comment count (target: >30 substantive)
- Time on front page (target: >4 hours)

**Reddit:**
- r/embedded upvotes (target: >100)
- r/embedded comments (target: >20)
- Cross-post performance (r/rust, r/esp32)

**Email Campaign:**
- Alpha user announcement open rate (target: >60%)
- Click-through rate (target: >30%)

## Business Metrics (Future)

### Revenue (Post-HIL Farm Launch)

**Monthly Recurring Revenue (MRR):**
- Month 12 target: $10K MRR (PRD target)
- Track by tier (if multiple pricing tiers)
- Track growth rate (target: 15% MoM)

**Paying Customers:**
- Month 12 target: 200 paying customers (PRD target)
- Average Revenue Per User (ARPU): $50/month (from $10K/200)

**Free→Paid Conversion:**
- Target: 10% (PRD target)
- Track conversion time (days from install to paid)
- Track conversion triggers (what drives upgrade?)

**Churn Rate:**
- Target: <5% monthly churn (PRD target)
- Track reasons for churn (survey on cancellation)
- Track reactivation rate (churned users who return)

### Customer Satisfaction

**Support Tickets:**
- Track ticket volume
- Track response time (target: <4 hours)
- Track resolution time (target: <24 hours)
- Track customer satisfaction (CSAT) on resolution

**Feature Adoption:**
- Track usage of new features
- Identify unused features (remove or improve)
- Track power users vs casual users

## Tracking Implementation

### Tools

**GitHub:**
- GitHub Insights (built-in analytics)
- Manual tracking via Issues API
- GitHub Actions for release download counts

**Analytics (if website exists):**
- Google Analytics or Plausible (privacy-focused)
- Track docs page views, referrals, search queries

**Telemetry (opt-in):**
```bash
# User opts in during first run
percepta config set telemetry enabled

# Track anonymized usage:
# - Command executed (observe, assert, generate, etc.)
# - Success/failure
# - Execution time
# - Platform (Linux, macOS, Windows)
# - Version number

# Never track:
# - Generated code content
# - Hardware video/images
# - Personal information
# - File paths or device names
```

**Surveys:**
- Post-install survey (after 1 week)
- NPS survey (quarterly)
- Churn survey (on cancellation, if paid)

### Metrics Dashboard

**Manual Spreadsheet (Launch Week):**
```
Date | GH Stars | Downloads | HN Points | Sentiment | Notes
-----|----------|-----------|-----------|-----------|------
Feb 13 | 50 | 12 | 45 | Positive | Launch day
Feb 14 | 124 | 38 | 89 | Mixed | Front page 6h
Feb 15 | 203 | 67 | 112 | Positive | Reddit traction
...
```

**Automated Dashboard (Month 1+):**
- Build simple dashboard (Grafana or custom)
- Pull from GitHub API, telemetry database
- Display key metrics: WAU, stars, downloads, sentiment
- Weekly email report to team

## Reporting Cadence

### Daily (Launch Week)
- GitHub stars and download count
- HN/Reddit engagement
- Critical bug count
- Sentiment snapshot

### Weekly (Month 1-3)
- WAU (from telemetry)
- GitHub stars and clones
- Issue/PR count and resolution rate
- Feature request summary
- Blog post performance
- Sentiment analysis

### Monthly (Ongoing)
- All primary metrics
- NPS survey results
- Roadmap progress
- Community highlights (best contributions, discussions)
- Comparison to targets (on track / behind / ahead)

### Quarterly (Long-term)
- Business metrics (MRR, paying customers, churn)
- Strategic review (product-market fit, roadmap adjustment)
- Competitive analysis (Embedder updates, new competitors)

## Success Milestones

### Week 1 Milestones
- [ ] 500 GitHub stars
- [ ] 200 binary downloads
- [ ] 100 WAU (if telemetry live)
- [ ] HN front page (any rank)
- [ ] Positive overall sentiment (>70%)

### Month 1 Milestones
- [ ] 1000 GitHub stars
- [ ] 1000 binary downloads
- [ ] 300 WAU
- [ ] 10+ merged community PRs
- [ ] First "success story" (user validates real bug with Percepta)

### Month 3 Milestones
- [ ] 2000 GitHub stars
- [ ] 3000 downloads
- [ ] 500 WAU
- [ ] Cloud HIL farm beta launched
- [ ] First paying customer (HIL farm)

### Month 6 Milestones
- [ ] 3500 GitHub stars
- [ ] 750 WAU
- [ ] 50 paying customers
- [ ] $2.5K MRR
- [ ] Featured in embedded podcast or blog

### Month 12 Milestones (PRD Targets)
- [ ] 1,500 WAU
- [ ] 200 paying customers
- [ ] $10K MRR
- [ ] 10% free→paid conversion
- [ ] <5% churn rate
- [ ] NPS 60+

## Red Flags

Watch for these warning signs:

**Growth Issues:**
- Stars plateauing after week 1 (<100 new stars/week)
- Download count not translating to usage (low WAU relative to downloads)
- High bounce rate on docs (>70%)

**Quality Issues:**
- Bug report ratio >50% of all issues
- Community PRs not merging (lack of contribution)
- Negative sentiment increasing (>30% negative)

**Product-Market Fit:**
- Low engagement (DAU/WAU < 0.2)
- High churn (if paid product live)
- Feature requests all over the map (no clear pattern)
- Users asking "why not just use X?" repeatedly

**Sustainability:**
- Maintainer burnout (can't keep up with issues)
- Community toxicity (negative interactions)
- Competitors launching better alternatives

## Action Plans

### If growth slows (Month 2-3)
- Double down on content (blog posts, tutorials, case studies)
- Create video tutorials for YouTube
- Engage more actively in embedded communities
- Launch referral program or GitHub Sponsors
- Host office hours or live coding sessions

### If quality complaints increase
- Prioritize bug fixes over features
- Improve documentation (especially troubleshooting)
- Add more test coverage
- Consider beta testing program for new features
- Publish quality metrics (test coverage, bug resolution time)

### If product-market fit unclear
- Survey users extensively (what do you use Percepta for?)
- Analyze command usage (which features actually used?)
- Talk to churned users (why did they stop?)
- Consider pivoting positioning or target audience
- Validate assumptions with user interviews

## Notes

**Privacy:**
- All telemetry must be opt-in
- Anonymous by default (no PII)
- Clear documentation on what's tracked
- Easy opt-out mechanism
- Open source telemetry code for auditability

**Transparency:**
- Publish metrics publicly (if comfortable)
- Share roadmap based on feedback
- Acknowledge failures and learnings
- Celebrate community contributions

**Sustainability:**
- Metrics are means to end (building great tool)
- Don't chase vanity metrics (stars don't equal success)
- Focus on user value and satisfaction
- Long-term thinking over short-term spikes

---

**Next Steps:**
1. Implement opt-in telemetry in CLI
2. Set up manual tracking spreadsheet for launch week
3. Configure GitHub Insights monitoring
4. Create post-install survey (1 week delay)
5. Build automated dashboard (Month 1 goal)
