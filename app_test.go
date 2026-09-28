package cli_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/titpetric/cli"
	"github.com/titpetric/cli/tests/assert"
)

var NewApp = cli.NewApp

// TestApp_AddCommand tests that commands are registered and can be found.
func TestApp_AddCommand(t *testing.T) {
	app := NewApp("testapp")

	app.AddCommand("hello", "Say hello", func() *Command {
		cmd := &Command{
			Name: "hello",
			Run: func(ctx context.Context, args []string) error {
				return nil
			},
		}
		return cmd
	})

	cmd, err := app.FindCommand([]string{"hello"}, "")
	assert.NoError(t, err)
	assert.Equal(t, "hello", cmd.Name)
}

// TestApp_FindCommand_NotFound tests that findCommand returns an error for unknown commands.
func TestApp_FindCommand_NotFound(t *testing.T) {
	app := NewApp("testapp")

	_, err := app.FindCommand([]string{"nonexistent"}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}

// TestApp_FindCommand_NoCommand tests that findCommand returns an error when no command is specified.
func TestApp_FindCommand_NoCommand(t *testing.T) {
	app := NewApp("testapp")

	_, err := app.FindCommand([]string{}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no command specified")
}

// TestParseCommands tests that parseCommands extracts command names up to the first pflag.
func TestParseCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "single command",
			args:     []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "command with flags",
			args:     []string{"hello", "--name", "world"},
			expected: []string{"hello"},
		},
		{
			name:     "multiple commands",
			args:     []string{"sub", "cmd", "--flag"},
			expected: []string{"sub", "cmd"},
		},
		{
			name:     "flag at start",
			args:     []string{"--flag", "hello"},
			expected: []string{"run"},
		},
		{
			name:     "empty args",
			args:     []string{},
			expected: []string{"run"},
		},
	}

	app := NewApp("testapp")
	app.DefaultCommand = "run"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.ParseCommands(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestApp_RunWithArgs_Help tests that --help and -h show help without error.
func TestApp_RunWithArgs_Help(t *testing.T) {
	newApp := func() *cli.App {
		app := NewApp("testapp")
		app.AddCommand("test", "Test command", func() *Command {
			return &Command{
				Bind: func(fs *cli.FlagSet) {
					var msg string
					fs.StringVar(&msg, "msg", "", "a message")
				},
				Run: func(ctx context.Context, args []string) error {
					t.Fatal("command should not execute when help is passed")
					return nil
				},
			}
		})
		app.DefaultCommand = "test"
		return app
	}

	for _, flag := range []string{"--help", "-h"} {
		t.Run("app "+flag, func(t *testing.T) {
			err := newApp().RunWithArgs([]string{flag})
			assert.NoError(t, err)
		})
		t.Run("command "+flag, func(t *testing.T) {
			err := newApp().RunWithArgs([]string{"test", flag})
			assert.NoError(t, err)
		})
	}
}

// TestApp_RunWithArgs_Integration tests the full command execution flow.
func TestApp_RunWithArgs_Integration(t *testing.T) {
	app := NewApp("testapp")
	executed := false

	app.AddCommand("test", "Test command", func() *Command {
		var msg string

		cmd := &Command{
			Name: "test",
			Bind: func(fs *cli.FlagSet) {
				fs.StringVar(&msg, "msg", "default", "")
			},
			Run: func(ctx context.Context, args []string) error {
				executed = true
				assert.Equal(t, "hello", msg)
				return nil
			},
		}

		return cmd
	})

	// Simulate command line with flag
	err := app.RunWithArgs([]string{"test", "--msg", "hello"})
	assert.NoError(t, err)
	assert.True(t, executed)
}

// TestApp_RunWithArgs_CommandErrorDoesNotPrintUsage ensures runtime failures
// are not presented as command-line usage errors.
func TestApp_RunWithArgs_CommandErrorDoesNotPrintUsage(t *testing.T) {
	app := NewApp("testapp")
	wantErr := errors.New("command failed")
	app.AddCommand("test", "Test command", func() *Command {
		return &Command{
			Run: func(ctx context.Context, args []string) error {
				return wantErr
			},
		}
	})

	var stdout, stderr bytes.Buffer
	app.Stdout, app.Stderr = &stdout, &stderr

	assert.ErrorIs(t, app.RunWithArgs([]string{"test"}), wantErr)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}

// TestApp_RunWithArgs_UsageStreams tests that requested help goes to Stdout
// while the usage printed for a rejected command or flag goes to Stderr.
func TestApp_RunWithArgs_UsageStreams(t *testing.T) {
	newApp := func(stdout, stderr *bytes.Buffer) *cli.App {
		app := NewApp("testapp")
		app.Stdout, app.Stderr = stdout, stderr
		app.AddCommand("test", "Test command", func() *Command {
			return &Command{
				Bind: func(fs *cli.FlagSet) {
					var msg string
					fs.StringVar(&msg, "msg", "", "a message")
				},
				Run: func(ctx context.Context, args []string) error {
					return nil
				},
			}
		})
		return app
	}

	t.Run("help to stdout", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		assert.NoError(t, newApp(&stdout, &stderr).RunWithArgs([]string{"test", "--help"}))
		assert.Contains(t, stdout.String(), "Usage: testapp test [--flags]")
		assert.Contains(t, stdout.String(), "--msg")
		assert.Empty(t, stderr.String())
	})

	t.Run("unknown flag to stderr", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := newApp(&stdout, &stderr).RunWithArgs([]string{"test", "--nope"})
		assert.Error(t, err)
		assert.Contains(t, stderr.String(), "Usage: testapp test [--flags]")
		assert.Empty(t, stdout.String())
	})

	t.Run("unknown command to stderr", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := newApp(&stdout, &stderr).RunWithArgs([]string{"nope"})
		assert.Error(t, err)
		assert.Contains(t, stderr.String(), "Usage: testapp (command) [--flags]")
		assert.Empty(t, stdout.String())
	})
}

