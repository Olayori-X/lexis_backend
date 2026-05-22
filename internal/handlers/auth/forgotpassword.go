package authhandlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Olayori-X/notes/api"
	"github.com/Olayori-X/notes/functions"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	log "github.com/sirupsen/logrus"
)

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Email string `json:"email"`
	}

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Error("Invalid request body: ", err)
		api.RequestErrorHandler(w, errors.New("invalid request"))
		return
	}

	email := strings.TrimSpace(params.Email)
	if email == "" {
		api.RequestErrorHandler(w, errors.New("email cannot be empty"))
		return
	}

	// Connect to DB
	dbi, err := sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Generate OTP
	otp, err := functions.GenerateOTPCode(6)
	if err != nil {
		log.Error("Failed to generate OTP: ", err)
		api.InternalErrorHandler(w)
		return
	}

	otpHash, err := functions.HashString(otp)
	if err != nil {
		log.Error("Failed to hash OTP: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Add forgot password record
	if err := (*dbi).AddForgotPasswordRecord(params.Email, otpHash); err != nil {
		log.Error("Failed to save forgot password record: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Send email
	err = functions.SendMail([]string{email}, "Reset Password Code", otp)
	if err != nil {
		log.Error("Failed to send email: ", err)
		api.RequestErrorHandler(w, errors.New("could not send reset code"))
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reset code sent to your email.",
	})
}
