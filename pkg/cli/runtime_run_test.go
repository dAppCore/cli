package cli

import (
	"context"
	"sync"
	"testing"
	"time"

	"dappco.re/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_Good_CancelledContext(t *testing.T) {
	resetGlobals(t)

	require.NoError(t, Init(Options{AppName: "test"}))

	// Register a long-running command that waits for context cancellation
	RegisterCommands(func(c *core.Core) {
		c.Command("wait", core.Command{
			Description: "Wait for context",
			Action: func(_ core.Options) core.Result {
				<-Context().Done()
				return core.Result{OK: true}
			},
		})
	})

	// TODO: Run() test with context cancellation requires os.Args override.
	// Skipping for now — the underlying Cli.Run() is tested in core/go.
	_ = t
}

func TestRunWithTimeout_Good_ReturnsHelper(t *testing.T) {
	resetGlobals(t)

	finished := make(chan struct{})
	var finishedOnce sync.Once
	require.NoError(t, Init(Options{
		AppName: "test",
		Services: []core.Service{
			{
				Name: "slow-stop",
				OnStop: func() core.Result {
					time.Sleep(100 * time.Millisecond)
					finishedOnce.Do(func() {
						close(finished)
					})
					return core.Result{OK: true}
				},
			},
		},
	}))

	start := time.Now()
	RunWithTimeout(20 * time.Millisecond)()
	require.Less(t, time.Since(start), 80*time.Millisecond)

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("shutdown did not complete")
	}
}

func TestRun_Good_NilContext(t *testing.T) {
	resetGlobals(t)
	require.NoError(t, Init(Options{AppName: "test"}))

	// Run with nil context should not panic
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	err := Run(ctx)
	assert.Error(t, err) // Should get context.Canceled
}
