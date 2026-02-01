# CLI package

[![Go Reference](https://pkg.go.dev/badge/github.com/titpetric/cli.svg)](https://pkg.go.dev/github.com/titpetric/cli)
[![Coverage](https://img.shields.io/badge/coverage-59.40%25-brightgreen.svg)](https://github.com/titpetric/cli)

This package contains the implementation for a minimal opinionated flags
framework similar to spf13/cobra. It all centers around the
`cli.Command` type but provides less functionality.

The package is low-dependency, only relying on the pflag package, to
provide `--` unix flag options.

This package is used in:

- [github.com/go-bridget/mig](https://github.com/go-bridget/mig) - database migration tooling,
- [github.com/titpetric/atkins](https://github.com/titpetric/atkins) - a local command runner for CI,
- [github.com/titpetric/vuego-cli](https://github.com/titpetric/vuego-cli) - a vuego template engine docs server, CLI,
- [github.com/titpetric/etl](https://github.com/titpetric/etl) - database agnostic tooling to interface databases
- [github.com/titpetric/exp](https://github.com/titpetric/exp) - experimental CLI tooling, notably `go-fsck`

It's existed in some form for many years and I kept rewriting the CLIs
to the next best thing, this exists for me to stop worring about flags
like not an already solved problem. Reuse what works.

- no environment handling (put it in default value with `os.Getenv`)
- provide defaults with or without env support
- scoped flagset bindings with Bind()
- error handling, usage, help

The benchmark are tools like `git`, `docker`, `docker compose`, `go`
where there is a mixture of command arguments (`git status`, `git pull`,
...) and flags like `-d`, `--force`. This package is the minimal shed in
the back, the "cli package we have at home".