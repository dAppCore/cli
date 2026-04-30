package cli

import (
	"context"
	"sync"

	"time"

	"dappco.re/go"
)

func TestRun_Good_CancelledContext(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "test"})))

	// Register a long-running command that waits for context cancellation
	RegisterCommands(func(c *core.Core) {
		c.Command("wait", core.Command{
			Description: "Wait for context",
			Action: func(_ core.Options) core.Result {
				<-Context().Done()
				return core.Ok(nil)
			},
		})
	})

	// TODO: Run() test with context cancellation requires os.Args override.
	// Skipping for now — the underlying Cli.Run() is tested in core/go.
	_ = t
}

func TestRunWithTimeout_Good_ReturnsHelper(t *core.T) {
	resetGlobals(t)

	finished := make(chan struct{})
	var finishedOnce sync.Once
	core.RequireNoError(t, cliResultError(Init(Options{
		AppName: "test",
		Services: []core.Service{
			{
				Name: "slow-stop",
				OnStop: func() core.Result {
					time.Sleep(100 * time.Millisecond)
					finishedOnce.Do(func() {
						close(finished)
					})
					return core.Ok(nil)
				},
			},
		},
	})))

	start := time.Now()
	RunWithTimeout(20 * time.Millisecond)()
	core.AssertLess(t, time.Since(start), 80*time.Millisecond)

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("shutdown did not complete")
	}
}

func TestRun_Good_NilContext(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "test"})))

	// Run with nil context should not panic
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	err := cliResultError(Run(ctx))
	core.AssertError(t, err) // Should get context.Canceled
}
