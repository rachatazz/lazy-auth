package repository

import (
	"lazy-auth/app/model"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          string `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name        string `gorm:"uniqueIndex:idx_name"`
	Description string
	Users       []User
}

type RoleRepository interface {
	Create(role Role) (*Role, error)
	GetAll(query model.QueryRole) ([]Role, int, error)
	GetById(id string) (*Role, error)
	GetByName(name string) (*Role, error)
	Update(role Role) (*Role, error)
	DeleteById(id string) error
}
