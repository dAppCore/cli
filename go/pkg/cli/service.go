// SPDX-License-Identifier: EUPL-1.2

// Service registration for the cli package. This package is primarily a
// library consumed inside a binary's func main() — the Service surface
// is intentionally minimal (version + app metadata introspection) so
// long-running cli-hosted binaries can be queried over IPC under
// go-process supervision.
//
//	c, _ := core.New(
//	    core.WithName("cli", cli.NewService(cli.CliConfig{})),
//	)
//	r := c.Action("cli.version").Run(ctx, core.Options{})

package cli

import (
	"context"

	core "dappco.re/go"
)

// CliConfig configures the cli service. Empty config uses the package-
// global AppName + AppVersion + Build* values populated by ldflags.
//
// Usage example: `cfg := cli.CliConfig{}`
type CliConfig struct{}

// Service is the registerable handle for the cli package — embeds
// *core.ServiceRuntime[CliConfig] for typed options access. The cli
// framework itself is consumed via the package-level cli.Main entry
// point; this Service exposes introspection actions for diagnostics
// when a cli-hosted binary is supervised by go-process.
//
// Usage example: `svc := core.MustServiceFor[*cli.Service](c, "cli"); _ = svc`
type Service struct {
	*core.ServiceRuntime[CliConfig]
	registrations core.Once
}

// NewService returns a factory that produces a *Service ready for
// c.Service() registration.
//
// Usage example: `c, _ := core.New(core.WithName("cli", cli.NewService(cli.CliConfig{})))`
func NewService(config CliConfig) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Ok(&Service{
			ServiceRuntime: core.NewServiceRuntime(c, config),
		})
	}
}

// Register builds the cli service with default CliConfig and returns
// the service Result directly — the imperative-style alternative to
// NewService for consumers wiring services without WithName options.
//
// Usage example: `r := cli.Register(c); svc := r.Value.(*cli.Service)`
func Register(c *core.Core) core.Result {
	return NewService(CliConfig{})(c)
}

// OnStartup registers the cli action handlers on the attached Core.
// Implements core.Startable. Idempotent via core.Once.
//
// Usage example: `r := svc.OnStartup(ctx)`
func (s *Service) OnStartup(context.Context) core.Result {
	if s == nil {
		return core.Ok(nil)
	}
	s.registrations.Do(func() {
		c := s.Core()
		if c == nil {
			return
		}
		c.Action("cli.version", s.handleVersion)
		c.Action("cli.app_name", s.handleAppName)
		c.Action("cli.build_info", s.handleBuildInfo)
	})
	return core.Ok(nil)
}

// OnShutdown is a no-op — the cli package holds no closable resources.
// Implements core.Stoppable.
//
// Usage example: `r := svc.OnShutdown(ctx)`
func (s *Service) OnShutdown(context.Context) core.Result {
	return core.Ok(nil)
}

// handleVersion — `cli.version` action handler. Returns the SemVer
// 2.0.0 version string in r.Value.
//
// Usage example: `r := c.Action("cli.version").Run(ctx, core.Options{})`
func (s *Service) handleVersion(_ core.Context, _ core.Options) core.Result {
	return core.Ok(SemVer())
}

// handleAppName — `cli.app_name` action handler. Returns the
// configured AppName (set via WithAppName before Main) in r.Value.
//
// Usage example: `r := c.Action("cli.app_name").Run(ctx, core.Options{})`
func (s *Service) handleAppName(_ core.Context, _ core.Options) core.Result {
	return core.Ok(AppName)
}

// handleBuildInfo — `cli.build_info` action handler. Returns a
// map[string]string with version, commit, date, and pre-release fields
// in r.Value for diagnostics.
//
// Usage example: `r := c.Action("cli.build_info").Run(ctx, core.Options{})`
func (s *Service) handleBuildInfo(_ core.Context, _ core.Options) core.Result {
	return core.Ok(map[string]string{
		"version":     AppVersion,
		"commit":      BuildCommit,
		"date":        BuildDate,
		"pre_release": BuildPreRelease,
		"semver":      SemVer(),
	})
}
