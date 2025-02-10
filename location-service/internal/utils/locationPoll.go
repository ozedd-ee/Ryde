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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error", err)
			break
		}

		var loc models.Location
		if err := json.Unmarshal(message, &loc); err != nil {
			fmt.Println("JSON Parse error", err)
			continue
		}
		UpdateChannel <- loc
	}
}
