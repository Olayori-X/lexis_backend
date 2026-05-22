package authhandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Olayori-X/notes/api"
	"github.com/Olayori-X/notes/functions"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	"github.com/Olayori-X/notes/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func VerifyOtpHandler(w http.ResponseWriter, r *http.Request) {
	var params = api.VerifyOtpParams{}
	// var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if strings.TrimSpace(params.UserID) == "" {
		log.Error("user_id cannot be empty")
		api.RequestErrorHandler(w, errors.New("user_id cannot be empty"))
		return
	}

	if strings.TrimSpace(params.OTP) == "" {
		log.Error("otp cannot be empty")
		api.RequestErrorHandler(w, errors.New("otp cannot be empty"))
		return
	}

	var database *sqltools.DatabaseInterface
	database, err = sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	var userDetails *models.User = (*database).GetUserDetails(params.UserID)

	err = bcrypt.CompareHashAndPassword([]byte(*userDetails.Code), []byte(params.OTP))
	if err != nil {
		log.Warn("Invalid login attempt for:", userDetails.Email)
		api.UnAuthorizedError(w)
		return
	}

	token, err := functions.HashString(params.OTP)
	if err != nil {
		log.Error("Hashing authorization code failed", err)
		api.InternalErrorHandler(w)
		return
	}

	fmt.Print(token)

	errChan := make(chan error, 2)

	go func() {
		if err := (*database).VerifyUser(params.UserID); err != nil {
			log.Error("Failed to verify user: ", err)
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		if err := (*database).UpsertLoggedInUser(userDetails.UserID, token); err != nil {
			log.Error("Failed to update logged-in user: ", err)
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Wait for both to finish
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			api.InternalErrorHandler(w)
			return
		}
	}

	var response = api.VerifyOtpResponse{
		Code:     http.StatusOK,
		Verified: true,
		UserID:   userDetails.UserID,
		Token:    params.OTP,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}
