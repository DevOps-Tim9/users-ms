package repository

import (
	"user-ms/src/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type IUserRepository interface {
	AddUser(*model.User) (int, error)
	DeleteUser(int) error
}

func NewUserRepository(database *gorm.DB) IUserRepository {
	return &UserRepository{
		database,
	}
}

type UserRepository struct {
	Database *gorm.DB
}

func (repo *UserRepository) AddUser(user *model.User) (int, error) {
	result := repo.Database.Create(user)

	if result.Error != nil {
		return -1, result.Error
	}

	return user.ID, nil
}

func (repo *UserRepository) DeleteUser(id int) error {
	result := repo.Database.Delete(&model.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
