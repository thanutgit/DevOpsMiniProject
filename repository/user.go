package repository

import (
	"DevOpsMiniProject/entity"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ent entity.User) (entity.User, error)
	DeleteUser(req string) error
	GetAllUser() ([]entity.User, error)
	GetStatusDB() string
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
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("Row found:", result.RowsAffected) // เพิ่มบรรทัดนี้
	return users, nil
}

func (u userRepository) GetStatusDB() string {
	sqlDB, err := u.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		return "🔴disconnected"
	} else {
		return "🟢connect"
	}
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
