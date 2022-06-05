package repository

import (
	"user-ms/src/dto"
	"user-ms/src/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
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