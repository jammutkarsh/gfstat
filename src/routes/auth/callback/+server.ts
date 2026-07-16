import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { exchangeCode } from '$lib/server/github';

export async function GET({ url, cookies }) {
	const code = url.searchParams.get('code');
	const state = url.searchParams.get('state');
	const savedState = cookies.get('oauth_state');

	if (!code) return new Response('missing code', { status: 400 });
	if (!state || state !== savedState) return new Response('state mismatch', { status: 400 });

	cookies.delete('oauth_state', { path: '/' });

	const token = await exchangeCode(code, env.GITHUB_CLIENT_ID, env.GITHUB_CLIENT_SECRET);
	cookies.set('gh_token', token, {
		httpOnly: true,
		secure: true,
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 30,
	});

	throw redirect(302, '/dashboard');
}
