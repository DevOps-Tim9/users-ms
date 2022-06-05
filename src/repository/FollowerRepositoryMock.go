package repository

import (
	"github.com/stretchr/testify/mock"
	"user-ms/src/model"
)

type FollowerRepositoryMock struct {
	mock.Mock
}

func (f FollowerRepositoryMock) AddFollower(follower *model.Follower) (int, error) {
	args := f.Called(follower)
	if args.Get(1) == nil {
		return args.Get(0).(int), nil
	}
	return -1, args.Get(1).(error)
}

func (f FollowerRepositoryMock) DeleteFollower(i int) error {
	args := f.Called(i)
	return args.Error(0)
}

func (f FollowerRepositoryMock) GetFollowing(i int) []model.Follower {
	panic("implement me")
}

func (f FollowerRepositoryMock) GetFollowers(i int) []model.Follower {
	panic("implement me")
}

func (f FollowerRepositoryMock) RemoveFollowing(i int, i2 int) error {
	args := f.Called(i, i2)
	return args.Error(0)
}
