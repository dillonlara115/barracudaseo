<script>
  import { onMount } from 'svelte';
  import { user } from '../lib/auth.js';
  import { supabase } from '../lib/supabase.js';
  import { X, Loader, Mail, UserPlus, Trash2, Check, AlertCircle, Send } from 'lucide-svelte';

  let loading = true;
  let members = [];
  let teamSizeLimit = 1;
  let activeCount = 0;
  let isOwner = false;
  let error = null;
  
  let inviteEmail = '';
  let inviting = false;
  let inviteError = null;
  let showInviteModal = false;
  let resendingInviteId = null;

  const API_URL = import.meta.env.VITE_CLOUD_RUN_API_URL || 'http://localhost:8080';

  onMount(() => {
    if ($user) {
      loadTeamMembers();
    }
  });

  async function getValidAccessToken() {
    const { data: sessionData, error: sessionError } = await supabase.auth.getSession();
    if (sessionError) {
      throw new Error('Not authenticated. Please sign in again.');
    }

    let currentSession = sessionData.session;
    if (!currentSession) {
      const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
      if (refreshError || !refreshed.session) {
        throw new Error('Session expired. Please sign in again.');
      }
      currentSession = refreshed.session;
    }

    const expiresAt = currentSession?.expires_at;
    if (expiresAt && expiresAt * 1000 < Date.now() + 60000) {
      const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
      if (refreshError || !refreshed.session) {
        throw new Error('Session expired. Please sign in again.');
      }
      currentSession = refreshed.session;
    }

    const token = currentSession?.access_token;
    if (!token) {
      throw new Error('Not authenticated');
    }

    return token;
  }

  async function loadTeamMembers() {
    if (!$user) return;
    
    loading = true;
    error = null;
    
    try {
      const token = await getValidAccessToken();
      const response = await fetch(`${API_URL}/api/v1/team/members`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        const errorMessage = errorData?.error || `Failed to load team members (${response.status})`;
        throw new Error(errorMessage);
      }

      const data = await response.json();
      members = data.members || [];
      // Use team_size_limit from API response, default to 1 if undefined/null/0
      teamSizeLimit = (data.team_size_limit != null && data.team_size_limit > 0) ? data.team_size_limit : 1;
      activeCount = data.active_count || 0;
      isOwner = data.is_owner || false;
      
      // Debug logging
      console.log('Team members loaded:', { 
        teamSizeLimit, 
        activeCount, 
        membersCount: members.length,
        isOwner,
        apiData: data 
      });
    } catch (err) {
      error = err.message || 'Failed to load team members';
      console.error('Failed to load team members:', err);
    } finally {
      loading = false;
    }
  }

  async function inviteMember() {
    if (!inviteEmail.trim()) return;
    
    inviting = true;
    inviteError = null;
    
    try {
      const token = await getValidAccessToken();
      const response = await fetch(`${API_URL}/api/v1/team/invite`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ email: inviteEmail.trim() }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to send invite');
      }

      const result = await response.json();
      
      // Success! Close modal and reload
      inviteEmail = '';
      showInviteModal = false;
      
      // Show success message
      const successMsg = document.createElement('div');
      successMsg.className = 'alert alert-success fixed bottom-4 right-4 w-auto z-50';
      successMsg.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg><span>Invite sent! Share this link: ${result.invite_url || 'Check your email'}</span>`;
      document.body.appendChild(successMsg);
      setTimeout(() => successMsg.remove(), 5000);
      
      await loadTeamMembers();
    } catch (err) {
      inviteError = err.message;
    } finally {
      inviting = false;
    }
  }

  async function resendInvite(memberId) {
    resendingInviteId = memberId;
    
    try {
      const token = await getValidAccessToken();
      const response = await fetch(`${API_URL}/api/v1/team/${memberId}/resend`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to resend invite');
      }

      const result = await response.json();
      
      // Show success message
      const successMsg = document.createElement('div');
      successMsg.className = 'alert alert-success fixed bottom-4 right-4 w-auto z-50';
      successMsg.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg><span>Invitation resent successfully!</span>`;
      document.body.appendChild(successMsg);
      setTimeout(() => successMsg.remove(), 3000);

      await loadTeamMembers();
    } catch (err) {
      error = err.message || 'Failed to resend invite';
      console.error('Failed to resend invite:', err);
    } finally {
      resendingInviteId = null;
    }
  }

  async function removeMember(memberId) {
    if (!confirm('Are you sure you want to remove this team member?')) {
      return;
    }
    
    try {
      const token = await getValidAccessToken();
      const response = await fetch(`${API_URL}/api/v1/team/${memberId}/remove`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to remove member');
      }

      await loadTeamMembers();
    } catch (err) {
      error = err.message || 'Failed to remove member';
      console.error('Failed to remove member:', err);
    }
  }

  function getStatusBadge(status) {
    switch (status) {
      case 'active':
        return 'badge-success';
      case 'pending':
        return 'badge-warning';
      case 'removed':
        return 'badge-error';
      default:
        return 'badge-ghost';
    }
  }

  function getRoleBadge(role) {
    return role === 'admin' ? 'badge-primary' : 'badge-ghost';
  }
</script>

