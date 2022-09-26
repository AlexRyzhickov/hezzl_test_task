package service

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"hezzl_test_task/campaigns_sevice/internal/models"
	"strconv"
)

const NotFoundError = "errors.item.notFound"

type Service struct {
	db   *gorm.DB
	repo Repository
}

type Repository interface {
	Load(string) (models.Item, bool)
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
	err := s.db.Transaction(func(tx *gorm.DB) error {
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
		return nil
	})
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
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New(NotFoundError)
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
	tx := s.db.Model(&item).Updates(values)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New(NotFoundError)
	}
	err := storeDataInRedis(ctx, id, campaignId, &(s.repo), item)
	if err != nil {
		return nil, err
	}
	return &item, err
}

func (s Service) ReadItem(ctx context.Context, id, campaignId int) (*models.Item, error) {
	item, ok := s.repo.Load(calculateKey(id, campaignId))
	if ok == true {
		return &item, nil
	}
	item = models.Item{Id: id, CampaignId: campaignId}
	tx := s.db.Find(&item)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New(NotFoundError)
	}
	err := s.repo.Store(ctx, calculateKey(id, campaignId), item)
	if err != nil {
		return nil, err
	}
	return &item, nil
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
