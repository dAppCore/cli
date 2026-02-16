# MCP Integration

The Core CLI includes a Model Context Protocol (MCP) server that provides file operations, RAG (Retrieval-Augmented Generation), metrics tracking, and process management tools for AI assistants like Claude Code.

## Overview

MCP is a protocol that allows AI assistants to interact with external tools and services. The Core CLI's MCP server exposes a rich set of tools for:

- **File Operations** - Read, write, edit, delete files and directories
- **RAG System** - Semantic search using Qdrant + Ollama embeddings
- **Metrics** - Record and query application metrics
- **Process Management** - Start, stop, and monitor processes
- **Language Detection** - Detect programming languages
- **WebSocket** - Real-time communication
- **WebView/CDP** - Browser automation via Chrome DevTools Protocol

## Quick Start

### Start MCP Server

**Stdio mode (for Claude Code):**
```bash
core mcp serve
```

**TCP mode:**
```bash
MCP_ADDR=localhost:9999 core mcp serve
```

**With workspace restriction:**
```bash
core mcp serve --workspace /path/to/project
```

### Configure in Claude Code

Add to your Claude Code MCP configuration (`~/.claude/mcp_config.json`):

```json
{
  "mcpServers": {
    "core-mcp": {
      "command": "core",
      "args": ["mcp", "serve"],
      "env": {}
    }
  }
}
```

With workspace restriction:

```json
{
  "mcpServers": {
    "core-mcp": {
      "command": "core",
      "args": ["mcp", "serve", "--workspace", "/home/user/projects/myapp"],
      "env": {}
    }
  }
}
```

## Transport Modes

### Stdio (Default)

Best for integration with Claude Code and other AI assistants that communicate via standard input/output.

```bash
core mcp serve
```

**Pros:**
- Simple setup
- Secure (no network exposure)
- Works with Claude Code out of the box

**Cons:**
- Single client only
- Process-coupled (client must manage server lifecycle)

### TCP

Best for network-based integrations or when you need multiple clients.

```bash
MCP_ADDR=localhost:9999 core mcp serve
```

**Pros:**
- Multiple clients can connect
- Server runs independently
- Can connect from remote machines (if exposed)

**Cons:**
- Requires network configuration
- No built-in authentication (use firewall/SSH tunnel)
- More complex setup

## Available Tools

### File Operations

#### file_read
Read contents of a file.

**Parameters:**
- `path` (string, required) - File path to read

**Example:**
```json
{
  "name": "file_read",
  "arguments": {
    "path": "/path/to/file.txt"
  }
}
```

#### file_write
Write content to a file (creates or overwrites).

**Parameters:**
- `path` (string, required) - File path to write
- `content` (string, required) - Content to write

**Example:**
```json
{
  "name": "file_write",
  "arguments": {
    "path": "/path/to/file.txt",
    "content": "Hello, world!"
  }
}
```

#### file_edit
Edit a file by replacing old content with new content.

**Parameters:**
- `path` (string, required) - File path to edit
- `old_content` (string, required) - Content to find
- `new_content` (string, required) - Replacement content

**Example:**
```json
{
  "name": "file_edit",
  "arguments": {
    "path": "/path/to/file.txt",
    "old_content": "old text",
    "new_content": "new text"
  }
}
```

#### file_delete
Delete a file.

**Parameters:**
- `path` (string, required) - File path to delete

#### file_rename
Rename or move a file.

**Parameters:**
- `old_path` (string, required) - Current file path
- `new_path` (string, required) - New file path

#### file_exists
Check if a file or directory exists.

**Parameters:**
- `path` (string, required) - Path to check

**Returns:**
- `exists` (boolean)
- `is_dir` (boolean)

#### dir_list
List contents of a directory.

**Parameters:**
- `path` (string, required) - Directory path
- `recursive` (boolean, optional) - Recursively list subdirectories

#### dir_create
Create a directory (including parent directories).

**Parameters:**
- `path` (string, required) - Directory path to create

### RAG (Retrieval-Augmented Generation)

#### rag_ingest
Ingest documents into a collection for semantic search.

**Parameters:**
- `collection` (string, required) - Collection name
- `documents` (array, required) - Documents to ingest
  - `id` (string) - Document ID
  - `text` (string) - Document content
  - `metadata` (object, optional) - Additional metadata

**Example:**
```json
{
  "name": "rag_ingest",
  "arguments": {
    "collection": "documentation",
    "documents": [
      {
        "id": "doc1",
        "text": "This is a documentation page about MCP",
        "metadata": {"type": "docs", "topic": "mcp"}
      }
    ]
  }
}
```

