# CLI Reference

This document provides a comprehensive reference for all available commands in the Core CLI.

## Overview

The Core CLI is a unified development tool that provides commands for:
- AI/ML operations and agent management
- Multi-repository development workflows
- Infrastructure and deployment management
- Security monitoring and analysis
- Data collection and RAG systems
- Cryptographic utilities
- Testing and quality assurance

## Quick Command Index

| Command | Category | Description |
|---------|----------|-------------|
| [ai](#ai) | AI/ML | Agentic task management and agent configuration |
| [collect](#collect) | Data | Collect data from external sources |
| [config](#config) | Config | Configuration management |
| [crypt](#crypt) | Security | Cryptographic utilities |
| [daemon](#daemon) | Service | Background service management |
| [deploy](#deploy) | DevOps | Coolify deployment management |
| [dev](#dev) | Development | Multi-repo development workflow |
| [docs](#docs) | Documentation | Documentation management |
| [doctor](#doctor) | Utilities | Environment verification |
| [forge](#forge) | Git | Forgejo instance management |
| [git](#git) | Git | Root-level git workflow commands |
| [gitea](#gitea) | Git | Gitea instance management |
| [go](#go) | Development | Go development commands |
| [help](#help) | Utilities | Help documentation |
| [lab](#lab) | Monitoring | Homelab monitoring dashboard |
| [mcp](#mcp) | AI/ML | Model Context Protocol server |
| [ml](#ml) | AI/ML | ML inference and training pipeline |
| [monitor](#monitor) | Security | Security monitoring |
| [plugin](#plugin) | Utilities | Plugin management |
| [prod](#prod) | DevOps | Production infrastructure management |
| [qa](#qa) | Development | Quality assurance workflow |
| [rag](#rag) | AI/ML | RAG system for embeddings and semantic search |
| [security](#security) | Security | Security management |
| [session](#session) | Utilities | Session recording and replay |
| [setup](#setup) | Utilities | Workspace setup and bootstrap |
| [test](#test) | Development | Test running |
| [unifi](#unifi) | Networking | UniFi network management |
| [update](#update) | Utilities | CLI self-update |
| [vm](#vm) | DevOps | LinuxKit VM management |
| [workspace](#workspace) | Config | Workspace configuration |

## Detailed Command Reference

### ai

**Agentic task management and agent configuration**

Manage AI agents, tasks, and integrations with external services.

**Subcommands:**

- `agent` - Manage AgentCI dispatch targets
  - `add` - Add a new dispatch target
  - `list` - List all configured targets
  - `status` - Show target status
  - `logs` - View target logs
  - `setup` - Initial target setup
  - `remove` - Remove a target
- `tasks` - Agentic task management
  - `list` - List all tasks
  - `view` - View task details
  - `update` - Update task status
  - `complete` - Mark task as complete
- `git` - Git integration for tasks
  - `commit` - Create commits for tasks
  - `PR` - Create pull requests
- `metrics` - Task metrics tracking
- `commands` - AI command management
- `updates` - Task updates
- `ratelimits` - Rate limiting configuration
- `dispatch` - Task dispatch management

**Examples:**

```bash
# List all AI tasks
core ai tasks list

# View task details
core ai tasks view TASK-123

# Add dispatch target
core ai agent add --name my-agent --url https://api.example.com
```

---

### collect

**Collect data from various sources**

Data collection tools for gathering information from external sources like GitHub, forums, and research papers.

**Persistent Flags:**

- `--output`, `-o` - Output directory
- `--verbose`, `-v` - Verbose output
- `--dry-run` - Show what would happen without executing

**Subcommands:**

- `github` - Collect data from GitHub repositories
- `bitcointalk` - Collect from BitcoinTalk forum
- `market` - Collect market data
- `papers` - Collect research papers
- `excavate` - Deep excavation of sources
- `process` - Process collected data
- `dispatch` - Dispatch collection jobs

**Examples:**

```bash
# Collect GitHub data
core collect github --org myorg --output ./data

# Collect research papers
core collect papers --query "machine learning" --output ./papers
```

---

### config

**Manage configuration**

Manage Core CLI configuration values stored in `~/.core/config.yaml`.

**Subcommands:**

- `get` - Get a config value
- `set` - Set a config value
- `list` - List all config values
- `path` - Show config file path

**Examples:**

```bash
# Get a config value
core config get dev.editor

# Set a config value
core config set dev.editor vim

# List all config
core config list

# Show config file location
core config path
```

**See also:** [Configuration Reference](configuration.md)

---

### crypt

**Cryptographic utilities**

Encrypt, decrypt, hash, and checksum files and data.

**Subcommands:**

- `hash` - Hash files or data
- `encrypt` - Encrypt files or data
- `keygen` - Generate cryptographic keys
- `checksum` - Generate and verify checksums

**Examples:**

```bash
# Hash a file
core crypt hash --file myfile.txt

# Generate a PGP key
core crypt keygen --name "John Doe" --email john@example.com

# Encrypt a file
core crypt encrypt --file secret.txt --recipient john@example.com

# Generate checksums
core crypt checksum --file release.tar.gz
```

---

### daemon

**Start the core daemon**

Start Core daemon for long-running services like MCP.

**Flags:**

- `--mcp-transport`, `-t` - MCP transport (stdio, tcp, socket)
- `--mcp-addr`, `-a` - Listen address for TCP/socket
- `--health-addr` - Health check endpoint address
- `--pid-file` - PID file path

**Environment Variables:**

- `CORE_MCP_TRANSPORT` - MCP transport type
- `CORE_MCP_ADDR` - MCP listen address
- `CORE_HEALTH_ADDR` - Health check address
- `CORE_PID_FILE` - PID file location

**Examples:**

```bash
# Start daemon with MCP on stdio
core daemon

# Start daemon with TCP MCP server
core daemon --mcp-transport tcp --mcp-addr :9100

# Start with health endpoint
core daemon --health-addr :8080
```

---

### deploy

**Deployment infrastructure management**

Manage Coolify deployments and infrastructure.

**Flags:**

- `--url` - Coolify instance URL
- `--token` - API token
- `--json` - JSON output format

**Subcommands:**

- `servers` - List Coolify servers
- `projects` - List projects
- `apps` - List applications
- `databases` / `dbs` / `db` - List databases
- `services` - List services
- `team` - Show team info
- `call` - Call arbitrary Coolify API operations

**Examples:**

```bash
# List all servers
core deploy servers

# List applications
core deploy apps

# Call custom API endpoint
core deploy call --method GET --path /api/v1/status
```

---

### dev

**Development workflow commands**

Multi-repository development workflows with health checks, commits, and CI monitoring.

**Subcommands:**

- `work` - Combined status, commit, and push workflow
- `health` - Quick health check across repos
- `commit` - Claude-assisted commit message generation
- `push` - Push repos with unpushed commits
- `pull` - Pull repos behind remote
- `issues` - List open issues
- `reviews` - List PRs needing review
- `ci` - Check GitHub Actions CI status
- `impact` - Analyze dependency impact
- `workflow` - CI/workflow management
- `api` - API synchronization
- `vm` - VM management
  - `install`, `boot`, `stop`, `status`, `shell`, `serve`, `test`, `claude`, `update`
- `file-sync` - Sync files across repos (safe for AI agents)
- `apply` - Run commands across repos (safe for AI agents)

**Examples:**

```bash
# Full workflow: status, commit, push
core dev work

# Health check all repos
core dev health

# Claude-assisted commit
core dev commit

# List open issues
core dev issues

# Check CI status
core dev ci
```

---

### docs

**Documentation management**

Manage project documentation.

**Subcommands:**

- `sync` - Sync documentation
- `list` - List documentation
- `commands` - Command documentation
- `scan` - Scan for documentation

**Examples:**

```bash
# Sync documentation
core docs sync

# List available docs
core docs list
```

---

### doctor

**Check development environment**

Verify required tools and dependencies are installed and configured.

**Flags:**

- `--verbose` - Show detailed check results

**Checks:**

- Required tools validation (Go, Git, Task, etc.)
- Optional tools validation
- GitHub SSH access
- GitHub CLI authentication
- Workspace validation

**Examples:**

```bash
# Quick environment check
core doctor

# Verbose check with details
core doctor --verbose
```

---

### forge

**Forgejo instance management**

Manage Forgejo repositories, issues, and pull requests.

**Subcommands:**

- `config` - Configure Forgejo connection
- `status` - Show instance status and version
- `repos` - List repositories
- `issues` - List and create issues
- `prs` - List pull requests
- `migrate` - Migrate repos from external services
- `sync` - Sync GitHub repos to upstream branches
- `orgs` - List organizations
- `labels` - List and create labels

**Examples:**

```bash
# Configure Forgejo
core forge config --url https://forge.example.com --token YOUR_TOKEN

# List repositories
core forge repos

# List issues
core forge issues

# Create an issue
core forge issues create --title "Bug fix" --body "Description"
```

---

### git

**Root-level git workflow commands**

Git operations wrapper (imports commands from dev package).

**Subcommands:**

- `health` - Show repo status
- `commit` - Claude-assisted commits
- `push` - Push repositories
- `pull` - Pull repositories
- `work` - Combined workflow
- `file-sync` - Safe file synchronization
- `apply` - Safe command execution

**Examples:**

```bash
# Check repo health
core git health

# Claude-assisted commit
core git commit
```

---

### gitea

**Gitea instance management**

Manage Gitea repositories and integrations.

**Subcommands:**

- `config` - Configure Gitea connection
- `repos` - List repositories
- `issues` - List and create issues
- `prs` - List pull requests
- `mirror` - Create GitHub-to-Gitea mirrors
- `sync` - Sync GitHub repos to upstream branches

**Examples:**

```bash
# Configure Gitea
core gitea config --url https://gitea.example.com

# Create GitHub mirror
core gitea mirror --repo owner/repo
```

---

### go

**Go development utilities**

Go-specific development tools for testing, coverage, and quality assurance.

**Subcommands:**

- `qa` - Quality assurance checks
- `test` - Run tests
- `cov` - Code coverage analysis
- `fmt` - Code formatting
- `lint` - Linting
- `install` - Install Go tools
- `mod` - Module management
- `work` - Workspace management
- `fuzz` - Fuzz testing

**Examples:**

```bash
# Run QA checks
core go qa

# Run tests with coverage
core go cov

# Format code
core go fmt

# Run specific test
core go test --run TestMyFunction
```

---

### help

**Display help documentation**

Access built-in help documentation.

**Flags:**

- `--search`, `-s` - Search help topics

**Examples:**

```bash
# Show help
core help

# Search help topics
core help --search "configuration"
```

---

### lab

**Homelab monitoring dashboard**

Real-time monitoring of machines, training runs, models, and services.

**Subcommands:**

- `serve` - Start lab dashboard web server
  - `--bind` - HTTP listen address

**Examples:**

```bash
# Start dashboard on default port
core lab serve

# Start on custom port
core lab serve --bind :8080
```

---

### mcp

**Model Context Protocol server**

MCP server providing file operations, RAG, and metrics tools for AI assistants.

**Subcommands:**

- `serve` - Start the MCP server
  - `--workspace`, `-w` - Restrict file operations to directory

**Environment Variables:**

- `MCP_ADDR` - TCP address (e.g., `localhost:9999`). If not set, uses stdio.

**Transport:**

- **stdio** (default) - For Claude Code integration
- **TCP** - Set via `MCP_ADDR` environment variable

**Examples:**

```bash
# Start with stdio (for Claude Code)
core mcp serve

# Start with workspace restriction
core mcp serve --workspace /path/to/project

# Start TCP server
MCP_ADDR=localhost:9999 core mcp serve
```

**See also:** [MCP Integration](mcp-integration.md)

---

### ml

**ML inference, scoring, and training pipeline**

Commands for model scoring, probe evaluation, data export, and format conversion.

**Persistent Flags:**

- `--api-url` - OpenAI-compatible API URL
- `--judge-url` - Judge model API URL (Ollama)
- `--judge-model` - Judge model name
- `--influx` - InfluxDB URL
- `--influx-db` - InfluxDB database
- `--db` - DuckDB database path
- `--model` - Model name for API

**Subcommands:**

- `score` - Score responses with heuristic and LLM judges
- `probe` - Run capability and content probes against a model
- `export` - Export golden set to training formats
- `expand` - Generate expansion responses
- `status` - Show training and generation progress
- `gguf` - Convert MLX LoRA adapter to GGUF format
- `convert` - Convert MLX LoRA adapter to PEFT format
- `agent` - Run the scoring agent daemon
- `worker` - Run a distributed worker node
- `serve` - Start OpenAI-compatible inference server
- `inventory` - Show DuckDB table inventory
- `query` - Run ad-hoc SQL against DuckDB
- `metrics` - Push golden set stats to InfluxDB
- `ingest` - Ingest benchmark scores to InfluxDB
- `normalize` - Deduplicate seeds into expansion prompts
- `seed-influx` - Migrate golden set from DuckDB to InfluxDB
- `consolidate` - Pull and merge response JSONL files
- `import-all` - Import all LEM data into DuckDB
- `approve` - Filter scored expansions into training JSONL
- `publish` - Upload Parquet dataset to HuggingFace Hub
- `coverage` - Analyze seed coverage by region and domain

**Examples:**

```bash
# Score model responses
core ml score --model gpt-4

# Run capability probes
core ml probe --api-url http://localhost:8000

# Export training data
core ml export --format jsonl
```

---

### monitor

**Security finding aggregation**

Aggregate GitHub code scanning, Dependabot, and secret scanning alerts.

**Flags:**

- `--repo` - Specific repository
- `--all` - All repositories
- `--severity` - Filter by severity
- `--json` - JSON output

**Examples:**

```bash
# Monitor all repos
core monitor --all

# Monitor specific repo with high severity
core monitor --repo owner/repo --severity high
```

---

### plugin

**Manage plugins**

Install and manage CLI plugins.

**Subcommands:**

- `install` - Install a plugin from GitHub
- `list` - List installed plugins
- `info` - Show detailed plugin information
- `update` - Update a plugin or all plugins
- `remove` - Remove an installed plugin

**Examples:**

```bash
# Install plugin
core plugin install owner/repo

# List installed
core plugin list

# Update all
core plugin update
```

---

### prod

**Production infrastructure management**

Manage Host UK production infrastructure.

**Flags:**

- `--config` - Path to infra.yaml

**Subcommands:**

- `status` - Show infrastructure health
- `setup` - Phase 1: discover topology, create LB, configure DNS
- `dns` - Manage DNS records via CloudNS
- `lb` - Manage Hetzner load balancer
- `ssh` - SSH into a production host

**Examples:**

```bash
# Check infrastructure status
core prod status

# Setup infrastructure
core prod setup

# SSH to host
core prod ssh web1
```

---

### qa

**Quality assurance workflow**

Verify work (CI status, reviews, issues) - complements 'dev' which is about doing work.

**Subcommands:**

- `watch` - Monitor GitHub Actions after push
- `review` - PR review status with actionable steps
- `health` - Aggregate CI health across repos
- `issues` - Intelligent issue triage
- `docblock` - Documentation block validation

**Examples:**

```bash
# Watch CI after push
core qa watch

# Check PR review status
core qa review

# Health check across repos
core qa health
```

---

### rag

**RAG system for embeddings and semantic search**

Retrieval-Augmented Generation system using Qdrant and Ollama.

**Persistent Flags:**

- `--qdrant-host` - Qdrant host (default: localhost)
- `--qdrant-port` - Qdrant port (default: 6334)
- `--ollama-host` - Ollama host (default: localhost)
- `--ollama-port` - Ollama port (default: 11434)
- `--model` - Embedding model (default: nomic-embed-text)
- `--verbose`, `-v` - Verbose output

**Subcommands:**

- `ingest` - Ingest documents into collections
- `query` - Query collections
- `collections` - Manage collections (list, stats, delete)

**Examples:**

```bash
# Ingest documents
core rag ingest --collection docs --path ./documentation

# Query collection
core rag query --collection docs --query "how to install"

# List collections
core rag collections list
```

---

### security

**Security management**

Security operations and scanning.

**Subcommands:**

- `jobs` - Security job management
- `secrets` - Secrets management
- `alerts` - Security alerts
- `scan` - Security scanning

**Examples:**

```bash
# List security alerts
core security alerts

# Run security scan
core security scan
```

---

### session

**Session recording and replay**

Record and replay CLI sessions with timeline generation.

**Subcommands:**

- `list` - List recent sessions
- `replay` - Generate HTML timeline (and optional MP4) from session
  - `--mp4` - Generate MP4 video
  - `--output` - Output directory
- `search` - Search across session transcripts

**Examples:**

```bash
# List sessions
core session list

# Replay session as HTML
core session replay SESSION_ID

# Replay with MP4
core session replay SESSION_ID --mp4
```

---

### setup

**Workspace setup and bootstrap**

Set up development workspace from registry.

**Flags:**

- `--registry` - Registry path
- `--only` - Specific setup targets
- `--dry-run` - Show what would happen
- `--all` - Setup all targets
- `--name` - Target name
- `--build` - Build after setup

**Subcommands:**

- `wizard` - Interactive setup wizard
- `github` - GitHub setup
- `registry` - Registry setup
- `ci` - CI setup
- `repo` - Repository setup
- `bootstrap` - Bootstrap new workspace

**Examples:**

```bash
# Interactive setup
core setup wizard

# Setup from registry
core setup --registry repos.yaml --all

# Bootstrap workspace
core setup bootstrap
```

---

### test

**Run tests**

Execute Go tests with various options.

**Flags:**

- `--verbose` - Verbose output
- `--coverage` - Generate coverage report
- `--short` - Run short tests only
- `--pkg` - Specific package
- `--run` - Test filter pattern
- `--race` - Enable race detector
- `--json` - JSON output

**Examples:**

```bash
# Run all tests
core test

# Run with coverage
core test --coverage

# Run specific test
core test --run TestMyFunction

# Race detection
core test --race
```

---

### unifi

**UniFi network management**

Manage UniFi network devices and configuration.

**Subcommands:**

- `config` - Configure UniFi connection
- `clients` - List connected clients
- `devices` - List infrastructure devices
- `sites` - List controller sites
- `networks` - List network segments and VLANs
- `routes` - List gateway routing table

**Examples:**

```bash
# Configure UniFi
core unifi config --url https://192.168.1.1

# List clients
core unifi clients

# List devices
core unifi devices
```

---

### update

**Update core CLI to latest version**

Self-update the Core CLI binary.

**Flags:**

- `--channel` - Release channel (stable, beta, alpha, dev)
- `--force` - Force update even if on latest
- `--check` - Check without applying

**Subcommands:**

- `check` - Check for available updates

**Channels:**

- `stable` - Stable releases
- `beta` - Beta releases
- `alpha` - Alpha releases
- `dev` - Rolling development releases

**Examples:**

```bash
# Update to latest stable
core update

# Check for updates
core update check

# Update to beta channel
core update --channel beta
```

---

### vm

**Virtual machine management**

Manage LinuxKit VMs for testing and development.

**Subcommands:**

- `run` - Run a VM
- `ps` - List running VMs
- `stop` - Stop a VM
- `logs` - View VM logs
- `exec` - Execute command in VM
- `templates` - Manage VM templates

**Examples:**

```bash
# Run VM
core vm run my-vm

# List running VMs
core vm ps

# Execute command
core vm exec my-vm -- ls -la
```

---

### workspace

**Manage workspace configuration**

Configure workspace-specific settings.

**Subcommands:**

- `active` - Show or set active package
- `task` - Workspace task management
- `agent` - Workspace agent configuration

**Examples:**

```bash
# Show active workspace
core workspace active

# Set active package
core workspace active --set ./pkg/mypackage
```

---

## Common Patterns

### Flags and Options

**Common Persistent Flags:**

- `--verbose`, `-v` - Verbose output
- `--dry-run` - Show what would happen without executing
- `--output`, `-o` - Output file or directory
- `--json` - JSON output format
- `--help`, `-h` - Show help for command

### Environment Variables

Many commands respect environment variables. See [Environment Variables](environment-variables.md) for a complete list.

### Configuration Files

Commands load configuration from:

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file (`~/.core/config.yaml`)
4. Default values (lowest priority)

See [Configuration Reference](configuration.md) for details.

## Getting Help

For command-specific help:

```bash
core <command> --help
core <command> <subcommand> --help
```

For general help:

```bash
core help
core --help
```
