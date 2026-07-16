import type { Follow } from './types';

export function mutuals(followers: Follow[], following: Follow[]): Follow[] {
	const fSet = new Set(followers.map((f) => f.login));
	return following.filter((f) => fSet.has(f.login));
}

export function notFollowingMeBack(followers: Follow[], following: Follow[]): Follow[] {
	const fSet = new Set(followers.map((f) => f.login));
	return following.filter((f) => !fSet.has(f.login));
}

export function iDontFollowBack(followers: Follow[], following: Follow[]): Follow[] {
	const wSet = new Set(following.map((f) => f.login));
	return followers.filter((f) => !wSet.has(f.login));
}
