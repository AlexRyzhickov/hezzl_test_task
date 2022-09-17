package service

import (
	"context"
	"hezzl_test_task/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) CreateItem(ctx context.Context, model *models.Item) error {
	return s.db.FirstOrCreate(model).Error
}

func (s Service) DeleteItem(ctx context.Context, id, campaignId int) error {
	return s.db.Delete(&models.Item{}).Error
}

func (s Service) ReadItem(ctx context.Context, id, campaignId int) (*models.Item, error) {
	return nil, nil
}

func (s Service) ReadItems(ctx context.Context) (*[]models.Item, error) {
	return nil, nil
}

func (s Service) UpdateItem(ctx context.Context, model *models.Item) error {
	return nil
}
