<script lang="ts">
	import { enhance } from '$app/forms';
	import { SvelteMap } from 'svelte/reactivity';
	import Autocomplete from '$lib/Autocomplete.svelte';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();

	type Tab = 'mutuals' | 'notFollowingMeBack' | 'iDontFollowBack';
	const tabs: { id: Tab; label: string }[] = [
		{ id: 'mutuals', label: 'mutuals' },
		{ id: 'notFollowingMeBack', label: "don't follow you back" },
		{ id: 'iDontFollowBack', label: "you don't follow" }
	];

	let activeTab: Tab = $state('mutuals');
	let filter = $state('');
	let locFilter = $state('');
	let minFollowers = $state('');
	let minFollowing = $state('');
	let minRepos = $state('');
	let submitting = $state<string | null>(null);
	let refreshing = $state(false);
	let scrolled = $state(false);

	type Detail = {
		name: string | null;
		bio: string | null;
		followers: number;
		following: number;
		publicRepos: number;
		publicGists: number;
		location: string | null;
		company: string | null;
		blog: string | null;
		twitterUsername: string | null;
		email: string | null;
		hireable: boolean;
		createdAt: string | null;
	};
	const details = new SvelteMap<string, Detail | 'loading' | 'error'>();

	// count/location filters need per-user detail — load it for the whole tab when one is active
	const detailFilterOn = $derived(
		locFilter.trim() !== '' || minFollowers !== '' || minFollowing !== '' || minRepos !== ''
	);

	$effect(() => {
		if (detailFilterOn) for (const f of data[activeTab]) loadDetail(f.login);
	});

	function currentList() {
		const q = filter.toLowerCase();
		const loc = locFilter.trim().toLowerCase();
		return data[activeTab].filter((f) => {
			if (q && !f.login.toLowerCase().includes(q)) return false;
			if (!detailFilterOn) return true;
			const d = details.get(f.login);
			if (!d || d === 'loading' || d === 'error') return false;
			if (minFollowers !== '' && d.followers < +minFollowers) return false;
			if (minFollowing !== '' && d.following < +minFollowing) return false;
			if (minRepos !== '' && d.publicRepos < +minRepos) return false;
			if (loc) {
				// "Indore, India" matches either "indore" or "india"
				const tokens = (d.location ?? '').toLowerCase().split(',').map((t) => t.trim());
				if (!tokens.some((t) => t.includes(loc))) return false;
			}
			return true;
		});
	}

	function relativeTime(ts: number) {
		const s = Math.floor((Date.now() - ts) / 1000);
		if (s < 60) return 'just now';
		if (s < 3600) return `${Math.floor(s / 60)}m ago`;
		if (s < 86400) return `${Math.floor(s / 3600)}h ago`;
		return `${Math.floor(s / 86400)}d ago`;
	}

	// user detail rarely changes — cache in localStorage for a day
	const DETAIL_TTL = 24 * 60 * 60 * 1000;

	const CACHE_KEY = 'gfstat:u2:'; // bump suffix when Detail shape changes

	function cacheGet(login: string): Detail | null {
		try {
			const raw = localStorage.getItem(CACHE_KEY + login);
			if (!raw) return null;
			const { d, t } = JSON.parse(raw);
			return Date.now() - t < DETAIL_TTL ? d : null;
		} catch {
			return null;
		}
	}

	function cachePut(login: string, d: Detail) {
		try {
			localStorage.setItem(CACHE_KEY + login, JSON.stringify({ d, t: Date.now() }));
		} catch {
			/* storage full/blocked — cache is best-effort */
		}
	}

	async function loadDetail(login: string) {
		if (details.has(login)) return;
		const cached = cacheGet(login);
		if (cached) {
			details.set(login, cached);
			return;
		}
		details.set(login, 'loading');
		try {
			// authenticated via our server proxy — anon quota (60/hr/IP) can't cover a full grid
			const res = await fetch(`/api/user/${login}`);
			if (!res.ok) throw new Error('fetch failed');
			const u = await res.json();
			const d: Detail = {
				name: u.name ?? null,
				bio: u.bio ?? null,
				followers: u.followers ?? 0,
				following: u.following ?? 0,
				publicRepos: u.publicRepos ?? u.public_repos ?? 0,
				publicGists: u.publicGists ?? u.public_gists ?? 0,
				location: u.location ?? null,
				company: u.company ?? null,
				blog: u.blog || null,
				twitterUsername: u.twitterUsername ?? u.twitter_username ?? null,
				email: u.email ?? null,
				hireable: !!u.hireable,
				createdAt: u.createdAt ?? u.created_at ?? null
			};
			details.set(login, d);
			cachePut(login, d);
		} catch {
			details.set(login, 'error');
		}
	}

	// lazy load a card's detail when it scrolls near the viewport
	function lazyDetail(node: HTMLElement, login: string) {
		const io = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting) {
					io.disconnect();
					loadDetail(login);
				}
			},
			{ rootMargin: '200px' }
		);
		io.observe(node);
		return { destroy: () => io.disconnect() };
	}

	function joinedYear(iso: string | null) {
		return iso ? new Date(iso).getFullYear() : null;
	}

	// unique location tokens from loaded cards, for the location filter's suggestions
	const locOptions = $derived.by(() => {
		const set = new Set<string>();
		for (const d of details.values()) {
			if (d && d !== 'loading' && d !== 'error' && d.location) {
				for (const t of d.location.split(',')) {
					const s = t.trim();
					if (s) set.add(s);
				}
			}
		}
		return [...set].sort();
	});
