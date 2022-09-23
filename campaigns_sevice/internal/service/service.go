package service

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"hezzl_test_task/internal/models"
)

type Service struct {
	db   *gorm.DB
	repo Repository
}

type Repository interface {
	Load(string) ([]models.Item, error)
	Store(context.Context, string, []models.Item) error
}

func NewService(db *gorm.DB, repo Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

func (s Service) CreateItem(ctx context.Context, model *models.Item) error {
	priority := 0
	row := s.db.Table("items").Select("max(priority)").Row()
	err := row.Scan(&priority)
	if err != nil {
		priority = 0
	}
	model.Priority = priority + 1
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
	items, err := s.repo.Load("items")
	if err != nil {
		return &items, nil
	}
	items = []models.Item{}
	err = s.db.Find(&items).Error
	if err == nil {
		err = s.repo.Store(ctx, "items", items)
	}
	return &items, err
}

func (s Service) UpdateItem(ctx context.Context, id, campaignId int, values map[string]interface{}) (*models.Item, error) {
	item := models.Item{Id: id, CampaignId: campaignId}
	return &item, s.db.Model(&item).Updates(values).Error
}
