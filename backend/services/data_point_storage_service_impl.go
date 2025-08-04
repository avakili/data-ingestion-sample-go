package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/avakili/data-ingestion/backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dataPointDb struct {
	DataPointId string    `gorm:"column:data_point_id;type:varchar(255);primaryKey"`
	DeviceId    string    `gorm:"column:device_id;type:varchar(255);not null"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;not null"`
	DataPayload string    `gorm:"column:data_payload;type:text;not null"`
}

// Overriding the table name:
func (dataPointDb) TableName() string {
	return "data_points"
}

func (d *dataPointDb) ToDatapoint() models.DataPoint {
	var dataPayload map[string]interface{}
	err := json.Unmarshal([]byte(d.DataPayload), &dataPayload)
	if err != nil {
		return models.DataPoint{}
	}

	return models.DataPoint{
		DataPointId: d.DataPointId,
		AddDataPointRequest: models.AddDataPointRequest{
			DeviceId:    d.DeviceId,
			Timestamp:   d.Timestamp,
			DataPayload: dataPayload,
		},
	}
}

type DataPointStorageServiceImpl struct {
	db *gorm.DB
}

func NewDataPointStorageServiceImpl(db *gorm.DB) *DataPointStorageServiceImpl {
	return &DataPointStorageServiceImpl{db: db}
}

func (s *DataPointStorageServiceImpl) SaveDataPoint(datapoint models.AddDataPointRequest) (string, error) {
	dataPointId := uuid.New().String()

	// Convert the data payload to a JSON string
	jsonData, err := json.Marshal(datapoint.DataPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data payload: %w", err)
	}

	result := s.db.Create(&dataPointDb{
		DataPointId: dataPointId,
		DeviceId:    datapoint.DeviceId,
		Timestamp:   datapoint.Timestamp,
		DataPayload: string(jsonData),
	})

	if result.Error != nil {
		return "", fmt.Errorf("failed to save datapoint: %w", result.Error)
	}

	return dataPointId, nil
}

func (s *DataPointStorageServiceImpl) GetDataPointsForDeviceId(deviceId string) ([]models.DataPoint, error) {
	var dataPointDbRows []dataPointDb
	result := s.db.Where("device_id = ?", deviceId).Find(&dataPointDbRows)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get data points for device id %s: %w", deviceId, result.Error)
	}

	resultDataPoints := make([]models.DataPoint, len(dataPointDbRows))
	for i, d := range dataPointDbRows {
		resultDataPoints[i] = d.ToDatapoint()
	}

	return resultDataPoints, nil
}
