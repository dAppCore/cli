# Environment Variables

The Core CLI uses environment variables for configuration, authentication, and runtime behavior. This document lists all supported environment variables and their purposes.

## Configuration System

Environment variables are part of Core's layered configuration system:

1. **Command-line flags** (highest priority)
2. **Environment variables** ← This document
3. **Configuration file** (`~/.core/config.yaml`)
4. **Default values** (lowest priority)

## Configuration File Locations

### User Configuration

| File | Path | Purpose |
|------|------|---------|
| Main config | `~/.core/config.yaml` | User preferences and framework settings |
| Agentic config | `~/.core/agentic.yaml` | AI agent service configuration |
| MCP config | `~/.claude/mcp_config.json` | Claude Code MCP server settings |
| Keybindings | `~/.claude/keybindings.json` | Claude Code keyboard shortcuts |

### Project Configuration

Located in `.core/` directory of project root:

| File | Path | Purpose |
|------|------|---------|
| Build config | `.core/build.yaml` | Build targets and flags |
| Release config | `.core/release.yaml` | Release automation settings |
| CI config | `.core/ci.yaml` | CI pipeline configuration |

### Registry Configuration

Searched in order:

1. Current directory: `./repos.yaml`
2. Parent directories (walking up)
3. Home Code directory: `~/Code/host-uk/repos.yaml`
4. Config directory: `~/.config/core/repos.yaml`

## General Configuration

### CORE_CONFIG_*

Map configuration values to the YAML hierarchy using dot notation.

**Format:** `CORE_CONFIG_<KEY>=<value>`

After stripping the `CORE_CONFIG_` prefix, the remaining variable name is converted to lowercase and underscores are replaced with dots.

**Examples:**

```bash
# dev.editor: vim
export CORE_CONFIG_DEV_EDITOR=vim

# log.level: debug
export CORE_CONFIG_LOG_LEVEL=debug

# ai.model: gpt-4
export CORE_CONFIG_AI_MODEL=gpt-4
```

### NO_COLOR

Disable ANSI color output.

**Values:** Any value (presence disables colors)

**Example:**
```bash
export NO_COLOR=1
core doctor  # No colored output
```

### CORE_DAEMON

Run application in daemon mode.

**Values:** `1` or `true`

**Example:**
```bash
export CORE_DAEMON=1
core mcp serve  # Runs as background daemon
```

## MCP Server

### MCP_ADDR

TCP address for MCP server. If not set, MCP uses stdio transport.

**Format:** `host:port` or `:port`

**Default:** (stdio mode)

**Examples:**
```bash
# Listen on localhost:9999
export MCP_ADDR=localhost:9999
core mcp serve

# Listen on all interfaces (⚠️ not recommended)
export MCP_ADDR=:9999
core mcp serve
```

### CORE_MCP_TRANSPORT

MCP transport mode.

**Values:** `stdio`, `tcp`, `socket`

**Default:** `stdio`

**Example:**
```bash
export CORE_MCP_TRANSPORT=tcp
export CORE_MCP_ADDR=:9100
core daemon
```

### CORE_MCP_ADDR

Alternative to `MCP_ADDR` for daemon mode.

**Format:** Same as `MCP_ADDR`

### CORE_HEALTH_ADDR

Health check endpoint address for daemon mode.

**Format:** `host:port` or `:port`

**Default:** None (disabled)

**Example:**
```bash
export CORE_HEALTH_ADDR=:8080
core daemon
# Health check: curl http://localhost:8080/health
```

### CORE_PID_FILE

PID file location for daemon mode.

**Default:** (platform-specific temp directory)

**Example:**
```bash
export CORE_PID_FILE=/var/run/core-mcp.pid
core daemon
```

## Service APIs

### COOLIFY_TOKEN

API token for Coolify deployments.

**Example:**
```bash
export COOLIFY_TOKEN=your-api-token
core deploy servers
```

### AGENTIC_TOKEN

API token for Agentic AI services.

**Example:**
```bash
export AGENTIC_TOKEN=your-token
core ai tasks list
```

## Networking

### UNIFI_URL

UniFi controller URL.

**Format:** `https://host` or `https://host:port`

**Example:**
```bash
export UNIFI_URL=https://192.168.1.1
core unifi clients
```

### UNIFI_INSECURE

