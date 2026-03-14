package cli_test

import (
	"context"
	"fmt"

	"github.com/titpetric/cli"
)

func ExampleApp() {
	app := cli.NewApp("mig")

	app.AddCommand("version", "Print version information", func() *cli.Command {
		var verbose bool

		return &cli.Command{
			Name:  "version",
			Title: "Print version information",
			Bind: func(fs *cli.FlagSet) {
				fs.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
			},
			Run: func(ctx context.Context, args []string) error {
				if verbose {
					fmt.Println("mig version 1.2.3 (commit abcdef)")
				} else {
					fmt.Println("mig version 1.2.3")
				}
				return nil
			},
		}
	})

	// In a real program you would call:
	//     _ = app.Run()
	//
	// For examples/tests, invoke RunWithArgs directly.
	_ = app.RunWithArgs([]string{"version"})

	// Output:
	// mig version 1.2.3
}
