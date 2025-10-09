.PHONY: build clean install test unittest

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
	echo ""; \
	read -p "Select authentication type [1-5, default: 1]: " AUTH_CHOICE; \
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

# Install to system (optional)
install: build
	@echo "Installing clog to /usr/local/bin..."
	@sudo cp clog /usr/local/bin/clog
	@echo "✓ Installed to /usr/local/bin/clog"

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
