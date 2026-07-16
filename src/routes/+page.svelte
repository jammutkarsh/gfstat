<script lang="ts">
	const SITE = 'https://gfstat.utkarshchourasia.in';
	const title = 'gfstat — who unfollowed you, in cold hard JSON';
	const description =
		"Clout-chasing, but for devs. See your GitHub mutuals, who's ghosting you, and whom you don't follow back — then fix the ratio one petty click at a time.";
	const ogImage = `${SITE}/meme.png`;
	const ogImageAlt = 'drake meme: rejecting checking who liked your Instagram post, approving checking who unfollowed you on GitHub';

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

	// escape "<" so nothing in the payload can break out of the tag (script close, comment open, etc.)
	// the closing tag is concatenated, not written literally, so it can't be mistaken for real markup
	const jsonLdScript =
		'<script type="application/ld+json">' +
		JSON.stringify(jsonLd).replace(/</g, '\\u003c') +
		'<' + '/script>';
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
	<meta property="og:image:width" content="1200" />
	<meta property="og:image:height" content="630" />
	<meta property="og:image:alt" content={ogImageAlt} />
	<meta property="og:url" content={SITE} />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:site" content="@jammutkarsh" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content={ogImage} />
	<meta name="twitter:image:alt" content={ogImageAlt} />

	{@html jsonLdScript}
</svelte:head>

<div class="hero">
	<h6>github follow stats</h6>
	<h1>Who unfollowed you,<br />in cold hard JSON.</h1>
	<p class="lead">
		Clout-chasing, but for devs. See who follows back, who's ghosting you, and fix the
		ratio one petty click at a time.
	</p>

	<img class="meme" src="/meme.png" alt="drake meme about checking github followers" width="1200" height="630" />

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
