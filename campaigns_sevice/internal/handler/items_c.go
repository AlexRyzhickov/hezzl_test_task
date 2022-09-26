package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"hezzl_test_task/campaigns_sevice/internal/models"

	"net/http"
	"strconv"
)

type CreateItemHandler struct {
	Service CreateItemService
	logger  *zerolog.Logger
}

func NewCreateItemHandler(s CreateItemService, logger *zerolog.Logger) *CreateItemHandler {
	return &CreateItemHandler{
		Service: s,
		logger:  logger,
	}
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
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	values := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	name, ok := values["name"].(string)
	if !ok {
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	item := &models.Item{}
	item.CampaignId = campaignId
	item.Name = name
	err = h.Service.CreateItem(r.Context(), item)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, item)
}
