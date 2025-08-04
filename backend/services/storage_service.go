package services

import "github.com/avakili/data-ingestion/backend/models"

type DataPointStorageService interface {
	SaveDataPoint(dataPoint models.AddDataPointRequest) (string, error)
	GetDataPointsForDeviceId(deviceId string) ([]models.DataPoint, error)
}
