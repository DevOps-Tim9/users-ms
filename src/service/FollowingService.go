package service

import (
	"errors"
	"strconv"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/model"
	"user-ms/src/repository"
)

type FollowingService struct {
	FollowerRepository         repository.IFollowerRepository
	FollowingRequestRepository repository.IFollowingRequestRepository
}

type IFollowingService interface {
	CreateRequest(*dto.FollowingRequestDTO) (int, error)
}

func NewFollowingService(followerRepository repository.IFollowerRepository, followingRequestRepository repository.IFollowingRequestRepository) IFollowingService {
	return &FollowingService{
		followerRepository,
		followingRequestRepository,
	}
}

func (service *FollowingService) CreateRequest(request *dto.FollowingRequestDTO) (int, error) {
	followingRequestId, err := service.FollowingRequestRepository.AddFollowingRequest(mapper.FollowingDTOToRequestFollower(request))
	if err != nil {
		return -1, errors.New("can't create the request")
	}
	return followingRequestId, nil
}

func (service *FollowingService) UpdateRequest(reqId int, request *dto.FollowingRequestDTO) (string, error) {
	followingRequestId, err := service.FollowingRequestRepository.UpdateFollowingRequest(reqId, mapper.FollowingDTOToRequestFollower(request))
	if err != nil {
		return string(-1), errors.New("can't create the request")
	}
	status := model.RequestStatus(request.RequestStatus)
	if model.ACCEPTED == status {
		_, _ = service.FollowerRepository.AddFollower(mapper.FollowingDTOToFollower(request))
	}
	return "Updated request " + strconv.Itoa(followingRequestId), nil
}

func (service *FollowingService) GetRequests() ([]model.FollowingRequest, error) {
	requests := service.FollowingRequestRepository.GetRequests()
	return requests, nil
}

func (service *FollowingService) GetRequestsByFollowingID(id int) ([]model.FollowingRequest, error) {
	requests := service.FollowingRequestRepository.GetRequestsByFollowingID(id)
	return requests, nil
}

func (service *FollowingService) CreateFollower(request *dto.FollowingRequestDTO) (int, error) {
	followerId, err := service.FollowerRepository.AddFollower(mapper.FollowingDTOToFollower(request))
	if err != nil {
		return -1, errors.New("can't create the request")
	}
	return followerId, nil
}

func (service *FollowingService) GetFollowers(id int) ([]model.Follower, error) {
	followers := service.FollowerRepository.GetFollowers(id)
	return followers, nil
}

func (service *FollowingService) GetFollowing(id int) ([]model.Follower, error) {
	following := service.FollowerRepository.GetFollowing(id)
	return following, nil
}

func (service *FollowingService) RemoveFollowing(id int, followingId int) error {
	service.FollowerRepository.RemoveFollowing(id, followingId)
	return nil
}
