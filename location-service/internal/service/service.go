package service

import (
	"fmt"
	"ryde/internal/data"
	"ryde/internal/models"
	"ryde/internal/utils"

	"github.com/gin-gonic/gin"
)

type LocationService struct {
	DataStore *data.DataStore
}

func NewLocationService(dataStore *data.DataStore) *LocationService {
	return &LocationService{
		DataStore: dataStore,
	}
}

func (s *LocationService) UpdateDriverLocation(c *gin.Context) {
	utils.PollLocation(c)

	for update := range utils.UpdateChannel {
		if err := s.DataStore.UpdateDriverLocation(c.Request.Context(), update.DriverID, update.Latitude, update.Longitude); err != nil {
			fmt.Println("Error updating location:", err)
			continue
		}
	}
}

func (s *LocationService) GetDriverLocation(c *gin.Context, driverID string) (*models.Location, error) {
	return s.DataStore.GetDriverLocation(c, driverID)
}
