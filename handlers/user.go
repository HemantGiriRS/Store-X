package handlers

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"storex/database/dbHelper"
	"storex/middlewares"
	"storex/models"
	"storex/utils"
	"strconv"
	"strings"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		logrus.Warnf("Failed to parse refresh token: %v", err)
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	tokenResponse, err := utils.GenerateTokenPair(claims.UserID, claims.Role)
	if err != nil {
		logrus.Errorf("Failed to generate token pair during refresh: %v", err)
		http.Error(w, "could not generate new tokens", http.StatusInternalServerError)
		return
	}

	if err := utils.EncodeResponse(w, http.StatusOK, tokenResponse); err != nil {
		logrus.Errorf("Failed to encode refresh token response: %v", err)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !utils.IsValidEmail(req.Email) {
		http.Error(w, "invalid email format or domain", http.StatusBadRequest)
		return
	}

	// This flag will help us determine the correct HTTP status code
	isNewUser := false

	user, err := dbHelper.FindUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // User does not exist, so we register them
			isNewUser = true
			user, err = dbHelper.RegisterUser(req.Email)
			if err != nil {
				logrus.Errorf("Failed to register new user '%s': %v", req.Email, err)
				http.Error(w, "error creating user profile", http.StatusInternalServerError)
				return
			}
		} else {
			logrus.Errorf("Failed to find user by email '%s': %v", req.Email, err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}

	tokenResponse, err := utils.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		logrus.Errorf("Failed to generate token pair for user '%s': %v", user.ID, err)
		http.Error(w, "could not generate tokens", http.StatusInternalServerError)
		return
	}

	finalResponse := models.SignInResponse{
		Message:      "Login successful",
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
	}

	statusCode := http.StatusOK // Default to 200 OK for existing users
	if isNewUser {
		statusCode = http.StatusCreated // Use 201 Created for new users
		finalResponse.Message = "Registration successful"
	}

	if err := utils.EncodeResponse(w, statusCode, finalResponse); err != nil {
		logrus.Errorf("Failed to encode sign-in response: %v", err)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	actorIDVal := r.Context().Value(middlewares.UserIDKey)
	createdBy, ok := actorIDVal.(uuid.UUID)
	if !ok {
		logrus.Error("Could not retrieve valid user ID from context in SignUp")
		http.Error(w, "Invalid user ID in token context", http.StatusInternalServerError)
		return
	}

	var req models.SignUpRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "Error in decoding request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.PhoneNo == "" || req.Type == "" {
		http.Error(w, "Missing required fields: email, phone_no, type", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(req.Email) {
		http.Error(w, "Invalid email format or domain", http.StatusBadRequest)
		return
	}
	if check := utils.ValidEmpTypes[req.Type]; !check {
		http.Error(w, "Invalid employee type", http.StatusBadRequest)
	}

	const defaultRole = "employee"

	newUserID, err := dbHelper.CreateUser(req, defaultRole, createdBy)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			logrus.Warnf("Attempt to create user with duplicate email or phone: %s", req.Email)
			http.Error(w, "An employee with this email or phone number already exists.", http.StatusConflict)
			return
		}

		logrus.Errorf("Failed to create user in database: %v", err)
		http.Error(w, "Could not create user.", http.StatusInternalServerError)
		return
	}

	response := models.SignUpResponse{
		Message: "Employee created successfully",
		UserID:  newUserID,
	}
	if err := utils.EncodeResponse(w, http.StatusCreated, response); err != nil {
		logrus.Errorf("Failed to encode sign-up response: %v", err)
	}
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	adminIDVal := r.Context().Value(middlewares.UserIDKey) // Use the key defined in your middleware
	if adminIDVal == nil {
		http.Error(w, "Unauthorized: Could not identify the admin user.", http.StatusUnauthorized)
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		http.Error(w, "Internal server error: user ID in context is of an invalid type.", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	employeeID, ok := vars["employee_id"] // This key must match your route definition
	if !ok {
		http.Error(w, "employee_id is missing in URL path", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(employeeID); err != nil {
		http.Error(w, "Invalid employee_id format. Must be a UUID.", http.StatusBadRequest)
		return
	}

	var req models.UpdateRoleRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	normalizedRole := strings.ToLower(req.Role)
	if !utils.ValidRoles[normalizedRole] {
		http.Error(w, "Invalid role provided. Must be one of: admin, asset_manager, employee_manager, employee.", http.StatusBadRequest)
		return
	}

	if err := dbHelper.UpdateEmployeeRole(employeeID, normalizedRole, adminID.String()); err != nil {
		// Check if the error is "not found"
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Employee with the specified ID was not found.", http.StatusNotFound)
		} else {
			logrus.Errorf("Failed to update employee role in database: %v", err)
			http.Error(w, "An internal server error occurred.", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "Employee role updated successfully"}
	if err := utils.EncodeResponse(w, http.StatusOK, response); err != nil {
		logrus.Errorf("Failed to send success response: %v", err)
	}
}

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filters := make(map[string]string)

	if q := queryParams.Get("q"); q != "" {
		filters["q"] = q
	}

	if role := strings.ToLower(queryParams.Get("role")); role != "" {
		if utils.ValidRoles[role] {
			filters["role"] = role
		}
	}

	if empType := strings.ToLower(queryParams.Get("type")); empType != "" {
		if utils.ValidEmpTypes[empType] {
			filters["type"] = empType
		}
	}

	if status := strings.ToLower(queryParams.Get("assignment_status")); status != "" {
		if utils.ValidAssignmentStatus[status] {
			filters["assignment_status"] = status
		}
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(queryParams.Get("pageSize"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	employees, totalRecords, err := dbHelper.GetEmployeesFiltered(filters, page, pageSize)
	if err != nil {
		logrus.Errorf("Failed to get employees from database: %v", err)
		http.Error(w, "Failed to retrieve employee data", http.StatusInternalServerError)
		return
	}

	totalPages := 0
	if totalRecords > 0 {
		totalPages = int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	}

	response := models.GetEmployeesResponse{
		Data: employees,
		Pagination: models.PaginationInfo{
			CurrentPage:  page,
			PageSize:     pageSize,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	}
	utils.EncodeResponse(w, http.StatusOK, response)
}
