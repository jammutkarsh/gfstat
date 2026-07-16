<script lang="ts">
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	let stars = $state<number | null>(null);
	let menuOpen = $state(false);

	// repo star count — anonymous, cached 1h; a nice-to-have, never blocks render
	$effect(() => {
		const KEY = 'gfstat:stars';
		try {
			const c = JSON.parse(localStorage.getItem(KEY) ?? 'null');
			if (c && Date.now() - c.t < 3600_000) {
				stars = c.n;
				return;
			}
		} catch {
			/* ignore */
		}
		fetch('https://api.github.com/repos/JammUtkarsh/gfstat')
			.then((r) => (r.ok ? r.json() : null))
			.then((j) => {
				if (j) {
					stars = j.stargazers_count;
					try {
						localStorage.setItem(KEY, JSON.stringify({ n: stars, t: Date.now() }));
					} catch {
						/* ignore */
					}
				}
			})
			.catch(() => {});
	});
</script>

<script lang="ts" module>
	// number formatter shared across renders
	const fmt = new Intl.NumberFormat();
</script>

<nav class="nav" class:open={menuOpen}>
	<div class="nav-inner">
		<a href="/" class="nav-brand"><span class="nav-prefix">~/</span>jammutkarsh/gfstat</a>
		<button
			class="hamburger"
			aria-label="toggle menu"
			aria-expanded={menuOpen}
			onclick={() => (menuOpen = !menuOpen)}
		>
			<span></span><span></span><span></span>
		</button>
		<div class="nav-links" class:open={menuOpen}>
			<a
				href="https://github.com/JammUtkarsh/gfstat"
				class="nav-link"
				target="_blank"
				rel="noopener"
				onclick={() => (menuOpen = false)}
			>
				★ github{#if stars !== null}<span class="star-count">{fmt.format(stars)}</span>{/if}
			</a>
			{#if data.loggedIn}
				<span class="nav-sep">|</span>
				<a
					href="/logout"
					class="nav-link"
					data-sveltekit-preload-data="off"
					data-sveltekit-reload
					onclick={() => (menuOpen = false)}>logout</a
				>
			{/if}
		</div>
	</div>
</nav>

<main>
	{@render children()}
</main>

<footer>
	<div class="footer-inner">
		<span class="meta">built with <a href="https://svelte.dev">svelte</a> + <a href="https://github.com/jammutkarsh/design-system">utc-ds</a></span>
	</div>
</footer>

<style>
	main {
		min-height: calc(100vh - 3.5rem - 6rem);
	}

	footer {
		border-top: 1px dashed var(--ds-border);
		margin-top: var(--ds-space-16);
	}

	.footer-inner {
		max-width: 92vw;
		margin-inline: auto;
		padding: var(--ds-space-6) var(--ds-space-4);
		display: flex;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: var(--ds-space-4);
	}

	.star-count {
		margin-left: var(--ds-space-1);
		font-family: var(--ds-font-mono);
		color: var(--ds-text-tertiary);
	}
</style>
