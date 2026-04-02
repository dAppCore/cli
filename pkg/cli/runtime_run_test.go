package cli

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_Good_ReturnsCommandError(t *testing.T) {
	resetGlobals(t)

	require.NoError(t, Init(Options{AppName: "test"}))

	RootCmd().AddCommand(NewCommand("boom", "Boom", "", func(_ *Command, _ []string) error {
		return errors.New("boom")
	}))
	RootCmd().SetArgs([]string{"boom"})

	err := Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestRun_Good_CancelledContext(t *testing.T) {
	resetGlobals(t)

	require.NoError(t, Init(Options{AppName: "test"}))

	RootCmd().AddCommand(NewCommand("wait", "Wait", "", func(_ *Command, _ []string) error {
		<-Context().Done()
		return nil
	}))
	RootCmd().SetArgs([]string{"wait"})

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(25*time.Millisecond, cancel)

	err := Run(ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
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
