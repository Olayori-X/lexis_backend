package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Olayori-X/notes/api"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	log "github.com/sirupsen/logrus"
)

func DeleteStatementHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userid")
	statementID := r.PathValue("id")

	if userID == "" {
		api.RequestErrorHandler(w, fmt.Errorf("user_id is required"))
		return
	}

	if statementID == "" {
		api.RequestErrorHandler(w, fmt.Errorf("statement_id is required"))
		return
	}

	var db *sqltools.DatabaseInterface
	db, err := sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	err = (*db).DeleteStatement(statementID, userID)
	if err != nil {
		log.Error("Failed to delete statement: ", err)
		api.RequestErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.DeleteStatementResponse{
		Code:    http.StatusOK,
		Message: "Statement deleted successfully",
	})
}
