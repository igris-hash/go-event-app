package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/internal/model"
	"github.com/igris-hash/go-event-app/internal/service"
	"github.com/igris-hash/go-event-app/internal/utils"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	eventService *service.EventService
}

// NewEventHandler creates a new EventHandler instance
func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// CreateEvent handles event creation
// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body model.Event true "Event details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var input model.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequest(c, "Invalid input", err)
		return
	}

	// Set creator ID from authenticated user
	input.CreatorID = c.GetInt("userID")

	event, err := h.eventService.Create(input)
	if err != nil {
		utils.SendInternalError(c, "Failed to create event", err)
		return
	}

	utils.SendCreated(c, "Event created successfully", event)
}

// GetEvent handles retrieving a single event
// @Summary Get event details
// @Description Get details of a specific event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequest(c, "Invalid event ID", err)
		return
	}

	event, err := h.eventService.GetByID(id)
	if err != nil {
		utils.SendNotFound(c, "Event not found")
		return
	}

	utils.SendSuccess(c, "Event retrieved successfully", event)
}

// ListEvents handles retrieving all events
// @Summary List all events
// @Description Get a list of all events
// @Tags events
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	events, err := h.eventService.GetAll()
	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve events", err)
		return
	}

	utils.SendSuccess(c, "Events retrieved successfully", events)
}

// UpdateEvent handles updating an event
// @Summary Update event details
// @Description Update details of a specific event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param event body model.Event true "Event details"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequest(c, "Invalid event ID", err)
		return
	}

	var input model.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequest(c, "Invalid input", err)
		return
	}

	// Verify ownership
	userID := c.GetInt("userID")
	event, err := h.eventService.GetByID(id)
	if err != nil {
		utils.SendNotFound(c, "Event not found")
		return
	}
	if event.CreatorID != userID {
		utils.SendForbidden(c, "You don't have permission to update this event")
		return
	}

	input.ID = id
	input.CreatorID = userID
	updatedEvent, err := h.eventService.Update(input)
	if err != nil {
		utils.SendInternalError(c, "Failed to update event", err)
		return
	}

	utils.SendSuccess(c, "Event updated successfully", updatedEvent)
}

// DeleteEvent handles deleting an event
// @Summary Delete an event
// @Description Delete a specific event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequest(c, "Invalid event ID", err)
		return
	}

	// Verify ownership
	userID := c.GetInt("userID")
	event, err := h.eventService.GetByID(id)
	if err != nil {
		utils.SendNotFound(c, "Event not found")
		return
	}
	if event.CreatorID != userID {
		utils.SendForbidden(c, "You don't have permission to delete this event")
		return
	}

	err = h.eventService.Delete(id)
	if err != nil {
		utils.SendInternalError(c, "Failed to delete event", err)
		return
	}

	utils.SendSuccess(c, "Event deleted successfully", nil)
}

// RegisterForEvent handles event registration
// @Summary Register for an event
// @Description Register the current user for a specific event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /events/{id}/register [post]
func (h *EventHandler) RegisterForEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequest(c, "Invalid event ID", err)
		return
	}

	userID := c.GetInt("userID")
	err = h.eventService.RegisterUser(eventID, userID)
	if err != nil {
		utils.SendInternalError(c, "Failed to register for event", err)
		return
	}

	utils.SendSuccess(c, "Successfully registered for event", nil)
}

// UnregisterFromEvent handles event unregistration
// @Summary Unregister from an event
// @Description Unregister the current user from a specific event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /events/{id}/unregister [post]
func (h *EventHandler) UnregisterFromEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequest(c, "Invalid event ID", err)
		return
	}

	userID := c.GetInt("userID")
	err = h.eventService.UnregisterUser(eventID, userID)
	if err != nil {
		utils.SendInternalError(c, "Failed to unregister from event", err)
		return
	}

	utils.SendSuccess(c, "Successfully unregistered from event", nil)
}

// GetRegisteredUsers handles retrieving all users registered for an event
// @Summary Get registered users
// @Description Get a list of users registered for a specific event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /events/{id}/registrations [get]
func (h *EventHandler) GetRegisteredUsers(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequest(c, "Invalid event ID", err)
		return
	}

	users, err := h.eventService.GetRegisteredUsers(eventID)
	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve registered users", err)
		return
	}

	utils.SendSuccess(c, "Successfully retrieved registered users", users)
}
