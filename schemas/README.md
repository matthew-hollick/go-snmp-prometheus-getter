# Configuration Schemas

This directory contains JSON Schema definitions for configuration structures used in the SNMP metrics collector.

## Device Configuration Schema

The `device-config.schema.json` defines the structure of device configuration documents stored in Elasticsearch. Each document represents a single device to be monitored.

### Key Components

1. **Core Device Information**
   - Unique identifier
   - Device name
   - Device type
   - Enabled/disabled status

2. **SNMP Settings**
   - Host and port
   - SNMP version
   - Community string
   - Timeout and retry settings

3. **Collector Settings**
   - Collector hostname
   - Software version
   - SNMP modules
   - Collection interval
   - Metric inclusion/exclusion lists

4. **Tags and Metadata**
   - Environment tag
   - Location tag
   - Role tag
   - Custom tags
   - Creation and update timestamps

### Example Configuration

```json
{
  "id": "switch-01",
  "name": "Core Switch 01",
  "type": "network-device",
  "enabled": true,
  "snmp_settings": {
    "host": "switch-01.network.hedgehog.internal",
    "port": 161,
    "version": "2c",
    "community": "public",
    "timeout": "5s",
    "retries": 3
  },
  "collector_settings": {
    "hostname": "snmp.collector.hedgehog.internal",
    "version": "1.0.0",
    "modules": ["net-snmp"],
    "collection_interval": "5m",
    "metrics": {
      "include": ["sysUpTime", "sysDescr", "ifInOctets", "ifOutOctets"],
      "exclude": []
    }
  },
  "tags": {
    "environment": "production",
    "location": "london",
    "role": "network-switch"
  },
  "created_at": "2025-02-18T22:05:21Z",
  "updated_at": "2025-02-18T22:05:21Z"
}
```

### Validation

You can validate your configuration documents against this schema using tools like:
- [JSON Schema Validator](https://www.jsonschemavalidator.net/)
- [ajv-cli](https://github.com/ajv-validator/ajv-cli)
- Various IDE extensions that support JSON Schema

### Notes

1. All hostnames must use the `.hedgehog.internal` domain
2. Time durations use Go-style format (e.g., "5s", "1m", "2h")
3. Timestamps must be in ISO 8601 format
4. The schema enforces required fields and value constraints
5. Custom tags are allowed under the `tags` object
