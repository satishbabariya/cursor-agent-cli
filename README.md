# Cursor Background Agent CLI

A command-line interface tool for managing Cursor Background Agents built with Go and Cobra.

## Features

- 🚀 **Easy Setup**: Initialize with your API key using `cursor-cli init`
- 📋 **List Agents**: View all your background agents with pagination support
- 🔍 **Agent Status**: Get detailed status and information about specific agents
- 💬 **Conversation History**: View the conversation history of any agent
- 📤 **Follow-up Instructions**: Send additional instructions to running agents
- 🔑 **API Key Management**: View information about your current API key
- ⚙️ **Configuration**: Stores API key securely in `~/.cursor-cli.yaml`

## Installation

### From Source

```bash
git clone https://github.com/satishbabariya/cursor-background-agent-cli.git
cd cursor-background-agent-cli
go build -o cursor-cli
sudo mv cursor-cli /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/satishbabariya/cursor-background-agent-cli@latest
```

## Getting Started

1. **Initialize the CLI with your API key:**
   ```bash
   cursor-cli init
   ```
   
   Get your API key from [Cursor Dashboard → Integrations](https://cursor.com/dashboard).

2. **List your background agents:**
   ```bash
   cursor-cli list
   ```

3. **Get status of a specific agent:**
   ```bash
   cursor-cli status bc_abc123
   ```

## Commands

### `cursor-cli init`
Initialize cursor-cli with your API key. This will prompt you to enter your API key and validate it.

### `cursor-cli list [flags]`
List background agents associated with your account. By default, only active agents (running, completed, failed, cancelled) are shown, excluding expired ones.

**Flags:**
- `-l, --limit int`: Number of agents to return (1-100, default: 20)
- `-c, --cursor string`: Pagination cursor from previous response
- `-a, --all`: Show all agents including expired ones

**Examples:**
```bash
cursor-cli list                    # Show only active agents
cursor-cli list --all              # Show all agents including expired
cursor-cli list --limit 50         # Show 50 active agents
cursor-cli list --cursor bc_def456 # Get next page
```

### `cursor-cli status <agent-id>`
Get the current status and detailed information about a specific background agent.

**Example:**
```bash
cursor-cli status bc_abc123
```

### `cursor-cli conversation <agent-id>`
Retrieve the conversation history of a background agent.

**Example:**
```bash
cursor-cli conversation bc_abc123
```

### `cursor-cli followup <agent-id> <prompt>`
Send an additional instruction to a running background agent.

**Example:**
```bash
cursor-cli followup bc_abc123 "Also add a section about troubleshooting"
```

### `cursor-cli keyinfo`
Display information about your current API key.

**Example:**
```bash
cursor-cli keyinfo
```

## Configuration

The CLI stores your API key in `~/.cursor-cli.yaml`. You can also set the API key using:

- Environment variable: `CURSOR_API_KEY`
- Command-line flag: `--api-key`

## API Reference

This CLI is built on top of the [Cursor Background Agents API](https://docs.cursor.com/en/background-agent/api/overview). The following endpoints are supported:

- `GET /v0/agents` - [List Agents](https://docs.cursor.com/en/background-agent/api/list-agents)
- `GET /v0/agents/{id}` - [Agent Status](https://docs.cursor.com/en/background-agent/api/agent-status)
- `GET /v0/agents/{id}/conversation` - [Agent Conversation](https://docs.cursor.com/en/background-agent/api/agent-conversation)
- `POST /v0/agents/{id}/followup` - [Add Follow-up](https://docs.cursor.com/en/background-agent/api/add-followup)
- `GET /v0/me` - User/API Key Info

## Error Handling

The CLI provides clear error messages for common issues:

- **API key not set**: Run `cursor-cli init` to set up your API key
- **Invalid API key**: Check that your API key is correct and active
- **Network errors**: Check your internet connection
- **Agent not found**: Verify the agent ID is correct

## Development

### Prerequisites

- Go 1.21 or later
- Access to Cursor Background Agents API

### Building from Source

```bash
git clone https://github.com/satishbabariya/cursor-background-agent-cli.git
cd cursor-background-agent-cli
go mod download
go build -o cursor-cli
```

### Project Structure

```
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command and configuration
│   ├── init.go            # API key initialization
│   ├── list.go            # List agents command
│   ├── status.go          # Agent status command
│   ├── conversation.go    # Agent conversation command
│   ├── followup.go        # Add follow-up command
│   └── keyinfo.go         # API key info command
├── internal/
│   ├── client/            # API client
│   │   └── client.go      # HTTP client and API methods
│   └── config/            # Configuration management
│       └── config.go      # Config file handling
├── main.go                # Entry point
├── go.mod                 # Go module file
└── README.md              # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- 📖 [Cursor Background Agents Documentation](https://docs.cursor.com/en/background-agent/api/overview)
- 🐛 [Report Issues](https://github.com/satishbabariya/cursor-background-agent-cli/issues)
- 💬 [Cursor Discord](https://discord.gg/cursor) - #background-agent channel
