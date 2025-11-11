<div align="center">
  <img src="clog.png" alt="clog logo" width="166" height="166">

  # clog - Claude Log Publisher for NATS

  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
  [![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
  [![NATS](https://img.shields.io/badge/NATS-Compatible-27AAE1?logo=nats.io)](https://nats.io/)
</div>

A lightweight Go binary for publishing Claude Code's task progress and questions to NATS. This tool enables real-time monitoring and integration of Claude Code workflows with your NATS-based systems.

Create the binary and put it in your project directory and tell Claude about it.

**Perfect for:**
- 🤖 Tracking Claude Code automation workflows
- 📊 Real-time monitoring of AI-assisted development
- 🔔 Event-driven notifications from Claude sessions
- 📈 Building analytics around Claude Code usage
- 🔗 Integrating Claude Code with existing NATS infrastructure

## Features

- ✅ **5 Authentication Methods**: None, Username/Password, Token, NKey, or Decentralized (JWT + Seed)
- ✅ **Zero Configuration**: Works out of the box with local NATS server
- ✅ **Environment Variable Support**: Override credentials at runtime
- ✅ **Secure by Default**: Credentials baked into binary, never stored in source
- ✅ **Comprehensive Event Types**: Tasks, questions, progress, and sessions
- ✅ **Claude Code Integration**: Easy hook setup for automatic tracking
- ✅ **Well-Tested**: Comprehensive test suite with multiple authentication scenarios

## Table of Contents

- [License](#license)
- [Quick Start with NATS](#quick-start-with-nats)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [NATS Subjects](#nats-subjects)
- [Message Format](#message-format)
- [Development](#development)
- [Contributing](#contributing)

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.

## Quick Start with NATS

If you don't have a NATS server running, here's how to quickly set one up:

### 1. Install NATS Server

**macOS (Homebrew):**
```bash
brew install nats-server
```

**Linux:**
```bash
# Download the latest release
curl -L https://github.com/nats-io/nats-server/releases/latest/download/nats-server-linux-amd64.zip -o nats-server.zip
unzip nats-server.zip
sudo mv nats-server /usr/local/bin/
```

**Or use Docker:**
```bash
docker run -d --name nats -p 4222:4222 nats:latest
```

### 2. Start NATS Server

```bash
# Simple start (no authentication)
nats-server

# Or with a custom port
nats-server -p 4222
```

You should see output like:
```
[1] 2025/10/09 18:00:00.000000 [INF] Starting nats-server
[1] 2025/10/09 18:00:00.000000 [INF]   Version:  2.10.0
[1] 2025/10/09 18:00:00.000000 [INF] Listening for client connections on 0.0.0.0:4222
```

### 3. Build and Use clog

```bash
# Clone and build
git clone https://github.com/davedotdev/clog.git
cd clog
go build -o clog ./cmd/main.go

# Send a test message
./clog -type=task -state=in_progress -message="Hello from clog!" -session="test-$(date +%s)"
```

### 4. Monitor NATS Messages (Optional)

Install the NATS CLI to see your messages in real-time:

```bash
# Install NATS CLI
brew install nats-io/nats-tools/nats  # macOS
# or download from: https://github.com/nats-io/natscli/releases

# Subscribe to all Claude events
nats sub "claude.>"

# Or subscribe to specific subjects
nats sub "claude.tasks.*"
nats sub "claude.questions.*"
```

Now when you use `clog`, you'll see the messages appear in your NATS subscriber!

### 5. Integrate with Claude Code (Global Setup)

Set up Claude Code to automatically use `clog` across **all your projects**:

```bash
# Run the setup command
make setup-claude
```

This will:
- ✅ Create `~/.claude/CLAUDE.md` with global instructions for Claude Code
- ✅ Configure `~/.claude/settings.json` to allow `clog` commands system-wide
- ✅ Make Claude Code aware of `clog` in every conversation

**What this does:**
- Claude Code will know about `clog` in **all your projects**, not just this one
- Claude will send NATS notifications before asking you questions
- No manual approval needed for `clog` commands
- Works globally since `clog` is installed in your system PATH

Now when Claude Code needs to ask you a question in **any project**, it will automatically send a NATS notification first!

## Installation

### Prerequisites

- Go 1.21 or later
- NATS server access with JWT/Seed credentials
- (Optional) Make for simplified building

### Building from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/davedotdev/clog.git
   cd clog
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build with configuration (recommended):**
   ```bash
   make build
   ```

   The build process will interactively prompt you for:
   - NATS URL (default: `nats://localhost:4222`)
   - Authentication type:
     1. **None** - No authentication
     2. **Username/Password** - Traditional username and password
     3. **Token** - Token-based authentication
     4. **NKey** - NKey authentication
     5. **Decentralized** - Decentralized authentication (JWT + Seed)

   Your configuration will be baked into the binary and then removed from the source code automatically.

4. **Set up Claude Code integration (optional, global):**
   ```bash
   make setup-claude
   ```

   This creates global configuration files in `~/.claude/` so Claude Code can automatically use `clog` across all your projects to send notifications when asking questions or tracking tasks.

5. **Install to system PATH:**
   ```bash
   make install
   ```

   This copies the binary to `/usr/local/bin/clog` (requires sudo).

### Alternative: Manual Build

If you prefer to manually edit credentials:

```bash
# Edit cmd/main.go and replace placeholders with your credentials
# Then build manually
go build -o clog ./cmd/main.go

# Copy to your PATH
sudo cp clog /usr/local/bin/clog
```

### Adding to PATH (for Claude Code)

For Claude Code to use `clog` in hooks and scripts, ensure it's in your PATH:

1. **Option 1: System-wide installation (recommended)**
   ```bash
   make install
   ```

2. **Option 2: User bin directory**
   ```bash
   mkdir -p ~/bin
   cp clog ~/bin/
   echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc  # or ~/.zshrc
   source ~/.bashrc  # or ~/.zshrc
   ```

3. **Verify installation:**
   ```bash
   which clog
   clog -h
   ```

## Configuration

### Baked-in Configuration (Default)

The recommended approach is to use `make build`, which:
1. Prompts for NATS URL and authentication type interactively
2. Bakes configuration into the binary at compile time
3. Keeps the source code clean with template placeholders

This ensures credentials are never committed to version control.

### Environment Variables (Runtime Override)

Override baked-in configuration at runtime using environment variables. The priority order is:

1. **Credentials file** (highest priority)
   ```bash
   export NATS_URL="nats://localhost:4222"
   export NATS_CREDS="/path/to/creds.file"
   ```

2. **Username/Password authentication**
   ```bash
   export NATS_URL="nats://localhost:4222"
   export NATS_USERNAME="myuser"
   export NATS_PASSWORD="mypassword"
   ```

3. **Token authentication**
   ```bash
   export NATS_URL="nats://localhost:4222"
   export NATS_TOKEN="mytoken"
   ```

4. **NKey authentication**
   ```bash
   export NATS_URL="nats://localhost:4222"
   export NATS_NKEY="SUABC..."
   ```

5. **Decentralized authentication (JWT + Seed)**
   ```bash
   export NATS_URL="nats://localhost:4222"
   export NATS_JWT="eyJ0eXAiOiJKV1Q..."
   export NATS_SEED="SUAK7SG5BVF..."
   ```

6. **Baked-in credentials** (lowest priority - from build time)

### Using with Claude Code

#### Automatic Setup (Recommended - Global)

Run the setup command to configure Claude Code to automatically use `clog` across all your projects:

```bash
make setup-claude
```

This will:
- Create `~/.claude/CLAUDE.md` - Global memory file that instructs Claude Code when and how to use clog
- Create or update `~/.claude/settings.json` - Global permissions file that allows clog commands without approval
- Ensure seamless integration with Claude Code in every project

**What Claude Code will do automatically:**
- Know about `clog` in all your projects (not just this one)
- Send NATS notifications before asking you questions
- Track task progress when working on multi-step tasks
- No manual intervention required
- Works anywhere since `clog` is in your system PATH

#### Manual Integration

If you prefer manual setup, you can also use `clog` in custom hooks or scripts.

**Example usage in scripts:**
```bash
#!/bin/bash
SESSION_ID="my-session-$(date +%s)"

# Start session
clog -type=session -message="Starting deployment process" -session="$SESSION_ID"

# Log task progress
clog -type=task -state=in_progress -message="Building application" -task-num="1/5" -session="$SESSION_ID"

# Task completed
clog -type=task -state=completed -message="Build successful" -task-num="1/5" -session="$SESSION_ID"

# Ask question
clog -type=question -state=blocked -message="Deploy to staging or production?" -session="$SESSION_ID"
```

## Usage

See help:
```bash
./clog -h
```

Examples:
```bash
# Task started
./clog -type=task -state=in_progress -message="Adding VAT breakdown" -task-num="3/15" -session="nye-api"

# Task completed
./clog -type=task -state=completed -message="VAT breakdown added" -task-num="3/15" -session="nye-api"

# Question
./clog -type=question -state=blocked -message="Should VAT be inclusive?" -session="nye-api"

# Progress
./clog -type=progress -message="50% complete" -session="nye-api"
```

## NATS Subjects

The tool publishes to these hardwired subjects based on type and state:

- `claude.tasks.started` - Task in progress
- `claude.tasks.completed` - Task completed
- `claude.tasks.blocked` - Task blocked
- `claude.questions.asked` - Question asked
- `claude.questions.waiting` - Waiting for answer
- `claude.progress.update` - Progress update
- `claude.session.started` - Session started
- `claude.session.completed` - Session completed

## Message Format

Messages are published as JSON:

```json
{
  "event": "claude.tasks.completed",
  "timestamp": "2025-10-09T14:30:00Z",
  "session_id": "nye-api-1696854321-a4f9",
  "message": "VAT breakdown added",
  "state": "completed",
  "task_num": "3/15"
}
```

## Exit Codes

- `0` - Success
- `1` - Invalid arguments
- `2` - NATS connection failed

## Development

### Available Make Targets

```bash
make build         # Build clog with interactive configuration
make setup-claude  # Configure Claude Code integration
make install       # Install clog to /usr/local/bin
make unittest      # Run unit tests
make test          # Build and run help output
make clean         # Clean build artifacts
make tidy          # Run go mod tidy
```

### Running Tests

```bash
# Run unit tests
make unittest

# Run all tests with coverage
go test -v -cover ./cmd/
```

### Building for Development

```bash
# Build without credential prompts (uses placeholders)
go build -o clog ./cmd/main.go

# Run with environment variables
export NATS_URL="nats://localhost:4222"
export NATS_CREDS="/path/to/test.creds"
./clog -type=task -message="Test" -session="dev"
```

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and best practices.

### Security Notes

- Configuration and credentials are baked into the binary at build time
- Source code uses template placeholders to prevent accidental credential commits
- `.gitignore` is configured to exclude binaries and credential files
- Never commit actual NATS credentials to version control
- Review the build process in `Makefile` to understand credential handling
- Supports 5 authentication methods: none, username/password, token, NKey, and decentralized (JWT + Seed)
