package api

type UpdateStatementParams struct {
	Statement   string `json:"statement"`
	Association string `json:"association"`
}

type UpdateStatementResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
