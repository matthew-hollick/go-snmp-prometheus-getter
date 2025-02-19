#!/bin/bash

# Create admin user
superset fab create-admin \
    --username hedgehog_admin \
    --firstname Hedgehog \
    --lastname Admin \
    --email hedgehog_admin@hedgehog.internal \
    --password "${SUPERSET_ADMIN_PASSWORD:-hedgehog_admin_password}"

# Initialize the database
superset db upgrade

# Create default roles and permissions
superset init

# Create Elasticsearch database connection
superset set-database-uri \
    --database-name "Elasticsearch Metrics" \
    --uri "elasticsearch+http://${ELASTICSEARCH_HOST}:${ELASTICSEARCH_PORT}/?fetch_size=10000&http_compress=True"