Skip TLS certificate verification for UniFi controller.

**Values:** `1`, `true`, `yes`

**Example:**
```bash
export UNIFI_INSECURE=1
export UNIFI_URL=https://192.168.1.1
core unifi devices
```

## Git/Forge Integration

### GITHUB_TOKEN

GitHub API token for operations requiring authentication.

**Note:** Usually set by GitHub CLI (`gh auth login`)

**Example:**
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxx
core dev issues
```

### FORGEJO_URL

Forgejo instance URL.

**Example:**
```bash
export FORGEJO_URL=https://forge.example.com
core forge repos
```

### FORGEJO_TOKEN

Forgejo API token.

**Example:**
```bash
export FORGEJO_TOKEN=your-token
core forge issues
```

### GITEA_URL

Gitea instance URL.

**Example:**
```bash
export GITEA_URL=https://gitea.example.com
core gitea repos
```

### GITEA_TOKEN

Gitea API token.

**Example:**
```bash
export GITEA_TOKEN=your-token
core gitea issues
```

## RAG System

### QDRANT_HOST

Qdrant vector database host.

**Default:** `localhost`

**Example:**
```bash
export QDRANT_HOST=localhost
export QDRANT_PORT=6334
core rag collections list
```

### QDRANT_PORT

Qdrant gRPC port.

**Default:** `6334`

### OLLAMA_HOST

Ollama API host for embeddings.

**Default:** `localhost`

**Example:**
```bash
export OLLAMA_HOST=localhost
export OLLAMA_PORT=11434
core rag ingest --collection docs --path ./documentation
```

### OLLAMA_PORT

Ollama API port.

**Default:** `11434`

## ML Pipeline

### ML_API_URL

OpenAI-compatible API endpoint for ML operations.

**Example:**
```bash
export ML_API_URL=http://localhost:8000/v1
core ml score
```

### ML_JUDGE_URL

Judge model API URL (typically Ollama).

**Example:**
```bash
export ML_JUDGE_URL=http://localhost:11434
core ml score --judge-model llama2
```

### ML_JUDGE_MODEL

Judge model name for scoring.

**Default:** (varies by command)

**Example:**
```bash
export ML_JUDGE_MODEL=llama2:latest
core ml score
```

### INFLUX_URL

InfluxDB URL for metrics storage.

**Example:**
```bash
export INFLUX_URL=http://localhost:8086
export INFLUX_DB=metrics
core ml metrics
```

### INFLUX_DB

InfluxDB database name.

**Example:**
```bash
export INFLUX_DB=ml_metrics
core ml ingest
```

## Build System

### CGO_ENABLED

Enable or disable CGO for builds.

**Values:** `0` (disabled), `1` (enabled)

**Default:** `0` (Go default)

**Example:**
```bash
export CGO_ENABLED=1
task cli:build
```

### GOOS

Target operating system for cross-compilation.

**Values:** `linux`, `darwin`, `windows`, `freebsd`, etc.

**Example:**
```bash
export GOOS=linux
export GOARCH=amd64
task cli:build
```

### GOARCH

Target architecture for cross-compilation.

**Values:** `amd64`, `arm64`, `386`, `arm`, etc.

**Example:**
```bash
export GOOS=darwin
export GOARCH=arm64
task cli:build  # Build for Apple Silicon
```

## Testing

### TEST_VERBOSE

Enable verbose test output.

**Values:** `1`, `true`

**Example:**
```bash
export TEST_VERBOSE=1
core test
```

### TEST_SHORT

Run short tests only.

**Values:** `1`, `true`

**Example:**
```bash
export TEST_SHORT=1
core test
```

### TEST_RACE

Enable race detector.

**Values:** `1`, `true`

**Example:**
```bash
export TEST_RACE=1
core test
```

## Infrastructure

### HETZNER_TOKEN

Hetzner Cloud API token.

**Example:**
```bash
export HETZNER_TOKEN=your-token
core prod status
```

### CLOUDNS_AUTH_ID

CloudNS authentication ID for DNS management.

**Example:**
```bash
export CLOUDNS_AUTH_ID=your-id
export CLOUDNS_AUTH_PASSWORD=your-password
core prod dns
```

### CLOUDNS_AUTH_PASSWORD

CloudNS authentication password.

## Virtual Machines

### LINUXKIT_PATH

Path to LinuxKit binary.

**Default:** `linuxkit` (from PATH)

**Example:**
```bash
export LINUXKIT_PATH=/usr/local/bin/linuxkit
core vm run
```

### VM_MEMORY

Default memory allocation for VMs.

**Format:** Integer (MB) or string with suffix (e.g., "2G")

**Default:** 2048 (2GB)

**Example:**
```bash
export VM_MEMORY=4096
core vm run my-vm
```

### VM_CPUS

Default CPU count for VMs.

**Default:** 2

**Example:**
```bash
export VM_CPUS=4
core vm run my-vm
```

## Development

### DEBUG

Enable debug logging.

**Values:** `1`, `true`

**Example:**
```bash
export DEBUG=1
core dev work
```

### LOG_LEVEL

Global log level.

**Values:** `debug`, `info`, `warn`, `error`

**Default:** `info`

**Example:**
```bash
export LOG_LEVEL=debug
core mcp serve
```

### EDITOR

Default text editor for interactive commands.

**Default:** `vim` or system default

**Example:**
```bash
export EDITOR=nvim
core dev commit  # Uses neovim for editing
```

### PAGER

Default pager for command output.

**Default:** `less` or system default

**Example:**
```bash
export PAGER=bat
core help
```

## Best Practices

### 1. Use Config File for Persistent Settings

Environment variables are temporary. For persistent configuration, use `~/.core/config.yaml`:

```yaml
dev:
  editor: nvim
