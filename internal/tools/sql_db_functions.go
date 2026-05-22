package sqltools

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Olayori-X/notes/functions"
	"github.com/Olayori-X/notes/models"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

// type mockDB struct{}

type RealDB struct {
	DB *sql.DB
}

func UserExists(db *RealDB, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := db.DB.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *RealDB) UserLoggedIn(userid string) *AuthenticatedUser {
	query := `
	SELECT user_id, code, created_at
	FROM loggedin_users
	WHERE user_id = $1;`

	var userID, code string
	var loginTime time.Time

	err := db.DB.QueryRow(query, userid).Scan(&userID, &code, &loginTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}

	return &AuthenticatedUser{
		UserID:    userID,
		Code:      code,
		LoginTime: loginTime,
	}
}

func (db *RealDB) GetUserLoginDetails(email string) *LoginDetails {
	query := `
	SELECT user_id, password, verified
	FROM users
	WHERE email = $1`

	var user_id, password string
	var verified bool

	err := db.DB.QueryRow(query, email).Scan(&user_id, &password, &verified)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Login attempt with non-existing user: %s\n", email)
			return nil
		}
		fmt.Printf("Error fetching login details for email '%s': %v\n", email, err)
		return nil
	}
	log.Warn("Login attempt with non-existing user:", user_id+" "+password+" "+fmt.Sprintf("%t", verified))
	fmt.Printf("Login attempt with non-existing user: %s %s %t\n", user_id, password, verified)
	return &LoginDetails{
		UserID:   user_id,
		Password: password,
		Verified: verified,
	}
}

func (db *RealDB) GetUserDetails(userid string) *models.User {
	query := `
	SELECT user_id, name, email, password, code,
	       created_at, updated_at
	FROM users
	WHERE user_id = $1;`

	var user models.User

	err := db.DB.QueryRow(query, userid).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Code,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}
	return &user
}

func (db *RealDB) AddUser(user models.User) (string, error) {
	// Check if user already exists
	exists, err := UserExists(db, user.Email)
	if err != nil {
		return "", fmt.Errorf("error checking user: %w", err)
	}
	if exists {
		return "", fmt.Errorf("user '%s' already exists", user.Email)
	}

	query := `
	INSERT INTO users (user_id, name, email, password, code)
	VALUES ($1, $2, $3, $4, $5) RETURNING user_id;`

	userid, err := functions.GenerateUserID()
	if err != nil {
		// return fmt.Errorf("could not generate user ID: %w", err)
		return "", err
	}

	var pk string
	err = db.DB.QueryRow(query, userid, user.Name, user.Email, user.Password, user.Code).Scan(&pk)

	if err != nil {
		return "", err
	}

	return pk, nil
}

func (db *RealDB) VerifyUser(userid string) error {
	query := `
		UPDATE users
		SET verified = TRUE,
			code = NULL
		WHERE user_id = $1;
	`

	_, err := db.DB.Exec(query, userid)
	if err != nil {
		return fmt.Errorf("failed to update user verification: %w", err)
	}

	return nil
}

func (db *RealDB) UpsertLoggedInUser(userID string, code string) error {
	query := `
	INSERT INTO loggedin_users (user_id, code)
	VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET code = EXCLUDED.code, created_at = CURRENT_TIMESTAMP;
	`

	_, err := db.DB.Exec(query, userID, code)
	if err != nil {
		return err
	}
	return nil
}

// SQL implementation
func (db *RealDB) UpdateUserCode(userID string, hashedCode string) error {
	_, err := db.DB.Exec(`
		UPDATE users
		SET code = $1, updated_at = NOW()
		WHERE user_id = $2
	`, hashedCode, userID)
	if err != nil {
		return fmt.Errorf("error updating user code: %w", err)
	}
	return nil
}

