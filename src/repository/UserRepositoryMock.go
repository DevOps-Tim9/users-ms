package repository

import (
	"user-ms/src/dto"
	"user-ms/src/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) GetBySearchParam(param string) ([]*dto.UserResponseDTO, error) {
	args := u.Called(param)
	if args.Get(1) == nil {
		return args.Get(0).([]*dto.UserResponseDTO), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserRepositoryMock) GetAll() ([]*dto.UserResponseDTO, error) {
	args := u.Called()
	if args.Get(1) == nil {
		return args.Get(0).([]*dto.UserResponseDTO), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserRepositoryMock) AddUser(user *model.User) (int, error) {
	args := u.Called(user)
	if args.Get(1) == nil {
		return args.Get(0).(int), nil
	}
	return -1, args.Get(1).(error)
}

func (u *UserRepositoryMock) DeleteUser(id int) error {
	panic("implement me")
}

func (u *UserRepositoryMock) Update(user *model.User) (*dto.UserResponseDTO, error) {
	args := u.Called(user)
	if args.Get(1) == nil {
		return args.Get(0).(*dto.UserResponseDTO), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserRepositoryMock) GetByEmail(email string) (*dto.UserResponseDTO, error) {
	args := u.Called(email)
	if args.Get(1) == nil {
		return args.Get(0).(*dto.UserResponseDTO), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserRepositoryMock) GetByID(id int) (*model.User, error) {
	args := u.Called(id)
	if args.Get(1) == nil {
		return args.Get(0).(*model.User), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserRepositoryMock) GetByAuth0ID(id string) (*model.User, error) {
	args := u.Called(id)
	if args.Get(1) == nil {
		return args.Get(0).(*model.User), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserRepositoryMock) UnblockUser(blockingID int, userID int) error {
	args := u.Called(blockingID, userID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (u *UserRepositoryMock) GetBlockedUsers(userID int) []model.User {
	args := u.Called(userID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]model.User)
}

func (u *UserRepositoryMock) CreateAdmin([]model.User) {
}

func (u *UserRepositoryMock) GetByUsername(username string) []model.User {
	panic("")
}
