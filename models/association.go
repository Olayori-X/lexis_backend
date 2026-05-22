package models

import "time"

type Association struct {
	ID          uint      `db:"primaryKey" json:"id"`
	StatementID uint      `db:"not null;index" json:"statement_id"`
	UserID      uint      `db:"not null;index" json:"user_id"`
	Content     string    `db:"not null;type:text" json:"content"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
