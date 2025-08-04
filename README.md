# Sample Go project

This repo contains sample service in Go that uses Postgres for storage.
It provides
the following routes:
* `POST /data_point`: to add a new data_point. It returns the UUID of the stored data point.
* `GET /data_point?device_id=...` allows retrieving all data points associated with a device.

## How to run

```bash
docker-compose up -d
```

The service would be exposed on [http://localhost:8080](http://localhost:8080).

The following adds a sample data point for device abc123 

```bash
curl -X POST -H "Content-Type: application/json" -d '{"device_id": "abc123", "timestamp": "2025-05-01T13:45:30Z", "data_payload": {"gene_count": 20456, "sample_quality": 98.6}}' http://localhost:8080/data_point
```

The following retrieves data points for device `abc123`:
```bash
curl http://localhost:8080/data_point?device_id=abc123
```

## Disclosure
* I've used Cursor to generate unit tests.