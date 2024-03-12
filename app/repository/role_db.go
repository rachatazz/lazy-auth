package repository

import (
	"fmt"

	"lazy-auth/app/model"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return roleRepository{db}
}

func (r roleRepository) Create(role Role) (*Role, error) {
	tx := r.db.Create(&role)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &role, nil
}

func (r roleRepository) GetAll(query model.QueryRole) ([]Role, int, error) {
	tx := r.db.Model(&Role{})

	sortBy := "created_at"
	if query.SortBy != nil {
		sortBy = *query.SortBy
	}

	orderBy := "DESC"
	if query.OrderBy != nil {
		orderBy = *query.OrderBy
	}
	tx = tx.Order(fmt.Sprintf("%v %v", sortBy, orderBy))

	if query.RoleType != nil {
		tx = tx.Where("role_type = ?", *query.RoleType)
	}

	limit := 100
	if query.Limit != nil {
		limit = *query.Limit
	}
	tx = tx.Limit(limit)

	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}
	tx = tx.Offset(offset)

	var roles []Role
	tx.Find(&roles)

	var total int64
	tx.Limit(-1).Count(&total)

	if tx.Error != nil {
		return nil, int(total), tx.Error
	}
	return roles, int(total), nil
}

func (r roleRepository) GetById(id string) (*Role, error) {
	var role Role
	tx := r.db.Where("id = ?", id).Take(&role)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &role, nil
}

func (r roleRepository) GetByName(name string) (*Role, error) {
	var role Role
	tx := r.db.Where("name = ?", name).Take(&role)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &role, nil
}

func (r roleRepository) Update(role Role) (*Role, error) {
	tx := r.db.Save(&role)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &role, nil
}

func (r roleRepository) DeleteById(id string) error {
	tx := r.db.Delete(&Role{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
