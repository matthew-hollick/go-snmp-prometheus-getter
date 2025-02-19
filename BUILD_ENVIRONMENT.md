# Build Environment Setup

This project uses [mise](https://mise.jdx.dev/) to manage tool versions and [uv](https://github.com/astral-sh/uv) for Python package management.

## Setting Up Your Development Environment

1. Install `mise` (if not already installed):
   ```bash
   curl https://mise.run | sh
   ```

2. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-snmp-prometheus-getter
   ```

3. Run the setup task:
   ```bash
   mise run setup
   ```
   This will:
   - Download all Go dependencies
   - Verify the dependency integrity
   - Build the project

4. Start the development services:
   ```bash
   mise run services-start
   ```

## Common Development Tasks

### Dependency Management
- `mise run deps-tidy` - Tidy up the module's dependencies
- `mise run deps-verify` - Verify dependencies have not been modified
- `mise run deps-update` - Update all dependencies to their latest versions
- `mise run deps-download` - Download all dependencies
- `mise run deps-clean` - Clean the module cache

### Building and Testing
- `mise run build` - Build the project
- `mise run test` - Run all tests
- `mise run test-coverage` - Run tests with coverage
- `mise run show-coverage` - Show coverage report in browser

### Service Management
- `mise run services-start` - Start all services
- `mise run services-stop` - Stop all services
- `mise run services-logs` - View service logs
- `mise run services-status` - Check service status

## Development Tools

### Go Tools
- golangci-lint: Comprehensive Go linter
- gosec: Security checker for Go code

These are configured in the `.pre-commit-config.yaml` file and will run automatically on git commit.

### Pre-commit Hooks
The following checks run automatically before each commit:
- Go formatting (gofmt)
- Go linting (golangci-lint)
- Security checks (gosec)
- YAML formatting
- Trailing whitespace removal
- EOF newline fixes

To run pre-commit checks manually:
```bash
uv run pre-commit run --all-files
```

## Updating Tools

To update tool versions:
```bash
mise upgrade
```

To update Python dependencies:
```bash
uv sync --upgrade
```

To update pre-commit hooks:
```bash
uv run pre-commit autoupdate