func (db *RealDB) UpdateUserProfile(user models.User) error {
	query := `
		UPDATE users
		SET name = $1,
			email = $3,
			phone = $4,
			updated_at = $5
		WHERE user_id = $6;
	`

	_, err := db.DB.Exec(query,
		user.Name,
		user.Email,
		time.Now(),
		user.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

func (db *RealDB) AddForgotPasswordRecord(email, code string) error {
	exists, err := UserExists(db, email)
	if err != nil {
		return fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		return fmt.Errorf("'%s' does not exist", email)
	}

	query := `
		INSERT INTO forgotpassword (user_id, code)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET
			code = EXCLUDED.code,
			created_at = CURRENT_TIMESTAMP;`

	_, err = db.DB.Exec(query, email, code)
	if err != nil {
		return fmt.Errorf("failed to add forgot password record: %w", err)
	}
	return nil
}

func (db *RealDB) ChangeUserPassword(userid string, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE user_id = $2;
	`

	_, err := db.DB.Exec(query, hashedPassword, userid)
	if err != nil {
		return fmt.Errorf("could not change password: %w", err)
	}

	// Optionally clear the forgot password code
	clearCodeQuery := `
		UPDATE forgotpassword
		SET code = NULL
		WHERE user_id = $1;
	`
	_, err = db.DB.Exec(clearCodeQuery, userid)
	// if err != nil {
	// 	log.Printf("Error clearing forgot password code for user %s: %v", email, err)
	// 	// Don't return this as critical — just log it.
	// }

	return nil
}

func (db *RealDB) GetUsers() ([]models.User, error) {
	query := `
	SELECT 
		user_id, 
		name, 
		email, 
		password, 
		code, 
		created_at, 
		updated_at
	FROM users;
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.UserID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.Code,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *RealDB) AddStatementWithAssociation(statement models.Statement) (string, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	INSERT INTO statements (user_id, statement_id, content, association)
		VALUES ($1, $2, $3, $4)
		RETURNING statement_id;`

	statementID, err := functions.GenerateUserID()
	if err != nil {
		return "", err
	}

	var pk string
	err = db.DB.QueryRow(query, statement.UserID, statementID, statement.Content, statement.Association).Scan(&pk)

	if err != nil {
		return "", fmt.Errorf("could not add statement: %w", err)
	}

	return statementID, nil
}

func (db *RealDB) SearchStatements(userID string, searchTerm string) ([]models.Statement, error) {
	query := `
	SELECT id, user_id, content, association, created_at, updated_at
	FROM statements
	WHERE user_id = $1
	AND content ILIKE $2
	ORDER BY created_at DESC;`

	rows, err := db.DB.Query(query, userID, "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statements []models.Statement
	for rows.Next() {
		var s models.Statement
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Content,
			&s.Association,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

func (db *RealDB) GetStatements(userID string) ([]models.Statement, error) {
	query := `
	SELECT id, statement_id, user_id, content, association, created_at, updated_at
	FROM statements
	WHERE user_id = $1
	ORDER BY created_at ASC;`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statements []models.Statement
	for rows.Next() {
		var s models.Statement
		err := rows.Scan(
			&s.ID,
			&s.StatementID,
			&s.UserID,
			&s.Content,
			&s.Association,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

func (db *RealDB) DeleteStatement(statementID string, userID string) error {
	query := `
	DELETE FROM statements
	WHERE statement_id = $1 AND user_id = $2;`

	result, err := db.DB.Exec(query, statementID, userID)
	if err != nil {
		return fmt.Errorf("could not delete statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("statement not found or does not belong to user")
	}

	return nil
}

func (db *RealDB) UpdateStatement(statementID string, userID string, content string, association string) error {
	query := `
	UPDATE statements
	SET content = $1,
		association = $2,
		updated_at = NOW()
	WHERE statement_id = $3 AND user_id = $4;`

	result, err := db.DB.Exec(query, content, association, statementID, userID)
	if err != nil {
		return fmt.Errorf("could not update statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("statement not found or does not belong to user")
	}

	return nil
}
