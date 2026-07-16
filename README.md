# gfstat — GitHub Follow Stats

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

Know your GitHub follow stats: mutuals, who doesn't follow you back, whom you don't follow back — with follow/unfollow right from the app.

## Tech

SvelteKit (TypeScript, server-side), Cloudflare Pages, [utc-ds](https://github.com/jammutkarsh/design-system) design system. Plain `fetch` against GitHub REST API.

## Local Dev

```bash
npm install
cp .env.example .env   # edit with your OAuth app credentials
npm run dev
```

### Required env vars

```env
GITHUB_CLIENT_ID=your_oauth_client_id
GITHUB_CLIENT_SECRET=your_oauth_client_secret
```

### OAuth App Setup

1. Register a new OAuth app at https://github.com/settings/developers
2. Homepage URL: `http://localhost:5173` (dev)
3. Authorization callback URL: `http://localhost:5173/auth/callback`
4. Copy Client ID and Secret to `.env`

## Deploy

```bash
npm run build
wrangler pages dev   # smoke test
wrangler pages deploy .svelte-kit/cloudflare
```

Prerequisites: Cloudflare KV namespace named `SESSIONS` (update `id` in `wrangler.toml`). Set secrets in Cloudflare dashboard: `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`.

Update OAuth app callback URL to `https://your-domain.pages.dev/auth/callback`.

## Commands

| Command | What |
|---------|------|
| `npm run dev` | Dev server |
| `npm run build` | Production build |
| `npm test` | Vitest |
| `npm run check` | Type-check |
