package cli

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"dappco.re/go"
	"time"
)

func TestDaemon_StartStop(t *core.T) {
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
	core.RequireNoError(t, daemon.Start(context.Background()))
	defer func() {
		core.RequireNoError(t, daemon.Stop(context.Background()))
	}()

	rawPID, err := os.ReadFile(pidFile)
	core.RequireNoError(t, err)
	core.AssertEqual(t, strconv.Itoa(os.Getpid()), strings.TrimSpace(string(rawPID)))

	addr := daemon.HealthAddr()
	core.RequireNotEmpty(t, addr)

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get("http://" + addr + "/health")
	core.RequireNoError(t, err)
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	core.RequireNoError(t, err)
	core.AssertEqual(t, http.StatusOK, resp.StatusCode)
	core.AssertEqual(t, "ok\n", string(body))

	resp, err = client.Get("http://" + addr + "/ready")
	core.RequireNoError(t, err)
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	core.RequireNoError(t, err)
	core.AssertEqual(t, http.StatusServiceUnavailable, resp.StatusCode)
	core.AssertEqual(t, "unhealthy\n", string(body))

	ready = true

	resp, err = client.Get("http://" + addr + "/ready")
	core.RequireNoError(t, err)
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	core.RequireNoError(t, err)
	core.AssertEqual(t, http.StatusOK, resp.StatusCode)
	core.AssertEqual(t, "ok\n", string(body))
}

func TestDaemon_StopRemovesPIDFile(t *core.T) {
	tmp := t.TempDir()
	pidFile := filepath.Join(tmp, "daemon.pid")

	daemon := NewDaemon(DaemonOptions{PIDFile: pidFile})
	core.RequireNoError(t, daemon.Start(context.Background()))

	_, err := os.Stat(pidFile)
	core.RequireNoError(t, err)
	core.RequireNoError(t, daemon.Stop(context.Background()))

	_, err = os.Stat(pidFile)
	core.RequireTrue(t, err != nil, "RequireError")
	core.AssertTrue(t, os.IsNotExist(err))
}

func TestStopPIDFile_Good(t *core.T) {
	tmp := t.TempDir()
	pidFile := filepath.Join(tmp, "daemon.pid")
	core.RequireNoError(t, os.WriteFile(pidFile, []byte("1234\n"), 0o644))

	originalSignal := processSignal
	originalAlive := processAlive
	originalNow := processNow
	originalSleep := processSleep
	originalPoll := processPollInterval
	originalShutdownWait := processShutdownWait
	t.Cleanup(func() {
		processSignal = originalSignal
		processAlive = originalAlive
		processNow = originalNow
		processSleep = originalSleep
		processPollInterval = originalPoll
		processShutdownWait = originalShutdownWait
	})

	var mu sync.Mutex
	var signals []syscall.Signal
	processSignal = func(pid int, sig syscall.Signal) error {
		mu.Lock()
		signals = append(signals, sig)
		mu.Unlock()
		return nil
	}
	processAlive = func(pid int) bool {
		mu.Lock()
		defer mu.Unlock()
		if len(signals) == 0 {
			return true
		}
		return signals[len(signals)-1] != syscall.SIGTERM
	}
	processPollInterval = 0
	processShutdownWait = 0
	core.RequireNoError(t, StopPIDFile(pidFile, time.Second))

	mu.Lock()
	defer mu.Unlock()
	core.RequireTrue(t, core.DeepEqual([]syscall.Signal{syscall.SIGTERM}, signals), core.Sprintf("want=%#v got=%#v", []syscall.Signal{syscall.SIGTERM}, signals))

	_, err := os.Stat(pidFile)
	core.RequireTrue(t, err != nil, "RequireError")
	core.AssertTrue(t, os.IsNotExist(err))
}

func TestStopPIDFile_Bad_Escalates(t *core.T) {
	tmp := t.TempDir()
	pidFile := filepath.Join(tmp, "daemon.pid")
	core.RequireNoError(t, os.WriteFile(pidFile, []byte("4321\n"), 0o644))

	originalSignal := processSignal
	originalAlive := processAlive
	originalNow := processNow
	originalSleep := processSleep
	originalPoll := processPollInterval
	originalShutdownWait := processShutdownWait
	t.Cleanup(func() {
		processSignal = originalSignal
		processAlive = originalAlive
		processNow = originalNow
		processSleep = originalSleep
		processPollInterval = originalPoll
		processShutdownWait = originalShutdownWait
	})

	var mu sync.Mutex
	var signals []syscall.Signal
	current := time.Unix(0, 0)
	processNow = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return current
	}
	processSleep = func(d time.Duration) {
		mu.Lock()
		current = current.Add(d)
		mu.Unlock()
	}
	processSignal = func(pid int, sig syscall.Signal) error {
		mu.Lock()
		signals = append(signals, sig)
		mu.Unlock()
		return nil
	}
	processAlive = func(pid int) bool {
		mu.Lock()
		defer mu.Unlock()
		if len(signals) == 0 {
			return true
		}
		return signals[len(signals)-1] != syscall.SIGKILL
	}
	processPollInterval = 10 * time.Millisecond
	processShutdownWait = 0
	core.RequireNoError(t, StopPIDFile(pidFile, 15*time.Millisecond))

	mu.Lock()
	defer mu.Unlock()
	core.RequireTrue(t, core.DeepEqual([]syscall.Signal{syscall.SIGTERM, syscall.SIGKILL}, signals), core.Sprintf("want=%#v got=%#v", []syscall.Signal{syscall.SIGTERM, syscall.SIGKILL}, signals))
}
