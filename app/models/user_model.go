package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	gorm.Model
	ID           uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	FirstName    string    `db:"firstname" json:"firstname" validate:"required,lte=255"`
	LastName     string    `db:"lastname" json:"lastname" validate:"required,lte=255"`
	Email        string    `db:"email" json:"email" validate:"required,email,lte=255"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int       `db:"user_status" json:"user_status" validate:"required,len=1"`
	UserRole     string    `db:"user_role" json:"user_role" validate:"required,lte=25"`
}
