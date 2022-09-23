package service

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"hezzl_test_task/internal/models"
	"strconv"
)

type Service struct {
	db   *gorm.DB
	repo Repository
}

type Repository interface {
	Load(string) (models.Item, error)
	Store(context.Context, string, models.Item) error
	Delete(ctx context.Context, key string) error
}

func NewService(db *gorm.DB, repo Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

func (s Service) CreateItem(ctx context.Context, item *models.Item) error {
	priority := 0
	row := s.db.Table("items").Select("max(priority)").Row()
	err := row.Scan(&priority)
	if err != nil {
		priority = 0
	}
	item.Priority = priority + 1
	err = s.db.Create(item).Error
	if err != nil {
		return err
	}
	err = storeDataInRedis(ctx, item.Id, item.CampaignId, &(s.repo), *item)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) DeleteItem(ctx context.Context, id, campaignId int) error {
	item := models.Item{Id: id, CampaignId: campaignId}
	tx := s.db.Delete(&item)
	if tx.RowsAffected == 0 {
		return errors.New("Object for removing was not found")
	}
	if tx.Error != nil {
		return nil
	}
	key := calculateKey(id, campaignId)
	err := s.repo.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) ReadItems(ctx context.Context) (*[]models.Item, error) {
	items := []models.Item{}
	err := s.db.Find(&items).Error
	return &items, err
}

func (s Service) UpdateItem(ctx context.Context, id, campaignId int, values map[string]interface{}) (*models.Item, error) {
	item := models.Item{Id: id, CampaignId: campaignId}
	err := s.db.Model(&item).Updates(values).Error
	if err != nil {
		return nil, err
	}
	err = storeDataInRedis(ctx, id, campaignId, &(s.repo), item)
	if err != nil {
		return nil, err
	}
	return &item, err
}

func calculateKey(id, campaignId int) string {
	return strconv.Itoa(id) + "&" + strconv.Itoa(campaignId)
}

func storeDataInRedis(ctx context.Context, id, campaignId int, repo *Repository, item models.Item) error {
	key := calculateKey(id, campaignId)
	err := (*repo).Store(ctx, key, item)
	if err != nil {
		return err
	}
	return nil
}