log:
  level: info
ai:
  model: gpt-4
```

Set via CLI:
```bash
core config set dev.editor nvim
```

### 2. Use .env Files for Project-Specific Settings

Create `.env` in project root:

```bash
# .env
MCP_ADDR=localhost:9100
QDRANT_HOST=localhost
OLLAMA_HOST=localhost
```

Load with:
```bash
source .env
core mcp serve
```

Or use `direnv` for automatic loading.

### 3. Secure Sensitive Tokens

Never commit tokens to version control:

```bash
# .gitignore
.env
.env.local
*.token
```

Use environment-specific files:
```bash
.env.development
.env.staging
.env.production
```

### 4. Override in CI/CD

Set environment variables in CI/CD pipelines:

```yaml
# .github/workflows/test.yml
env:
  CORE_CONFIG_LOG_LEVEL: debug
  NO_COLOR: 1
```

### 5. Document Project Requirements

Create `.env.example` with required variables:

```bash
# .env.example
MCP_ADDR=localhost:9100
QDRANT_HOST=localhost
QDRANT_PORT=6334
OLLAMA_HOST=localhost
# FORGEJO_TOKEN=your-token-here
```

## Checking Configuration

### View Current Config

```bash
# All config values
core config list

# Specific value
core config get dev.editor
```

### View Config File Location

```bash
core config path
```

### Debug Environment Variables

```bash
# Show all CORE_* variables
env | grep CORE_

# Show all MCP-related variables
env | grep MCP
```

### Test Configuration

```bash
# Verify environment setup
core doctor --verbose
```

## Troubleshooting

### Variable Not Taking Effect

1. **Check priority:** Command flags override environment variables
   ```bash
   # This flag takes priority
   core mcp serve --workspace /path
   ```

2. **Verify variable is set:**
   ```bash
   echo $MCP_ADDR
   ```

3. **Check for typos:**
   ```bash
   # Wrong
   export CORE_CONFIG_dev_editor=vim

   # Correct
   export CORE_CONFIG_DEV_EDITOR=vim
   ```

4. **Restart shell or reload:**
   ```bash
   source ~/.bashrc  # or ~/.zshrc
   ```

### Config File vs Environment Variables

If config file and environment variable conflict, environment variable wins:

```yaml
# ~/.core/config.yaml
log:
  level: info
```

```bash
# This overrides the config file
export CORE_CONFIG_LOG_LEVEL=debug
core mcp serve  # Uses debug level
```

### Case Sensitivity

Variable names are case-sensitive:

```bash
# Wrong
export mcp_addr=:9100

# Correct
export MCP_ADDR=:9100
```

## See Also

- [Configuration Reference](configuration.md) - Complete config file documentation
- [CLI Reference](cli-reference.md) - Command documentation
- [MCP Integration](mcp-integration.md) - MCP server setup
- [Getting Started](getting-started.md) - Installation and setup
