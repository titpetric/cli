# Package cli

```go
import (
	"github.com/titpetric/cli"
}
```
# CLI package

This package contains the implementation for a minimal opinionated flags
framework similar to spf13/cobra. It all centers around the `cli.Command` type
but provides less functionality.

To create a new CLI application:

```go
app := cli.NewApp("mig")
app.AddCommand("version", version.Name, version.New)

	if err := app.Run(); err != nil {
	        return err
	}

```

The `version.New` is a `func() *cli.Command`.

The Command type defines Name and Title as strings, equivallent to cobra
`Command.Use` (Name) and `Command.Long` (Title). There is no equivalent
of `Command.Short`.

The API choices are different, cobra's `AddCommand` took a command, and the
command type was passed into Run().

The cli package creates a `CommandInfo` with AddCommand, and then calls
the constructor of the `*Command` type. The type must have Run filled, and
can implement Bind(*FlagSet) to read in CLI flags.

The Run function is context aware, supporting observability.

## Types

```go
// App is the cli entrypoint.
type App struct {
	Name	string

	commands	map[string]CommandInfo
	commandOrder	[]string
}
```

```go
// Command and CommandInfo types for CLI command handling.
type (
	// FlagSet is here to prevent pflag leaking to imports.
	FlagSet	= pflag.FlagSet

	// Command is an individual command
	Command	struct {
		Name, Title	string

		Bind	func(*FlagSet)
		Run	func(context.Context, []string) error
	}

	// CommandInfo is the constructor info for a command
	CommandInfo	struct {
		Name	string
		Title	string
		New	func() *Command
	}
)
```

## Vars

```go
// Flag variable binding functions from spf13/pflag.
var (
	BoolVar		= pflag.BoolVar
	DurationVar	= pflag.DurationVar
	Int64Var	= pflag.Int64Var
	IntVar		= pflag.IntVar
	StringVar	= pflag.StringVar
	Uint64Var	= pflag.Uint64Var
	UintVar		= pflag.UintVar
	StringSliceVar	= pflag.StringSliceVar

	BoolVarP	= pflag.BoolVarP
	DurationVarP	= pflag.DurationVarP
	Int64VarP	= pflag.Int64VarP
	IntVarP		= pflag.IntVarP
	StringVarP	= pflag.StringVarP
	Uint64VarP	= pflag.Uint64VarP
	UintVarP	= pflag.UintVarP
	StringSliceVarP	= pflag.StringSliceVarP

	PrintDefaults	= pflag.PrintDefaults
)
```

## Function symbols

- `func NewApp (name string) *App`
- `func ParseCommands (args []string) []string`
- `func ParseWithFlagSet (fs *FlagSet, args []string) error`
- `func (*App) AddCommand (name,title string, constructor func() *Command)`
- `func (*App) FindCommand (commands []string, fallback string) (*Command, error)`
- `func (*App) Help ()`
- `func (*App) HelpCommand (fs *FlagSet, command *Command)`
- `func (*App) Run () error`
- `func (*App) RunWithArgs (args []string) error`

### NewApp

NewApp creates a new App instance.

```go
func NewApp (name string) *App
```

### ParseCommands

ParseCommands cleans up args[], returning only commands.

It looks inside args[] up until the first parameter that starts with "-", a
flag parameter. We asume all the parameters before are command names.

Example: [a, b, -c, d, e] becomes [a, b].

```go
func ParseCommands (args []string) []string
```

### ParseWithFlagSet

ParseWithFlagSet parses flags and environment variables for a scoped FlagSet.

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

### Run

Run passes os.Args without the command name to RunWithArgs().

```go
func (*App) Run () error
```

### RunWithArgs

RunWithArgs is a cli entrypoint which sets up a cancellable context for the command.

```go
func (*App) RunWithArgs (args []string) error
```


