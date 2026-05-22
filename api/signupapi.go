package api

type SignupParams struct {
	Email    string `db:"email" json:"emailPhone"`
	Name     string `db:"name" json:"name"`
	Password string `db:"password" json:"password"`
}

type SignupResponse struct {
	//success code, usually 200
	Code int `json:"code"`

	//Username of the user
	Username string `json:"username"`

	//Message to be displayed
	Message string `json:"message"`
}
