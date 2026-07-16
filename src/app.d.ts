declare global {
	namespace App {
		interface Platform {
			env?: {
				SESSIONS: KVNamespace;
			};
		}
	}
}

interface KVNamespace {
	get(key: string, options?: { type: 'text' }): Promise<string | null>;
	put(key: string, value: string, options?: { expirationTtl?: number }): Promise<void>;
	delete(key: string): Promise<void>;
}

export {};
