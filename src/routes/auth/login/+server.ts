import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export function GET({ cookies }) {
	const state = crypto.randomUUID();
	cookies.set('oauth_state', state, {
		httpOnly: true,
		secure: true,
		sameSite: 'lax',
		path: '/',
		maxAge: 600,
	});
	const url = new URL('https://github.com/login/oauth/authorize');
	url.searchParams.set('client_id', env.GITHUB_CLIENT_ID);
	url.searchParams.set('scope', 'read:user user:follow');
	url.searchParams.set('state', state);
	throw redirect(302, url.toString());
}
