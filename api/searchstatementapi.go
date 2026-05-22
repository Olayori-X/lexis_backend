package api

import "github.com/Olayori-X/notes/models"

type SearchStatementsResponse struct {
	Code       int                `json:"code"`
	Statements []models.Statement `json:"statements"`
}
