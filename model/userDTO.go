package model

import (
	"github.com/dan1M/gorm-user-auth-tuto/model"
)

type UserCreateDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateDTO struct {
	Email string `json:"email"`
}

func (s *UserService) CreateUser(data *model.UserCreateDTO) (*model.User, error) {

}
