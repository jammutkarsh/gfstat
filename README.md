# gfStat (GitHub Follow Stat)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub Release](https://img.shields.io/github/release/JammUtkarsh/gfstat.svg?style=flat)

gfStat is a tool for GitHub users that provides insights into your GitHub followers and following. With gfStat, you can easily discover:

- **Mutual followers**: See who follows you and whom you follow back.
- **Followers you don't follow**: Find out who follows you, but you don't follow back.
- **Following that don't follow you**: Identify GitHub users you follow, but they don't follow you in return.

## Why gfStat?

It's a side project to learn a bunch of new technologies and tools and practice writing. Go again since I've been doing data structures and algorithms for a long time now.

Also, in this online world, where one can get paid to tweet, it's fun to have some statistics about your profile's growth to build better stuff.
Or it's simply a tool to stalk your followers and following. 😜

The tech stack I plan to use in this project might not be shipped in `V1.0`, but they will definitely be released in future versions. The tech stack includes:

- [Go](https://go.dev/)
- [React](https://react.dev/) and/or [HTMX](https://htmx.org/)
- [GraphQL](https://graphql.org/) and [REST](https://restfulapi.net/)

A few more things that don't come under the tech stack, but I plan to do:

- It will be a CLI-first tool that can also be used to run the server. Something like `gfstat -u <username>` will generate a CLI output, and `gfstat serve` will start an API server.
- This API can be further used to build a frontend using React or HTMX.
- The initial launch of the API server will support REST with limited fields and later will be upgraded to GraphQL so that users can query the data they want.

## High-Level Roadmap

- `V0.x` CLI support to ingest a single user and output as JSON.
- `V1.0` Option to output as JSON or markdown in CLI.
- `V1.x` CLI supports ingesting multiple users in one go. Like having an input file and an output file.
- `V2.0` Support for starting an API server with GET requests and query parameters for usernames (*subject to change*)
- `V2.x` Support for paginated responses
- `V3.0` GraphQL support

Note: The roadmap is subject to change and will be updated as and when required.
