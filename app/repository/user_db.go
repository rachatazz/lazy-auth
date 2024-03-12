package repository

import (
	"strings"

	"lazy-auth/app/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{db}
}

func (r userRepository) GetMany(query model.QueryUser) ([]User, int, error) {
	tx := r.db.Model(&User{})

	if query.OrderBy != nil {
		orderStr := *query.OrderBy + " "
		if query.SortBy != nil {
			orderStr = orderStr + strings.ToUpper(*query.SortBy)
		} else {
			orderStr = orderStr + "DESC"
		}
		tx = tx.Order(orderStr)
	}

	if query.ID != nil {
		tx = tx.Where("id = ?", *query.ID)
	}

	if query.Keyword != nil {
		tx = tx.Where(
			"display_name ? OR first_name ILIKE ? OR last_name ILIKE ?",
			"%"+*query.Keyword+"%",
			"%"+*query.Keyword+"%",
			"%"+*query.Keyword+"%",
		)
	}

	if query.Limit != nil {
		tx = tx.Limit(*query.Limit)
	}

	if query.Offset != nil {
		tx = tx.Offset(*query.Offset)
	}

	var users []User
	tx.Find(&users)

	var total int64
	tx.Limit(-1).Count(&total)
	if tx.Error != nil {
		return nil, int(total), tx.Error
	}
	return users, int(total), nil
}

func (r userRepository) GetById(id string) (*User, error) {
	var user User
	tx := r.db.Preload("Role").Where("id = ?", id).Take(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (r userRepository) GetByUsername(username string) (*User, error) {
	var user User
	tx := r.db.Where("username = ?", username).Take(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (r userRepository) Create(user *User) error {
	tx := r.db.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r userRepository) Update(user *User) error {
	tx := r.db.Save(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r userRepository) DaleteById(id string) error {
	tx := r.db.Delete(&User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
