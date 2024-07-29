package repository

import (
	"time"

	"lazy-auth/app/model"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID               string `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	RoleID           string
	Role             Role
	Username         string `gorm:"uniqueIndex:idx_username"`
	Email            string `gorm:"uniqueIndex:idx_email"`
	PasswordHash     string
	DisplayName      string
	FirstName        string
	LastName         string
	VerifyFlag       bool `gorm:"default:false"`
	Ticket           string
	TicketExpiresAt  time.Time
	LastAccessAt     time.Time
	ChangePasswordAt time.Time
}

type UserRepository interface {
	GetMany(query model.QueryUser) ([]User, int, error)
	GetById(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByTicket(ticket string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	DaleteById(id string) error
}
