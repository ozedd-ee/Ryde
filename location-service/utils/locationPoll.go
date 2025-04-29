package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ryde/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

var UpdateChannel = make(chan models.Location, 10)

func PollLocation(c *gin.Context) {
	token := c.Query("token")
	claims, err := ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error", err)
		return
	}
	defer conn.Close()
	defer close(UpdateChannel)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error", err)
			break
		}
		user_id := claims.UserID
		var update struct {
			lat float64
			lon float64
		}

		if err := json.Unmarshal(message, &update); err != nil {
			fmt.Println("JSON Parse error", err)
			continue
		}
		var location models.Location
		location.DriverID = user_id
		location.Latitude = update.lat
		location.Longitude = update.lon

		UpdateChannel <- location
	}
}
