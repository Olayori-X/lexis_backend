package middleware

import (
	"errors"
	"net/http"

	"github.com/Olayori-X/notes/api"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var UnAuthorizedError = errors.New("Invalid Username or Token")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		username := r.Header.Get("userid")
		token := r.Header.Get("Authorization")

		if username == "" || token == "" {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		var database *sqltools.DatabaseInterface
		database, err = sqltools.NewDatabase()

		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		loginDetails := (*database).UserLoggedIn(username)

		if loginDetails == nil {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		log.Printf("Login details for user %s: %+v", token, loginDetails)
		err = bcrypt.CompareHashAndPassword([]byte((*loginDetails).Code), []byte(token))
		if err != nil {
			log.Print("Invalid token for user ", username, ": ", err)
			log.Warn("Invalid login attempt for:", username)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
