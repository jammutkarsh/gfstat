import { redirect } from '@sveltejs/kit';

export function GET({ cookies }) {
	cookies.delete('gh_token', { path: '/' });
	throw redirect(302, '/');
}
