package handlers

import (
	"mental-health-companion/internal/database"
	"mental-health-companion/internal/models"
	"mental-health-companion/internal/utils"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func sendError(c *gin.Context, err *utils.APIError) {
	c.JSON(err.HTTPStatus, err)
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, utils.NewInvalidRequestFormatError())
		return
	}

	// Required fields validation
	if req.Email == "" {
		sendError(c, utils.NewMissingRequiredFieldError("email"))
		return
	}
	if req.Password == "" {
		sendError(c, utils.NewMissingRequiredFieldError("password"))
		return
	}

	// Email format validation
	if _, err := mail.ParseAddress(req.Email); err != nil {
		sendError(c, utils.NewInvalidEmailError())
		return
	}

	// Email to lowercase
	req.Email = strings.ToLower(req.Email)

	// Check if email exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		sendError(c, utils.NewEmailAlreadyExistsError())
		return
	}

	// Password validation
	if passwordErrors := utils.ValidatePassword(req.Password); len(passwordErrors) > 0 {
		sendError(c, utils.NewInvalidPasswordError(passwordErrors))
		return
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Surname:  req.Surname,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		sendError(c, utils.NewDatabaseError(result.Error))
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		sendError(c, utils.NewInternalServerError(err))
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, utils.NewInvalidRequestFormatError())
		return
	}

	// Required fields validation
	if req.Email == "" {
		sendError(c, utils.NewMissingRequiredFieldError("email"))
		return
	}
	if req.Password == "" {
		sendError(c, utils.NewMissingRequiredFieldError("password"))
		return
	}

	req.Email = strings.ToLower(req.Email)

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		sendError(c, utils.NewInvalidCredentialsError())
		return
	}

	if err := user.ComparePassword(req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			sendError(c, utils.NewInvalidCredentialsError())
			return
		}
		sendError(c, utils.NewInternalServerError(err))
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		sendError(c, utils.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendError(c, utils.NewInvalidTokenError())
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		sendError(c, utils.NewRecordNotFoundError("User"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"surname":   user.Surname,
			"createdAt": user.CreatedAt,
		},
	})
}

// Şifre güvenlik kontrolü
func isPasswordSecure(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*(),.?\":{}|<>", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
