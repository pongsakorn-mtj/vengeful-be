package handlers

import (
	"math"
	"net/http"
	"strconv"
	"vengeful-be/internal/models"
	"vengeful-be/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetUsersHandler struct {
	logger   *logrus.Logger
	userRepo *repository.UserRepository
}

func NewGetUsersHandler(logger *logrus.Logger, userRepo *repository.UserRepository) *GetUsersHandler {
	return &GetUsersHandler{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (h *GetUsersHandler) GetAll(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get users from repository
	users, total, err := h.userRepo.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get users")
		c.JSON(http.StatusInternalServerError, models.GetAllUsersResponse{
			Status:  "error",
			Message: "Failed to get users",
		})
		return
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	// Convert repository users to response users
	var responseUsers []models.UserResponse
	for _, user := range users {
		responseUsers = append(responseUsers, models.UserResponse{
			FirstName:             user.FirstName,
			LastName:              user.LastName,
			PhoneNo:               user.PhoneNo,
			Email:                 user.Email,
			IsAcceptTnc:           user.IsAcceptTnc,
			IsAcceptPrivacyPolicy: user.IsAcceptPrivacyPolicy,
			CreatedAt:             user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, models.GetAllUsersResponse{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data: &models.GetUsersData{
			Users: responseUsers,
			Pagination: &models.Pagination{
				CurrentPage:  page,
				TotalPages:   totalPages,
				TotalRecords: total,
				Limit:        limit,
			},
		},
	})
}
