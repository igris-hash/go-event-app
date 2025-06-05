package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/internal/model"
	"github.com/igris-hash/go-event-app/internal/service"
	"github.com/igris-hash/go-event-app/internal/utils"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserRegistration true "User registration details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var input model.UserRegistration
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequest(c, "Invalid input", err)
		return
	}

	user, err := h.userService.Register(input)
	if err != nil {
		utils.SendInternalError(c, "Failed to register user", err)
		return
	}

	utils.SendCreated(c, "User registered successfully", user)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body model.UserLogin true "User login credentials"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var input model.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequest(c, "Invalid input", err)
		return
	}

	token, err := h.userService.Login(input)
	if err != nil {
		utils.SendUnauthorized(c, "Invalid credentials")
		return
	}

	utils.SendSuccess(c, "Login successful", gin.H{"token": token})
}

// GetProfile handles getting the current user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	user, err := h.userService.GetByID(userID)
	if err != nil {
		utils.SendNotFound(c, "User not found")
		return
	}

	utils.SendSuccess(c, "Profile retrieved successfully", user)
}

// UpdateProfile handles updating the current user's profile
// @Summary Update user profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body model.UserUpdate true "User update details"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	var input model.UserUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequest(c, "Invalid input", err)
		return
	}

	user, err := h.userService.Update(userID, input)
	if err != nil {
		utils.SendInternalError(c, "Failed to update profile", err)
		return
	}

	utils.SendSuccess(c, "Profile updated successfully", user)
}
