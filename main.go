package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	// Commands via self-registration
	_ "forge.lthn.ai/core/cli/cmd/ai"
	_ "forge.lthn.ai/core/cli/cmd/api"
	_ "forge.lthn.ai/core/cli/cmd/collect"
	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/crypt"
	_ "forge.lthn.ai/core/cli/cmd/daemon"
	_ "forge.lthn.ai/core/cli/cmd/deploy"
	_ "forge.lthn.ai/core/cli/cmd/dev"
	_ "forge.lthn.ai/core/cli/cmd/docs"
	_ "forge.lthn.ai/core/cli/cmd/doctor"
	_ "forge.lthn.ai/core/cli/cmd/forge"
	_ "forge.lthn.ai/core/cli/cmd/gitcmd"
	_ "forge.lthn.ai/core/cli/cmd/go"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/lab"
	_ "forge.lthn.ai/core/cli/cmd/mcpcmd"
	_ "forge.lthn.ai/core/cli/cmd/ml"
	_ "forge.lthn.ai/core/cli/cmd/module"
	_ "forge.lthn.ai/core/cli/cmd/monitor"
	_ "forge.lthn.ai/core/cli/cmd/pkgcmd"
	_ "forge.lthn.ai/core/cli/cmd/plugin"
	_ "forge.lthn.ai/core/cli/cmd/prod"
	_ "forge.lthn.ai/core/cli/cmd/qa"
	_ "forge.lthn.ai/core/cli/cmd/rag"
	_ "forge.lthn.ai/core/cli/cmd/security"
	_ "forge.lthn.ai/core/cli/cmd/session"
	_ "forge.lthn.ai/core/cli/cmd/setup"
	_ "forge.lthn.ai/core/cli/cmd/test"
	_ "forge.lthn.ai/core/cli/cmd/unifi"
	_ "forge.lthn.ai/core/cli/cmd/updater"
	_ "forge.lthn.ai/core/cli/cmd/vm"
	_ "forge.lthn.ai/core/cli/cmd/workspace"
	_ "forge.lthn.ai/core/go-devops/build/buildcmd"

	// Variant repos (optional — comment out to exclude)
	// _ "forge.lthn.ai/core/php"
	// _ "forge.lthn.ai/core/ci"
)

func main() {
	cli.Main()
}
