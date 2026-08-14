// Package cli implements a minimal, opinionated command and flag framework
// built on spf13/pflag.
//
// An App registers command constructors with [App.AddCommand]. The selected
// constructor creates a [Command], whose Bind callback defines scoped flags and
// whose Run callback executes with a signal-aware context.
//
// A minimal application looks like this:
//
//	app := cli.NewApp("mig")
//	app.AddCommand("version", "Print version information", func() *cli.Command {
//		return &cli.Command{
//			Run: func(ctx context.Context, args []string) error {
//				fmt.Println("mig version 1.2.3")
//				return nil
//			},
//		}
//	})
//	if err := app.Run(); err != nil {
//		return err
//	}
//
// [App.DefaultCommand] selects a command when no explicit command is present.
// Flags are defined in [Command.Bind]. Before argument parsing, matching
// environment variables are applied to unchanged flags: names are lowercased
// and underscores become hyphens, so DB_DSN maps to --db-dsn.
//
// The -h and --help flags print help and return nil. Command lookup and flag
// parsing errors print relevant usage before being returned. Errors produced by
// [Command.Run] are returned without printing usage.
package cli
