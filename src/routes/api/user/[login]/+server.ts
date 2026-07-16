import type { RequestHandler } from '@sveltejs/kit';
import { getUserDetail } from '$lib/server/github';
import { AuthError } from '$lib/server/github';

const LOGIN_RE = /^[a-zA-Z0-9-]{1,39}$/;

export const GET: RequestHandler = async ({ params, cookies }) => {
	const token = cookies.get('gh_token');
	if (!token) return new Response('unauthorized', { status: 401 });

	const login = params.login;
	if (!login || !LOGIN_RE.test(login)) return new Response('invalid login', { status: 400 });

	try {
		const detail = await getUserDetail(token, login);
		return new Response(JSON.stringify(detail), {
			headers: {
				'content-type': 'application/json',
				'cache-control': 'private, max-age=3600',
			},
		});
	} catch (e) {
		if (e instanceof AuthError) return new Response('token expired', { status: 401 });
		console.error('proxy error', e);
		return new Response('proxy error', { status: 502 });
	}
};
