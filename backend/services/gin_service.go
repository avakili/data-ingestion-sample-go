package services

import (
	"net/http"

	"github.com/avakili/data-ingestion/backend/models"
	"github.com/gin-gonic/gin"
)

type DataPointController struct {
	dataPointStorageService DataPointStorageService
}

func NewDataPointController(dataPointStorageService DataPointStorageService) *DataPointController {
	return &DataPointController{dataPointStorageService: dataPointStorageService}
}

func (controller *DataPointController) AddDatapoint(c *gin.Context) {
	datapointRequest := models.AddDataPointRequest{}
	c.ShouldBindJSON(&datapointRequest)

	dataPointId, err := controller.dataPointStorageService.SaveDataPoint(datapointRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data_point_id": dataPointId})
}

func (controller *DataPointController) GetDataPointsForDeviceId(c *gin.Context) {
	deviceId := c.Query("device_id")

	dataPoints, err := controller.dataPointStorageService.GetDataPointsForDeviceId(deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data_points": dataPoints})
}
