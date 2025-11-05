<script>
  import { signIn, signUp, signOut } from '../lib/auth.js';
  import { user } from '../lib/auth.js';

  let isSignUp = false;
  let email = '';
  let password = '';
  let displayName = '';
  let loading = false;
  let error = null;
  let success = null;

  $: isAuthenticated = $user !== null;

  async function handleSubmit() {
    loading = true;
    error = null;
    success = null;

    try {
      if (isSignUp) {
        const { data, error: signUpError } = await signUp(email, password, displayName);
        if (signUpError) throw signUpError;
        success = 'Account created! Please check your email to verify your account.';
      } else {
        const { data, error: signInError } = await signIn(email, password);
        if (signInError) throw signInError;
        success = 'Signed in successfully!';
      }
    } catch (err) {
      error = err.message || 'An error occurred';
    } finally {
      loading = false;
    }
  }

  async function handleSignOut() {
    loading = true;
    error = null;
    try {
      const { error: signOutError } = await signOut();
      if (signOutError) throw signOutError;
    } catch (err) {
      error = err.message || 'An error occurred';
    } finally {
      loading = false;
    }
  }
</script>

{#if isAuthenticated}
  <div class="dropdown dropdown-end">
    <label tabindex="0" class="btn btn-ghost">
      <div class="avatar placeholder">
        <div class="bg-neutral text-neutral-content rounded-full w-8">
          <span class="text-xs">{$user?.email?.charAt(0).toUpperCase() || 'U'}</span>
        </div>
      </div>
      <span class="ml-2">{$user?.email || 'User'}</span>
    </label>
    <ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
      <li>
        <a on:click={handleSignOut} class="text-error">
          Sign Out
        </a>
      </li>
    </ul>
  </div>
{:else}
  <div class="card bg-base-100 shadow-xl w-full max-w-md mx-auto">
    <div class="card-body">
      <h2 class="card-title justify-center">
        {isSignUp ? 'Sign Up' : 'Sign In'}
      </h2>

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

      <form on:submit|preventDefault={handleSubmit}>
        {#if isSignUp}
          <div class="form-control w-full">
            <label class="label">
              <span class="label-text">Display Name</span>
            </label>
            <input
              type="text"
              placeholder="Your name"
              class="input input-bordered w-full"
              bind:value={displayName}
            />
          </div>
        {/if}

        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">Email</span>
          </label>
          <input
            type="email"
            placeholder="email@example.com"
            class="input input-bordered w-full"
            bind:value={email}
            required
          />
        </div>

        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">Password</span>
          </label>
          <input
            type="password"
            placeholder="password"
            class="input input-bordered w-full"
            bind:value={password}
            required
            minlength="6"
          />
        </div>

        <div class="form-control mt-6">
          <button
            type="submit"
            class="btn btn-primary"
            disabled={loading}
          >
            {#if loading}
              <span class="loading loading-spinner loading-sm"></span>
            {:else}
              {isSignUp ? 'Sign Up' : 'Sign In'}
            {/if}
          </button>
        </div>
      </form>

      <div class="divider">OR</div>

      <button
        class="btn btn-ghost"
        on:click={() => isSignUp = !isSignUp}
      >
        {isSignUp ? 'Already have an account? Sign in' : "Don't have an account? Sign up"}
      </button>
    </div>
  </div>
{/if}

