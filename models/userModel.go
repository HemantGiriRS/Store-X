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

type SignUpRequest struct {
	Email   string `json:"email" db:"email"`
	PhoneNo string `json:"phone_no" db:"phone_no"`
	Type    string `json:"type" db:"type"`
}

type SignUpResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}
