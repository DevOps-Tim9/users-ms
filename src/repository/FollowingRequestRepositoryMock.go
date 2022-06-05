package repository

import (
	"github.com/stretchr/testify/mock"
	"user-ms/src/model"
)

type FollowingRequestRepositoryMock struct {
	mock.Mock
}

func (f FollowingRequestRepositoryMock) AddFollowingRequest(request *model.FollowingRequest) (int, error) {
	args := f.Called(request)
	if args.Get(1) == nil {
		return args.Get(0).(int), nil
	}
	return -1, args.Get(1).(error)
}

func (f FollowingRequestRepositoryMock) UpdateFollowingRequest(i int, request *model.FollowingRequest) (*model.FollowingRequest, error) {
	args := f.Called(i, request)
	if args.Get(1) == nil {
		return args.Get(0).(*model.FollowingRequest), nil
	}
	return nil, args.Get(1).(error)
}

func (f FollowingRequestRepositoryMock) DeleteFollowingRequest(i int) error {
	args := f.Called(i)
	return args.Error(0)
}

func (f FollowingRequestRepositoryMock) GetRequests() []model.FollowingRequest {
	panic("implement me")
}

func (f FollowingRequestRepositoryMock) GetRequestsByFollowingID(i int) []model.FollowingRequest {
	panic("implement me")
}