// TestApp_Run tests that Run dispatches os.Args without the program name.
func TestApp_Run(t *testing.T) {
	app := NewApp("testapp")
	var got []string

	app.AddCommand("test", "Test command", func() *Command {
		return &Command{
			Run: func(ctx context.Context, args []string) error {
				got = args
				return nil
			},
		}
	})

	args := os.Args
	os.Args = []string{"testapp", "test", "alpha"}
	t.Cleanup(func() {
		os.Args = args
	})

	assert.NoError(t, app.Run())
	assert.Equal(t, []string{"alpha"}, got)
}

// TestApp_Help tests that Help lists every registered command with its title.
func TestApp_Help(t *testing.T) {
	app := NewApp("testapp")
	app.AddCommand("first", "The first command", func() *Command { return &Command{} })
	app.AddCommand("second", "The second command", func() *Command { return &Command{} })

	var output bytes.Buffer
	app.Stdout = &output
	app.Help()

	assert.Contains(t, output.String(), "Usage: testapp (command) [--flags]")
	assert.Contains(t, output.String(), "first     The first command")
	assert.Contains(t, output.String(), "second    The second command")
}

// TestApp_HelpCommand tests that HelpCommand prints the usage line, the
// command description and the flags bound by the command.
func TestApp_HelpCommand(t *testing.T) {
	app := NewApp("testapp")
	var msg string

	command := &Command{
		Name:  "test",
		Usage: func() string { return "  Does a thing.  " },
		Bind: func(fs *cli.FlagSet) {
			fs.StringVar(&msg, "msg", "", "a message")
		},
	}
	fs := cli.NewFlagSet(command.Name, cli.ContinueOnError)
	command.Bind(fs)

	var output bytes.Buffer
	app.Stdout = &output
	app.HelpCommand(fs, command)

	assert.Contains(t, output.String(), "Usage: testapp test [--flags]")
	assert.Contains(t, output.String(), "Does a thing.")
	assert.Contains(t, output.String(), "Available options:")
	assert.Contains(t, output.String(), "--msg")
}

// TestApp_HelpCommand_Default tests that a default command is not named in
// its own usage line.
func TestApp_HelpCommand_Default(t *testing.T) {
	app := NewApp("testapp")
	command := &Command{Name: "test", Default: true}
	fs := cli.NewFlagSet(command.Name, cli.ContinueOnError)

	var output bytes.Buffer
	app.Stdout = &output
	app.HelpCommand(fs, command)

	assert.Contains(t, output.String(), "Usage: testapp [--flags]")
	assert.NotContains(t, output.String(), "Available options:")
}

// TestApp_HasCommand tests that HasCommand reports registered commands only.
func TestApp_HasCommand(t *testing.T) {
	app := NewApp("testapp")
	assert.False(t, app.HasCommand("test"))

	app.AddCommand("test", "Test command", func() *Command { return &Command{} })
	assert.True(t, app.HasCommand("test"))
	assert.False(t, app.HasCommand("other"))
}
