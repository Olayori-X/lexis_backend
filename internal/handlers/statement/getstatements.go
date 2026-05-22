package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Olayori-X/notes/api"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	log "github.com/sirupsen/logrus"
)

func GetStatementsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userid")

	if userID == "" {
		api.RequestErrorHandler(w, fmt.Errorf("user_id is required"))
		return
	}

	var db *sqltools.DatabaseInterface
	db, err := sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	statements, err := (*db).GetStatements(userID)
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.SearchStatementsResponse{
		Code:       http.StatusOK,
		Statements: statements,
	})
}
