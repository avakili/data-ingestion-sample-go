package models

import "time"

type AddDataPointRequest struct {
	DeviceId    string                 `json:"device_id"`
	Timestamp   time.Time              `json:"timestamp"`
	DataPayload map[string]interface{} `json:"data_payload"`
}

type DataPoint struct {
	DataPointId string `json:"data_point_id"`
	AddDataPointRequest
}
