# Project Backlog

## Monitoring and Metrics
- [ ] Implement service metrics endpoint
- [ ] Add scrape duration tracking
- [ ] Add success/failure rate metrics
- [ ] Add Elasticsearch write latency metrics
- [ ] Add error rate monitoring

## Instrumentation
- [ ] Add tracing support
- [ ] Add detailed performance metrics
- [ ] Implement health check endpoints

## Error Handling Improvements
- [ ] Implement sophisticated backoff mechanism
- [ ] Add configurable retry conditions
- [ ] Add circuit breaker pattern
- [ ] Add error rate monitoring

## Development Infrastructure
- [ ] Set up CI pipeline
- [ ] Add example Prometheus exporter configurations
- [ ] Create development environment setup scripts
- [ ] Add automated integration tests
- [ ] Set up pre-commit hooks
- [ ] Add devbox.json for development tools
- [ ] Identify and mark unused files
  - Review all files in the codebase
  - Check for files not referenced by imports or build system
  - Rename unused files with ".DELETEME" extension
  - Document findings for team review
  - Remove files after team approval
- [ ] Implement proper SNMP simulator healthcheck
  - Research best practices for SNMP service health monitoring
  - Consider alternatives to snmpget for health checks
  - Options to explore:
    - Using netcat to check UDP port availability
    - Creating a lightweight health endpoint
    - Building a custom health check script
  - Implement and test chosen solution

## Documentation
- [ ] Add API documentation
- [ ] Add configuration guide
- [ ] Add troubleshooting guide
- [ ] Add development setup guide
- [ ] Add architecture diagrams
- [ ] Add deployment guide

## Security Enhancements
- [ ] Add support for additional authentication methods
- [ ] Implement credential rotation
- [ ] Add audit logging
- [ ] Add security hardening guide

## Performance Optimizations
- [ ] Optimize metric batching
- [ ] Implement metric buffering
- [ ] Add performance tuning guide
- [ ] Optimize memory usage

## Data Schema and Storage
- [ ] Implement ECS-aligned JSON schema for SNMP data
  - Core fields to implement:
    - host.* (device details, hostname, IP)
    - metrics.* (SNMP metric values)
    - event.* (collection metadata, timestamps)
    - observer.* (SNMP collector information)
    - network.* (interface metrics, network stats)
  - Tasks:
    - Create JSON schema definition
    - Implement transformation from Prometheus to ECS format
    - Add validation for schema compliance
    - Document field mappings and usage
    - Create example queries for common use cases

## SNMP Exporter Configuration
- Set up SNMP exporter generator with proper configuration
  - Investigate issues with generator configuration format
- [ ] Configure SNMP simulator to emulate a network switch
  - Research and obtain Cisco network switch MIB files
  - Create SNMP record file with switch metrics:
    - Interface statistics (in/out bytes, errors)
    - Port status
    - System resources (CPU, memory)
    - Temperature sensors
  - Update SNMP exporter configuration for switch metrics
  - Add documentation for switch metric collection
  - Test with standard network monitoring queries
  - Ensure MIBs are properly loaded and accessible
  - Test generated configuration with different SNMP devices
  - Document the generator process for future updates

## Feature Requests
- [ ] Add support for metric filtering
- [ ] Add support for metric transformation
- [ ] Add support for metric aggregation
- [ ] Add support for alerting
