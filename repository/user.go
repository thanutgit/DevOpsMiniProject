package repository

import (
	"DevOpsMiniProject/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ent entity.User) (entity.User, error)
	DeleteUser(req string) error
	GetAllUser() ([]entity.User, error)
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

func (u userRepository) GetAllUser() ([]entity.User, error) {
	var users []entity.User
	result := u.db.Find(&users)
	if result != nil {
		return nil, result.Error
	}

	return users, nil
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
