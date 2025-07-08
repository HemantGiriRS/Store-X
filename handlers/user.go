package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"storex/database/dbHelper"
	"storex/models"
	"storex/utils"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "error in decoding request", http.StatusInternalServerError)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}

	token, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil || token == nil {
		http.Error(w, "invalid refresh token", http.StatusBadRequest)
		return
	}

	newAccessToken, er := utils.GenerateAccessToken(token.UserID, token.Role)
	if er != nil {
		http.Error(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, er := utils.GenerateRefreshToken(token.UserID, token.Role)
	if er != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	res := models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	if err := utils.EncodeResponse(w, http.StatusOK, &res); err != nil {
		http.Error(w, "error in sending response", http.StatusBadRequest)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest

	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "error in decoding request", http.StatusBadRequest)
		return
	}

	if isValidEmail := utils.IsValidEmail(req.Email); !isValidEmail {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	//check user in db
	var user *models.UserContext
	user, err := dbHelper.FindUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //if not exist
			user, err = dbHelper.CreateUser(req.Email)
			if err != nil {
				http.Error(w, "error in creating user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "error in finding user", http.StatusInternalServerError)
			return
		}
	}

	newAccessToken, er := utils.GenerateAccessToken(user.ID, user.Role)
	if er != nil {
		http.Error(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, er := utils.GenerateRefreshToken(user.ID, user.Role)
	if er != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	res := models.SignInResponse{
		Message:      "Login successful",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	if err := utils.EncodeResponse(w, http.StatusCreated, &res); err != nil {
		http.Error(w, "error in sending response", http.StatusBadRequest)
	}

}
