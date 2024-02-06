package models

import "gorm.io/gorm"

// SignUp struct to describe register a new user.
type SignUp struct {
	gorm.Model
	FirstName string `json:"firstname" validate:"required,lte=255"`
	LastName  string `json:"lastname" validate:"required,lte=255"`
	Email     string `json:"email" validate:"required,email,lte=255"`
	Password  string `json:"password" validate:"required,lte=255"`
	// UserRole string `json:"user_role" validate:"required,lte=25"`
}

// SignIn struct to describe login user.
type SignIn struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}
