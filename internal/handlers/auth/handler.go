package auth

import (
	"net/http"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models/auth"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Login to the application
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.LoginRequest true "Login credentials"
// @Success      200  {object}  common.Response{data=auth.LoginResponse}  "Login successful"
// @Failure      400  {object}  common.Response  "Invalid request"
// @Failure      401  {object}  common.Response  "Invalid credentials"
// @Failure      404  {object}  common.Response  "User not found"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	user, token, err := h.service.Auth.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case common.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token, "user": user})
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.RegisterRequest true "User registration details"
// @Success      201  {object}  common.Response{data=auth.RegisterResponse}  "User registered successfully"
// @Failure      400  {object}  common.Response  "Invalid request"
// @Failure      409  {object}  common.Response  "User already exists"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, password and full_name are required"})
		return
	}

	user, token, err := h.service.Auth.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		switch err {
		case common.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
		"token":   token,
		"user":    user,
	})
}