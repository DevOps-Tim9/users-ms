package model

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

type Follower struct {
	ID          int `json:"id"`
	FollowerId  int `json:"followers_id" gorm:"TYPE:integer REFERENCES users" validate:"required"`
	FollowingId int `json:"following_id" gorm:"TYPE:integer REFERENCES users" validate:"required"`
}

func (f *Follower) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func (f *Follower) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(f)
}

func (f *Follower) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}
