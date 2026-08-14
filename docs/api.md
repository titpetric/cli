# Package cli

```go
import (
	"github.com/titpetric/cli"
}
```
Package cli implements a minimal, opinionated command and flag framework
built on spf13/pflag.

An App registers command constructors with [App.AddCommand]. The selected
constructor creates a [Command], whose Bind callback defines scoped flags and
whose Run callback executes with a signal-aware context.

A minimal application looks like this:

	app := cli.NewApp("mig")
	app.AddCommand("version", "Print version information", func() *cli.Command {
		return &cli.Command{
			Run: func(ctx context.Context, args []string) error {
				fmt.Println("mig version 1.2.3")
				return nil
			},
		}
	})
	if err := app.Run(); err != nil {
		return err
	}

[App.DefaultCommand] selects a command when no explicit command is present.
Flags are defined in [Command.Bind]. Before argument parsing, matching
environment variables are applied to unchanged flags: names are lowercased
and underscores become hyphens, so DB_DSN maps to --db-dsn.

The -h and --help flags print help and return nil. Command lookup and flag
parsing errors print relevant usage before being returned. Errors produced by
[Command.Run] are returned without printing usage.

## Types

```go
// App is the cli entrypoint.
type App struct {
	// Name is the executable name shown in usage output.
	Name	string
	// DefaultCommand is selected when no explicit command is provided.
	DefaultCommand	string

	commands	map[string]CommandInfo
	commandOrder	[]string
}
```

```go
// Command and CommandInfo types for CLI command handling.
type (
	// FlagSet is here to prevent pflag leaking to imports.
	FlagSet	= pflag.FlagSet

	// Command is an individual command.
	Command	struct {
		// Name defaults to the name registered with App.AddCommand.
		Name	string
		// Title defaults to the title registered with App.AddCommand.
		Title	string
		// Default omits the command name from this command's usage line.
		Default	bool
		// Usage returns optional descriptive text printed before flag defaults.
		Usage	func() string
		// Bind defines this command's flags.
		Bind	func(*FlagSet)
		// Run executes the command with its remaining positional arguments.
		Run	func(context.Context, []string) error

		// Flags is populated by App.RunWithArgs with the flags defined by Bind.
		// HelpCommand uses it to print command flag defaults.
		Flags	*FlagSet
	}

	// CommandInfo is the constructor info for a command
	CommandInfo	struct {
		Name	string
		Title	string
		New	func() *Command
	}
)
```

## Function symbols

- `func NewApp (name string) *App`
- `func ParseWithFlagSet (fs *FlagSet, args []string) error`
- `func (*App) AddCommand (name,title string, constructor func() *Command)`
- `func (*App) FindCommand (commands []string, fallback string) (*Command, error)`
- `func (*App) HasCommand (name string) bool`
- `func (*App) Help ()`
- `func (*App) HelpCommand (fs *FlagSet, command *Command)`
- `func (*App) ParseCommands (args []string) []string`
- `func (*App) Run () error`
- `func (*App) RunWithArgs (args []string) error`

### NewApp

NewApp creates a new App instance.

```go
func NewApp (name string) *App
```

### ParseWithFlagSet

ParseWithFlagSet applies environment variables and parses args for a scoped
FlagSet. Environment names are lowercased and underscores become hyphens;
variables without an underscore or a matching flag are ignored. Argument
values take precedence over environment values.

```go
func ParseWithFlagSet (fs *FlagSet, args []string) error
```

### AddCommand

AddCommand adds a command to the app.

```go
func (*App) AddCommand (name,title string, constructor func() *Command)
```

### FindCommand

FindCommand finds a command for the app.

```go
func (*App) FindCommand (commands []string, fallback string) (*Command, error)
```

### HasCommand

HasCommand checks if a command exists in the app.

```go
func (*App) HasCommand (name string) bool
```

### Help

Help prints out registered commands for app.

```go
func (*App) Help ()
```

### HelpCommand

HelpCommand prints out help for a specific command.

```go
func (*App) HelpCommand (fs *FlagSet, command *Command)
```

### ParseCommands

ParseCommands cleans up args[], returning only commands.
If no commands are detected, DefaultCommand is returned.

```go
func (*App) ParseCommands (args []string) []string
```

### Run

Run passes os.Args without the command name to RunWithArgs().

```go
func (*App) Run () error
```

### RunWithArgs

RunWithArgs selects and executes a command with a context canceled by SIGINT
or SIGTERM. Help requests return nil; lookup and flag parsing errors print
usage, while errors returned by Command.Run do not.

```go
func (*App) RunWithArgs (args []string) error
```


