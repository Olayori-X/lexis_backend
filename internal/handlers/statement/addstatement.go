package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Olayori-X/notes/api"
	sqltools "github.com/Olayori-X/notes/internal/tools"
	"github.com/Olayori-X/notes/models"
	log "github.com/sirupsen/logrus"
)

func AddStatementHandler(w http.ResponseWriter, r *http.Request) {
	var params api.AddStatementParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	if params.UserID == "" || params.Statement == "" || params.Association == "" {
		api.RequestErrorHandler(w, fmt.Errorf("user_id, statement and association are required"))
		return
	}

	newStatement := models.Statement{
		UserID:      params.UserID,
		Content:     params.Statement,
		Association: params.Association,
	}

	var db *sqltools.DatabaseInterface
	db, err = sqltools.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		api.InternalErrorHandler(w)
		return
	}

	statementID, err := (*db).AddStatementWithAssociation(newStatement)
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(api.AddStatementResponse{
		Code:        http.StatusCreated,
		StatementID: statementID,
		Message:     "Statement added successfully",
	})
}
