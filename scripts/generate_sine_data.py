import math
import json
from datetime import datetime, timedelta
import subprocess

# Constants
AMPLITUDE = 100
PERIOD = 3600  # 1 hour in seconds
TOTAL_HOURS = 24
POINTS_PER_SECOND = 1

# Create index mapping
index_mapping = {
    "mappings": {
        "properties": {
            "timestamp": {"type": "date"},
            "value": {"type": "float"},
            "metadata": {
                "properties": {
                    "type": {"type": "keyword"},
                    "period_seconds": {"type": "integer"}
                }
            }
        }
    },
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    }
}

# Create index
subprocess.run([
    "curl", "-X", "PUT",
    "http://localhost:9200/test_sin",
    "-H", "Content-Type: application/json",
    "-d", json.dumps(index_mapping)
])

# Generate data
start_time = datetime.now()
bulk_data = []

for second in range(TOTAL_HOURS * 3600):
    timestamp = start_time + timedelta(seconds=second)
    # Calculate sine wave value between -100 and 100
    value = AMPLITUDE * math.sin(2 * math.pi * second / PERIOD)
    
    # Create Elasticsearch bulk action
    bulk_data.append(json.dumps({
        "index": {"_index": "test_sin"}
    }))
    bulk_data.append(json.dumps({
        "timestamp": timestamp.isoformat(),
        "value": value,
        "metadata": {
            "type": "sine_wave",
            "period_seconds": PERIOD
        }
    }))
    
    # Upload in batches of 1000 documents
    if len(bulk_data) >= 2000:
        bulk_data_str = "\n".join(bulk_data) + "\n"
        subprocess.run([
            "curl", "-X", "POST",
            "http://localhost:9200/_bulk",
            "-H", "Content-Type: application/x-ndjson",
            "-d", bulk_data_str
        ])
        bulk_data = []

# Upload any remaining documents
if bulk_data:
    bulk_data_str = "\n".join(bulk_data) + "\n"
    subprocess.run([
        "curl", "-X", "POST",
        "http://localhost:9200/_bulk",
        "-H", "Content-Type: application/x-ndjson",
        "-d", bulk_data_str
    ])

print("Data generation complete")
