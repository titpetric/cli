package cli

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// This package builds on spf13/pflag functionality.
//
// We omit exposing functions which return a pointer from the spf13/pflag
// public API, so we can encourage defining the flag values into structs.
//
// We also don't expose a lot of the spf13/pflag functionality here, since we
// expect that it will be wrapped in cli.App, and it makes little sense to use
// spf13/pflag "primitives" when we're creating a higher level abstraction.
//
// That being said, it's still possible to use the spf13/pflag API, but there
// should be little reason to do that.

// Command and CommandInfo types for CLI command handling.
type (
	// FlagSet is here to prevent pflag leaking to imports.
	FlagSet = pflag.FlagSet

	// Command is an individual command.
	Command struct {
		// Name defaults to the name registered with App.AddCommand.
		Name string
		// Title defaults to the title registered with App.AddCommand.
		Title string
		// Default omits the command name from this command's usage line.
		Default bool
		// Usage returns optional descriptive text printed before flag defaults.
		Usage func() string
		// Bind defines this command's flags.
		Bind func(*FlagSet)
		// Run executes the command with its remaining positional arguments.
		Run func(context.Context, []string) error

		// Flags is populated by App.RunWithArgs with the flags defined by Bind.
		// HelpCommand uses it to print command flag defaults.
		Flags *FlagSet
	}

	// CommandInfo is the constructor info for a command
	CommandInfo struct {
		Name  string
		Title string
		New   func() *Command
	}
)

// ParseWithFlagSet applies environment variables and parses args for a scoped
// FlagSet. Environment names are lowercased and underscores become hyphens;
// variables without an underscore or a matching flag are ignored. Argument
// values take precedence over environment values.
func ParseWithFlagSet(fs *FlagSet, args []string) error {
	// FlagSets are optional, but generally filled.
	if fs == nil {
		return nil
	}

	// parse environment variables and set on FlagSet
	for _, v := range os.Environ() {
		vals := strings.SplitN(v, "=", 2)
		if len(vals) != 2 {
			continue
		}

		flagName := vals[0]

		// only consider scoped envs
		if !strings.Contains(flagName, "_") {
			continue
		}

		flagName = strings.ToLower(flagName)
		flagName = strings.Replace(flagName, "_", "-", -1)

		// check if destination flag exists or modified
		fn := fs.Lookup(flagName)
		if fn == nil || fn.Changed {
			continue
		}
		if err := fn.Value.Set(vals[1]); err != nil {
			return err
		}
	}
	return fs.Parse(args)
}
