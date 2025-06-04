package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/models"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")

	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}
	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	err = event.RegisterNewEvent(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register new event."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "New event registered."})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")

	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}
	var event models.Event
	event.ID = eventID

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to cancel event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event has been cancelled"})
}
