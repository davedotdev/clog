.PHONY: build clean install install-system test unittest setup-claude

# Setup Claude Code memory file (user-level, global)
setup-claude:
	@echo "Setting up Claude Code configuration (global)..."
	@mkdir -p ~/.claude
	@printf '%s\n' \
		'# clog - Claude Log Publisher' \
		'' \
		'## Startup Instructions' \
		'' \
		'You have access to `clog`, a tool for sending NATS messages when working with users.' \
		'' \
		'When starting any conversation, you should:' \
		'' \
		'1. Run `clog -h` to see available commands and usage' \
		'2. Remember that `clog` is globally available in your PATH' \
		'' \
		'## When to Use clog' \
		'' \
		'**IMPORTANT**: Before asking the user any question, you should send a NATS notification using:' \
		'' \
		'```bash' \
		'clog -type=question -state=blocked -message="<your question here>" -session="<session-id>"' \
		'```' \
		'' \
		'Then ask your question normally in the conversation.' \
		'' \
		'## Available Event Types' \
		'' \
		'- **question**: Use before asking the user a question (state: blocked)' \
		'- **task**: Track task progress (states: in_progress, completed, blocked)' \
		'- **progress**: Send progress updates' \
		'- **session**: Mark session start/end' \
		'' \
		'## Example Usage' \
		'' \
		'```bash' \
		'# Before asking a question' \
		'clog -type=question -state=blocked -message="Should I deploy to staging or production?" -session="deploy-task"' \
		'' \
		'# Starting a task' \
		'clog -type=task -state=in_progress -message="Running tests" -task-num="1/5" -session="test-suite"' \
		'' \
		'# Completing a task' \
		'clog -type=task -state=completed -message="Tests passed" -task-num="5/5" -session="test-suite"' \
		'```' \
		> ~/.claude/CLAUDE.md
	@echo "✓ Created ~/.claude/CLAUDE.md"
	@if [ ! -f ~/.claude/settings.json ]; then \
		echo "Creating ~/.claude/settings.json..."; \
		printf '%s\n' \
			'{' \
			'  "permissions": {' \
			'    "allow": [' \
			'      "Bash(clog:*)"' \
			'    ],' \
			'    "deny": [],' \
			'    "ask": []' \
			'  }' \
			'}' \
			> ~/.claude/settings.json; \
		echo "✓ Created ~/.claude/settings.json with clog permissions"; \
	else \
		if ! grep -q "Bash(clog:\*)" ~/.claude/settings.json 2>/dev/null; then \
			echo "⚠ Warning: ~/.claude/settings.json exists but may not include 'Bash(clog:*)' permission"; \
			echo "  Please manually add \"Bash(clog:*)\" to the allow list in ~/.claude/settings.json"; \
		else \
			echo "✓ ~/.claude/settings.json already configured with clog permissions"; \
		fi \
	fi

