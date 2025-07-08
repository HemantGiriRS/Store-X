package models

type SignInRequest struct {
	Email string `json:"email" db:"email"`
}

type SignInResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserContext struct {
	ID   string `json:"id" db:"id"`
	Role string `json:"role" db:"role"`
}
