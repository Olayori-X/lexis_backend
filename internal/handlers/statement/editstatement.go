package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Olayori-X/notes/api"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	log "github.com/sirupsen/logrus"
)

func UpdateStatementHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("UpdateStatementHandler called")
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

	var params api.UpdateStatementParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	if params.Statement == "" || params.Association == "" {
		api.RequestErrorHandler(w, fmt.Errorf("statement and association are required"))
		return
	}

	var db *sqltools.DatabaseInterface
	db, err = sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	err = (*db).UpdateStatement(statementID, userID, params.Statement, params.Association)
	if err != nil {
		log.Error("Failed to update statement: ", err)
		api.RequestErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.UpdateStatementResponse{
		Code:    http.StatusOK,
		Message: "Statement updated successfully",
	})
}
