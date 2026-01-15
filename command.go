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

// Flag variable binding functions from spf13/pflag.
var (
	BoolVar        = pflag.BoolVar
	DurationVar    = pflag.DurationVar
	Int64Var       = pflag.Int64Var
	IntVar         = pflag.IntVar
	StringVar      = pflag.StringVar
	Uint64Var      = pflag.Uint64Var
	UintVar        = pflag.UintVar
	StringSliceVar = pflag.StringSliceVar

	BoolVarP        = pflag.BoolVarP
	DurationVarP    = pflag.DurationVarP
	Int64VarP       = pflag.Int64VarP
	IntVarP         = pflag.IntVarP
	StringVarP      = pflag.StringVarP
	Uint64VarP      = pflag.Uint64VarP
	UintVarP        = pflag.UintVarP
	StringSliceVarP = pflag.StringSliceVarP

	PrintDefaults = pflag.PrintDefaults
)

// Command and CommandInfo types for CLI command handling.
type (
	// FlagSet is here to prevent pflag leaking to imports.
	FlagSet = pflag.FlagSet

	// Command is an individual command
	Command struct {
		Name, Title string

		Bind func(*FlagSet)
		Run  func(context.Context, []string) error
	}

	// CommandInfo is the constructor info for a command
	CommandInfo struct {
		Name  string
		Title string
		New   func() *Command
	}
)

// ParseWithFlagSet parses flags and environment variables for a scoped FlagSet.
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