#### rag_query
Query a collection using semantic search.

**Parameters:**
- `collection` (string, required) - Collection name
- `query` (string, required) - Search query
- `limit` (integer, optional) - Maximum results (default: 5)

**Returns:**
Array of matching documents with scores.

#### rag_collections
List all RAG collections.

**Returns:**
Array of collection names and their document counts.

### Metrics

#### metrics_record
Record a metric value.

**Parameters:**
- `name` (string, required) - Metric name
- `value` (number, required) - Metric value
- `tags` (object, optional) - Metric tags/labels
- `timestamp` (string, optional) - ISO 8601 timestamp

**Example:**
```json
{
  "name": "metrics_record",
  "arguments": {
    "name": "api.requests",
    "value": 1,
    "tags": {"endpoint": "/api/users", "status": "200"}
  }
}
```

#### metrics_query
Query recorded metrics.

**Parameters:**
- `name` (string, optional) - Filter by metric name
- `tags` (object, optional) - Filter by tags
- `from` (string, optional) - Start time (ISO 8601)
- `to` (string, optional) - End time (ISO 8601)

### Language Detection

#### lang_detect
Detect programming language from code or filename.

**Parameters:**
- `content` (string, optional) - Code content
- `filename` (string, optional) - Filename

**Returns:**
- `language` (string) - Detected language
- `confidence` (number) - Detection confidence (0-1)

#### lang_list
List all supported languages.

**Returns:**
Array of supported language identifiers.

### Process Management

#### process_start
Start a new process.

**Parameters:**
- `id` (string, required) - Process identifier
- `command` (string, required) - Command to execute
- `args` (array, optional) - Command arguments
- `env` (object, optional) - Environment variables
- `cwd` (string, optional) - Working directory

**Example:**
```json
{
  "name": "process_start",
  "arguments": {
    "id": "build-process",
    "command": "go",
    "args": ["build", "-o", "bin/app"],
    "cwd": "/path/to/project"
  }
}
```

#### process_stop
Stop a running process (SIGTERM).

**Parameters:**
- `id` (string, required) - Process identifier

#### process_kill
Force kill a process (SIGKILL).

**Parameters:**
- `id` (string, required) - Process identifier

#### process_list
List all managed processes.

**Returns:**
Array of process info (ID, command, status, PID).

#### process_output
Get output from a process.

**Parameters:**
- `id` (string, required) - Process identifier
- `tail` (integer, optional) - Number of recent lines

#### process_input
Send input to a process.

**Parameters:**
- `id` (string, required) - Process identifier
- `input` (string, required) - Input to send

### WebSocket

#### ws_start
Start a WebSocket hub server.

**Parameters:**
- `addr` (string, optional) - Listen address (default: ":8080")

**Returns:**
- `url` (string) - WebSocket server URL

#### ws_info
Get WebSocket server information.

**Returns:**
- `running` (boolean)
- `url` (string)
- `connections` (integer)

### WebView/CDP (Chrome DevTools Protocol)

#### webview_connect
Connect to Chrome/Chromium via CDP.

**Parameters:**
- `url` (string, required) - Debug URL (e.g., "http://localhost:9222")

#### webview_navigate
Navigate to a URL.

**Parameters:**
- `url` (string, required) - URL to navigate to

#### webview_click
Click an element.

**Parameters:**
- `selector` (string, required) - CSS selector

#### webview_type
Type text into an element.

**Parameters:**
- `selector` (string, required) - CSS selector
- `text` (string, required) - Text to type

#### webview_query
Query DOM for elements.

**Parameters:**
- `selector` (string, required) - CSS selector

**Returns:**
Array of matching elements.

#### webview_console
Get console logs.

**Returns:**
Array of console messages.

#### webview_eval
Evaluate JavaScript.

**Parameters:**
- `script` (string, required) - JavaScript code

**Returns:**
Evaluation result.

#### webview_screenshot
Take a screenshot.

**Parameters:**
- `full_page` (boolean, optional) - Full page screenshot

**Returns:**
- `data` (string) - Base64-encoded PNG

#### webview_wait
Wait for condition.

**Parameters:**
- `type` (string, required) - Condition type ("selector", "navigation", "timeout")
- `value` (string/number, required) - Condition value

#### webview_disconnect
Disconnect from CDP.

## Security Considerations

### Workspace Restriction

**Always** use `--workspace` flag when exposing MCP server to AI assistants:

