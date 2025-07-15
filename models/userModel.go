package models

import "time"

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

type EmployeeDetail struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	PhoneNo     string    `json:"phone_no" db:"phone_no"`
	Type        string    `json:"type" db:"type"`
	Role        string    `json:"role" db:"role"`
	AssetStatus int       `json:"asset_status" db:"asset_status"` // <-- ADD THIS LINE
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type GetEmployeesResponse struct {
	Data       []EmployeeDetail `json:"data"`
	Pagination PaginationInfo   `json:"pagination"`
}
