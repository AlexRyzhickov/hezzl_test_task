package service

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"hezzl_test_task/internal/models"
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
	return s.db.Create(model).Error
}

func (s Service) DeleteItem(ctx context.Context, id, campaignId int) error {
	item := models.Item{Id: id, CampaignId: campaignId}
	tx := s.db.Delete(&item)
	if tx.RowsAffected == 0 {
		return errors.New("Object for removing was not found")
	}
	return tx.Error
}

func (s Service) ReadItems(ctx context.Context) (*[]models.Item, error) {
	findContacts := []models.Item{}
	err := s.db.Find(&findContacts).Error
	return &findContacts, err
}

func (s Service) UpdateItem(ctx context.Context, id, campaignId int, values map[string]interface{}) error {
	item := models.Item{Id: id, CampaignId: campaignId}
	return s.db.Model(&item).Updates(values).Error
}