<div class="card bg-base-100 shadow-lg">
  <div class="card-body">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="card-title text-xl">Team Members</h2>
        <p class="text-sm text-base-content/70 mt-1">
          {activeCount} of {teamSizeLimit} seats used
        </p>
      </div>
      
      {#if isOwner && activeCount < teamSizeLimit}
        <button 
          class="btn btn-primary btn-sm"
          on:click={() => showInviteModal = true}
        >
          <UserPlus class="w-4 h-4" />
          Invite Member
        </button>
      {/if}
    </div>

    {#if loading}
      <div class="flex items-center justify-center py-8">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    {:else if error}
      <div class="alert alert-error">
        <AlertCircle class="w-5 h-5" />
        <span>{error}</span>
      </div>
    {:else}
      {#if members.length === 0}
        <div class="text-center py-8 text-base-content/70">
          <p>No team members yet.</p>
          {#if isOwner}
            <button 
              class="btn btn-primary btn-sm mt-4"
              on:click={() => showInviteModal = true}
            >
              <UserPlus class="w-4 h-4" />
              Invite Your First Member
            </button>
          {/if}
        </div>
      {:else}
        <div class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
              <tr>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Joined</th>
                {#if isOwner}
                  <th>Actions</th>
                {/if}
              </tr>
            </thead>
            <tbody>
              {#each members as member (member.id)}
                <tr>
                  <td>
                    <div class="flex items-center gap-2">
                      <Mail class="w-4 h-4 text-base-content/50" />
                      <span>{member.email}</span>
                    </div>
                  </td>
                  <td>
                    <span class="badge {getRoleBadge(member.role)} badge-sm">
                      {member.role || 'member'}
                    </span>
                  </td>
                  <td>
                    <span class="badge {getStatusBadge(member.status)} badge-sm">
                      {member.status}
                    </span>
                  </td>
                  <td>
                    {#if member.joined_at}
                      {new Date(member.joined_at).toLocaleDateString()}
                    {:else if member.invited_at}
                      <span class="text-base-content/50">Invited {new Date(member.invited_at).toLocaleDateString()}</span>
                    {:else}
                      <span class="text-base-content/50">—</span>
                    {/if}
                  </td>
                  {#if isOwner}
                    <td>
                      <div class="flex gap-2">
                        {#if member.status === 'pending'}
                          {#if member.user_id !== $user?.id}
                            <button 
                              class="btn btn-xs btn-primary btn-outline"
                              on:click={() => resendInvite(member.id)}
                              disabled={resendingInviteId === member.id}
                              title="Resend invitation email"
                            >
                              {#if resendingInviteId === member.id}
                                <Loader class="w-3 h-3 animate-spin" />
                              {:else}
                                <Send class="w-3 h-3" />
                              {/if}
                              Resend
                            </button>
                            <button 
                              class="btn btn-xs btn-error"
                              on:click={() => removeMember(member.id)}
                              title="Remove team member"
                            >
                              <Trash2 class="w-3 h-3" />
                              Remove
                            </button>
                          {:else}
                            <span class="text-xs text-base-content/50">You</span>
                          {/if}
                        {:else if member.status === 'active'}
                          {#if member.user_id !== $user?.id}
                            <button 
                              class="btn btn-xs btn-error"
                              on:click={() => removeMember(member.id)}
                              title="Remove team member"
                            >
                              <Trash2 class="w-3 h-3" />
                              Remove
                            </button>
                          {:else}
                            <span class="text-xs text-base-content/50">You</span>
                          {/if}
                        {/if}
                      </div>
                    </td>
                  {/if}
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    {/if}
  </div>
</div>

<!-- Invite Modal -->
{#if showInviteModal}
  <input type="checkbox" id="invite-modal" class="modal-toggle" checked />
  <div class="modal" role="dialog">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4">Invite Team Member</h3>
      
      <div class="form-control w-full">
        <label class="label">
          <span class="label-text">Email Address</span>
        </label>
        <input 
          type="email" 
          placeholder="colleague@example.com" 
          class="input input-bordered w-full"
          bind:value={inviteEmail}
          disabled={inviting}
        />
        <label class="label">
          <span class="label-text-alt">They'll receive an invite link to join your team.</span>
        </label>
      </div>

      {#if inviteError}
        <div class="alert alert-error mt-4">
          <AlertCircle class="w-5 h-5" />
          <span>{inviteError}</span>
        </div>
      {/if}

      {#if activeCount >= teamSizeLimit}
        <div class="alert alert-warning mt-4">
          <AlertCircle class="w-5 h-5" />
          <span>Team size limit reached. Upgrade your plan to add more members.</span>
        </div>
      {/if}

      <div class="modal-action">
        <label 
          for="invite-modal" 
          class="btn"
          on:click={() => { showInviteModal = false; inviteEmail = ''; inviteError = null; }}
        >
          Cancel
        </label>
        <button 
          class="btn btn-primary" 
          on:click={inviteMember}
          disabled={inviting || !inviteEmail.trim() || activeCount >= teamSizeLimit}
        >
          {#if inviting}
            <Loader class="w-4 h-4 animate-spin" />
          {:else}
            <Mail class="w-4 h-4" />
          {/if}
          Send Invite
        </button>
      </div>
    </div>
  </div>
{/if}

