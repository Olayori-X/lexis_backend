package sqltools

import (
	"time"

	"github.com/Olayori-X/notes/models"
	log "github.com/sirupsen/logrus"
)

type LoginDetails struct {
	Password string
	UserID   string
	Verified bool
}

type AuthenticatedUser struct {
	Code      string
	UserID    string
	LoginTime time.Time
}

type CoinDetails struct {
	Coins    int64
	Username string
}

type DatabaseInterface interface {
	GetUserLoginDetails(username string) *LoginDetails
	GetUserDetails(userID string) *models.User
	AddUser(user models.User) (string, error)
	SetupDatabase() error
	GetUsers() ([]models.User, error)
	UpsertLoggedInUser(userID string, code string) error
	UpdateUserCode(userID string, hashedCode string) error
	UserLoggedIn(userid string) *AuthenticatedUser
	VerifyUser(userID string) error
	UpdateUserProfile(user models.User) error
	AddForgotPasswordRecord(userID, code string) error
	ChangeUserPassword(email string, hashedPassword string) error
	AddStatementWithAssociation(statement models.Statement) (string, error)
	SearchStatements(userID string, searchTerm string) ([]models.Statement, error)
	GetStatements(userID string) ([]models.Statement, error)
	DeleteStatement(statementID string, userID string) error
	UpdateStatement(statementID string, userID string, content string, association string) error
}

func NewDatabase() (*DatabaseInterface, error) {
	var database DatabaseInterface = &RealDB{}
	var err error = database.SetupDatabase()

	if err != nil {
		log.Error("Failed to set up database connection: ", err)
		return nil, err
	}

	return &database, nil
}
