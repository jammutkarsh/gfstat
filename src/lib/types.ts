export type Follow = { login: string; avatarUrl: string; htmlUrl: string };

export type Profile = {
	id: number;
	login: string;
	name: string | null;
	bio: string | null;
	avatarUrl: string;
	htmlUrl: string;
	followers: number;
	following: number;
	publicRepos: number;
	location: string | null;
	blog: string | null;
	company: string | null;
};

export type CacheEntry = {
	profile: Profile;
	followers: Follow[];
	following: Follow[];
	fetchedAt: number;
};
