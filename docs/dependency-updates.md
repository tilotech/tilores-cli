# Dependency updates

The Go modules of this project are updated every week by the Go updater of
[tilotech/dependency-bot](https://gitlab.com/tilotech/dependency-bot), through
its pipeline schedule `tilores-cli (Golang)`. Its
[README](https://gitlab.com/tilotech/dependency-bot/-/blob/main/README.md),
section *Go updates*, has the details. Its policy for this project is
`policies/tilores-cli.json` there.

## What arrives

Merge requests by the `dependabot` account, one per group:

| Branch | Contains | Merges by itself |
|---|---|---|
| `dependabot-gomod-go` | minor and patch updates, the Tilotech modules (`gitlab.com/tilotech/*`, `github.com/tilotech/*`) included | yes, once `gatekeeper` approves it and the pipeline is green |
| `dependabot-gomod-go-major` | major updates | never |

Each run builds these branches again from `main` with `go get`, `go mod tidy`
and `go mod vendor`, so `vendor/` changes with them. The branch pipeline is
the only test of an update.

## A merge request that waits

A major update or a red pipeline waits for a developer:

1. Branch from the bot's branch, or from `main` to take a single major. Make
   what the update needs, and open a merge request of your own. It needs an
   approval by another person.
2. Do not push to the bot's branch. The push takes the bot's approval away,
   and the updater then leaves the branch alone.
3. Once your merge request is merged, the next run closes the bot's one.

A closed merge request comes back at the next run.

## Keeping a module back

Add a hold to `.gitlab/go-updates.json`. All fields are required, and `scope`
is `major` or `all`:

```json
{
  "holds": [
    {
      "module": "github.com/example/lib",
      "scope": "major",
      "breaks": "v2 removes Client.Do, which internal/x calls",
      "removeWhen": "internal/x uses Client.Send",
      "added": "2026-11-25",
      "ref": "tilotech/tilores-cli!123"
    }
  ]
}
```

## Rules this relies on

- `main` stays tidy and vendored: `go mod tidy` and `go mod vendor` change
  nothing. Otherwise the updater stops for this project.
- The `go` line of `go.mod` stays as it is. A module that needs a newer Go is
  listed under *Not included* in the merge request, and the bot's job fails
  until someone raises the `go` line.
- `main` requires one approval by someone other than the author. For the
  bot's merge requests, `gatekeeper` gives it when the update passes the
  policy.
- The bot's merge request merges by itself only after a successful pipeline:
  GitLab's auto-merge waits for it.
