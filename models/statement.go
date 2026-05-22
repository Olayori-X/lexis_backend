package models

import "time"

type Statement struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	StatementID string    `db:"statement_id" json:"statement_id"`
	Content     string    `db:"not null;type:text" json:"content"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Association string    `db:"association;type:text" json:"associations,omitempty"`
}
