package authhandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Olayori-X/notes/api"
	"github.com/Olayori-X/notes/functions"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	"github.com/Olayori-X/notes/models"
	log "github.com/sirupsen/logrus"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("SignupHandler called")
	var params = api.SignupParams{}
	// var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if strings.TrimSpace(params.Email) == "" {
		api.RequestErrorHandler(w, errors.New("email cannot be empty"))
		return
	}
	if strings.TrimSpace(params.Password) == "" {
		api.RequestErrorHandler(w, errors.New("password cannot be empty"))
		return
	}

	fmt.Printf("Received signup request: %+v\n", params)

	var database *sqltools.DatabaseInterface
	database, err = sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	hashedPassword, err := functions.HashString(params.Password)
	if err != nil {
		log.Error("Failed to hash password:", err)
		api.InternalErrorHandler(w)
		return
	}

	code, err := functions.GenerateOTPCode(6)
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	hashedCode, err := functions.HashString(code)
	if err != nil {
		log.Error("Failed to hash password:", err)
		api.InternalErrorHandler(w)
		return
	}

	newUser := models.User{
		UserID:    "userID",
		Name:      "None",
		Email:     params.Email,
		Code:      &hashedCode,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	errChan := make(chan error, 2)
	var userid string

	// Send email concurrently
	go func() {
		messageID, err := functions.SendSimpleMessage(params.Email, "OTP code", code)
		// err = functions.SendMail([]string{params.Email}, "OTP Code", code)
		if err != nil {
			log.Error("Failed to send email: ", err)
			log.Error("MessageId: ", messageID)
			errChan <- errors.New("an error occurred while sending the mail")
			return
		}
		errChan <- nil
	}()

	// Add user concurrently
	go func() {
		var errAdd error
		userid, errAdd = (*database).AddUser(newUser)
		if errAdd != nil {
			log.Error("Failed to add user: ", errAdd)
			errChan <- errors.New("user already exists or could not be added")
			return
		}
		errChan <- nil
	}()

	// Wait for both operations to complete
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			api.RequestErrorHandler(w, err)
			return
		}
	}

	var response = api.SignupResponse{
		Code:     http.StatusOK,
		Message:  "Signup successful",
		Username: userid,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}
