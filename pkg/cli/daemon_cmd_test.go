package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddDaemonCommand_RegistersSubcommands(t *testing.T) {
	root := &Command{Use: "test"}

	AddDaemonCommand(root, DaemonCommandConfig{
		Name:       "daemon",
		PIDFile:    "/tmp/test-daemon.pid",
		HealthAddr: "127.0.0.1:0",
	})

	// Should have the daemon command
	daemonCmd, _, err := root.Find([]string{"daemon"})
	require.NoError(t, err)
	require.NotNil(t, daemonCmd)

	// Should have subcommands
	var subNames []string
	for _, sub := range daemonCmd.Commands() {
		subNames = append(subNames, sub.Name())
	}
	assert.Contains(t, subNames, "start")
	assert.Contains(t, subNames, "stop")
	assert.Contains(t, subNames, "status")
	assert.Contains(t, subNames, "run")
}

func TestDaemonCommandConfig_DefaultName(t *testing.T) {
	root := &Command{Use: "test"}

	AddDaemonCommand(root, DaemonCommandConfig{})

	// Should default to "daemon"
	daemonCmd, _, err := root.Find([]string{"daemon"})
	require.NoError(t, err)
	require.NotNil(t, daemonCmd)
}