```bash
core mcp serve --workspace /path/to/safe/directory
```

This restricts file operations to the specified directory and its subdirectories, preventing the AI from accessing sensitive files outside the project.

### File Operations

Without workspace restriction, the MCP server has full filesystem access. This is powerful but dangerous:

❌ **Never** run unrestricted in production:
```bash
core mcp serve  # Full filesystem access!
```

✅ **Always** restrict to project directory:
```bash
core mcp serve --workspace /home/user/projects/myapp
```

### Network Exposure

When using TCP mode:

- **Localhost only** by default: `MCP_ADDR=localhost:9999`
- **Never** expose to public internet without authentication
- Use SSH tunnels for remote access
- Consider firewall rules to restrict access

### Process Management

Processes started via `process_start` run with the same permissions as the MCP server. Be cautious when:

- Running as root/administrator
- Executing untrusted commands
- Exposing process management to external clients

## Running as Daemon

Use the `daemon` command to run MCP server as a background service:

```bash
core daemon --mcp-transport tcp --mcp-addr :9100
```

With health endpoint:

```bash
core daemon --mcp-transport tcp --mcp-addr :9100 --health-addr :8080
```

Check health:

```bash
curl http://localhost:8080/health
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MCP_ADDR` | TCP listen address | (stdio) |
| `CORE_MCP_TRANSPORT` | Transport mode (stdio, tcp) | stdio |
| `CORE_MCP_ADDR` | Alternative to MCP_ADDR | - |
| `CORE_HEALTH_ADDR` | Health check endpoint | - |

## Use Cases

### 1. Claude Code Integration

Enable Claude Code to read and modify files in your project:

```json
{
  "mcpServers": {
    "core": {
      "command": "core",
      "args": ["mcp", "serve", "--workspace", "${workspaceFolder}"]
    }
  }
}
```

### 2. Semantic Code Search

Ingest codebase into RAG system for semantic search:

```bash
# Start MCP server
core mcp serve --workspace /path/to/project

# Ingest code (via AI assistant)
# Use: rag_ingest with collection="code" and documents from source files

# Query code
# Use: rag_query with query="authentication logic"
```

### 3. Build Automation

Monitor build processes with real-time output:

```bash
# Start MCP with WebSocket
# Use: ws_start
# Use: process_start with command="go build"
# Subscribe to process output via WebSocket
```

### 4. Web Automation Testing

Automate browser testing:

```bash
# Start Chrome with remote debugging:
# google-chrome --remote-debugging-port=9222

# Use CDP tools:
# webview_connect url="http://localhost:9222"
# webview_navigate url="https://example.com"
# webview_click selector="#login-button"
# webview_screenshot
```

## Troubleshooting

### Server Won't Start

**Error:** Port already in use

```bash
# Check what's using the port
lsof -i :9999

# Kill the process or use different port
MCP_ADDR=:9100 core mcp serve
```

**Error:** Permission denied

```bash
# Don't use privileged ports (<1024) without sudo
# Use high ports instead:
MCP_ADDR=:9999 core mcp serve
```

### File Operations Fail

**Error:** File outside workspace

Make sure file paths are within the workspace:

```bash
# Workspace: /home/user/project
# ✅ OK: /home/user/project/src/main.go
# ❌ Fail: /etc/passwd
```

**Error:** Permission denied

Check file/directory permissions and MCP server user.

### RAG Not Working

**Error:** Qdrant connection failed

Ensure Qdrant is running:

```bash
# Start Qdrant (Docker)
docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

**Error:** Ollama not available

Ensure Ollama is running:

```bash
# Start Ollama
ollama serve

# Pull embedding model
ollama pull nomic-embed-text
```

### Process Management Issues

**Error:** Process not found

List processes to verify ID:

```json
{"name": "process_list"}
```

**Error:** Process already exists

Use unique process IDs or stop existing process first.

## Performance Tips

1. **Use workspace restriction** - Reduces filesystem traversal overhead
2. **Batch RAG ingestion** - Ingest multiple documents in one call
3. **Limit query results** - Use `limit` parameter in rag_query
4. **Process output streaming** - Use WebSocket for real-time output
5. **Reuse connections** - Keep WebView CDP connection open for multiple operations

## See Also

- [CLI Reference](cli-reference.md) - Complete command documentation
- [Daemon Mode](daemon-mode.md) - Running MCP as a service
- [RAG System](rag-system.md) - Detailed RAG documentation
- [Environment Variables](environment-variables.md) - Configuration options
