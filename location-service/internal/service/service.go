package service

import (
	"fmt"
	"ryde/internal/data"
	"ryde/internal/models"
	"ryde/internal/utils"

	"github.com/gin-gonic/gin"
)

type LocationService struct {
	LocationStore *data.LocationStore
}

func NewLocationService(dataStore *data.LocationStore) *LocationService {
	return &LocationService{
		LocationStore: dataStore,
	}
}

func (s *LocationService) UpdateDriverLocation(c *gin.Context) {
	utils.PollLocation(c)

	for update := range utils.UpdateChannel {
		if err := s.LocationStore.UpdateDriverLocation(c.Request.Context(), update.DriverID, update.Latitude, update.Longitude); err != nil {
			fmt.Println("Error updating location:", err)
			continue
		}
	}
}

func (s *LocationService) GetDriverLocation(c *gin.Context, driverID string) (*models.Location, error) {
	return s.LocationStore.GetDriverLocation(c.Request.Context(), driverID)
}

func (s *LocationService) FindNearbyDrivers(c *gin.Context, lat, lon, radius float64) ([]string, error) {
	return s.LocationStore.FindNearbyDrivers(c.Request.Context(), lat, lon, radius)
}
