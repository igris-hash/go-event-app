package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/models"
)

func main() {
	server := gin.Default()
	server.GET("/events", getEvents)
	server.POST("/events", createEvent)
	server.Run(":8000")
}

func getEvents(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Hello world"})
}

func createEvent(context *gin.Context) {
	var request models.Event
	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	request.ID = int(time.Millisecond)
	request.UserID = 1
	request.Save()
	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": request})
}
