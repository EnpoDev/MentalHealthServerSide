package utils

import "net/http"

// APIError represents a standardized error response
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	HTTPStatus int    `json:"-"`
}

// Common error codes
const (
	// Authentication Errors (1xxx)
	ErrCodeInvalidCredentials = "ERR_1001"
	ErrCodeTokenExpired       = "ERR_1002"
	ErrCodeInvalidToken       = "ERR_1003"
	ErrCodeMissingToken       = "ERR_1004"
	ErrCodeInvalidTokenFormat = "ERR_1005"

	// Validation Errors (2xxx)
	ErrCodeInvalidEmail         = "ERR_2001"
	ErrCodeInvalidPassword      = "ERR_2002"
	ErrCodeEmailAlreadyExists   = "ERR_2003"
	ErrCodeInvalidRequestFormat = "ERR_2004"
	ErrCodeMissingRequiredField = "ERR_2005"
	ErrCodeInvalidField         = "ERR_2006"

	// Database Errors (3xxx)
	ErrCodeDatabaseError  = "ERR_3001"
	ErrCodeRecordNotFound = "ERR_3002"
	ErrCodeDuplicateEntry = "ERR_3003"

	// Server Errors (5xxx)
	ErrCodeInternalServer     = "ERR_5001"
	ErrCodeServiceUnavailable = "ERR_5002"
)

// Common validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error constructors
func NewInvalidCredentialsError() *APIError {
	return &APIError{
		Code:       ErrCodeInvalidCredentials,
		Message:    "Invalid email or password",
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewTokenExpiredError() *APIError {
	return &APIError{
		Code:       ErrCodeTokenExpired,
		Message:    "Token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewInvalidTokenError() *APIError {
	return &APIError{
		Code:       ErrCodeInvalidToken,
		Message:    "Invalid token",
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewMissingTokenError() *APIError {
	return &APIError{
		Code:       ErrCodeMissingToken,
		Message:    "Authorization token is missing",
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewInvalidTokenFormatError() *APIError {
	return &APIError{
		Code:       ErrCodeInvalidTokenFormat,
		Message:    "Invalid token format",
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewInvalidEmailError() *APIError {
	return &APIError{
		Code:       ErrCodeInvalidEmail,
		Message:    "Invalid email format",
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewInvalidPasswordError(details []ValidationError) *APIError {
	return &APIError{
		Code:       ErrCodeInvalidPassword,
		Message:    "Password does not meet security requirements",
		Details:    details,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewEmailAlreadyExistsError() *APIError {
	return &APIError{
		Code:       ErrCodeEmailAlreadyExists,
		Message:    "Email is already registered",
		HTTPStatus: http.StatusConflict,
	}
}

func NewInvalidRequestFormatError() *APIError {
	return &APIError{
		Code:       ErrCodeInvalidRequestFormat,
		Message:    "Invalid request format",
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewMissingRequiredFieldError(field string) *APIError {
	return &APIError{
		Code:    ErrCodeMissingRequiredField,
		Message: "Required field is missing",
		Details: ValidationError{
			Field:   field,
			Message: "This field is required",
		},
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewInvalidFieldError(field, message string) *APIError {
	return &APIError{
		Code:    ErrCodeInvalidField,
		Message: "Invalid field value",
		Details: ValidationError{
			Field:   field,
			Message: message,
		},
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewDatabaseError(err error) *APIError {
	return &APIError{
		Code:       ErrCodeDatabaseError,
		Message:    "Database operation failed",
		Details:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewRecordNotFoundError(resource string) *APIError {
	return &APIError{
		Code:       ErrCodeRecordNotFound,
		Message:    resource + " not found",
		HTTPStatus: http.StatusNotFound,
	}
}

func NewInternalServerError(err error) *APIError {
	return &APIError{
		Code:       ErrCodeInternalServer,
		Message:    "Internal server error",
		Details:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// Helper function to validate password and return detailed errors
func ValidatePassword(password string) []ValidationError {
	var errors []ValidationError

	if len(password) < 8 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		})
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
		case isSpecialChar(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter",
		})
	}

	if !hasLower {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter",
		})
	}

	if !hasNumber {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one number",
		})
	}

	if !hasSpecial {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one special character (!@#$%^&*(),.?\":{}|<>)",
		})
	}

	return errors
}

func isSpecialChar(c rune) bool {
	specialChars := "!@#$%^&*(),.?\":{}|<>"
	for _, sc := range specialChars {
		if c == sc {
			return true
		}
	}
	return false
}
