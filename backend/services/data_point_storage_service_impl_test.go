package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/avakili/data-ingestion/backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	// Create the table schema
	err = db.AutoMigrate(&dataPointDb{})
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func TestSaveDataPoint_Success(t *testing.T) {
	db := setupTestDB(t)
	service := NewDataPointStorageServiceImpl(db)

	payload := map[string]interface{}{
		"temperature": 22.5,
		"humidity":    60.9,
	}
	req := models.AddDataPointRequest{
		DeviceId:    "device123",
		Timestamp:   time.Now().UTC().Truncate(time.Second),
		DataPayload: payload,
	}

	id, err := service.SaveDataPoint(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Check that the data was actually saved
	var dp dataPointDb
	result := db.First(&dp, "data_point_id = ?", id)
	assert.NoError(t, result.Error)
	assert.Equal(t, req.DeviceId, dp.DeviceId)
	assert.WithinDuration(t, req.Timestamp, dp.Timestamp, time.Second)

	var savedPayload map[string]interface{}
	err = json.Unmarshal([]byte(dp.DataPayload), &savedPayload)
	assert.NoError(t, err)
	assert.Equal(t, payload["temperature"], savedPayload["temperature"])
	assert.Equal(t, payload["humidity"], savedPayload["humidity"])
}

func TestSaveDataPoint_InvalidPayload(t *testing.T) {
	db := setupTestDB(t)
	service := NewDataPointStorageServiceImpl(db)

	// DataPayload contains a value that cannot be marshaled to JSON (e.g., a channel)
	payload := map[string]interface{}{
		"bad": make(chan int),
	}
	req := models.AddDataPointRequest{
		DeviceId:    "device123",
		Timestamp:   time.Now(),
		DataPayload: payload,
	}

	id, err := service.SaveDataPoint(req)
	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestGetDataPointsForDeviceId_Empty(t *testing.T) {
	db := setupTestDB(t)
	service := NewDataPointStorageServiceImpl(db)

	points, err := service.GetDataPointsForDeviceId("nonexistent")
	assert.NoError(t, err)
	assert.Empty(t, points)
}

func TestGetDataPointsForDeviceId_Success(t *testing.T) {
	db := setupTestDB(t)
	service := NewDataPointStorageServiceImpl(db)

	payload := map[string]interface{}{
		"temperature": 25.0,
	}
	req := models.AddDataPointRequest{
		DeviceId:    "deviceABC",
		Timestamp:   time.Now().UTC().Truncate(time.Second),
		DataPayload: payload,
	}

	id, err := service.SaveDataPoint(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	points, err := service.GetDataPointsForDeviceId("deviceABC")
	assert.NoError(t, err)
	assert.Len(t, points, 1)
	assert.Equal(t, id, points[0].DataPointId)
	assert.Equal(t, req.DeviceId, points[0].DeviceId)
	assert.WithinDuration(t, req.Timestamp, points[0].Timestamp, time.Second)
	assert.Equal(t, payload["temperature"], points[0].DataPayload["temperature"])
}

func TestToDatapoint_InvalidJSON(t *testing.T) {
	// Simulate a dataPointDb with invalid JSON in DataPayload
	dp := dataPointDb{
		DataPointId: "id1",
		DeviceId:    "dev1",
		Timestamp:   time.Now(),
		DataPayload: "{invalid json",
	}
	result := dp.ToDatapoint()
	assert.Empty(t, result.DataPointId)
	assert.Empty(t, result.DeviceId)
	assert.Nil(t, result.DataPayload)
}
