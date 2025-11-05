<script>
  export let error = null;
</script>

{#if error}
  <div class="alert alert-error max-w-2xl mx-auto mt-4">
    <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <div>
      <h3 class="font-bold">Configuration Error</h3>
      <div class="text-sm">{error}</div>
      <div class="text-sm mt-2">
        {#if import.meta.env.DEV}
          <p><strong>Local Development:</strong> Create a <code>.env.local</code> file in the <code>web/</code> directory:</p>
          <pre class="bg-base-200 p-2 rounded mt-1 text-xs overflow-x-auto"><code>PUBLIC_SUPABASE_URL=http://127.0.0.1:54321
PUBLIC_SUPABASE_ANON_KEY=your-anon-key-here</code></pre>
          <p class="mt-2">Then restart the dev server (<code>npm run dev</code>).</p>
        {:else}
          <p>Please set the following environment variables in Vercel:</p>
          <ul class="list-disc list-inside mt-1">
            <li><code>PUBLIC_SUPABASE_URL</code></li>
            <li><code>PUBLIC_SUPABASE_ANON_KEY</code></li>
          </ul>
          <p class="mt-2">After setting them, redeploy your application.</p>
        {/if}
      </div>
    </div>
  </div>
{/if}

