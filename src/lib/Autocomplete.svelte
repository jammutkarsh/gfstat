<script lang="ts">
	let {
		value = $bindable(''),
		options,
		placeholder
	}: { value: string; options: string[]; placeholder: string } = $props();

	let open = $state(false);
	let hi = $state(-1); // keyboard-highlighted index

	const matches = $derived.by(() => {
		const q = value.trim().toLowerCase();
		const list = q ? options.filter((o) => o.toLowerCase().includes(q)) : options;
		return list.slice(0, 50); // ponytail: cap render at 50; list is scrollable anyway
	});

	function pick(o: string) {
		value = o;
		open = false;
	}

	function onkeydown(e: KeyboardEvent) {
		if (!open) return;
		if (e.key === 'ArrowDown') {
			hi = Math.min(hi + 1, matches.length - 1);
			e.preventDefault();
		} else if (e.key === 'ArrowUp') {
			hi = Math.max(hi - 1, 0);
			e.preventDefault();
		} else if (e.key === 'Enter' && hi >= 0) {
			pick(matches[hi]);
			e.preventDefault();
		} else if (e.key === 'Escape') {
			open = false;
		}
	}
</script>

<div class="ac">
	<input
		type="text"
		{placeholder}
		bind:value
		onfocus={() => (open = true)}
		oninput={() => {
			open = true;
			hi = -1;
		}}
		onblur={() => (open = false)}
		{onkeydown}
	/>
	{#if open && matches.length}
		<ul class="ac-list" role="listbox">
			{#each matches as m, i (m)}
				<li>
					<!-- mousedown fires before the input's blur, so the click always lands -->
					<button
						type="button"
						class:hi={i === hi}
						role="option"
						aria-selected={i === hi}
						onmousedown={(e) => {
							e.preventDefault();
							pick(m);
						}}>{m}</button
					>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.ac {
		position: relative;
		width: 100%;
	}

	input {
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

	input::placeholder {
		color: var(--ds-text-tertiary);
	}

	input:focus {
		border: 1px solid var(--ds-primary);
	}

	.ac-list {
		position: absolute;
		top: calc(100% + 2px);
		left: 0;
		right: 0;
		z-index: 200;
		margin: 0;
		padding: 0;
		list-style: none;
		background: var(--ds-bg-elevated);
		border: 1px dashed var(--ds-border);
		border-radius: var(--ds-radius-sm);
		max-height: 220px;
		overflow-y: auto;
	}

	.ac-list button {
		display: block;
		width: 100%;
		text-align: left;
		font-family: var(--ds-font-mono);
		font-size: var(--ds-text-sm);
		color: var(--ds-text-secondary);
		background: transparent;
		border: none;
		padding: var(--ds-space-2) var(--ds-space-3);
		cursor: pointer;
		transition:
			background var(--ds-transition),
			color var(--ds-transition);
	}

	.ac-list button:hover,
	.ac-list button.hi {
		background: var(--ds-primary-muted);
		color: var(--ds-primary);
	}
</style>
