# Contributing to clog

Thank you for your interest in contributing to clog!

## Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/yourusername/clog.git
   cd clog
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run tests:**
   ```bash
   make unittest
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes to the code

3. Run tests to ensure everything works:
   ```bash
   make unittest
   ```

4. Build and test the binary:
   ```bash
   make build
   make test
   ```

### Code Style

- Follow idiomatic Go conventions
- Run `go fmt` on all code before committing
- Add comments for exported functions and types
- Keep functions focused and testable

### Testing

- Write unit tests for all new functionality
- Ensure existing tests pass before submitting
- Aim for good test coverage
- Use table-driven tests where appropriate

### Submitting Changes

1. Commit your changes with descriptive commit messages:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```

2. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create a Pull Request on GitHub

## Security Considerations

### Credential Handling

- **NEVER commit actual NATS credentials** to the repository
- The `make build` process is designed to keep credentials out of source code
- `.gitignore` is configured to prevent credential files from being committed
- Review your changes before committing to ensure no secrets are included

### Testing with Credentials

When testing with real credentials:

1. Use environment variables:
   ```bash
   export NATS_URL="your-test-url"
   export NATS_CREDS="/path/to/test-creds"
   ```

2. Or use a local `.env` file (which is gitignored):
   ```bash
   NATS_URL=nats://localhost:4222
   NATS_CREDS=/path/to/creds
   ```

## Project Structure

```
clog/
├── cmd/
│   ├── main.go       # Main application code
│   └── main_test.go  # Unit tests
├── Makefile          # Build and development tasks
├── README.md         # User documentation
├── CONTRIBUTING.md   # This file
└── go.mod            # Go module definition
```

## Makefile Targets

- `make build` - Build the binary with interactive credential prompts
- `make unittest` - Run all unit tests
- `make test` - Build and run manual tests
- `make clean` - Remove build artifacts
- `make install` - Install binary to `/usr/local/bin`
- `make tidy` - Run `go mod tidy`

## Questions or Issues?

If you have questions or encounter issues, please open an issue on GitHub.

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
