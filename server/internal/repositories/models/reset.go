package models

import "time"

type ResetToken struct {
	ID        int       `db:"id"`
	Token     string    `db:"token"`
	UserID    int       `db:"id_user"`
	ExpiresAt time.Time `db:"expires_at"`
}
