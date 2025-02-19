# Next Steps

## SNMP Simulator Improvements
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

## Documentation
### Sequence Diagram
```mermaid
sequenceDiagram
    participant ES as Elasticsearch<br/>(TCP/9200)
    participant Getter as SNMP Getter
    participant Exporter as SNMP Exporter<br/>(TCP/9116)
    participant Simulator as SNMP Simulator<br/>(UDP/11161)

    Note over Getter: Bootstrap Phase
    Getter->>ES: HTTP GET /service_configuration/_search
    ES-->>Getter: Return device configurations (JSON)
    
    Note over Getter: Collection Phase
    loop Every collection interval
        Getter->>Exporter: HTTP GET /snmp?target=simulator&module=net-snmp
        Exporter->>Simulator: SNMP v2c GET/WALK (UDP/11161)
        Simulator-->>Exporter: SNMP v2c Response
        Exporter-->>Getter: HTTP 200 (Prometheus format)
        Getter->>ES: HTTP POST /metrics/_doc (JSON)
    end
```

### Component Flowchart
```mermaid
flowchart LR
    ES[(Elasticsearch)]
    Getter[SNMP Getter]
    Exporter[SNMP Exporter]
    Simulator[SNMP Simulator]

    ES <--> |Device configs\nMetrics storage| Getter
    Getter --> |HTTP metrics query| Exporter
    Exporter --> |SNMP v2c\nport 11161| Simulator
```

## Testing
- [ ] Create integration tests for the full metrics collection flow
- [ ] Add unit tests for URL construction and query parameters
- [ ] Test error scenarios and recovery
