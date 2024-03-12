package repository

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID           string `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID       string
	User         User
	RefreshToken string
	ExpiresAt    time.Time
}

type SessionRepository interface {
	Create(session *Session) error
	GetById(id string) (*Session, error)
	GetByRefreshToken(refreshToken string) (*Session, error)
	Update(session *Session) error
	DeleteById(id string) error
	DeleteByUserId(id string) error
}
