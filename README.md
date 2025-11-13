# Overview

`talosctlenv` is a CLI tool designed to manage multiple versions of `talosctl` easily.
It provides commands to list available versions, install specific versions, and switch between them.
The tool integrates with GitHub to fetch official releases and supports filtering and sorting using semantic versioning.

# Sub-commands

## list

TODO

## list-remote

The `list-remote` command retrieves all available `talosctl` versions from the official GitHub repository (`siderolabs/talos`).

### Behavior

- Fetches releases from GitHub.
- By default, **only stable versions** are shown (pre-release versions like alpha, beta, rc are excluded).

### Flags

- `--devel`: Include pre-release versions (alpha, beta, rc) in the output.
  
- `--limit <number>`: Limit the number of versions displayed. The tool will stop fetching once the limit is reached.

### Authentication

To avoid hitting GitHub’s unauthenticated rate limit (60 requests/hour), you can set the `GITHUB_TOKEN` environment variable:

```bash
export GITHUB_TOKEN=your_personal_access_token
```

This increases the limit to 5000 requests/hour and ensures smoother operation.

## install

TODO

## use

TODO
