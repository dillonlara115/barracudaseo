# Frontend Guidelines (Svelte + DaisyUI + Tailwind)

## Component Structure

Follow this standard structure for all Svelte components:

```svelte
<script>
	// 1. Imports (lifecycle, components, utilities)
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';
	import { push } from 'svelte-spa-router';

	// 2. Props (exports)
	export let projectId = null;
	export let issues = [];

	// 3. Event dispatcher
	const dispatch = createEventDispatcher();

	// 4. Local state
	let loading = false;
	let error = null;

	// 5. Reactive statements
	$: filteredData = data.filter((d) => d.visible);

	// 6. Functions
	function handleClick() {}

	// 7. Lifecycle hooks
	onMount(async () => {});
	onDestroy(() => {});
</script>

<!-- Template -->
{#if loading}
	<LoadingState />
{:else if error}
	<ErrorState message={error} />
{:else}
	<Content />
{/if}

<style>
	/* Scoped CSS only when Tailwind classes insufficient */
</style>
```

## DaisyUI Components

### Buttons
```svelte
<button class="btn btn-primary">Primary Action</button>
<button class="btn btn-ghost">Secondary</button>
<button class="btn btn-sm btn-circle btn-ghost" aria-label="Close">✕</button>
<button class="btn btn-outline btn-error">Destructive</button>
```

### Cards
```svelte
<div class="card bg-base-200 shadow-lg">
	<div class="card-body">
		<h3 class="card-title">Title</h3>
		<p>Content</p>
		<div class="card-actions justify-end">
			<button class="btn btn-primary">Action</button>
		</div>
	</div>
</div>
```

### Modals
```svelte
{#if showModal}
	<dialog class="modal modal-open">
		<div class="modal-box bg-base-200">
			<h3 class="text-lg font-bold">Modal Title</h3>
			<p class="py-4">Content</p>
			<div class="modal-action">
				<button class="btn btn-ghost" on:click={handleCancel}>Cancel</button>
				<button class="btn btn-primary" on:click={handleSubmit}>Submit</button>
			</div>
		</div>
		<form method="dialog" class="modal-backdrop">
			<button on:click={handleCancel}>close</button>
		</form>
	</dialog>
{/if}
```

### Form Controls
```svelte
<div class="form-control w-full">
	<label class="label" for="input-id">
		<span class="label-text">Label</span>
	</label>
	<input
		id="input-id"
		type="text"
		class="input input-bordered w-full"
		bind:value={inputValue}
	/>
	<label class="label">
		<span class="label-text-alt text-base-content/70">Helper text</span>
	</label>
</div>

<select class="select select-bordered w-full" bind:value={selected}>
	<option value="all">All</option>
</select>
```

### Alerts
```svelte
<div class="alert alert-info">
	<span class="loading loading-spinner loading-sm"></span>
	<span>Loading...</span>
</div>
<div class="alert alert-success">Success</div>
<div class="alert alert-warning">Warning</div>
<div class="alert alert-error">Error</div>
```

### Badges
```svelte
<span class="badge badge-primary">Primary</span>
<span class="badge badge-error">Error</span>
<span class="badge badge-lg badge-ghost">Large Ghost</span>
```

## Tailwind Conventions

### Theme Colors (Custom Barracuda Theme)
- `bg-base-100`, `bg-base-200`, `bg-base-300` - Background layers
- `text-base-content` - Primary text
- `text-base-content/70` - Secondary text (with opacity)
- `text-error`, `text-warning`, `text-success`, `text-info` - Semantic colors
- `btn-primary` uses theme primary color (#8ec07c green)

### Layout Patterns
```svelte
<!-- Responsive grid -->
<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">

<!-- Flex container -->
<div class="flex items-center justify-between gap-4">

<!-- Stack with spacing -->
<div class="space-y-4">
```

### Responsive Design
- Mobile-first: Start with base styles, add `md:` and `lg:` breakpoints
- Common breakpoints: `md:` (768px), `lg:` (1024px)

## API Call Pattern

Always use the centralized data functions from `lib/data.js`:

```javascript
async function loadData() {
	loading = true;
	const { data, error } = await fetchProjects();
	loading = false;

	if (error) {
		errorMessage = error.message;
		return;
	}

	projects = data;
}
```

Return format is always `{ data, error }` - never throws.

## Event Communication

```javascript
// Child component emits events
const dispatch = createEventDispatcher();
dispatch('created', { id: newId });
dispatch('close');

// Parent listens
<ChildComponent on:created={(e) => handleCreated(e.detail)} on:close={handleClose} />
```

## Form Validation Pattern

```javascript
async function handleSubmit() {
	// 1. Trim and validate required fields
	if (!formData.name.trim()) {
		error = 'Name is required';
		return;
	}

	// 2. Format validation
	try {
		new URL(formData.url);
	} catch {
		error = 'Invalid URL format';
		return;
	}

	// 3. Submit
	loading = true;
	error = null;
	const result = await createItem(formData);
	loading = false;

	if (result.error) {
		error = result.error.message;
		return;
	}

	dispatch('created', result.data);
}
```

## Polling Pattern

```javascript
let pollInterval = null;

function startPolling() {
	if (pollInterval) return;
	pollInterval = setInterval(async () => {
		await refreshData();
	}, 1000);
}

onMount(() => {
	document.addEventListener('visibilitychange', handleVisibilityChange);
	startPolling();
});

onDestroy(() => {
	clearInterval(pollInterval);
	document.removeEventListener('visibilitychange', handleVisibilityChange);
});
```

## Accessibility

- Always use `for` attribute on labels matching input `id`
- Add `aria-label` to icon-only buttons
- Use semantic HTML elements (`<button>`, `<nav>`, `<main>`)

## Formatting

Run before committing:
```bash
cd web && npm run format
```

Config: tabs, single quotes, no trailing commas, 100 char line width.
