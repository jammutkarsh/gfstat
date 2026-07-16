<script lang="ts">
	const SITE = 'https://gfstat.utkarshchourasia.in';
	const title = 'gfstat — who unfollowed you, in cold hard JSON';
	const description =
		"Clout-chasing, but for devs. See your GitHub mutuals, who's ghosting you, and whom you don't follow back — then fix the ratio one petty click at a time.";
	const ogImage = `${SITE}/meme.jpg`;

	const jsonLd = {
		'@context': 'https://schema.org',
		'@type': 'WebApplication',
		name: 'gfstat',
		url: SITE,
		description,
		applicationCategory: 'DeveloperApplication',
		operatingSystem: 'Web',
		offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' }
	};
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<link rel="canonical" href={SITE} />

	<meta property="og:type" content="website" />
	<meta property="og:site_name" content="gfstat" />
	<meta property="og:locale" content="en_US" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content={ogImage} />
	<meta property="og:url" content={SITE} />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:site" content="@jammutkarsh" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content={ogImage} />

	{@html `<script type="application/ld+json">${JSON.stringify(jsonLd)}<\/script>`}
</svelte:head>

<div class="hero">
	<h6>github follow stats</h6>
	<h1>Who unfollowed you,<br />in cold hard JSON.</h1>
	<p class="lead">
		Clout-chasing, but for devs. See who follows back, who's ghosting you, and fix the
		ratio one petty click at a time.
	</p>

	<img class="meme" src="/meme.jpg" alt="drake meme about checking github followers" width="500" height="500" />

	<a href="/auth/login" class="btn btn-primary arrow" data-sveltekit-preload-data="off">Sign in with GitHub</a>
</div>

<style>
	.hero {
		max-width: var(--ds-content-width);
		margin: var(--ds-space-16) auto 0;
		padding: 0 var(--ds-space-4);
	}

	h6 {
		color: var(--ds-primary);
		margin-bottom: var(--ds-space-2);
	}

	.lead {
		margin: var(--ds-space-4) 0 var(--ds-space-6);
	}

	.meme {
		display: block;
		width: 100%;
		max-width: 320px;
		height: auto;
		margin: 0 0 var(--ds-space-8);
		border: 1px dashed var(--ds-border);
		border-radius: var(--ds-radius);
	}
</style>
