package repository

import (
	"DevOpsMiniProject/entity"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ent entity.User) (entity.User, error)
	DeleteUser(username string) error
	GetAllUser() ([]entity.User, error)
	PingDB(ctx context.Context) error
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

func (u userRepository) DeleteUser(username string) error {
	result := u.db.Delete(entity.User{}, "username = ?", username)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (u userRepository) GetAllUser() ([]entity.User, error) {
	var users []entity.User
	result := u.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("Row found:", result.RowsAffected) // เพิ่มบรรทัดนี้
	return users, nil
}

func (u userRepository) PingDB(ctx context.Context) error {
	sqlDB, err := u.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
