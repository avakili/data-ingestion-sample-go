package routes

import (
	"github.com/avakili/data-ingestion/backend/services"
	"github.com/gin-gonic/gin"
)

func DataPointRoutes(r *gin.Engine, dataStorage services.DataPointStorageService) {
	controller := services.NewDataPointController(dataStorage)
	datapointGroup := r.Group("/data_point")
	{
		datapointGroup.POST("", controller.AddDatapoint)
		datapointGroup.GET("", controller.GetDataPointsForDeviceId)
	}
}