# Build the binary
build:
	@echo "Building clog with configuration..."
	@echo ""
	@echo "=== NATS Configuration ==="
	@echo ""
	@read -p "NATS URL [default: nats://localhost:4222]: " NATS_URL; \
	NATS_URL=$${NATS_URL:-nats://localhost:4222}; \
	echo ""; \
	echo "Authentication types:"; \
	echo "  1) none           - No authentication"; \
	echo "  2) userpass       - Username and password"; \
	echo "  3) token          - Token-based authentication"; \
	echo "  4) nkey           - NKey authentication"; \
	echo "  5) decentralized  - Decentralized (JWT + Seed)"; \
	echo "  6) creds-file     - Extract from NATS credentials file"; \
	echo ""; \
	read -p "Select authentication type [1-6, default: 1]: " AUTH_CHOICE; \
	AUTH_CHOICE=$${AUTH_CHOICE:-1}; \
	case $$AUTH_CHOICE in \
		1) AUTH_TYPE="none"; USERNAME=""; PASSWORD=""; TOKEN=""; NKEY=""; JWT=""; SEED=""; ;; \
		2) AUTH_TYPE="userpass"; \
		   read -p "Username: " USERNAME; \
		   read -sp "Password: " PASSWORD; echo ""; \
		   TOKEN=""; NKEY=""; JWT=""; SEED=""; ;; \
		3) AUTH_TYPE="token"; \
		   read -sp "Token: " TOKEN; echo ""; \
		   USERNAME=""; PASSWORD=""; NKEY=""; JWT=""; SEED=""; ;; \
		4) AUTH_TYPE="nkey"; \
		   read -p "NKey: " NKEY; \
		   USERNAME=""; PASSWORD=""; TOKEN=""; JWT=""; SEED=""; ;; \
		5) AUTH_TYPE="decentralized"; \
		   read -p "NATS JWT: " JWT; \
		   read -p "NATS Seed: " SEED; \
		   USERNAME=""; PASSWORD=""; TOKEN=""; NKEY=""; ;; \
		6) AUTH_TYPE="decentralized"; \
		   read -p "Path to NATS credentials file: " CREDS_FILE; \
		   CREDS_FILE=$$(eval echo "$$CREDS_FILE"); \
		   if [ ! -f "$$CREDS_FILE" ]; then \
		       echo "Error: Credentials file not found: $$CREDS_FILE"; \
		       exit 1; \
		   fi; \
		   JWT=$$(sed -n '/-----BEGIN NATS USER JWT-----/,/------END NATS USER JWT------/p' "$$CREDS_FILE" | grep -v "BEGIN\|END" | tr -d '\n'); \
		   SEED=$$(sed -n '/-----BEGIN USER NKEY SEED-----/,/------END USER NKEY SEED------/p' "$$CREDS_FILE" | grep -v "BEGIN\|END" | tr -d '\n'); \
		   if [ -z "$$JWT" ] || [ -z "$$SEED" ]; then \
		       echo "Error: Failed to extract JWT or Seed from credentials file"; \
		       exit 1; \
		   fi; \
		   echo "✓ Successfully extracted JWT and Seed from credentials file"; \
		   USERNAME=""; PASSWORD=""; TOKEN=""; NKEY=""; ;; \
		*) echo "Invalid choice, defaulting to 'none'"; \
		   AUTH_TYPE="none"; USERNAME=""; PASSWORD=""; TOKEN=""; NKEY=""; JWT=""; SEED=""; ;; \
	esac; \
	echo ""; \
	echo "Backing up main.go..."; \
	cp cmd/main.go cmd/main.go.bak; \
	echo "Injecting configuration into code..."; \
	sed -i.tmp "s|defaultNATSURL  = \".*\"|defaultNATSURL  = \"$$NATS_URL\"|" cmd/main.go; \
	sed -i.tmp "s|defaultAuthType = \".*\" // none, userpass, token, nkey, decentralized|defaultAuthType = \"$$AUTH_TYPE\" // none, userpass, token, nkey, decentralized|" cmd/main.go; \
	sed -i.tmp "s|defaultUsername = \".*\"|defaultUsername = \"$$USERNAME\"|" cmd/main.go; \
	sed -i.tmp "s|defaultPassword = \".*\"|defaultPassword = \"$$PASSWORD\"|" cmd/main.go; \
	sed -i.tmp "s|defaultToken    = \".*\"|defaultToken    = \"$$TOKEN\"|" cmd/main.go; \
	sed -i.tmp "s|defaultNKey     = \".*\"|defaultNKey     = \"$$NKEY\"|" cmd/main.go; \
	sed -i.tmp "s|defaultNATSJWT  = \".*\"|defaultNATSJWT  = \"$$JWT\"|" cmd/main.go; \
	sed -i.tmp "s|defaultNATSSeed = \".*\"|defaultNATSSeed = \"$$SEED\"|" cmd/main.go; \
	rm -f cmd/main.go.tmp; \
	echo "Building binary..."; \
	go build -o clog ./cmd/main.go; \
	echo "Restoring template placeholders in main.go..."; \
	mv cmd/main.go.bak cmd/main.go; \
	echo ""; \
	echo "✓ Binary built: ./clog (with configuration baked in)"; \
	echo "✓ Source file restored to template state"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f clog
	@echo "✓ Clean complete"

# Install to user's local bin directory (recommended)
install: build
	@echo "Installing clog to ~/bin..."
	@mkdir -p ~/bin
	@cp clog ~/bin/clog
	@chmod a+x ~/bin/clog
	@echo "✓ Installed to ~/bin/clog with executable permissions"
	@echo ""
	@echo "IMPORTANT: Ensure ~/bin is in your PATH"
	@echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
	@echo '  export PATH="$$HOME/bin:$$PATH"'
	@echo ""
	@echo "Then reload your shell or run: source ~/.bashrc (or ~/.zshrc)"

# Install to system-wide location (requires sudo)
install-system: build
	@echo "Installing clog to system-wide location..."
	@echo ""
	@if [ "$$(uname)" = "Darwin" ]; then \
		read -p "Install path [default: /usr/local/bin]: " INSTALL_PATH; \
		INSTALL_PATH=$${INSTALL_PATH:-/usr/local/bin}; \
	else \
		INSTALL_PATH=/usr/local/bin; \
	fi; \
	echo "Installing to $$INSTALL_PATH..."; \
	sudo mkdir -p $$INSTALL_PATH; \
	sudo cp clog $$INSTALL_PATH/clog; \
	sudo chmod a+x $$INSTALL_PATH/clog; \
	echo "✓ Installed to $$INSTALL_PATH/clog with executable permissions"

# Run go mod tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

# Run unit tests
unittest:
	@echo "Running unit tests..."
	@go test -v ./cmd/
	@echo ""
	@echo "✓ All tests passed"

# Test the binary with example commands
test: build
	@echo "Testing clog..."
	@./clog -h
	@echo ""
	@echo "Test completed. Try running:"
	@echo "  ./clog -type=task -state=in_progress -message=\"Test message\" -session=\"test-123\""
