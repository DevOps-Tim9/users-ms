package mapper

import (
	"user-ms/src/dto"
	"user-ms/src/model"
)

func FollowingDTOToFollower(followingRequestDTO *dto.FollowingRequestDTO) *model.Follower {
	var follower model.Follower
	follower.ID = followingRequestDTO.ID
	follower.FollowerId = followingRequestDTO.FollowerId
	follower.FollowingId = followingRequestDTO.FollowingId
	return &follower
}

func FollowingDTOToRequestFollower(followingRequestDTO *dto.FollowingRequestDTO) *model.FollowingRequest {
	var followingRequest model.FollowingRequest
	followingRequest.ID = followingRequestDTO.ID
	followingRequest.FollowerId = followingRequestDTO.FollowerId
	followingRequest.FollowingId = followingRequestDTO.FollowingId
	followingRequest.RequestStatus = model.RequestStatus(followingRequestDTO.RequestStatus)
	return &followingRequest
}

func RequestToFollowingDTO(request *model.FollowingRequest) *dto.FollowingRequestDTO {
	var follower dto.FollowingRequestDTO
	follower.FollowingId = request.FollowerId
	follower.FollowerId = request.FollowingId
	follower.RequestStatus = int(request.RequestStatus)
	return &follower
}

func FollowingRequestToFollower(request *model.FollowingRequest) *model.Follower {
	var follower model.Follower
	follower.FollowingId = request.FollowerId
	follower.FollowerId = request.FollowingId
	return &follower
}
