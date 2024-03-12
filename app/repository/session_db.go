package repository

import (
	"time"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return sessionRepository{db}
}

func (r sessionRepository) Create(session *Session) error {
	tx := r.db.Create(&session)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r sessionRepository) GetById(id string) (*Session, error) {
	var session Session
	tx := r.db.Where("id = ? AND expires_at > ?", id, time.Now()).Take(&session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func (r sessionRepository) GetByRefreshToken(refreshToken string) (*Session, error) {
	var session Session
	tx := r.db.Where("refresh_token = ? ", refreshToken).Take(&session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func (r sessionRepository) Update(session *Session) error {
	tx := r.db.Save(&session)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r sessionRepository) DeleteById(id string) error {
	tx := r.db.Delete(&Session{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r sessionRepository) DeleteByUserId(id string) error {
	tx := r.db.Where("user_id = ?", id).Delete(&Session{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
