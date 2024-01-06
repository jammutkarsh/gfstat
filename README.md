# gfStat (GitHub Follow Stat)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub Release](https://img.shields.io/github/release/JammUtkarsh/gfstat.svg?style=flat)

gfStat is a tool for GitHub users that provides insights into your GitHub followers and following. With gfStat, you can easily discover:

- **Mutual followers**: See who follows you and whom you follow back.
- **Followers you don't follow**: Find out who follows you, but you don't follow back.
- **Following that don't follow you**: Identify GitHub users you follow, but they don't follow you in return.

## Why gfStat?

It's a side project to learn a bunch of new technologies and tools and practice writing. I am using Go again because I've been doing data structures and algorithms for a long enough time now which made me unable to write effective Go code.

## Tech Involved

- About 99% of the project is built in [Go](https://go.dev/) the rest of it is simple HTML templating.
- The only external dependency that I have used is [GitHub's Go SDK](https://pkg.go.dev/github.com/google/go-github/v56#section-readme), the rest of application is built standard library.


## Local Dev

To run and debug this program locally:

Prepare a `.env` file

```env
GH_BASIC_CLIENT_ID=1234abc
GH_BASIC_SECRET_ID=1234xyz
```

These two values which you get from [New OAuth App](https://github.com/settings/developers). Then:

```bash
export $(cat .env | xargs) # to set env vars from .env
go run . # you are up and running.
```

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%201.svg)](https://www.digitalocean.com/?refcode=93388f4a1ca0&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)
