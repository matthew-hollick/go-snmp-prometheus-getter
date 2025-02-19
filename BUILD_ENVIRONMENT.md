# Setting Up Your Workspace

This project uses [mise](https://mise.jdx.dev/) to manage tool versions and [uv](https://github.com/astral-sh/uv) for Python package management.

## Getting Started

1. Install `mise` (if you don't have it):
   ```bash
   curl https://mise.run | sh
   ```

2. Get the code:
   ```bash
   git clone <repository-url>
   cd go-snmp-prometheus-getter
   ```

3. Run the setup task:
   ```bash
   mise run setup
   ```
   This will:
   - Download all needed files
   - Check file integrity
   - Build the project

4. Start the services:
   ```bash
   mise run services-start
   ```

## Common Tasks

### Managing Dependencies
- `mise run deps-tidy` - Clean up the project's dependencies
- `mise run deps-verify` - Check dependencies haven't changed
- `mise run deps-update` - Get latest versions of dependencies
- `mise run deps-download` - Download all dependencies
- `mise run deps-clean` - Remove old dependencies

### Building and Testing
- `mise run build` - Build the project
- `mise run test` - Run all tests
- `mise run test-coverage` - Run tests and check coverage
- `mise run show-coverage` - Show coverage in your web browser

### Service Management
- `mise run services-start` - Start all services
- `mise run services-stop` - Stop all services
- `mise run services-logs` - View service information
- `mise run services-status` - Check if services are running

## Development Tools

### Code Quality Tools
- Go code checker: Looks for common mistakes
- Security checker: Checks for security issues

These run automatically when you save your work, as set up in `.pre-commit-config.yaml`.

### Automatic Checks
These checks run before saving your work:
- Code formatting
- Code quality checks
- Security checks
- File formatting
- Whitespace cleanup
- File ending fixes

To run checks manually:
```bash
uv run pre-commit run --all-files
```

## Updating Tools

To update tool versions:
```bash
mise upgrade
```

To update Python packages:
```bash
uv sync --upgrade
```

To update the automatic checks:
```bash
uv run pre-commit autoupdate
```

## Solving Common Problems

### Connection Issues

1. Device Connection Problems
   ```bash
   # Check if test device is running
   mise run device-status
   
   # Test connection
   mise run device-test
   ```

2. Database Connection Problems
   ```bash
   # Check database status
   mise run db-status
   
   # Test settings
   mise run db-test-config
   ```

3. Build Problems
   ```bash
   # Clean old files
   mise run clean
   
   # Start fresh
   mise run rebuild
   ```

### Finding Problems

- Service information: `mise run services-logs`
- Extra device information: Set `DEVICE_DEBUG=1` before starting
- Database information: Set `DB_DEBUG=1`

For more help, see the guides in `docs/`.