</script>

<svelte:head>
	<title>gfstat — {data.profile.login}</title>
</svelte:head>

<svelte:window onscroll={() => (scrolled = window.scrollY > 600)} />

{#if scrolled}
	<button class="to-top" onclick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}>
		↑ top
	</button>
{/if}

<div class="page">
	<section class="card profile">
		<div class="profile-head">
			<img class="profile-avatar" src={data.profile.avatarUrl} alt={data.profile.login} />
			<div class="profile-id">
				<h6>signed in as</h6>
				<h3>{data.profile.name ?? data.profile.login}</h3>
				<a class="profile-handle" href={data.profile.htmlUrl} target="_blank" rel="noopener">@{data.profile.login}</a>
			</div>
		</div>

		{#if data.profile.bio}
			<p class="profile-bio">{data.profile.bio}</p>
		{/if}

		<div class="profile-stats">
			<a class="stat" href="{data.profile.htmlUrl}?tab=followers" target="_blank" rel="noopener">
				<span class="stat-num">{data.profile.followers}</span>
				<span class="stat-label">followers</span>
			</a>
			<a class="stat" href="{data.profile.htmlUrl}?tab=following" target="_blank" rel="noopener">
				<span class="stat-num">{data.profile.following}</span>
				<span class="stat-label">following</span>
			</a>
			<a class="stat" href="{data.profile.htmlUrl}?tab=repositories" target="_blank" rel="noopener">
				<span class="stat-num">{data.profile.publicRepos}</span>
				<span class="stat-label">repos</span>
			</a>
		</div>

		{#if data.profile.location || data.profile.blog || data.profile.company}
			<p class="meta profile-meta">
				{#if data.profile.location}<span class="pm">◆ {data.profile.location}</span>{/if}
				{#if data.profile.company}<span class="pm">▲ {data.profile.company}</span>{/if}
				{#if data.profile.blog}
					<a class="pm" href={data.profile.blog.startsWith('http') ? data.profile.blog : `https://${data.profile.blog}`}
						target="_blank" rel="noopener">→ {data.profile.blog}</a>
				{/if}
			</p>
		{/if}
	</section>

	<div class="toolbar">
		<span class="meta">data as of {relativeTime(data.fetchedAt)}</span>
		<form
			method="POST"
			action="?/refresh"
			use:enhance={() => {
				refreshing = true;
				return async ({ update }) => {
					await update();
					refreshing = false;
				};
			}}
		>
			<button type="submit" class="btn btn-ghost btn-sm" disabled={refreshing}>
				<span class="spin" class:spinning={refreshing}>↻</span>
				{refreshing ? 'refreshing…' : 'refresh'}
			</button>
		</form>
	</div>

	<div class="tabs" role="tablist">
		{#each tabs as tab (tab.id)}
			<button
				class="tab"
				class:active={activeTab === tab.id}
				role="tab"
				aria-selected={activeTab === tab.id}
				onclick={() => (activeTab = tab.id)}
			>
				{tab.label} <span class="tab-count">[{data[tab.id].length}]</span>
			</button>
		{/each}
	</div>

	<div class="filters">
		<Autocomplete bind:value={filter} options={data[activeTab].map((f) => f.login)} placeholder="filter by login…" />
		<Autocomplete bind:value={locFilter} options={locOptions} placeholder="location…" />
		<input class="filter num" type="number" min="0" placeholder="min followers" bind:value={minFollowers} />
		<input class="filter num" type="number" min="0" placeholder="min following" bind:value={minFollowing} />
		<input class="filter num" type="number" min="0" placeholder="min repos" bind:value={minRepos} />
	</div>

	{#if form?.error}
		<div class="callout callout-caution">
			<p class="callout-title">error</p>
			<p>{form.error}</p>
		</div>
	{/if}

	{#key activeTab}
	<div class="grid">
		{#each currentList() as follow (follow.login)}
			{@const d = details.get(follow.login)}
			<article class="mini" use:lazyDetail={follow.login}>
				<div class="mini-top">
					<a href={follow.htmlUrl} target="_blank" rel="noopener">
						<img class="mini-avatar" src={follow.avatarUrl} alt={follow.login} loading="lazy" />
					</a>
					<a class="mini-name" href={follow.htmlUrl} target="_blank" rel="noopener">
						{d && d !== 'loading' && d !== 'error' && d.name ? d.name : follow.login}
					</a>
					<a class="mini-handle" href={follow.htmlUrl} target="_blank" rel="noopener">@{follow.login}</a>
					{#if activeTab === 'iDontFollowBack'}
						<form
							method="POST"
							action="?/follow"
							use:enhance={() => {
								submitting = follow.login;
								return async ({ update }) => {
									submitting = null;
									await update();
								};
							}}
						>
							<input type="hidden" name="login" value={follow.login} />
							<button type="submit" class="btn btn-ghost btn-sm" disabled={submitting === follow.login}>
								{submitting === follow.login ? '…' : 'follow'}
							</button>
						</form>
					{:else}
						<form
							method="POST"
							action="?/unfollow"
							use:enhance={() => {
								submitting = follow.login;
								return async ({ update }) => {
									submitting = null;
									await update();
								};
							}}
						>
							<input type="hidden" name="login" value={follow.login} />
							<button type="submit" class="btn btn-ghost btn-sm" disabled={submitting === follow.login}>
								{submitting === follow.login ? '…' : 'unfollow ✕'}
							</button>
						</form>
					{/if}
				</div>

				<div class="mini-body">
					{#if d === 'loading' || d === undefined}
						<span class="mini-dim">loading…</span>
					{:else if d === 'error'}
						<span class="mini-dim">could not load profile</span>
					{:else}
						<div class="mini-loaded">
							<p class="mini-bio">
								{#if d.bio}{d.bio}{:else}<span class="nil">nil</span>{/if}
							</p>
							<p class="mini-line">
								<a href="{follow.htmlUrl}?tab=followers" target="_blank" rel="noopener"><b>{d.followers}</b> followers</a>
								·
								<a href="{follow.htmlUrl}?tab=following" target="_blank" rel="noopener"><b>{d.following}</b> following</a>
							</p>
							<p class="mini-line">
								<a href="{follow.htmlUrl}?tab=repositories" target="_blank" rel="noopener"><b>{d.publicRepos ?? 0}</b> repos</a>
								·
								<a href="https://gist.github.com/{follow.login}" target="_blank" rel="noopener"><b>{d.publicGists ?? 0}</b> gists</a>
							</p>
							<p class="mini-line">◆ {#if d.location}{d.location}{:else}<span class="nil">nil</span>{/if}</p>
							<p class="mini-line">▲ {#if d.company}{d.company}{:else}<span class="nil">nil</span>{/if}</p>
							<p class="mini-line">⧗ joined {#if d.createdAt}{joinedYear(d.createdAt)}{:else}<span class="nil">nil</span>{/if}</p>
							<p class="mini-line">
								→
								{#if d.blog}
									<a href={d.blog.startsWith('http') ? d.blog : `https://${d.blog}`}
										target="_blank" rel="noopener">{d.blog}</a>
								{:else}<span class="nil">nil</span>{/if}
							</p>
							<p class="mini-line">
								𝕏
								{#if d.twitterUsername}
									<a href="https://x.com/{d.twitterUsername}" target="_blank" rel="noopener">@{d.twitterUsername}</a>
								{:else}<span class="nil">nil</span>{/if}
							</p>
							<p class="mini-line">
								✉
								{#if d.email}
									<a href="mailto:{d.email}">{d.email}</a>
								{:else}<span class="nil">nil</span>{/if}
							</p>
							{#if d.hireable}
								<p class="mini-line"><span class="badge badge-success">hireable</span></p>
							{/if}
						</div>
					{/if}
				</div>
			</article>
		{:else}
			<p class="empty meta">{filter ? '[ no matches ]' : '[ empty ]'}</p>
		{/each}
	</div>
	{/key}
</div>

<style>
	.page {
		max-width: var(--ds-max-width);
		margin: var(--ds-space-8) auto 0;
		padding: 0 var(--ds-space-4);
	}

	/* ---- profile card ---- */
	.profile {
		border-style: dashed;
		max-width: var(--ds-content-width);
		margin-inline: auto;
	}

	.profile-head {
		display: flex;
		gap: var(--ds-space-4);
		align-items: center;
	}

	.profile-avatar {
		width: 72px;
		height: 72px;
		border-radius: var(--ds-radius);
		border: 1px solid var(--ds-primary);
		flex-shrink: 0;
	}

	.profile-id h6 {
		color: var(--ds-primary);
		margin: 0 0 var(--ds-space-1);
	}

	.profile-id h3 {
		margin: 0;
	}

	.profile-handle {
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-tertiary);
		text-decoration: none;
		transition: color var(--ds-transition);
	}

	.profile-handle:hover {
		color: var(--ds-info);
	}

	.profile-bio {
		margin: var(--ds-space-4) 0 0;
		color: var(--ds-text-secondary);
	}

	.profile-stats {
		display: flex;
		gap: var(--ds-space-6);
		margin: var(--ds-space-4) 0 0;
		padding-top: var(--ds-space-4);
		border-top: 1px dashed var(--ds-border);
	}

	.stat {
		display: flex;
		flex-direction: column;
		text-decoration: none;
	}

	.stat-num {
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-xl);
		font-weight: var(--ds-weight-semibold);
		color: var(--ds-text-primary);
		transition: color var(--ds-transition);
	}

	.stat:hover .stat-num {
		color: var(--ds-primary);
	}

	.stat-label {
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-xs);
		color: var(--ds-text-tertiary);
	}

	.profile-meta {
		margin: var(--ds-space-4) 0 0;
		gap: var(--ds-space-4);
		color: var(--ds-text-tertiary);
	}

	.profile-meta .pm {
		color: var(--ds-text-secondary);
	}

	a.pm {
		text-decoration: none;
		transition: color var(--ds-transition);
	}

	a.pm:hover {
		color: var(--ds-primary);
	}

	/* ---- toolbar ---- */
	.toolbar {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin: var(--ds-space-6) 0 var(--ds-space-4);
	}

	.toolbar form {
		margin: 0;
	}

	.spin {
		display: inline-block;
	}

	.spin.spinning {
		animation: spin 800ms linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.spin.spinning {
			animation: none;
		}
	}

	/* ---- tabs ---- */
	.tabs {
		display: flex;
		flex-wrap: wrap;
		border-bottom: 1px dashed var(--ds-border);
	}

	.tab {
		flex: 1 1 0;
		text-align: center;
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-secondary);
		background: transparent;
		border: none;
		border-bottom: 1px solid transparent;
		padding: var(--ds-space-2) var(--ds-space-3);
		margin-bottom: -1px;
		cursor: pointer;
		transition: color var(--ds-transition);
	}

	.tab:hover {
		color: var(--ds-text-primary);
	}

	.tab.active {
		color: var(--ds-primary);
		border-bottom-color: var(--ds-primary);
	}

	.tab-count {
		color: var(--ds-text-tertiary);
	}

	.tab.active .tab-count {
		color: var(--ds-primary);
	}

	/* ---- filter inputs ---- */
	.filters {
		display: flex;
		flex-wrap: wrap;
		gap: var(--ds-space-2);
		margin: var(--ds-space-4) 0;
	}

	/* each filter grows equally to fill the row */
	.filters > :global(*) {
		flex: 1 1 0;
		min-width: 140px;
	}

	.filter.num {
		-moz-appearance: textfield;
		appearance: textfield;
	}

	.filter.num::-webkit-outer-spin-button,
	.filter.num::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	.filter {
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-primary);
		background: var(--ds-bg-code);
		border: 1px dashed var(--ds-border);
		border-radius: var(--ds-radius-sm);
		padding: var(--ds-space-2) var(--ds-space-3);
		width: 100%;
		outline: none;
		transition: border-color var(--ds-transition);
	}

	.filter::placeholder {
		color: var(--ds-text-tertiary);
	}

	.filter:focus {
		border: 1px solid var(--ds-primary);
	}

	/* ---- card grid ---- */
	/* flex, not grid: a partial last row (N % 4 != 0) centers itself */
	.grid {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: var(--ds-space-3);
		animation: fade-up 200ms ease both;
	}

	.grid > .mini {
		flex: 0 0 calc(25% - var(--ds-space-3) * 3 / 4);
	}

	@media (max-width: 900px) {
		.grid > .mini {
			flex-basis: calc(50% - var(--ds-space-3) / 2);
		}
	}

	@media (max-width: 520px) {
		.grid > .mini {
			flex-basis: 100%;
		}
	}

	.mini {
		display: flex;
		flex-direction: column;
		border: 1px dashed var(--ds-border);
		border-radius: var(--ds-radius);
		padding: var(--ds-space-4) var(--ds-space-3) var(--ds-space-3);
		background: transparent;
		transition:
			background var(--ds-transition),
			border-color var(--ds-transition),
			transform var(--ds-transition);
	}

	.mini:hover {
		background: var(--ds-bg-secondary);
		border-color: var(--ds-primary);
		transform: translateY(-2px);
	}

	.mini-loaded {
		animation: fade-up 250ms ease both;
	}

	@keyframes fade-up {
		from {
			opacity: 0;
			transform: translateY(4px);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.grid,
		.mini-loaded {
			animation: none;
		}
		.mini:hover {
			transform: none;
		}
	}

	/* ---- vertical ID-card header ---- */
	.mini-top {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--ds-space-1);
	}

	.mini-top form {
		margin: var(--ds-space-2) 0 0;
	}

	.mini-avatar {
		width: 64px;
		height: 64px;
		border-radius: var(--ds-radius);
		border: 1px solid var(--ds-border);
		margin-bottom: var(--ds-space-1);
	}

	.mini:hover .mini-avatar {
		border-color: var(--ds-primary);
	}

	.mini-name {
		font-weight: var(--ds-weight-medium);
		color: var(--ds-text-primary);
		text-decoration: none;
		overflow-wrap: anywhere; /* full name, wraps instead of truncating */
		transition: color var(--ds-transition);
	}

	.mini-name:hover {
		color: var(--ds-primary);
	}

	.mini-line a:hover b {
		color: var(--ds-primary);
	}

	.mini-handle {
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-secondary);
		text-decoration: none;
		overflow-wrap: anywhere;
		transition: color var(--ds-transition);
	}

	.mini-handle:hover {
		color: var(--ds-primary);
	}

	/* ---- body: one line per fact, nothing truncated ---- */
	.mini-body {
		margin-top: var(--ds-space-3);
		padding-top: var(--ds-space-3);
		border-top: 1px dashed var(--ds-border);
		flex: 1;
	}

	.mini-dim {
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-secondary);
	}

	.mini-bio {
		margin: 0 0 var(--ds-space-2);
		font-size: var(--ds-text-sm);
		line-height: 1.5;
		min-height: 3em; /* reserve two lines even when bio is nil */
		color: var(--ds-text-secondary);
		overflow-wrap: anywhere;
	}

	/* placeholder for fields the user hasn't filled in — dimmer than data, still legible */
	.nil {
		color: var(--ds-text-tertiary);
		font-style: italic;
	}

	.mini-line {
		margin: var(--ds-space-2) 0 0;
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		line-height: 1.6;
		color: var(--ds-text-secondary);
		overflow-wrap: anywhere;
	}

	.mini-line b {
		font-weight: var(--ds-weight-semibold);
		color: var(--ds-text-primary);
	}

	.mini-line a {
		color: var(--ds-text-primary);
		text-decoration: none;
		transition: color var(--ds-transition);
	}

	.mini-line a:hover {
		color: var(--ds-primary);
	}

	.empty {
		width: 100%;
		justify-content: center;
		padding: var(--ds-space-8) 0;
	}

	/* ---- scroll to top ---- */
	.to-top {
		position: fixed;
		right: var(--ds-space-6);
		bottom: var(--ds-space-6);
		z-index: 100;
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-primary);
		background: var(--ds-bg-secondary);
		border: 1px dashed var(--ds-border);
		border-radius: var(--ds-radius);
		padding: var(--ds-space-2) var(--ds-space-3);
		cursor: pointer;
		transition:
			color var(--ds-transition),
			border-color var(--ds-transition);
		animation: fade-up 200ms ease both;
	}

	.to-top:hover {
		color: var(--ds-primary);
		border-color: var(--ds-primary);
	}

	@media (prefers-reduced-motion: reduce) {
		.to-top {
			animation: none;
		}
	}

	.callout {
		margin: 0 0 var(--ds-space-4);
	}
</style>
