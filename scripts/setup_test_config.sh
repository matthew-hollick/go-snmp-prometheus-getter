#!/bin/bash
set -e

# Create the configuration index
curl -X PUT "http://elasticsearch:9200/service_configuration" \
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

# Add the example configuration
curl -X POST "http://elasticsearch:9200/service_configuration/_doc" \
     -H "Content-Type: application/json" \
     -d @/config/elasticsearch_device1_config.json

echo "Configuration setup completed successfully"
