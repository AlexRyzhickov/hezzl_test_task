package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"hezzl_test_task/internal/models"
	"log"
	"net/http"
	"strconv"
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
	item := &models.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	err = h.Service.CreateItem(r.Context(), item)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, item)
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
	return "/item/remove/id/{id}/campaignId/{campaignId}"
}

func (h *DeleteItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	campaignId, err := strconv.Atoi(chi.URLParam(r, "campaignId"))
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	err = h.Service.DeleteItem(r.Context(), id, campaignId)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, struct {
		Id         int  `json:"id"`
		CampaignId int  `json:"campaign_id"`
		Removed    bool `json:"removed"`
	}{Id: id, CampaignId: campaignId, Removed: true})
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
	return "/item/list"
}

func (h *ListItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.ReadItems(r.Context())
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, list)
}

type UpdateItemHandler struct {
	Service UpdateItemService
}

type UpdateItemService interface {
	UpdateItem(ctx context.Context, values map[string]interface{}) error
}

func (h *UpdateItemHandler) Method() string {
	return http.MethodPatch
}

func (h *UpdateItemHandler) Path() string {
	return "/item/update"
}

func (h *UpdateItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
