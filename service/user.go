package service

import (
	"github.com/dan1M/gorm-user-auth-tuto/model"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUser(id int) (*model.User, error) {
	var user model.User

	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUsers() ([]*model.User, error) {
	var users []*model.User
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteUser(id int) (*model.User, error) {
	var user model.User

	err := s.db.Delete(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	var users *model.User

	err := s.db.Where("email = ?", email).First(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) CreateUser(data *model.UserCreateDTO) (*model.User, error) {
	user := &model.User{
		Email:    data.Email,
		Password: data.Password,
	}
	err := s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(data *model.UserUpdateDTO) (*model.User, error) {
	user := &model.User{
		Email: data.Email,
	}
	err := s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
