package service

import (
	"github.com/dan1M/gorm-user-auth-tuto/model"
	"github.com/kjk/betterguid"
	"gorm.io/gorm"
)

type RtService struct {
	db *gorm.DB
}

func NewRtService(db *gorm.DB) *RtService {
	return &RtService{db: db}
}

func (s *RtService) GetRT(id int) (*model.RefreshToken, error) {
	var rt model.RefreshToken

	err := s.db.First(&rt, id).Error
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (rt *RtService) CreateRT(data *model.RtCreateDTO) (*model.RefreshToken, error) {
	hash := betterguid.New()

	token := &model.RefreshToken{
		UserID: data.UserID,
		Ip:     data.Ip,
		Hash:   hash,
	}

	err := rt.db.Save(token).Error
	if err != nil {
		return nil, err
	}

	return token, nil
}
