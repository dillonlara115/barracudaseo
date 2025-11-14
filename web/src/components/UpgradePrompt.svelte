<script>
  import { push } from 'svelte-spa-router';
  import { CreditCard, X } from 'lucide-svelte';

  export let feature = 'this feature';
  export let requiredTier = 'Pro';
  export let onClose = null;
</script>

<div class="card bg-base-100 shadow-lg border-2 border-primary">
  <div class="card-body">
    <div class="flex items-start justify-between">
      <div class="flex-1">
        <h2 class="card-title text-xl mb-2">
          <CreditCard class="w-6 h-6" />
          Upgrade Required
        </h2>
        <p class="text-base-content/70 mb-4">
          {feature} is only available for {requiredTier} plan subscribers. Upgrade to unlock this feature and more.
        </p>
        <div class="flex gap-2">
          <button 
            class="btn btn-primary"
            on:click={() => push('#/billing')}
          >
            View Plans
          </button>
          {#if onClose}
            <button 
              class="btn btn-ghost"
              on:click={onClose}
            >
              Close
            </button>
          {/if}
        </div>
      </div>
      {#if onClose}
        <button 
          class="btn btn-sm btn-ghost btn-circle"
          on:click={onClose}
        >
          <X class="w-4 h-4" />
        </button>
      {/if}
    </div>
  </div>
</div>

