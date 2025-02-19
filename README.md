# Elasticsearch Data Writer

A Go-based service for writing data to Elasticsearch.

## Project Requirements

### Core Functionality
- Scrape metrics from multiple Prometheus exporters
- Write data to Elasticsearch using the official Go client
- Read and validate configuration from a configuration file
- Validate Elasticsearch service using certificate hash
- Support structured data writing based on configurable schema
- Support configurable concurrent scraping and writing
- Implement simple retry mechanism for error handling

### Data Flow
1. Retrieve configuration from Elasticsearch
2. For each configured Prometheus exporter:
   - Validate connection and certificate
   - Scrape metrics using provided parameters
   - Transform metrics to match schema
3. Write transformed data to Elasticsearch

### Configuration Requirements
The service requires two configuration sources:

1. Initial Bootstrap Configuration File:
   - Elasticsearch connection details for configuration retrieval
   - Authentication credentials
   - Certificate hash for service validation
   - Configuration index name
   - Log level settings
   - Concurrency settings
   - Simple retry configuration
   - Configuration refresh period

2. Main Configuration (stored in Elasticsearch):
   - Target Elasticsearch connection details
   - Authentication credentials for configuration retrieval
   - Target index names:
     - Data write index
     - Configuration index
   - Write operation credentials
   - Certificate hash for service validation
   - Data schema configuration
   - List of Prometheus exporters, each with:
     - Unique identifier
     - URL endpoint
     - List of query parameters
     - Certificate hash for validation
     - Scrape interval
     - Optional metric prefix
     - Optional labels to add

### Error Handling
- All errors logged to stdout with configurable log level
- Simple retry mechanism with configurable:
  - Maximum attempts
  - Delay between attempts

### Configuration Updates
- Bootstrap configuration read once at startup
- Main configuration checked periodically (configurable interval)
- Configuration changes validated before applying
- Changes applied without service restart

## Development Setup

### Testing
- Unit tests for configuration parsing
- Integration tests with mock Prometheus exporters
- End-to-end tests with containerized services

### Local Development
- Docker Compose configuration for local services
- Mock Prometheus exporters for testing
- Local Elasticsearch instance

### Prerequisites
- Go 1.21 or later
- Docker and Docker Compose for local development
- Make for build automation

## Security Notes
- Configuration file should have restricted permissions
- Credentials should be stored securely
- Certificate validation must be implemented
- Basic authentication support for Prometheus exporters

## Network Device Metrics Collector

A service that collects metrics from network devices using SNMP and stores them in Elasticsearch.

## SNMP Testing

We use the Python-based `snmpsim` package for SNMP simulation during development and testing. The simulator configuration and data files are in the `/snmpsim` directory.

### Running the Simulator

```bash
# Using Docker
cd snmpsim
docker build -t snmpsim .
docker run -p 11161:11161/udp snmpsim

# Test the simulator
snmpget -v2c -c public localhost:11161 1.3.6.1.2.1.1.1.0
```

The simulator provides:
- System information (name, location, contact)
- Interface statistics
- Resource metrics (CPU, memory)

See `/snmpsim/README.md` for more details about the simulated device.

Note: There is an experimental Go-based simulator in `/internal/experimental/simulator`, but it is not currently in use.

## Development

### Prerequisites

1. Go 1.21 or later
2. Access to Elasticsearch
3. Root privileges for SNMP simulator (ports below 1024)

### Building

```bash
# Build all components
mise run build

# Build specific component
mise run build-simulator
mise run build-collector
```

### Testing

```bash
# Run all tests
mise run test

# Run specific component tests
mise run test-simulator
mise run test-collector
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests
4. Submit a pull request

## License

Copyright (c) 2025 Hedgehog Analytics. All rights reserved.

See [BACKLOG.md](BACKLOG.md) for planned features and improvements.
