// # CLI package
//
// This package contains the implementation for a minimal opinionated flags
// framework similar to spf13/cobra. It all centers around the `cli.Command` type
// but provides less functionality.
//
// To create a new CLI application:
//
// ```go
// app := cli.NewApp("mig")
// app.AddCommand("version", version.Name, version.New)
//
//	if err := app.Run(); err != nil {
//	        return err
//	}
//
// ```
//
// The `version.New` is a `func() *cli.Command`.
//
// The Command type defines Name and Title as strings, equivallent to cobra
// `Command.Use` (Name) and `Command.Long` (Title). There is no equivalent
// of `Command.Short`.
//
// The API choices are different, cobra's `AddCommand` took a command, and the
// command type was passed into Run().
//
// The cli package creates a `CommandInfo` with AddCommand, and then calls
// the constructor of the `*Command` type. The type must have Run filled, and
// can implement Bind(*FlagSet) to read in CLI flags.
//
// The Run function is context aware, supporting observability.
package cli
