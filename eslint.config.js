import js from '@eslint/js';
import ts from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import svelteParser from 'svelte-eslint-parser';
import globals from 'globals';

export default ts.config(
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		rules: {
			'@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
			'svelte/no-navigation-without-resolve': 'off',
			'svelte/prefer-svelte-reactivity': 'off',
		},
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parser: svelteParser,
			parserOptions: { parser: ts.parser },
			globals: { ...globals.browser },
		},
	},
	{
		ignores: [
			'.svelte-kit/**',
			'build/**',
			'node_modules/**',
			'.output/**',
			'.wrangler/**',
		],
	},
);
