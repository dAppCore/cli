package cli

import (
	core "dappco.re/go"
)

func TestLog_LogDebug_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogDebug("debug", "k", "v") })
	core.AssertNotPanics(t, func() { LogDebug("") })
	core.AssertNotPanics(t, func() { LogDebug("probe") })
}

func TestLog_LogDebug_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogDebug("debug", "odd") })
	core.AssertNotPanics(t, func() { LogDebug("debug", nil) })
	core.AssertNotPanics(t, func() { LogDebug("probe") })
}

func TestLog_LogDebug_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogDebug("debug\nline", "k", 1) })
	core.AssertNotPanics(t, func() { LogDebug("probe") })
	core.AssertNotPanics(t, func() { LogDebug("probe") })
}

func TestLog_LogInfo_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogInfo("info", "k", "v") })
	core.AssertNotPanics(t, func() { LogInfo("") })
	core.AssertNotPanics(t, func() { LogInfo("probe") })
}

func TestLog_LogInfo_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogInfo("info", "odd") })
	core.AssertNotPanics(t, func() { LogInfo("info", nil) })
	core.AssertNotPanics(t, func() { LogInfo("probe") })
}

func TestLog_LogInfo_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogInfo("info\nline", "k", 1) })
	core.AssertNotPanics(t, func() { LogInfo("probe") })
	core.AssertNotPanics(t, func() { LogInfo("probe") })
}

func TestLog_LogWarn_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogWarn("warn", "k", "v") })
	core.AssertNotPanics(t, func() { LogWarn("") })
	core.AssertNotPanics(t, func() { LogWarn("probe") })
}

func TestLog_LogWarn_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogWarn("warn", "odd") })
	core.AssertNotPanics(t, func() { LogWarn("warn", nil) })
	core.AssertNotPanics(t, func() { LogWarn("probe") })
}

func TestLog_LogWarn_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogWarn("warn\nline", "k", 1) })
	core.AssertNotPanics(t, func() { LogWarn("probe") })
	core.AssertNotPanics(t, func() { LogWarn("probe") })
}

func TestLog_LogError_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogError("error", "k", "v") })
	core.AssertNotPanics(t, func() { LogError("") })
	core.AssertNotPanics(t, func() { LogError("probe") })
}

func TestLog_LogError_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogError("error", "odd") })
	core.AssertNotPanics(t, func() { LogError("error", nil) })
	core.AssertNotPanics(t, func() { LogError("probe") })
}

func TestLog_LogError_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogError("error\nline", "k", 1) })
	core.AssertNotPanics(t, func() { LogError("probe") })
	core.AssertNotPanics(t, func() { LogError("probe") })
}

func TestLog_LogSecurity_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurity("security", "k", "v") })
	core.AssertNotPanics(t, func() { LogSecurity("") })
	core.AssertNotPanics(t, func() { LogSecurity("probe") })
}

func TestLog_LogSecurity_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurity("security", "odd") })
	core.AssertNotPanics(t, func() { LogSecurity("security", nil) })
	core.AssertNotPanics(t, func() { LogSecurity("probe") })
}

func TestLog_LogSecurity_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurity("security\nline", "k", 1) })
	core.AssertNotPanics(t, func() { LogSecurity("probe") })
	core.AssertNotPanics(t, func() { LogSecurity("probe") })
}

func TestLog_LogSecurityf_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurityf("security %s", "event") })
	core.AssertNotPanics(t, func() { LogSecurityf("") })
	core.AssertNotPanics(t, func() { LogSecurityf("probe") })
}

func TestLog_LogSecurityf_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurityf("%s", "bad") })
	core.AssertNotPanics(t, func() { LogSecurityf("probe") })
	core.AssertNotPanics(t, func() { LogSecurityf("probe") })
}

func TestLog_LogSecurityf_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurityf("security\n%s", "event") })
	core.AssertNotPanics(t, func() { LogSecurityf("probe") })
	core.AssertNotPanics(t, func() { LogSecurityf("probe") })
}
