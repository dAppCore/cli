# Core CLI

Unified interface for development, multi-repo management, deployment, AI/ML operations, and infrastructure management.

## Quick Links

- **[Complete CLI Reference](../cli-reference.md)** - Comprehensive documentation of all commands
- **[Build Variants](../build-variants.md)** - Create custom builds with different feature sets
- **[MCP Integration](../mcp-integration.md)** - Model Context Protocol server for AI assistants
- **[Environment Variables](../environment-variables.md)** - Configuration via environment

## Command Categories

### Development & Testing
| Command | Description |
|---------|-------------|
| [dev](dev/) | Multi-repo workflow + dev environment |
| [go](go/) | Go development tools |
| [test](test/) | Run Go tests with coverage |
| [qa](../cli-reference.md#qa) | Quality assurance workflow |
| [doctor](doctor/) | Check environment |

### AI & Machine Learning
| Command | Description |
|---------|-------------|
| [ai](ai/) | AI agent task management and Claude integration |
| AI ml | ML inference and training pipeline |
| AI mcp | Model Context Protocol server |
| AI rag | RAG system for semantic search |
| AI collect | Data collection from external sources |

### Infrastructure & Deployment
| Command | Description |
|---------|-------------|
| AI deploy | Coolify deployment management |
| AI prod | Production infrastructure management |
| [vm](vm/) | LinuxKit VM management |
| AI unifi | UniFi network management |
| AI monitor | Security monitoring |

### Git & Forge
| Command | Description |
|---------|-------------|
| AI git | Root-level git workflow |
| AI forge | Forgejo instance management |
| AI gitea | Gitea instance management |

### Configuration & Utilities
| Command | Description |
|---------|-------------|
| AI config | Configuration management |
| AI crypt | Cryptographic utilities |
| [docs](docs/) | Documentation management |
| AI help | Help documentation |
| AI plugin | Plugin management |
| AI workspace | Workspace configuration |

### Services
| Command | Description |
|---------|-------------|
| AI daemon | Background service management |
| AI session | Session recording and replay |
| AI lab | Homelab monitoring dashboard |

## Installation

```bash
go install forge.lthn.ai/core/cli/cmd/core@latest
```

Verify: `core doctor`

See [Getting Started](../getting-started.md) for all installation options.

## Build Variants

The Core CLI supports modular builds. You can create specialized variants:

- **Full** - All commands (default)
- **Minimal** - Essential commands only
- **Developer** - Development-focused
- **DevOps** - Infrastructure and deployment
- **AI/ML** - AI and machine learning operations

See [Build Variants](../build-variants.md) for details.

## Quick Start

```bash
# Check environment
core doctor

# Multi-repo workflow
core dev work

# Start MCP server for AI assistants
core mcp serve --workspace /path/to/project

# Run tests
core test

# Get help
core help
core <command> --help
```

## Documentation

- [CLI Reference](../cli-reference.md) - Complete command documentation
- [Build Variants](../build-variants.md) - Custom builds
- [MCP Integration](../mcp-integration.md) - MCP server setup
- [Environment Variables](../environment-variables.md) - Configuration options
- [Configuration](../configuration.md) - Config file reference
- [Getting Started](../getting-started.md) - Installation guide
