package repository

import (
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

func (repo *UserRepositoryMock) DeleteUser(id int) error {
	panic("implement me!")
}
