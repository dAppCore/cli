package cli

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemon_StartStop(t *testing.T) {
	tmp := t.TempDir()
	pidFile := filepath.Join(tmp, "daemon.pid")
	ready := false

	daemon := NewDaemon(DaemonOptions{
		PIDFile:    pidFile,
		HealthAddr: "127.0.0.1:0",
		HealthCheck: func() bool {
			return true
		},
		ReadyCheck: func() bool {
			return ready
		},
	})

	require.NoError(t, daemon.Start(context.Background()))
	defer func() {
		require.NoError(t, daemon.Stop(context.Background()))
	}()

	rawPID, err := os.ReadFile(pidFile)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(os.Getpid()), strings.TrimSpace(string(rawPID)))

	addr := daemon.HealthAddr()
	require.NotEmpty(t, addr)

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get("http://" + addr + "/health")
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok\n", string(body))

	resp, err = client.Get("http://" + addr + "/ready")
	require.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, "unhealthy\n", string(body))

	ready = true

	resp, err = client.Get("http://" + addr + "/ready")
	require.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok\n", string(body))
}

func TestDaemon_StopRemovesPIDFile(t *testing.T) {
	tmp := t.TempDir()
	pidFile := filepath.Join(tmp, "daemon.pid")

	daemon := NewDaemon(DaemonOptions{PIDFile: pidFile})
	require.NoError(t, daemon.Start(context.Background()))

	_, err := os.Stat(pidFile)
	require.NoError(t, err)

	require.NoError(t, daemon.Stop(context.Background()))

	_, err = os.Stat(pidFile)
	require.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
