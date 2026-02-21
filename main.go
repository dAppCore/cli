package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	// Commands via self-registration (local to CLI)
	_ "forge.lthn.ai/core/cli/cmd/ai"
	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/dev"
	_ "forge.lthn.ai/core/cli/cmd/docs"
	_ "forge.lthn.ai/core/cli/cmd/doctor"
	_ "forge.lthn.ai/core/cli/cmd/gitcmd"
	_ "forge.lthn.ai/core/cli/cmd/go"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/lab"
	_ "forge.lthn.ai/core/cli/cmd/module"
	_ "forge.lthn.ai/core/cli/cmd/monitor"
	_ "forge.lthn.ai/core/cli/cmd/pkgcmd"
	_ "forge.lthn.ai/core/cli/cmd/plugin"
	_ "forge.lthn.ai/core/cli/cmd/qa"
	_ "forge.lthn.ai/core/cli/cmd/session"
	_ "forge.lthn.ai/core/cli/cmd/setup"
	_ "forge.lthn.ai/core/cli/cmd/updater"
	_ "forge.lthn.ai/core/cli/cmd/workspace"

	// Commands via self-registration (external repos)
	_ "forge.lthn.ai/core/go-ai/cmd/daemon"
	_ "forge.lthn.ai/core/go-ai/cmd/mcpcmd"
	_ "forge.lthn.ai/core/go-ai/cmd/security"
	_ "forge.lthn.ai/core/go-api/cmd/api"
	_ "forge.lthn.ai/core/go-crypt/cmd/crypt"
	_ "forge.lthn.ai/core/go-crypt/cmd/testcmd"
	_ "forge.lthn.ai/core/go-devops/build/buildcmd"
	_ "forge.lthn.ai/core/go-devops/cmd/deploy"
	_ "forge.lthn.ai/core/go-devops/cmd/prod"
	_ "forge.lthn.ai/core/go-devops/cmd/vm"
	_ "forge.lthn.ai/core/go-ml/cmd"
	_ "forge.lthn.ai/core/go-netops/cmd/unifi"
	_ "forge.lthn.ai/core/go-scm/cmd/collect"
	_ "forge.lthn.ai/core/go-scm/cmd/forge"

	// Variant repos (optional — comment out to exclude)
	// _ "forge.lthn.ai/core/php"
	// _ "forge.lthn.ai/core/ci"
)

func main() {
	cli.Main()
}
