# Registry PoC
Relates to https://github.com/opentofu/opentofu/issues/741

This PoC demonstrates how GitHub repository could be used as the source of truth for the providers and modules of OpenTofu, and later be hosted in GH Pages for OpenTofu CLI to consume

## How this works:

- Any changes made in the `main` branch are immediately replicated in registry acceptable format in the `host` branch, using `go run ./cmd/main.go publish ...`
- Existing providers/modules are being checked periodically for newer versions, using `go run ./cmd/main.go update ...`

## How to add a new provider/module:

- Simply fork the repository, and run `go run ./cmd/main.go initialize providers/<NAMESPACE>/<PROVIDER>.json` or `go run ./cmd/main.go initialize modules/<NAMESPACE>/<NAME>/<SYSTEM>.json`, to locally create the requested JSON file. You need to set a `export GH_TOKEN=<MY_GITHUB_TOKEN>` for the API calls to work
- After creating the files, push them to your fork, and create a PR to the base repository. Once the files are in, they will be processed for hosting

Adding providers/modules to the registry will probably not happen often, and will require manual approval of the maintenance team. However, we could provide some means of automating the process in the future, if we'd like

## How to bump versions for an existing provider/module

There are a couple of possible approaches:

- **The approach this current PoC currently implements:** Periodic update of all providers and modules. To save on GH API calls, and to not be throttled by GitHub, we only look for new releases for repositories for which the latest semver tag is not part of the existing versions
  - There are some cases where latest tag being missing from the file doesn't really indicate there are new releases. Mainly, in cases of prerelease tags, but also in cases of broken tags that were created for some reason for the providers.
    - The logic here could be improved, to somehow decide which such tags are irrelevant 
  - We could decide to not auto-update all versions from the get-go, and only start with the more "legit" providers and modules
  - If needed, because of GH API throttling in the future, we could change the periodic bump (bump 1/6 of the providers every 6 hours, for example). We could also implement smarter backoffs depending on the GH API limit, if necessary
- Manual upgrade - Anyone could open a PR to add new versions to a provider/module
  - We could auto-approve and auto-merge based on some basic checks (check the SHAs are correct, and that the signatures are correct, for example)
  - We could provide a GH Action that could be incorporated in the GH release process, to automatically create a bump PR
  - We could also provide a CLI (possible under OpenTofu CLI itself) that could also create the bump PR for you
- A hybrid of both approaches is possible (Periodically updating all providers/modules + allowing manual update for faster version bump / Only periodically update some providers/modules + other providers/modules would require manual bumps)
