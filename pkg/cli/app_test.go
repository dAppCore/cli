package cli

import (
	core "dappco.re/go"
	"testing/fstest"
)

func cliResultError(r core.Result) error {
	if r.OK {
		return nil
	}
	if err, ok := r.Value.(error); ok {
		return err
	}
	return core.NewError(r.Error())
}

func cliPlainCLI(t *core.T) {
	t.Helper()
	originalTheme := currentTheme
	originalColor := ColorEnabled()
	UseASCII()
	SetColorEnabled(false)
	t.Cleanup(func() {
		currentTheme = originalTheme
		SetColorEnabled(originalColor)
		SetStdout(nil)
		SetStderr(nil)
		SetStdin(nil)
	})
}

func cliCaptureStdout(t *core.T, fn func()) string {
	t.Helper()
	out := core.NewBuilder()
	SetStdout(out)
	defer SetStdout(nil)
	fn()
	return out.String()
}

func cliCaptureStderr(t *core.T, fn func()) string {
	t.Helper()
	out := core.NewBuilder()
	SetStderr(out)
	defer SetStderr(nil)
	fn()
	return out.String()
}

func cliFakeCommands(t *core.T, scripts map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range scripts {
		path := core.Path(dir, name)
		data := []byte("#!/bin/sh\n" + body)
		core.RequireTrue(t, core.WriteFile(path, data, 0o755).OK)
	}
	t.Setenv("PATH", dir+string(core.PathListSeparator)+core.Getenv("PATH"))
	return dir
}

func cliRunSelf(t *core.T, envKey string) error {
	t.Helper()
	t.Setenv(envKey, "1")
	return cliResultError(runProcessOutput(core.Background(), core.Args()[0], "-test.run", "^"+t.Name()+"$"))
}

func resetGlobals(t *core.T) {
	t.Helper()
	if instance != nil {
		Shutdown()
		instance = nil
	}
	registeredCommandsMu.Lock()
	registeredCommands = nil
	registeredLocales = nil
	commandsAttached = false
	registeredCommandsMu.Unlock()
	SetStdin(nil)
	SetStdout(nil)
	SetStderr(nil)
	t.Cleanup(func() {
		if instance != nil {
			Shutdown()
			instance = nil
		}
		registeredCommandsMu.Lock()
		registeredCommands = nil
		registeredLocales = nil
		commandsAttached = false
		registeredCommandsMu.Unlock()
		SetStdin(nil)
		SetStdout(nil)
		SetStderr(nil)
	})
}

func TestApp_SemVer_Good(t *core.T) {
	oldVersion, oldPre, oldCommit, oldDate := AppVersion, BuildPreRelease, BuildCommit, BuildDate
	AppVersion, BuildPreRelease, BuildCommit, BuildDate = "1.2.3", "", "unknown", "unknown"
	defer func() { AppVersion, BuildPreRelease, BuildCommit, BuildDate = oldVersion, oldPre, oldCommit, oldDate }()

	core.AssertEqual(t, "1.2.3", SemVer())
}

func TestApp_SemVer_Bad(t *core.T) {
	oldVersion, oldPre, oldCommit, oldDate := AppVersion, BuildPreRelease, BuildCommit, BuildDate
	AppVersion, BuildPreRelease, BuildCommit, BuildDate = "0.0.0", "dev.1", "unknown", "unknown"
	defer func() { AppVersion, BuildPreRelease, BuildCommit, BuildDate = oldVersion, oldPre, oldCommit, oldDate }()

	core.AssertEqual(t, "0.0.0-dev.1", SemVer())
}

func TestApp_SemVer_Ugly(t *core.T) {
	oldVersion, oldPre, oldCommit, oldDate := AppVersion, BuildPreRelease, BuildCommit, BuildDate
	AppVersion, BuildPreRelease, BuildCommit, BuildDate = "1.0.0", "rc.1", "abc123", "20260428"
	defer func() { AppVersion, BuildPreRelease, BuildCommit, BuildDate = oldVersion, oldPre, oldCommit, oldDate }()

	core.AssertEqual(t, "1.0.0-rc.1+abc123.20260428", SemVer())
}

func TestApp_WithAppName_Good(t *core.T) {
	old := AppName
	defer func() { AppName = old }()
	WithAppName("codex")

	core.AssertEqual(t, "codex", AppName)
}

func TestApp_WithAppName_Bad(t *core.T) {
	old := AppName
	defer func() { AppName = old }()
	WithAppName("")

	core.AssertEqual(t, "", AppName)
}

func TestApp_WithAppName_Ugly(t *core.T) {
	old := AppName
	defer func() { AppName = old }()
	WithAppName("core dev")

	core.AssertEqual(t, "core dev", AppName)
}

func TestApp_WithLocales_Good(t *core.T) {
	fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
	src := WithLocales(fs, ".")

	core.AssertEqual(t, ".", src.Dir)
	core.AssertNotNil(t, src.FS)
}

func TestApp_WithLocales_Bad(t *core.T) {
	src := WithLocales(nil, ".")

	core.AssertEqual(t, ".", src.Dir)
	core.AssertNil(t, src.FS)
}

func TestApp_WithLocales_Ugly(t *core.T) {
	fs := fstest.MapFS{}
	src := WithLocales(fs, "")

	core.AssertEqual(t, "", src.Dir)
	core.AssertNotNil(t, src.FS)
}

func TestApp_Main_Good(t *core.T) {
	if core.Getenv("AX7_MAIN_GOOD") == "1" {
		Main()
		return
	}
	err := cliRunSelf(t, "AX7_MAIN_GOOD")
	core.AssertError(t, err)
}

func TestApp_Main_Bad(t *core.T) {
	if core.Getenv("AX7_MAIN_BAD") == "1" {
		Main(func(*core.Core) { panic("main setup failed") })
		return
	}
	err := cliRunSelf(t, "AX7_MAIN_BAD")
	core.AssertError(t, err)
}

func TestApp_Main_Ugly(t *core.T) {
	if core.Getenv("AX7_MAIN_UGLY") == "1" {
		Main(func(c *core.Core) { c.App().Name = "ugly" })
		return
	}
	err := cliRunSelf(t, "AX7_MAIN_UGLY")
	core.AssertError(t, err)
}

func TestApp_MainWithLocales_Good(t *core.T) {
	if core.Getenv("AX7_MAIN_WITH_LOCALES_GOOD") == "1" {
		fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
		MainWithLocales([]LocaleSource{WithLocales(fs, ".")})
		return
	}
	err := cliRunSelf(t, "AX7_MAIN_WITH_LOCALES_GOOD")
	core.AssertError(t, err)
}

func TestApp_MainWithLocales_Bad(t *core.T) {
	if core.Getenv("AX7_MAIN_WITH_LOCALES_BAD") == "1" {
		MainWithLocales([]LocaleSource{{}})
		return
	}
	err := cliRunSelf(t, "AX7_MAIN_WITH_LOCALES_BAD")
	core.AssertError(t, err)
}

func TestApp_MainWithLocales_Ugly(t *core.T) {
	if core.Getenv("AX7_MAIN_WITH_LOCALES_UGLY") == "1" {
		MainWithLocales(nil, func(c *core.Core) { c.App().Version = SemVer() })
		return
	}
	err := cliRunSelf(t, "AX7_MAIN_WITH_LOCALES_UGLY")
	core.AssertError(t, err)
}
