package model

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

type Gender int

const (
	Male Gender = iota
	Female
)

type User struct {
	ID             int     `json:"id"`
	Auth0ID        string  `json:"auth0_id"`
	FirstName      string  `json:"first_name" validate:"required"`
	LastName       string  `json:"last_name" validate:"required"`
	Email          string  `json:"email" validate:"required,email" gorm:"unique"`
	Password       string  `json:"password" validate:"required"`
	PhoneNumber    string  `json:"phone_number"`
	Gender         *Gender `json:"gender" validate:"required"`
	Username       string  `json:"user_name" gorm:"unique" validate:"required"`
	DateOfBirth    float32 `json:"date_od_birth"`
	Biography      string  `json:"biography"`
	Education      string  `json:"education"`
	WorkExperience string  `json:"work_experience"`
	Skills         string  `json:"skills"`
	Interests      string  `json:"interests"`
	Active         bool    `json:"active"`
	Public         bool    `json:"public"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}
