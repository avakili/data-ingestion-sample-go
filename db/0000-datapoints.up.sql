CREATE TABLE data_points (
    data_point_id VARCHAR(255) PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    data_payload TEXT NOT NULL
);

CREATE INDEX idx_device_id ON data_points (device_id);
