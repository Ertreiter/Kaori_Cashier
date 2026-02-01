package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kaori/backend/internal/middleware"
	"github.com/kaori/backend/internal/model"
	"github.com/kaori/backend/internal/service"
	"github.com/kaori/backend/pkg/response"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// LoginWithPIN handles POST /api/auth/login/pin
func (h *AuthHandler) LoginWithPIN(c *gin.Context) {
	var req model.LoginPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.service.LoginWithPIN(req.Email, req.PIN)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// Me handles GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "User not found in context")
		return
	}

	user, err := h.service.GetCurrentUser(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if user == nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "User not found in context")
		return
	}

	if err := h.service.Logout(userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}
