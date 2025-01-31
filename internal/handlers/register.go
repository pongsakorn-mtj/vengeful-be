package handlers

import (
	"net/http"
	"vengeful-be/internal/models"
	"vengeful-be/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RegisterHandler struct {
	logger   *logrus.Logger
	userRepo *repository.UserRepository
}

func NewRegisterHandler(logger *logrus.Logger, userRepo *repository.UserRepository) *RegisterHandler {
	return &RegisterHandler{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to bind request body")

		c.JSON(http.StatusBadRequest, models.RegisterResponse{
			Status:  "error",
			Message: "Invalid request format",
		})
		return
	}

	// Log the registration attempt
	h.logger.WithFields(logrus.Fields{
		"email":     request.Email,
		"firstName": request.FirstName,
		"lastName":  request.LastName,
	}).Info("Registration attempt")

	// Validate acceptance of terms and policies
	if !request.IsAcceptTnc || !request.IsAcceptPrivacyPolicy {
		c.JSON(http.StatusBadRequest, models.RegisterResponse{
			Status:  "error",
			Message: "Must accept terms and conditions and privacy policy",
		})
		return
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(c.Request.Context(), request.Email)
	if err != nil {
		h.logger.WithError(err).Error("Failed to check email existence")
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Status:  "error",
			Message: "Internal server error",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, models.RegisterResponse{
			Status:  "error",
			Message: "Email already registered",
		})
		return
	}

	// Create user in database
	if err := h.userRepo.Create(c.Request.Context(), &request); err != nil {
		h.logger.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Status:  "error",
			Message: "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, models.RegisterResponse{
		Status:  "success",
		Message: "Registration successful",
	})
}
