# Next Steps

## SNMP Simulator Improvements
- [ ] Switch to using `lextudio/snmpsim` from https://github.com/lextudio/snmpsim
  - Current simulator may be outdated or incompatible
  - `lextudio/snmpsim` is actively maintained and supports more features
  - Provides better SNMPv2c support
- [ ] Add more OIDs to provide a realistic network switch simulation
- [ ] Configure to listen on port 11161 to avoid network port collisions
- [ ] Ensure SNMPv2c is properly configured
- [ ] Verify community string is set to "public"

## SNMP Exporter Configuration
- [ ] Replace current configuration with the official Prometheus SNMP exporter configuration
  - Reference: https://github.com/prometheus/snmp_exporter/blob/main/snmp.yml
  - Only simplify if required for basic functionality
- [ ] Test all configured metrics

## SNMP Getter Enhancements
- [ ] Validate URL construction for SNMP exporter queries
- [ ] Improve error handling for SNMP exporter responses
- [ ] Add retry logic with exponential backoff

## Code Quality Improvements
- [ ] Fix type mismatches in collector package:
  - Update `SNMPSettings` struct to include missing fields
  - Fix test cases to use correct types
- [ ] Address configuration issues:
  - Fix undefined fields in `BootstrapConfiguration`
  - Update test cases to match new structure
- [ ] Fix Elasticsearch client issues:
  - Correct client initialization in tests
  - Update struct field usage
- [ ] Address security concerns:
  - Set appropriate TLS MinVersion in `service.go`
  - Fix memory aliasing in service package
- [ ] Fix style issues:
  - Apply `gofmt -s` to `service.go`
  - Fix whitespace and cuddling issues:
    - Correct assignment grouping in `exporter/client.go`
    - Fix switch statement placement in `schema/transformer.go`
    - Update return statement placement in multiple files
    - Correct goroutine launches in `service.go`

## Documentation
- [ ] Review and validate architecture diagrams in `ARCHITECTURE.md`
- [ ] Add implementation notes for each component
- [ ] Create troubleshooting guide

## Testing
- [ ] Create integration tests for the full metrics collection flow
- [ ] Add unit tests for URL construction and query parameters
- [ ] Test error scenarios and recovery

## Setup Tasks
- [ ] Recreate `snmp_exporter/snmp.yml` using the generator:
  ```bash
  cd snmp_exporter
  make generate
  # Or download from:
  # https://github.com/prometheus/snmp_exporter/blob/main/snmp.yml
  ```
  Note: This file is not tracked in git due to its size. Each developer should generate it locally.
