# Quality gates

## Tooling

Tests use Go's native `testing` package and the existing Testify assertions.
[gotestsum](https://github.com/gotestyourself/gotestsum) runs those tests unchanged,
prints failures and skips, and writes JUnit XML and Go test JSON. Failed tests are
not automatically retried.

[golangci-lint v2](https://golangci-lint.run/docs/welcome/install/) provides static
analysis, including govet, Staticcheck, gosec, and Testify checks. Its version is
pinned in `.golangci-lint-version`. Make installs the official release binary with
its checksum-verifying installer, pinned to an upstream commit. It does not use
a globally installed linter or build a potentially different dependency graph.
CI uses the official golangci-lint action for caching and annotations, reading
the same version file and configuration.

[govulncheck](https://go.dev/doc/tutorial/govulncheck) checks reachable vulnerable
functions in the application and its dependencies, including the Go standard
library. Findings fail the gate; they are not suppressed for an old compiler.

[actionlint](https://github.com/rhysd/actionlint) checks workflow syntax,
expressions, and embedded shell commands through ShellCheck when available.
CI's Ubuntu runner includes ShellCheck.

The Go tools are pinned with native `tool` directives in `tools/go.mod`, isolated
from the application's dependency graph. `go.mod` selects the supported Go
patch version; the tools module follows it. Dependency and action updates are
proposed weekly by Dependabot. Update the linter's version file deliberately
and run the gates before accepting a new linter release.

## Local commands

Go, Make, Bash, and curl are required on Linux or macOS. Race tests also need
a C compiler.

```sh
make tools       # prepare pinned tools; no global installs
make test        # offline unit, architecture, contract, and golden tests
make race        # the same tests with the race detector
make coverage    # race tests, JUnit/JSON reports, and atomic coverage
make lint        # lint the host platform
make vuln        # scan the host platform for reachable vulnerabilities
make check       # formatting, module consistency, workflows, lint, security, coverage
make dist        # Linux amd64 and macOS amd64/arm64 release binaries
```

`make check` needs network access for missing tools/modules and the vulnerability
database. Reports are written under `build/reports`, which is ignored by Git.
Coverage is reported rather than gated on an arbitrary percentage.

Run platform-specific analysis without cross-compiling the analysis tools:

```sh
make lint vuln TARGET_GOOS=linux
make lint vuln TARGET_GOOS=darwin
```

Live integration is separate from the offline suite. On Linux or macOS, with
Herdr 0.9.5 installed:

```sh
HERDR_SCOPE_HERDR=/absolute/path/to/herdr make integration
```

The test uses a private short-lived HOME, XDG directories, socket, and server;
it never attaches to the user's session. Plugin processes are race-instrumented;
private race reports and shutdown failures fail the test, including cleanup.
An explicitly requested integration run fails if Herdr is unavailable. CI downloads the tested Linux release and
verifies its pinned SHA256 before running it.

Use `make fmt` to apply the configured formatters, `make tidy` to update module
files, and `make golden` only for intended visual changes. The corresponding CI
checks do not rewrite source, module files, or expected output.

## GitHub Actions

`.github/workflows/ci.yml` runs on master pushes, pull requests, merge queues,
manual dispatch, and a weekly schedule. It is also reusable by the release
workflow. It provides:

- Native race tests and coverage on Linux only, with the module's Go version
  and the latest stable Go. CI disables automatic
  toolchain switching so each run really uses its selected compiler.
- Lint and source vulnerability analysis for Linux and macOS.
- Formatting, module tidy/verification, and workflow validation.
- A real, isolated Herdr 0.9.5 integration run on Linux.
- Release cross-compilation for Linux amd64 and macOS amd64/arm64, plus the
  existing installer checks.
- JUnit, JSON, and coverage artifacts, retained even when tests fail.

`CI passed` depends on every gate and fails if any gate fails, is cancelled, or
is skipped. Configure the master branch ruleset to require this status; workflow
files alone cannot enforce merge protection. The repository had no active
branch protection or rulesets at the time of this audit. No repository settings
are changed by this checkout.

Release builds and non-publishing release rehearsals call the same CI workflow
before packaging. Stable release tags must point to commits on master;
pre-release tags and non-publishing rehearsals may use a development branch.
Publishing still requires a version tag and happens in a separate write-enabled job. A test-only Ubuntu run can no longer bypass the
lint, race, security, and integration gates.

CI tests run on Linux only; macOS is cross-built for releases without native
macOS test jobs. Linux arm64 is not a release target.

Windows is excluded from automated platform testing, analysis, and releases.
Existing Windows source and legacy installer code are retained, but no Windows
binaries or ZIP archives are produced.
