package service

import (
	"testing"
	"user-ms/src/dto"
	"user-ms/src/model"
	"user-ms/src/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type FollowingTestsSuite struct {
	suite.Suite
	followerRepositoryMock         *repository.FollowerRepositoryMock
	followingRequestRepositoryMock *repository.FollowingRequestRepositoryMock
	userRepositoryMock             *repository.UserRepositoryMock
	service                        IFollowingService
}

func TestUserFollowingTestsSuite(t *testing.T) {
	suite.Run(t, new(FollowingTestsSuite))
}

func (suite *FollowingTestsSuite) SetupSuite() {
	suite.followerRepositoryMock = new(repository.FollowerRepositoryMock)
	suite.followingRequestRepositoryMock = new(repository.FollowingRequestRepositoryMock)
	suite.userRepositoryMock = new(repository.UserRepositoryMock)
	suite.service = NewFollowingService(suite.followerRepositoryMock, suite.followingRequestRepositoryMock, suite.userRepositoryMock, nil)
}

func (suite *FollowingTestsSuite) TestNewFollowingTestsService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *FollowingTestsSuite) TestUpdateFollowerRequest() {
	pendingStatus := model.PENDING
	acceptedStatus := model.ACCEPTED

	followingRequestDTO := dto.FollowingRequestDTO{
		FollowingId:   1234,
		FollowerId:    2222,
		RequestStatus: int(pendingStatus),
	}

	updatedReqDTO := dto.FollowingRequestDTO{
		FollowingId:   1234,
		FollowerId:    2222,
		RequestStatus: int(acceptedStatus),
	}

	updatedReq := model.FollowingRequest{
		FollowingId:   1234,
		FollowerId:    2222,
		RequestStatus: acceptedStatus,
	}
	suite.followingRequestRepositoryMock.On("UpdateFollowingRequest", mock.AnythingOfType("int"), mock.AnythingOfType("*model.FollowingRequest")).Return(&updatedReq, nil).Once()
	suite.followerRepositoryMock.On("AddFollower", mock.AnythingOfType("*model.Follower")).Return(1, nil).Once()

	reqDTO, err := suite.service.UpdateRequest(followingRequestDTO.FollowingId, &updatedReqDTO)

	assert.Equal(suite.T(), int(acceptedStatus), reqDTO.RequestStatus)
	assert.Equal(suite.T(), nil, err)

}
