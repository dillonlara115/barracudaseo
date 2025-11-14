<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { userProfile, loadSubscriptionData, isProOrTeam } from '../lib/subscription.js';
  import Integrations from './Integrations.svelte';
  import UpgradePrompt from '../components/UpgradePrompt.svelte';

  let loading = true;
  let hasAccess = false;

  onMount(async () => {
    await loadSubscriptionData();
    
    // Subscribe to profile changes
    const unsubscribe = userProfile.subscribe((profile) => {
      hasAccess = isProOrTeam(profile);
      loading = false;
    });

    return unsubscribe;
  });

  function handleClose() {
    push('/');
  }
</script>

{#if loading}
  <div class="flex items-center justify-center min-h-screen">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
{:else if hasAccess}
  <Integrations />
{:else}
  <div class="container mx-auto p-6 max-w-4xl">
    <div class="mb-6">
      <button 
        class="btn btn-ghost btn-sm mb-4"
        on:click={() => push('/')}
      >
        ← Back to Projects
      </button>
    </div>
    <UpgradePrompt 
      feature="Integrations" 
      requiredTier="Pro"
      onClose={handleClose}
    />
  </div>
{/if}

