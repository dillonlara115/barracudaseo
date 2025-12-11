<script>
  import { onMount } from 'svelte';
  import { supabase } from '../lib/supabase.js';
  import { push } from 'svelte-spa-router';

  let loading = true;
  let submitting = false;
  let error = null;
  let success = null;
  let password = '';
  let confirmPassword = '';
  let userEmail = '';
  let needsRecoverySession = false;

  onMount(async () => {
    loading = true;
    error = null;
    success = null;
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session || !session.user) {
        needsRecoverySession = true;
        return;
      }
      userEmail = session.user.email || '';
    } catch (err) {
      error = err.message || 'Unable to load reset session.';
    } finally {
      loading = false;
    }
  });

  async function handleReset(event) {
    event.preventDefault();
    if (submitting) return;
    error = null;
    success = null;

    if (!password || password.length < 8) {
      error = 'Password must be at least 8 characters.';
      return;
    }
    if (password !== confirmPassword) {
      error = 'Passwords do not match.';
      return;
    }

    submitting = true;
    try {
      const { error: updateError } = await supabase.auth.updateUser({ password });
      if (updateError) throw updateError;
      success = 'Password updated. You can now log in with your new password.';
      setTimeout(() => push('#/'), 1500);
    } catch (err) {
      error = err.message || 'Failed to update password.';
    } finally {
      submitting = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-base-200">
  <div class="card w-full max-w-md bg-base-100 shadow-xl">
    <div class="card-body space-y-4">
      <h1 class="text-2xl font-bold">Reset your password</h1>

      {#if loading}
        <div class="flex items-center justify-center py-6">
          <span class="loading loading-spinner loading-md"></span>
        </div>
      {:else if needsRecoverySession}
        <div class="alert alert-warning">
          <span>Reset link invalid or expired. Please start the reset from the login page.</span>
        </div>
        <button class="btn btn-primary w-full" on:click={() => push('#/auth')}>
          Back to Login
        </button>
      {:else}
        {#if error}
          <div class="alert alert-error">
            <span>{error}</span>
          </div>
        {/if}
        {#if success}
          <div class="alert alert-success">
            <span>{success}</span>
          </div>
        {/if}

        {#if userEmail}
          <p class="text-sm text-base-content/70">Account: {userEmail}</p>
        {/if}

        <form class="space-y-4" on:submit|preventDefault={handleReset}>
          <div class="form-control">
            <label class="label">
              <span class="label-text">New password</span>
            </label>
            <input
              type="password"
              class="input input-bordered w-full"
              bind:value={password}
              autocomplete="new-password"
              required
              minlength="8"
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Confirm password</span>
            </label>
            <input
              type="password"
              class="input input-bordered w-full"
              bind:value={confirmPassword}
              autocomplete="new-password"
              required
              minlength="8"
            />
          </div>
          <button class="btn btn-primary w-full" type="submit" disabled={submitting}>
            {#if submitting}
              <span class="loading loading-spinner loading-sm"></span>
              <span class="ml-2">Updating...</span>
            {:else}
              Update password
            {/if}
          </button>
          <button type="button" class="btn btn-ghost w-full" on:click={() => push('#/auth')}>
            Back to Login
          </button>
        </form>
      {/if}
    </div>
  </div>
</div>

