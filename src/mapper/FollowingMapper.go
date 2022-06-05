package mapper

import (
	"user-ms/src/dto"
	"user-ms/src/model"
)

func FollowingDTOToFollower(followingRequestDTO *dto.FollowingRequestDTO) *model.Follower {
	var follower model.Follower
	follower.FollowerId = followingRequestDTO.FollowerId
	follower.FollowingId = followingRequestDTO.FollowingId
	return &follower
}

func FollowingDTOToRequestFollower(followingRequestDTO *dto.FollowingRequestDTO) *model.FollowingRequest {
	var followingRequest model.FollowingRequest
	followingRequest.FollowerId = followingRequestDTO.FollowerId
	followingRequest.FollowingId = followingRequestDTO.FollowingId
	followingRequest.RequestStatus = model.RequestStatus(followingRequestDTO.RequestStatus)
	return &followingRequest
}

func FollowingRequestToFollower(request *model.FollowingRequest) *model.Follower {
	var follower model.Follower
	follower.FollowingId = request.FollowerId
	follower.FollowerId = request.FollowingId
	return &follower
}
