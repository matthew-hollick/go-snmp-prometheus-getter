# Network Device Information Collector

A service that collects information from network devices and stores it in a searchable database.

## Overview

This service provides:
- Network device information collection using standard protocols
- Information storage in a searchable format
- Integration with industry-standard tools
- Support for common network devices
- Central device management

## How It Works

The service has three main parts:
1. Information Gatherer: Manages device settings and coordinates collection
2. Data Translator: Converts device information into a standard format
3. Test Environment: Provides a safe space for development

For detailed diagrams showing how these parts work together, see `NEXT.md`.

## Development

### Before You Start
- Go programming language version 1.21 or later
- Docker and Docker Compose for running services
- mise for development tools

### Getting Started
1. See `BUILD_ENVIRONMENT.md` for setting up your workspace
2. Review `NEXT.md` for current work items
3. Check `BACKLOG.md` for planned features

### Testing Your Setup
We use the `lextudio/snmpsim` package to create a test environment.

```bash
# Start all services
mise run services-start

# View service information
mise run services-logs

# Stop all services
mise run services-stop
```

## Settings

### Device Settings
Devices are set up using documents in the database. For example:
```json
{
  "id": "switch-01",
  "name": "Test Switch",
  "type": "network-device",
  "enabled": true,
  "connection_settings": {
    "host": "switch.simulator.hedgehog.internal",
    "port": 11161,
    "version": "2c",
    "access_code": "public"
  }
}
```

### Service Settings
See `config.example.toml` for available options.

## Contributing
1. Review `BUILD_ENVIRONMENT.md` for workspace setup
2. Check `NEXT.md` for current tasks
3. See `BACKLOG.md` for planned features

## Security Notes
- Access codes should be treated as sensitive information
- All login details should be stored securely
- Use non-standard ports (11161) for testing
