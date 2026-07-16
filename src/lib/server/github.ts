import type { Follow, Profile } from '$lib/types';

const GH = 'https://api.github.com';
const HEADERS = {
	Accept: 'application/vnd.github+json',
	'User-Agent': 'gfstat',
};

export class AuthError extends Error {
	constructor(msg: string) {
		super(msg);
		this.name = 'AuthError';
	}
}

export async function exchangeCode(
	code: string,
	clientId: string,
	clientSecret: string,
): Promise<string> {
	const body = new URLSearchParams({ client_id: clientId, client_secret: clientSecret, code });
	const res = await fetch('https://github.com/login/oauth/access_token', {
		method: 'POST',
		headers: { ...HEADERS, 'Content-Type': 'application/x-www-form-urlencoded' },
		body,
	});
	const data = (await res.json()) as { access_token?: string; error?: string };
	if (!data.access_token) throw new Error(data.error ?? 'no access_token returned');
	return data.access_token;
}

export async function getProfile(token: string): Promise<Profile> {
	const res = await ghFetch('/user', token);
	ghCheckAuth(res);
	const u = (await res.json()) as Record<string, unknown>;
	return {
		id: u.id as number,
		login: u.login as string,
		name: (u.name as string) ?? null,
		bio: (u.bio as string) ?? null,
		avatarUrl: u.avatar_url as string,
		htmlUrl: u.html_url as string,
		followers: u.followers as number,
		following: u.following as number,
		publicRepos: u.public_repos as number,
		location: (u.location as string) ?? null,
		blog: (u.blog as string) ?? null,
		company: (u.company as string) ?? null,
	};
}

async function listPage(
	token: string,
	kind: 'followers' | 'following',
	page: number,
): Promise<{ rows: Follow[]; link: string | null }> {
	const url = `${GH}/user/${kind}?per_page=100&page=${page}`;
	const res = await ghFetchRaw(url, token);
	ghCheckAuth(res);
	const rows: Follow[] = ((await res.json()) as Array<Record<string, unknown>>).map((r) => ({
		login: r.login as string,
		avatarUrl: r.avatar_url as string,
		htmlUrl: r.html_url as string,
	}));
	return { rows, link: res.headers.get('Link') };
}

export async function listAll(token: string, kind: 'followers' | 'following'): Promise<Follow[]> {
	const { rows, link } = await listPage(token, kind, 1);
	const all: Follow[] = [...rows];
	const lastMatch = link?.match(/[?&]page=(\d+)>; rel="last"/);
	if (lastMatch) {
		const last = parseInt(lastMatch[1]);
		const pages = [];
		for (let p = 2; p <= last; p++) pages.push(listPage(token, kind, p));
		const results = await Promise.all(pages);
		for (const r of results) all.push(...r.rows);
	}
	return all;
}

export async function setFollow(token: string, login: string, follow: boolean): Promise<void> {
	const url = `${GH}/user/following/${login}`;
	const res = await ghFetchRaw(url, token, follow ? 'PUT' : 'DELETE');
	ghCheckAuth(res);
	if (res.status !== 204) throw new Error(`GitHub returned ${res.status} ${res.statusText}`);
}

export async function getUserDetail(token: string, login: string) {
	const res = await ghFetch(`/users/${login}`, token);
	ghCheckAuth(res);
	const u = (await res.json()) as Record<string, unknown>;
	return {
		name: (u.name as string) ?? null,
		bio: (u.bio as string) ?? null,
		followers: u.followers as number,
		following: u.following as number,
		publicRepos: u.public_repos as number,
		publicGists: u.public_gists as number,
		location: (u.location as string) ?? null,
		company: (u.company as string) ?? null,
		blog: (u.blog as string) || null,
		twitterUsername: (u.twitter_username as string) ?? null,
		email: (u.email as string) ?? null,
		hireable: !!u.hireable,
		createdAt: u.created_at as string,
	};
}

function ghFetch(path: string, token: string): Promise<Response> {
	return ghFetchRaw(`${GH}${path}`, token);
}

function ghFetchRaw(url: string, token: string, method = 'GET'): Promise<Response> {
	return fetch(url, {
		method,
		headers: { ...HEADERS, Authorization: `Bearer ${token}` },
		...(method === 'PUT' ? { body: null } : {}),
	});
}

export function ghCheckAuth(res: Response): void {
	if (res.status === 401) throw new AuthError('token expired or revoked');
}
