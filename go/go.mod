module github.com/AJayClauGermany/the-most-well-tested-checkout-system/go

go 1.26.0

toolchain go1.26.10

// The libraries the code and the tests use. Every version is pinned exactly, and
// go.sum holds their checksums (like package-lock.json in JavaScript).
require (
	github.com/cucumber/godog v0.16.0
	github.com/google/go-cmp v0.7.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3
	go.uber.org/mock v0.6.0
	golang.org/x/text v0.44.0
	modernc.org/sqlite v1.62.0
	pgregory.net/rapid v1.3.0
)

// The modules of the developer tools listed at the bottom (pinned the same way).
require (
	golang.org/x/vuln v1.9.0 // indirect
	gotest.tools/gotestsum v1.13.0 // indirect
	honnef.co/go/tools v0.8.2 // indirect
)

// Everything the libraries and tools above need in turn (filled in by `go mod tidy`).
require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/bitfield/gotestdox v0.2.2 // indirect
	github.com/cucumber/gherkin/go/v42 v42.0.0 // indirect
	github.com/cucumber/messages/go/v34 v34.2.0 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.5 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/telemetry v0.0.0-20261010193011-f71dc92da45e // indirect
	golang.org/x/term v0.35.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	modernc.org/libc v1.80.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)

// Developer tools (Go 1.24 "tool" directives). Run them with `go tool <name>`,
// e.g. `go tool staticcheck ./...` – nothing to install by hand, and everybody
// gets exactly the version pinned above.
tool (
	go.uber.org/mock/mockgen
	golang.org/x/vuln/cmd/govulncheck
	gotest.tools/gotestsum
	honnef.co/go/tools/cmd/staticcheck
)
