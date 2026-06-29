package cli

import (
	"context"

	core "dappco.re/go"
)

// serviceFor builds a *Service attached to a fresh Core, as a consumer
// would via core.WithName("cli", cli.NewService(...)).
func serviceFor(t *core.T) (*core.Core, *Service) {
	c := core.New()
	r := NewService(CliConfig{})(c)
	core.AssertTrue(t, r.OK)
	svc, ok := r.Value.(*Service)
	core.AssertTrue(t, ok)
	core.AssertNotNil(t, svc)
	return c, svc
}

func TestService_NewService_Good(t *core.T) {
	_, svc := serviceFor(t)

	core.AssertNotNil(t, svc.ServiceRuntime)
}

func TestService_NewService_Bad(t *core.T) {
	factory := NewService(CliConfig{})

	core.AssertNotNil(t, factory)
}

func TestService_NewService_Ugly(t *core.T) {
	c := core.New()
	first := NewService(CliConfig{})(c).Value.(*Service)
	second := NewService(CliConfig{})(c).Value.(*Service)

	core.AssertNotNil(t, first)
	core.AssertNotNil(t, second)
}

func TestService_Register_Good(t *core.T) {
	c := core.New()
	r := Register(c)

	core.AssertTrue(t, r.OK)
	_, ok := r.Value.(*Service)
	core.AssertTrue(t, ok)
}

func TestService_Register_Bad(t *core.T) {
	c := core.New()
	svc := Register(c).Value.(*Service)

	core.AssertNotNil(t, svc.ServiceRuntime)
}

func TestService_Register_Ugly(t *core.T) {
	c := core.New()
	first := Register(c).Value.(*Service)
	second := Register(c).Value.(*Service)

	core.AssertNotNil(t, first)
	core.AssertNotNil(t, second)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	c, svc := serviceFor(t)
	r := svc.OnStartup(context.Background())

	core.AssertTrue(t, r.OK)
	core.AssertTrue(t, c.Action("cli.version").Exists())
	core.AssertTrue(t, c.Action("cli.app_name").Exists())
	core.AssertTrue(t, c.Action("cli.build_info").Exists())
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	var svc *Service
	r := svc.OnStartup(context.Background())

	core.AssertTrue(t, r.OK)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	c, svc := serviceFor(t)
	first := svc.OnStartup(context.Background())
	second := svc.OnStartup(context.Background())

	core.AssertTrue(t, first.OK)
	core.AssertTrue(t, second.OK)
	core.AssertTrue(t, c.Action("cli.version").Exists())
}

func TestService_Service_OnShutdown_Good(t *core.T) {
	_, svc := serviceFor(t)
	r := svc.OnShutdown(context.Background())

	core.AssertTrue(t, r.OK)
}

func TestService_Service_OnShutdown_Bad(t *core.T) {
	var svc *Service
	r := svc.OnShutdown(context.Background())

	core.AssertTrue(t, r.OK)
}

func TestService_Service_OnShutdown_Ugly(t *core.T) {
	_, svc := serviceFor(t)
	first := svc.OnShutdown(context.Background())
	second := svc.OnShutdown(context.Background())

	core.AssertTrue(t, first.OK)
	core.AssertTrue(t, second.OK)
}

// handleVersion returns the SemVer string via the cli.version action.
func TestService_Service_handleVersion_Good(t *core.T) {
	c, svc := serviceFor(t)
	svc.OnStartup(context.Background())
	r := c.Action("cli.version").Run(c.Context(), core.Options{})

	core.AssertTrue(t, r.OK)
	core.AssertEqual(t, SemVer(), r.Value.(string))
}

func TestService_Service_handleVersion_Bad(t *core.T) {
	c := core.New()
	r := c.Action("cli.version").Run(c.Context(), core.Options{})

	core.AssertFalse(t, r.OK)
}

func TestService_Service_handleVersion_Ugly(t *core.T) {
	_, svc := serviceFor(t)
	r := svc.handleVersion(context.Background(), core.Options{})

	core.AssertTrue(t, r.OK)
	core.AssertNotEmpty(t, r.Value.(string))
}

// handleAppName returns the configured AppName via the cli.app_name action.
func TestService_Service_handleAppName_Good(t *core.T) {
	c, svc := serviceFor(t)
	svc.OnStartup(context.Background())
	r := c.Action("cli.app_name").Run(c.Context(), core.Options{})

	core.AssertTrue(t, r.OK)
	core.AssertEqual(t, AppName, r.Value.(string))
}

func TestService_Service_handleAppName_Bad(t *core.T) {
	c := core.New()
	r := c.Action("cli.app_name").Run(c.Context(), core.Options{})

	core.AssertFalse(t, r.OK)
}

func TestService_Service_handleAppName_Ugly(t *core.T) {
	_, svc := serviceFor(t)
	r := svc.handleAppName(context.Background(), core.Options{})

	core.AssertTrue(t, r.OK)
	core.AssertNotEmpty(t, r.Value.(string))
}

// handleBuildInfo returns a diagnostics map via the cli.build_info action.
func TestService_Service_handleBuildInfo_Good(t *core.T) {
	c, svc := serviceFor(t)
	svc.OnStartup(context.Background())
	r := c.Action("cli.build_info").Run(c.Context(), core.Options{})

	core.AssertTrue(t, r.OK)
	info := r.Value.(map[string]string)
	core.AssertEqual(t, AppVersion, info["version"])
	core.AssertEqual(t, SemVer(), info["semver"])
}

func TestService_Service_handleBuildInfo_Bad(t *core.T) {
	c := core.New()
	r := c.Action("cli.build_info").Run(c.Context(), core.Options{})

	core.AssertFalse(t, r.OK)
}

func TestService_Service_handleBuildInfo_Ugly(t *core.T) {
	_, svc := serviceFor(t)
	r := svc.handleBuildInfo(context.Background(), core.Options{})

	info := r.Value.(map[string]string)
	core.AssertContains(t, info, "commit")
	core.AssertContains(t, info, "date")
	core.AssertContains(t, info, "pre_release")
}
