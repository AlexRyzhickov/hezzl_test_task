package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"hezzl_test_task/campaigns_sevice/internal/models"
	"hezzl_test_task/campaigns_sevice/internal/service"
	"hezzl_test_task/campaigns_sevice/internal/utils"
	"net/http"
	"strconv"
)

type DeleteItemHandler struct {
	Service DeleteItemService
	logger  *zerolog.Logger
}

func NewDeleteItemHandler(s DeleteItemService, logger *zerolog.Logger) *DeleteItemHandler {
	return &DeleteItemHandler{
		Service: s,
		logger:  logger,
	}
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
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	campaignId, err := strconv.Atoi(chi.URLParam(r, "campaignId"))
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	err = h.Service.DeleteItem(r.Context(), id, campaignId)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		if err.Error() == service.NotFoundError {
			writeResponse(w, http.StatusNotFound, utils.NotFoundMsg())
			return
		}
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, struct {
		Id         int  `json:"id"`
		CampaignId int  `json:"campaign_id"`
		Removed    bool `json:"removed"`
	}{Id: id, CampaignId: campaignId, Removed: true})
}
