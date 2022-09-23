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
	return "/item/create/{campaignId}"
}

func (h *CreateItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	campaignId, err := strconv.Atoi(chi.URLParam(r, "campaignId"))
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	values := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	name, ok := values["name"].(string)
	if !ok {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	item := &models.Item{}
	item.CampaignId = campaignId
	item.Name = name
	err = h.Service.CreateItem(r.Context(), item)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, item)
}
