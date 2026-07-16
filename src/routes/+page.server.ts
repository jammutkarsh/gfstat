import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
	if (cookies.get('gh_token')) throw redirect(302, '/dashboard');
};
