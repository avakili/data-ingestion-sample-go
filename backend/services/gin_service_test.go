package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/avakili/data-ingestion/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of DataPointStorageService
type mockDataPointStorageService struct {
	SaveDataPointFunc            func(dataPoint models.AddDataPointRequest) (string, error)
	GetDataPointsForDeviceIdFunc func(deviceId string) ([]models.DataPoint, error)
}

func (m *mockDataPointStorageService) SaveDataPoint(dataPoint models.AddDataPointRequest) (string, error) {
	return m.SaveDataPointFunc(dataPoint)
}

func (m *mockDataPointStorageService) GetDataPointsForDeviceId(deviceId string) ([]models.DataPoint, error) {
	return m.GetDataPointsForDeviceIdFunc(deviceId)
}

func TestAddDatapoint_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockDataPointStorageService{
		SaveDataPointFunc: func(dataPoint models.AddDataPointRequest) (string, error) {
			return "test-id-123", nil
		},
	}
	controller := NewDataPointController(mockService)

	payload := map[string]interface{}{
		"device_id":    "dev1",
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"data_payload": map[string]interface{}{"temperature": 22.5},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/datapoints", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.AddDatapoint(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "test-id-123", resp["data_point_id"])
}

func TestAddDatapoint_SaveError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockDataPointStorageService{
		SaveDataPointFunc: func(dataPoint models.AddDataPointRequest) (string, error) {
			return "", errors.New("save failed")
		},
	}
	controller := NewDataPointController(mockService)

	payload := map[string]interface{}{
		"device_id":    "dev1",
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"data_payload": map[string]interface{}{"temperature": 22.5},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/datapoints", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.AddDatapoint(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "save failed")
}

func TestControllerGetDataPointsForDeviceId_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expectedPoints := []models.DataPoint{
		{
			DataPointId: "id1",
			AddDataPointRequest: models.AddDataPointRequest{
				DeviceId:    "dev1",
				Timestamp:   time.Now().UTC().Truncate(time.Second),
				DataPayload: map[string]interface{}{"temperature": 22.5},
			},
		},
	}
	mockService := &mockDataPointStorageService{
		GetDataPointsForDeviceIdFunc: func(deviceId string) ([]models.DataPoint, error) {
			assert.Equal(t, "dev1", deviceId)
			return expectedPoints, nil
		},
	}
	controller := NewDataPointController(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/datapoints?device_id=dev1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Request.URL.RawQuery = "device_id=dev1"

	controller.GetDataPointsForDeviceId(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	points, ok := resp["data_points"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, points, 1)
	point := points[0].(map[string]interface{})
	assert.Equal(t, "id1", point["data_point_id"])
	assert.Equal(t, "dev1", point["device_id"])
}

func TestGetDataPointsForDeviceId_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockDataPointStorageService{
		GetDataPointsForDeviceIdFunc: func(deviceId string) ([]models.DataPoint, error) {
			return nil, errors.New("db error")
		},
	}
	controller := NewDataPointController(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/datapoints?device_id=dev1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Request.URL.RawQuery = "device_id=dev1"

	controller.GetDataPointsForDeviceId(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "db error")
}
