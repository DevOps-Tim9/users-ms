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
