import { describe, expect, it } from 'vitest';
import { mutuals, notFollowingMeBack, iDontFollowBack } from './compare';
import type { Follow } from './types';

const input: { followers: Follow[]; following: Follow[] } = {
	followers: [
		{ login: 'JammUtkarsh', htmlUrl: 'https://github.com/JammUtkarsh', avatarUrl: '' },
		{ login: 'MadhaviGupta', htmlUrl: 'https://github.com/MadhaviGupta', avatarUrl: '' },
		{ login: 'Mishank24', htmlUrl: 'https://github.com/Mishank24', avatarUrl: '' },
		{ login: 'SparshGarg1', htmlUrl: 'https://github.com/SparshGarg1', avatarUrl: '' },
		{ login: 'SujalSamai', htmlUrl: 'https://github.com/SujalSamai', avatarUrl: '' },
		{ login: 'TheGameisYash', htmlUrl: 'https://github.com/TheGameisYash', avatarUrl: '' },
		{ login: 'dcdeepesh', htmlUrl: 'https://github.com/dcdeepesh', avatarUrl: '' },
		{ login: 'golemvincible', htmlUrl: 'https://github.com/golemvincible', avatarUrl: '' },
		{ login: 'shreyash2002', htmlUrl: 'https://github.com/shreyash2002', avatarUrl: '' },
		{ login: 'shristigupta12', htmlUrl: 'https://github.com/shristigupta12', avatarUrl: '' },
		{ login: 'sushantsharma08', htmlUrl: 'https://github.com/sushantsharma08', avatarUrl: '' },
		{ login: 'tanishjain158', htmlUrl: 'https://github.com/tanishjain158', avatarUrl: '' },
	],
	following: [
		{ login: 'JammUtkarsh', htmlUrl: 'https://github.com/JammUtkarsh', avatarUrl: '' },
		{ login: 'MadhaviGupta', htmlUrl: 'https://github.com/MadhaviGupta', avatarUrl: '' },
		{ login: 'Mishank24', htmlUrl: 'https://github.com/Mishank24', avatarUrl: '' },
		{ login: 'SparshGarg1', htmlUrl: 'https://github.com/SparshGarg1', avatarUrl: '' },
		{ login: 'SujalSamai', htmlUrl: 'https://github.com/SujalSamai', avatarUrl: '' },
		{ login: 'TheGameisYash', htmlUrl: 'https://github.com/TheGameisYash', avatarUrl: '' },
		{ login: 'dcdeepesh', htmlUrl: 'https://github.com/dcdeepesh', avatarUrl: '' },
		{ login: 'dxaman', htmlUrl: 'https://github.com/dxaman', avatarUrl: '' },
		{ login: 'golemvincible', htmlUrl: 'https://github.com/golemvincible', avatarUrl: '' },
		{ login: 'shreyash2002', htmlUrl: 'https://github.com/shreyash2002', avatarUrl: '' },
		{ login: 'shristigupta12', htmlUrl: 'https://github.com/shristigupta12', avatarUrl: '' },
		{ login: 'sushantsharma08', htmlUrl: 'https://github.com/sushantsharma08', avatarUrl: '' },
	],
};

describe('mutuals', () => {
	it('finds mutual followers', () => {
		const got = mutuals(input.followers, input.following)
			.map((f) => f.login)
			.sort();
		const want = [
			'JammUtkarsh',
			'MadhaviGupta',
			'Mishank24',
			'SparshGarg1',
			'SujalSamai',
			'TheGameisYash',
			'dcdeepesh',
			'golemvincible',
			'shreyash2002',
			'shristigupta12',
			'sushantsharma08',
		].sort();
		expect(got).toEqual(want);
	});
});

describe('notFollowingMeBack', () => {
	it('finds following who do not follow back', () => {
		const got = notFollowingMeBack(input.followers, input.following).map((f) => f.login);
		expect(got).toEqual(['dxaman']);
	});
});

describe('iDontFollowBack', () => {
	it('finds followers you do not follow back', () => {
		const got = iDontFollowBack(input.followers, input.following).map((f) => f.login);
		expect(got).toEqual(['tanishjain158']);
	});
});
