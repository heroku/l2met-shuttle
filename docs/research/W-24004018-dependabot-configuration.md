# W-24004018: Dependabot configuration research

## Decision

Implement [W-24004018](https://gus.lightning.force.com/lightning/r/ADM_Work__c/a07EE00002iq2M4YAI/view) with a single `.github/dependabot.yml` entry for the repository's root Go module:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
```

The GUS item reports `.github/dependabot.yml not found` and links the [GitHub Dependabot configuration options](https://docs.github.com/en/code-security/supply-chain-security/keeping-your-dependencies-updated-automatically/configuration-options-for-dependency-updates) and an internal [Golden Path Dependency Management Guide](https://salesforce.quip.com/ic3QAQH6kvI6). It says the item will auto-close within 24 hours after resolution.

## Repository evidence

- [`go.mod`](../../go.mod) declares `github.com/heroku/l2met-shuttle` as a root-level Go module and lists its dependencies, so GitHub's `gomod` ecosystem and `/` directory apply.
- [`go.sum`](../../go.sum) records the module checksums.
- The checked-in `vendor/` tree needs no explicit option because GitHub says Go vendored dependencies are automatically identified and maintained ([`vendor` option reference](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/dependabot-options-reference#vendor--)).
- The repository has no GitHub Actions workflows. [`.travis.yml`](../../.travis.yml) configures Travis, so there is no `github-actions` ecosystem to monitor.
- [`CONTRIBUTING.md`](../../CONTRIBUTING.md) identifies `master` as the pull-request target and `go test ./...` as the pre-PR check.

## GitHub requirements

GitHub requires `version: 2`, an `updates` list, `package-ecosystem`, `directory` or `directories`, and `schedule.interval` ([Dependabot options reference](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/dependabot-options-reference#required-keys)). GitHub identifies `gomod` as the package ecosystem value for Go modules ([supported ecosystems](https://docs.github.com/en/code-security/dependabot/ecosystems-supported-by-dependabot/supported-ecosystems-and-repositories)). The configuration must be stored as `.github/dependabot.yml` or `.github/dependabot.yaml` on the default branch ([configuration file location](https://docs.github.com/en/code-security/concepts/supply-chain-security/about-the-dependabot-yml-file#where-to-store-the-dependabotyml-file)).

`weekly` is a valid GitHub interval and is the conservative choice for this low-activity repository. It is a recommendation rather than a verified Salesforce policy because the authenticated Golden Path document could not be retrieved during this investigation. If that guide mandates another interval or organization-specific options, its requirements supersede this recommendation.

## Validation

1. Parse `.github/dependabot.yml` as YAML and verify the required keys and values.
2. Run `go test ./...` as documented in [`CONTRIBUTING.md`](../../CONTRIBUTING.md).
3. Confirm GitHub accepts the configuration after it reaches `master`. GitHub begins monitoring dependencies when the configuration is added or updated ([how the configuration works](https://docs.github.com/en/code-security/concepts/supply-chain-security/about-the-dependabot-yml-file#how-the-dependabotyml-file-works)).
4. Allow up to 24 hours for the GUS item to auto-close, per the work item.

## Out of scope

- Updating `go.mod`, `go.sum`, or `vendor/`.
- Adding GitHub Actions or replacing Travis.
- Changing the Go version.
- Adding labels, reviewers, registries, grouping, ignore rules, pull-request limits, or auto-merge behavior without an owning requirement.
- Changing repository-level Dependabot alert or security-update settings, which GitHub manages separately from this configuration file ([Dependabot configuration overview](https://docs.github.com/en/code-security/concepts/supply-chain-security/about-the-dependabot-yml-file)).
