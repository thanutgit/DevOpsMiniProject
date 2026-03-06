package repository

import (
	"DevOpsMiniProject/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ent entity.User) (entity.User, error)
	DeleteUser(req string) error
}

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) CreateUser(ent entity.User) (entity.User, error) {
	result := u.db.Create(&ent)
	if result.Error != nil {
		return entity.User{}, result.Error
	}
	return ent, nil
}

func (u userRepository) DeleteUser(req string) error {
	result := u.db.Delete(entity.User{}, "username = ?", req)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
