#!/bin/bash
set -e

# Wait for Elasticsearch to be ready
until curl -s "http://elasticsearch:9200/_cluster/health" > /dev/null; do
    echo "Waiting for Elasticsearch..."
    sleep 5
done

# Create the configuration index
curl -X PUT "http://elasticsearch:9200/snmp-config" \
     -H "Content-Type: application/json" \
     -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { "type": "keyword" },
      "description": { "type": "text" },
      "snmp_settings": {
        "properties": {
          "host": { "type": "keyword" },
          "port": { "type": "integer" },
          "version": { "type": "keyword" },
          "community": { "type": "keyword" },
          "oids": {
            "type": "nested",
            "properties": {
              "oid": { "type": "keyword" },
              "name": { "type": "keyword" },
              "description": { "type": "text" },
              "type": { "type": "keyword" }
            }
          }
        }
      },
      "prometheus_settings": {
        "properties": {
          "metric_prefix": { "type": "keyword" },
          "labels": { "type": "object" }
        }
      },
      "enabled": { "type": "boolean" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

# Add a test configuration
curl -X POST "http://elasticsearch:9200/snmp-config/_doc/test1" \
     -H "Content-Type: application/json" \
     -d '{
  "id": "test1",
  "name": "Test SNMP Device",
  "description": "Test device for verifying SNMP exporter connectivity",
  "snmp_settings": {
    "host": "snmp-simulator",
    "port": 161,
    "version": "2c",
    "community": "public",
    "oids": [
      {
        "oid": "1.3.6.1.2.1.1.1.0",
        "name": "sysDescr",
        "description": "System Description",
        "type": "string"
      }
    ]
  },
  "prometheus_settings": {
    "metric_prefix": "snmp_test",
    "labels": {
      "device": "test_device"
    }
  },
  "enabled": true,
  "created_at": "'"$(date -u +"%Y-%m-%dT%H:%M:%SZ")"'",
  "updated_at": "'"$(date -u +"%Y-%m-%dT%H:%M:%SZ")"'"
}'

echo "Test configuration created successfully"
