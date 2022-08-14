package users

import (
	"dialogv2/internal/database/entities"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r Repository) GetById(id string) (*entities.User, error) {
	var user entities.User
	result := r.DB.Take(&user, "uid = ?", id)

	return &user, result.Error
}

func (r Repository) Create(user *entities.User) (err error) {
	err = r.DB.Create(user).Error

	if err != nil {
		return
	}

	return
}

func (r Repository) Update(user *entities.User) (err error) {
	err = r.DB.Updates(user).Error

	if err != nil {
		return
	}

	return
}

func (r Repository) Delete(user *entities.User) (err error) {
	err = r.DB.Delete(user).Error

	if err != nil {
		return
	}

	return
}
