# cli - A minimal CLI command package using pflag

[![Go Reference](https://pkg.go.dev/badge/github.com/titpetric/cli.svg)](https://pkg.go.dev/github.com/titpetric/cli)

Package `cli` is a small, opinionated command and flag framework built on
[`spf13/pflag`](https://github.com/spf13/pflag). It provides command dispatch,
command-scoped flags, environment-backed flag values, usage text, and a
signal-aware `context.Context` without the larger API surface of Cobra.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/titpetric/cli"
)

func main() {
	app := cli.NewApp("greet")
	app.DefaultCommand = "hello"

	app.AddCommand("hello", "Print a greeting", func() *cli.Command {
		var name string

		return &cli.Command{
			Default: true,
			Usage: func() string {
				return "Print a greeting for NAME."
			},
			Bind: func(fs *cli.FlagSet) {
				fs.StringVarP(&name, "name", "n", "world", "name to greet")
			},
			Run: func(ctx context.Context, args []string) error {
				fmt.Printf("Hello, %s!\n", name)
				return nil
			},
		}
	})

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

The command can then be invoked explicitly or as the default:

```console
$ greet hello --name Ada
Hello, Ada!
$ greet -n Ada
Hello, Ada!
```

`AddCommand` supplies `Command.Name` and `Command.Title` when the constructor
leaves them empty. Set `Command.Default` when its help usage should omit the
command name.

## Flags and environment variables

`Bind` receives a command-scoped `*cli.FlagSet`, which is an alias for
`*pflag.FlagSet`. Flags may appear before or after positional arguments, and
`--` stops flag parsing. `Run` receives the remaining positional arguments
with an explicitly selected command name removed.

Before parsing arguments, `ParseWithFlagSet` applies matching environment
variables to flags that were not already changed. Environment names are
lowercased and underscores become hyphens, so `DB_DSN` supplies the value for
`--db-dsn`. Environment variables without an underscore are ignored. Invalid
environment values are returned as errors.

## Help and errors

- `-h` and `--help` print application or command help and return `nil`.
- An unknown or missing command prints application usage and returns an error.
- Invalid flags print command usage and return the parsing error.
- Errors returned by `Command.Run`, including cancellation, are returned
  without printing usage.

`Run` creates a context canceled by `SIGINT` or `SIGTERM` and passes it to the
selected command. Applications remain responsible for printing returned
errors and selecting an exit status.

## Projects using this package

- [github.com/go-bridget/mig](https://github.com/go-bridget/mig) - database migration tooling
- [github.com/titpetric/atkins](https://github.com/titpetric/atkins) - a local command runner for CI
- [github.com/titpetric/vuego-cli](https://github.com/titpetric/vuego-cli) - a Vuego template engine documentation server
- [github.com/titpetric/etl](https://github.com/titpetric/etl) - database-agnostic database tooling
- [github.com/titpetric/exp](https://github.com/titpetric/exp) - experimental CLI tooling, notably `go-fsck`
