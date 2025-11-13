# Overview

`talosctlenv` is a CLI tool designed to manage multiple versions of `talosctl` easily.
It provides commands to list available versions, install specific versions, and switch between them.
The tool integrates with GitHub to fetch official releases and supports filtering and sorting using semantic versioning.

See the [docs folder](./docs/) for more information.

## A note on Authentication

This tool makes use of the GitHub APIs.
You can use it without setting up anything if your use is infrequent.

But, if you plan on downloading a lot of versions (or frequently check which versions are available),
I strongly recommend setting up a `GITHUB_TOKEN` to get around the API rate-limiting.

Unauthenticated calls to the APIs are limited to 60 per hour.
Using a token you can get up to 5,000 per hour.

During development I observed, for example, that the `list-remote` can make up to at least 12 API calls to retrieve the whole list.
Do it 5 times and you're done.

### How to

To create a token:
1. go to https://github.com/settings/tokens/new
2. add a note
3. set the expiration you like
4. select scope: `repo > public_repo`.
5. scroll to the bottom of the page and click "Generate token"
6. copy the token that will appear

Now, up to you where you save it.
In the repo you have a `.env.template` file you can use: put your token there and rename the file to `.env` (don't worry, it's in the .gitignore).
Then, when you need to run the tool do:

```sh
source .env
```

Now run the tool and you'll have the limit upped to 5,000 calls an hour.
