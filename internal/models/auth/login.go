package auth

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response body for successful login
type LoginResponse struct {
	Code  int    `json:"code"`
	ID    string `json:"id"`
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// RegisterResponse represents the response body for successful registration
type RegisterResponse struct {
	Code  int    `json:"code"`
	ID    string `json:"id"`
	Token string `json:"token"`
}
