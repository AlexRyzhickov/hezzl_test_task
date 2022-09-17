package handler

import (
	"context"
	"hezzl_test_task/internal/models"
	"net/http"
)

type CreateItemHandler struct {
	Service CreateItemService
}

type CreateItemService interface {
	CreateItem(ctx context.Context, model *models.Item) error
}

func (h *CreateItemHandler) Method() string {
	return http.MethodPost
}

func (h *CreateItemHandler) Path() string {
	return "/item/create"
}

func (h *CreateItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type DeleteItemHandler struct {
	Service DeleteItemService
}

type DeleteItemService interface {
	DeleteItem(ctx context.Context, id, campaignId int) error
}

func (h *DeleteItemHandler) Method() string {
	return http.MethodDelete
}

func (h *DeleteItemHandler) Path() string {
	return "/item/remove/{id,campaignId}"
}

func (h *DeleteItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type ReadItemHandler struct {
	Service ReadItemService
}

type ReadItemService interface {
	ReadItem(ctx context.Context, id, campaignId int) (*models.Item, error)
}

func (h *ReadItemHandler) Method() string {
	return http.MethodGet
}

func (h *ReadItemHandler) Path() string {
	return "/item/{id,campaignId}"
}

func (h *ReadItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type ListItemsHandler struct {
	Service ListItemsService
}

type ListItemsService interface {
	ReadItems(ctx context.Context) (*[]models.Item, error)
}

func (h *ListItemsHandler) Method() string {
	return http.MethodGet
}

func (h *ListItemsHandler) Path() string {
	return "/item/{id,campaignId}"
}

func (h *ListItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type UpdateItemHandler struct {
	Service UpdateItemService
}

type UpdateItemService interface {
	UpdateItem(ctx context.Context, model *models.Item) error
}

func (h *UpdateItemHandler) Method() string {
	return http.MethodPatch
}

func (h *UpdateItemHandler) Path() string {
	return "/item/update/{id,campaignId}"
}

func (h *UpdateItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
