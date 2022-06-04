package mapper

import (
	"user-ms/src/dto"
	"user-ms/src/model"
)

func RegistrationRequestDTOToUser(registeredUserDto *dto.RegistrationRequestDTO) *model.User {

	var user model.User
	user.Username = registeredUserDto.Username
	user.FirstName = registeredUserDto.FirstName
	user.LastName = registeredUserDto.LastName
	user.DateOfBirth = registeredUserDto.DateOfBirth
	user.Gender = registeredUserDto.Gender
	user.Email = registeredUserDto.Email
	user.PhoneNumber = registeredUserDto.PhoneNumber
	user.Password = registeredUserDto.Password
	return &user
}

func UserToDTO(userEntity *model.User) *dto.UserResponseDTO {
	var user dto.UserResponseDTO

	user.ID = userEntity.ID
	user.Auth0ID = userEntity.Auth0ID
	user.FirstName = userEntity.FirstName
	user.LastName = userEntity.LastName
	user.Email = userEntity.Email
	user.PhoneNumber = userEntity.PhoneNumber
	user.Gender = userEntity.Gender
	user.Username = userEntity.Username
	user.DateOfBirth = userEntity.DateOfBirth
	user.Biography = userEntity.Biography
	user.Education = userEntity.Education
	user.WorkExperience = userEntity.WorkExperience
	user.Skills = userEntity.Skills
	user.Interests = userEntity.Interests
	user.Public = userEntity.Public

	return &user
}

func UserUpdateDTOToUser(userUpdateDTO *dto.UserUpdateDTO) *model.User {
	var user model.User

	user.ID = userUpdateDTO.ID
	user.FirstName = userUpdateDTO.FirstName
	user.LastName = userUpdateDTO.LastName
	user.Email = userUpdateDTO.Email
	user.PhoneNumber = userUpdateDTO.PhoneNumber
	user.Gender = userUpdateDTO.Gender
	user.Username = userUpdateDTO.Username
	user.DateOfBirth = userUpdateDTO.DateOfBirth
	user.Biography = userUpdateDTO.Biography
	user.Education = userUpdateDTO.Education
	user.WorkExperience = userUpdateDTO.WorkExperience
	user.Skills = userUpdateDTO.Skills
	user.Interests = userUpdateDTO.Interests
	user.Public = userUpdateDTO.Public

	return &user
}
