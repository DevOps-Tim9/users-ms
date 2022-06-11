package repository

import (
	"errors"
	"fmt"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type IUserRepository interface {
	AddUser(*model.User) (int, error)
	DeleteUser(int) error
	Update(*model.User) (*dto.UserResponseDTO, error)
	GetByID(int) (*model.User, error)
	GetByAuth0ID(string) (*model.User, error)
	GetByEmail(string) (*dto.UserResponseDTO, error)
	UnblockUser(int, int) error
	GetBlockedUsers(int) []model.User
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

func (repo *UserRepository) Update(user *model.User) (*dto.UserResponseDTO, error) {
	result := repo.Database.Save(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return mapper.UserToDTO(user), nil
}

func (repo *UserRepository) GetByID(id int) (*model.User, error) {
	userEntity := model.User{
		ID: id,
	}
	if err := repo.Database.Where("ID = ?", id).First(&userEntity).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("User with ID %d not found", id))
	}

	return &userEntity, nil
}

func (repo *UserRepository) GetByAuth0ID(id string) (*model.User, error) {
	userEntity := model.User{
		Auth0ID: id,
	}
	if err := repo.Database.Where("auth0_id = ?", id).First(&userEntity).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("User with auth0_ID %s not found", id))
	}

	return &userEntity, nil
}

func (repo *UserRepository) GetByEmail(email string) (*dto.UserResponseDTO, error) {
	userEntity := model.User{
		Email: email,
	}
	if err := repo.Database.Where("email = ?", email).First(&userEntity).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("User with email %s not found", email))
	}

	return mapper.UserToDTO(&userEntity), nil
}

func (repo *UserRepository) UnblockUser(blockingID int, userID int) error {
	result := repo.Database.Exec(fmt.Sprintf("delete from user_blocked where user_id = %d and blocked_id = %d", userID, blockingID))
	return result.Error
}

func (repo *UserRepository) GetBlockedUsers(userID int) []model.User {
	userEntity := model.User{
		ID: userID,
	}
	if err := repo.Database.Where("ID = ?", userID).Preload("Blocked").First(&userEntity).Error; err != nil {
	}

	return userEntity.Blocked
}
