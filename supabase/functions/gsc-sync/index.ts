// Supabase Edge Function: gsc-sync
// Schedules background synchronization of GSC data by calling the Barracuda API.
// Ensure the following environment variables are set in the function's environment:
// - CLOUD_RUN_API_URL (or API_BASE_URL): The base URL of the Barracuda API
// - GSC_SYNC_SECRET: Shared secret that authorizes access to the cron endpoint
// Optionally:
// - GSC_SYNC_LOOKBACK_DAYS: Override the number of days (default 30)

import 'https://deno.land/std@0.224.0/dotenv/load.ts';

const apiUrl =
  Deno.env.get('BARRACUDA_API_URL') ??
  Deno.env.get('CLOUD_RUN_API_URL') ??
  Deno.env.get('API_BASE_URL') ??
  'http://localhost:8080';

const cronSecret = Deno.env.get('GSC_SYNC_SECRET') ?? '';
const lookbackDays = Number(Deno.env.get('GSC_SYNC_LOOKBACK_DAYS') ?? '30');

export const config = {
  runtime: 'edge',
  verifyJWT: false,
};

export default async function handler(_request: Request): Promise<Response> {
  if (!cronSecret) {
    return new Response(
      JSON.stringify({ error: 'GSC_SYNC_SECRET not configured' }),
      { status: 500, headers: { 'Content-Type': 'application/json' } },
    );
  }

  const targetUrl = `${apiUrl.replace(/\/$/, '')}/api/internal/gsc/sync`;

  const response = await fetch(targetUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Cron-Secret': cronSecret,
    },
    body: JSON.stringify({
      lookback_days: Number.isFinite(lookbackDays) && lookbackDays > 0 ? lookbackDays : 30,
    }),
  });

  const text = await response.text();
  return new Response(text, {
    status: response.status,
    headers: {
      'Content-Type': 'application/json',
    },
  });
}
