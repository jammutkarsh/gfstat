import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getProfile, listAll, setFollow, AuthError } from '$lib/server/github';
import { mutuals, notFollowingMeBack, iDontFollowBack } from '$lib/compare';
import { getCache, putCache, deleteCache } from '$lib/server/cache';
import type { CacheEntry } from '$lib/types';

export const load: PageServerLoad = async ({ cookies, platform }) => {
	const token = cookies.get('gh_token');
	if (!token) throw redirect(302, '/');

	try {
		const profile = await getProfile(token);
		let cache = await getCache(platform, profile.id);
		let followers, following, fetchedAt: number;

		if (cache) {
			followers = cache.followers;
			following = cache.following;
			fetchedAt = cache.fetchedAt;
		} else {
			[followers, following] = await Promise.all([
				listAll(token, 'followers'),
				listAll(token, 'following'),
			]);
			cache = { profile, followers, following, fetchedAt: Date.now() };
			await putCache(platform, profile.id, cache);
			fetchedAt = cache.fetchedAt;
		}

		return {
			profile,
			mutuals: mutuals(followers, following),
			notFollowingMeBack: notFollowingMeBack(followers, following),
			iDontFollowBack: iDontFollowBack(followers, following),
			fetchedAt,
		};
	} catch (e) {
		if (e instanceof AuthError) throw redirect(302, '/logout');
		throw e;
	}
};

const LOGIN_RE = /^[a-zA-Z0-9-]{1,39}$/;

export const actions = {
	follow: async ({ request, cookies, platform }) => {
		const token = cookies.get('gh_token');
		if (!token) throw redirect(302, '/');
		const data = await request.formData();
		const login = data.get('login')?.toString() ?? '';
		if (!LOGIN_RE.test(login)) return { error: 'invalid login', status: 400 };
		try {
			await setFollow(token, login, true);
			const profile = await getProfile(token);
			const cache = await getCache(platform, profile.id);
			if (cache) {
				const f: CacheEntry = cache as CacheEntry;
				if (!f.following.some((x) => x.login === login)) {
					f.following.push({
						login,
						avatarUrl: `https://avatars.githubusercontent.com/${login}`,
						htmlUrl: `https://github.com/${login}`,
					});
				}
				await putCache(platform, profile.id, f);
			}
		} catch (e) {
			if (e instanceof AuthError) throw redirect(302, '/logout');
			return { error: 'follow failed', status: 502 };
		}
	},
	unfollow: async ({ request, cookies, platform }) => {
		const token = cookies.get('gh_token');
		if (!token) throw redirect(302, '/');
		const data = await request.formData();
		const login = data.get('login')?.toString() ?? '';
		if (!LOGIN_RE.test(login)) return { error: 'invalid login', status: 400 };
		try {
			await setFollow(token, login, false);
			const profile = await getProfile(token);
			const cache = await getCache(platform, profile.id);
			if (cache) {
				const f: CacheEntry = cache as CacheEntry;
				f.following = f.following.filter((x) => x.login !== login);
				await putCache(platform, profile.id, f);
			}
		} catch (e) {
			if (e instanceof AuthError) throw redirect(302, '/logout');
			return { error: 'unfollow failed', status: 502 };
		}
	},
	refresh: async ({ cookies, platform }) => {
		const token = cookies.get('gh_token');
		if (!token) throw redirect(302, '/');
		try {
			const profile = await getProfile(token);
			await deleteCache(platform, profile.id);
		} catch (e) {
			if (e instanceof AuthError) throw redirect(302, '/logout');
		}
	},
} satisfies Actions;
