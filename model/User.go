package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Role     string    `json:"role"` // vodič, turista, administrator
	IsBlocked bool		`json:"isBlocked"`
}

func (user *User) BeforeCreate() {
	user.ID = uuid.New()
}
