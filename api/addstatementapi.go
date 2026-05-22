package api

type AddStatementParams struct {
	UserID      string `json:"user_id"`
	Statement   string `json:"statement"`
	Association string `json:"association"`
}

type AddStatementResponse struct {
	Code        int    `json:"code"`
	StatementID string `json:"statement_id"`
	Message     string `json:"message"`
}
