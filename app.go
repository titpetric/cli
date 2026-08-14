package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/pflag"
)

var errNoCommand = errors.New("no command found")

// App is the cli entrypoint.
type App struct {
	// Name is the executable name shown in usage output.
	Name string
	// DefaultCommand is selected when no explicit command is provided.
	DefaultCommand string

	commands     map[string]CommandInfo
	commandOrder []string
}

// NewApp creates a new App instance.
func NewApp(name string) *App {
	return &App{
		Name:         name,
		commands:     make(map[string]CommandInfo),
		commandOrder: []string{},
	}
}

// Run passes os.Args without the command name to RunWithArgs().
func (app *App) Run() error {
	return app.RunWithArgs(os.Args[1:])
}

// RunWithArgs selects and executes a command with a context canceled by SIGINT
// or SIGTERM. Help requests return nil; lookup and flag parsing errors print
// usage, while errors returned by Command.Run do not.
func (app *App) RunWithArgs(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	commands := app.ParseCommands(args)
	explicitCommand := app.hasExplicitCommand(args)
	command, err := app.FindCommand(commands, app.DefaultCommand)

	// Create a scoped FlagSet for this command
	name := app.Name
	if command != nil {
		name = command.Name
	}
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.Usage = func() {
		if command != nil && explicitCommand {
			app.HelpCommand(fs, command)
		} else {
			app.Help()
		}
	}

	if err != nil {
		if errors.Is(fs.Parse(args), pflag.ErrHelp) {
			return nil
		}
		app.Help()
		return err
	}

	// bind command specific flags
	if command.Bind != nil {
		command.Bind(fs)
	}

	// build a separate FlagSet with only command-specific flags
	command.Flags = pflag.NewFlagSet(command.Name+"-flags", pflag.ContinueOnError)
	fs.VisitAll(func(f *pflag.Flag) {
		command.Flags.AddFlag(f)
	})

	// parse flags and set from environment
	if err := ParseWithFlagSet(fs, args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return nil
		}
		app.HelpCommand(fs, command)
		return err
	}

	// Run command if defined
	if command.Run != nil {
		remainingArgs := fs.Args()

		if len(remainingArgs) > 0 && remainingArgs[0] == command.Name {
			remainingArgs = remainingArgs[1:]
		}

		err = command.Run(ctx, remainingArgs)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Missing Run() for command")
}

// Help prints out registered commands for app.
func (app *App) Help() {
	fmt.Println("Usage:", app.Name, "(command) [--flags]")
	fmt.Println("Available commands:")
	fmt.Println()

	maxLen := 0
	for _, name := range app.commandOrder {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}
	pad := "   "
	format := pad + "%-" + fmt.Sprintf("%d", maxLen+3) + "s %s\n"
	for _, name := range app.commandOrder {
		command := app.commands[name]
		fmt.Printf(format, command.Name, command.Title)
	}
	fmt.Println()
}

// HelpCommand prints out help for a specific command.
func (app *App) HelpCommand(fs *FlagSet, command *Command) {
	usage := app.Name
	if !command.Default {
		usage += " " + command.Name
	}
	usage += " [--flags]"
	fmt.Println("Usage:", usage)
	fmt.Println()

	if command.Usage != nil {
		if text := strings.TrimSpace(command.Usage()); text != "" {
			fmt.Println(text)
			fmt.Println()
		}
	}

	// Print command-specific flags only (excludes --help)
	flags := fs
	if command.Flags != nil {
		flags = command.Flags
	}
	if flags.HasFlags() {
		fmt.Println("Available options:")
		fmt.Println()
		flags.PrintDefaults()
		fmt.Println()
	}
}

// AddCommand adds a command to the app.
func (app *App) AddCommand(name, title string, constructor func() *Command) {
	info := CommandInfo{
		Name:  name,
		Title: title,
		New:   constructor,
	}
	app.commands[name] = info
	app.commandOrder = append(app.commandOrder, name)
}

// HasCommand checks if a command exists in the app.
func (app *App) HasCommand(name string) bool {
	_, ok := app.commands[name]
	return ok
}

// hasExplicitCommand checks if args contain an explicit command name (not a flag).
func (app *App) hasExplicitCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := args[0]
	return len(first) > 0 && first[0] != '-' && app.HasCommand(first)
}

// FindCommand finds a command for the app.
func (app *App) FindCommand(commands []string, fallback string) (*Command, error) {
	spawn := func(info CommandInfo) (*Command, error) {
		command := info.New()
		if command.Name == "" {
			command.Name = info.Name
		}
		if command.Title == "" {
			command.Title = info.Title
		}
		return command, nil
	}

	// This is just fully naive, we use the first command we find
	// but we could be smarter and have sub commands? Maybe one day.
	if len(commands) > 0 {
		commandName := commands[0]
		if info, ok := app.commands[commandName]; ok {
			return spawn(info)
		}
	}
	if info, ok := app.commands[fallback]; ok {
		return spawn(info)
	}
	if len(commands) > 0 {
		return nil, fmt.Errorf("unknown command: %q", commands[0])
	}
	return nil, fmt.Errorf("no command specified")
}

// ParseCommands cleans up args[], returning only commands.
// If no commands are detected, DefaultCommand is returned.
func (app *App) ParseCommands(args []string) []string {
	result := []string{}
	for _, v := range args {
		if len(v) > 0 && v[0:1] == "-" {
			break
		}
		result = append(result, v)
	}
	if len(result) == 0 {
		result = append(result, app.DefaultCommand)
		return result
	}
	return result
}
