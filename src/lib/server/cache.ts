import type { CacheEntry } from '$lib/types';

const mem = new Map<string, string>();

// ponytail: KV-backed cache, in-memory Map falls back when platform.env absent (local dev)
export function kv(platform: App.Platform | undefined) {
	const sessions = platform?.env?.SESSIONS;
	if (sessions) return sessions;
	return {
		async get(key: string): Promise<string | null> {
			return mem.get(key) ?? null;
		},
		async put(key: string, value: string, opts?: { expirationTtl?: number }): Promise<void> {
			mem.set(key, value);
			if (opts?.expirationTtl) {
				setTimeout(() => mem.delete(key), opts.expirationTtl * 1000);
			}
		},
		async delete(key: string): Promise<void> {
			mem.delete(key);
		},
	};
}

export async function getCache(
	platform: App.Platform | undefined,
	userId: number,
): Promise<CacheEntry | null> {
	const raw = await kv(platform).get(`user:${userId}`);
	if (!raw) return null;
	return JSON.parse(raw) as CacheEntry;
}

export async function putCache(
	platform: App.Platform | undefined,
	userId: number,
	entry: CacheEntry,
): Promise<void> {
	await kv(platform).put(`user:${userId}`, JSON.stringify(entry), { expirationTtl: 86400 });
}

export async function deleteCache(
	platform: App.Platform | undefined,
	userId: number,
): Promise<void> {
	await kv(platform).delete(`user:${userId}`);
}
