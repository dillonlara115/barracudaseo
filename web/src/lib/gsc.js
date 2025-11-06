const severityWeights = {
  error: 10,
  warning: 5,
  info: 1,
};

export function normalizeUrlForGSC(url = '') {
  if (!url) return '';
  const trimmed = url.trim().toLowerCase();
  if (!trimmed) return '';
  return trimmed.endsWith('/') ? trimmed.slice(0, -1) : trimmed;
}

function getSeverityWeight(severity) {
  return severityWeights[severity] ?? 1;
}

function calculateEnrichedPriority(issue, perf) {
  const basePriority = getSeverityWeight(issue?.severity);

  const impressions = perf?.impressions ?? 0;
  const ctrRatio = perf?.ctr ?? 0;
  const position = perf?.position ?? 0;

  let trafficMultiplier = 1.0;
  if (impressions > 10000) {
    trafficMultiplier = 3.0;
  } else if (impressions > 1000) {
    trafficMultiplier = 2.0;
  } else if (impressions < 100) {
    trafficMultiplier = 0.5;
  }

  const ctrPercent = ctrRatio * 100;
  let ctrMultiplier = 1.0;
  if (impressions > 1000 && ctrPercent < 2.0) {
    ctrMultiplier = 1.5;
  }

  let positionMultiplier = 1.0;
  if (impressions > 500 && position > 10 && position < 20) {
    positionMultiplier = 1.3;
  }

  return Number((basePriority * trafficMultiplier * ctrMultiplier * positionMultiplier).toFixed(2));
}

function generateRecommendationReason(perf) {
  const impressions = perf?.impressions ?? 0;
  const ctrPercent = (perf?.ctr ?? 0) * 100;

  if (impressions > 10000) {
    return `This page has high search visibility (${Math.round(impressions).toLocaleString()} impressions). Fixing this issue could significantly impact performance.`;
  }

  if (impressions > 1000) {
    if (ctrPercent < 2.0) {
      return `This page has solid visibility (${Math.round(impressions).toLocaleString()} impressions) but a low CTR (${ctrPercent.toFixed(1)}%). Optimizing could improve clicks.`;
    }
    return `This page attracts a meaningful amount of searches (${Math.round(impressions).toLocaleString()} impressions). Improvements here are likely to move the needle.`;
  }

  if (impressions < 100) {
    return 'This page currently has limited visibility. Fixing this keeps the page healthy as traffic grows.';
  }

  return null;
}

function extractMetrics(row = {}) {
  const metrics = row.metrics || {};
  const impressions = Number(metrics.impressions ?? 0);
  const clicks = Number(metrics.clicks ?? 0);
  const ctr = Number(metrics.ctr ?? 0);
  const position = Number(metrics.position ?? 0);
  const pageUrl = typeof row.dimension_value === 'string' ? row.dimension_value : '';

  return {
    impressions: Math.max(0, Math.round(impressions)),
    clicks: Math.max(0, Math.round(clicks)),
    ctr,
    position,
    page_url: pageUrl,
    top_queries: Array.isArray(row.top_queries) ? row.top_queries : [],
  };
}

export function buildEnrichedIssues(issues = [], pageRows = []) {
  if (!issues.length || !pageRows.length) {
    return [];
  }

  const rowsByUrl = new Map();
  pageRows.forEach((row) => {
    const pageUrl = typeof row.dimension_value === 'string' ? row.dimension_value : '';
    const normalized = normalizeUrlForGSC(pageUrl);
    if (!normalized) return;
    rowsByUrl.set(normalized, extractMetrics(row));
  });

  return issues.map((issue) => {
    const issueUrl = issue?.url ?? '';
    const normalizedIssueUrl = normalizeUrlForGSC(issueUrl);
    const metrics = normalizedIssueUrl ? rowsByUrl.get(normalizedIssueUrl) : null;

    if (!metrics) {
      return {
        issue,
        enriched_priority: Number(getSeverityWeight(issue?.severity)),
      };
    }

    const enrichedPriority = calculateEnrichedPriority(issue, metrics);
    const recommendationReason = generateRecommendationReason(metrics);

    return {
      issue,
      gsc_performance: metrics,
      enriched_priority: enrichedPriority,
      recommendation_reason: recommendationReason ?? undefined,
    };
  });
}
