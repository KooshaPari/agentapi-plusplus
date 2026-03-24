// Package cli_test provides tests for the CLI utilities.
package cli_test

import (
	"testing"

	"github.com/coder/agentapi/internal/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateRootCommand tests the root command creation.
func TestCreateRootCommand(t *testing.T) {
	cmd := cli.CreateRootCommand("1.0.0")
	require.NotNil(t, cmd)
	assert.Equal(t, "agentapi", cmd.Use)
	assert.Equal(t, "AgentAPI CLI", cmd.Short)
	assert.Equal(t, "1.0.0", cmd.Version)
}

// TestCreateRootCommand_CallHelp tests that calling the root command shows help.
func TestCreateRootCommand_CallHelp(t *testing.T) {
	cmd := cli.CreateRootCommand("1.0.0")
	require.NotNil(t, cmd)
	require.NotNil(t, cmd.RunE)

	// Calling RunE with no args should call Help()
	err := cmd.RunE(cmd, []string{})
	// Help() returns nil on success
	assert.NoError(t, err)
}

// TestNewCommandBuilder tests creating a new command builder.
func TestNewCommandBuilder(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	require.NotNil(t, builder)

	cmd := builder.Build()
	require.NotNil(t, cmd)
	assert.Equal(t, "test-cmd", cmd.Use)
}

// TestCommandBuilder_Short tests setting the short description.
func TestCommandBuilder_Short(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	builder.Short("This is a test command")

	cmd := builder.Build()
	assert.Equal(t, "This is a test command", cmd.Short)
}

// TestCommandBuilder_Long tests setting the long description.
func TestCommandBuilder_Long(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	longDesc := "This is a long description for the test command"
	builder.Long(longDesc)

	cmd := builder.Build()
	assert.Equal(t, longDesc, cmd.Long)
}

// TestCommandBuilder_RunE tests setting the RunE function.
func TestCommandBuilder_RunE(t *testing.T) {
	called := false
	runFunc := func(cmd *cobra.Command, args []string) error {
		called = true
		return nil
	}

	builder := cli.NewCommandBuilder("test-cmd")
	builder.RunE(runFunc)

	cmd := builder.Build()
	require.NotNil(t, cmd.RunE)
	err := cmd.RunE(cmd, []string{})
	assert.NoError(t, err)
	assert.True(t, called)
}

// TestCommandBuilder_AddStringFlag tests adding a string flag.
func TestCommandBuilder_AddStringFlag(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	builder.AddStringFlag("output", "o", "default.txt", "Output file")

	cmd := builder.Build()
	require.NotNil(t, cmd.Flags().Lookup("output"))
	flag := cmd.Flags().Lookup("output")
	assert.Equal(t, "default.txt", flag.DefValue)
	assert.Equal(t, "o", flag.Shorthand)
}

// TestCommandBuilder_AddBoolFlag tests adding a boolean flag.
func TestCommandBuilder_AddBoolFlag(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	builder.AddBoolFlag("verbose", "v", false, "Verbose output")

	cmd := builder.Build()
	require.NotNil(t, cmd.Flags().Lookup("verbose"))
	flag := cmd.Flags().Lookup("verbose")
	assert.Equal(t, "false", flag.DefValue)
	assert.Equal(t, "v", flag.Shorthand)
}

// TestCommandBuilder_AddIntFlag tests adding an integer flag.
func TestCommandBuilder_AddIntFlag(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	builder.AddIntFlag("count", "c", 10, "Count of items")

	cmd := builder.Build()
	require.NotNil(t, cmd.Flags().Lookup("count"))
	flag := cmd.Flags().Lookup("count")
	assert.Equal(t, "10", flag.DefValue)
	assert.Equal(t, "c", flag.Shorthand)
}

// TestCommandBuilder_AddSubcommand tests adding a subcommand.
func TestCommandBuilder_AddSubcommand(t *testing.T) {
	rootBuilder := cli.NewCommandBuilder("root")
	subCmd := cli.NewCommandBuilder("sub").Build()

	rootBuilder.AddSubcommand(subCmd)

	root := rootBuilder.Build()
	require.Equal(t, 1, len(root.Commands()))
	assert.Equal(t, "sub", root.Commands()[0].Use)
}

// TestCommandBuilder_Chainable tests that CommandBuilder methods are chainable.
func TestCommandBuilder_Chainable(t *testing.T) {
	builder := cli.NewCommandBuilder("test-cmd")
	cmd := builder.
		Short("Test command").
		Long("This is a test command").
		AddStringFlag("name", "n", "", "Name").
		AddBoolFlag("debug", "d", false, "Debug mode").
		AddIntFlag("threads", "t", 1, "Number of threads").
		Build()

	require.NotNil(t, cmd)
	assert.Equal(t, "test-cmd", cmd.Use)
	assert.Equal(t, "Test command", cmd.Short)
	assert.Equal(t, "This is a test command", cmd.Long)
	assert.NotNil(t, cmd.Flags().Lookup("name"))
	assert.NotNil(t, cmd.Flags().Lookup("debug"))
	assert.NotNil(t, cmd.Flags().Lookup("threads"))
}
