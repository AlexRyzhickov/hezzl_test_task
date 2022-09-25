package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"hezzl_test_task/internal/models"
	"hezzl_test_task/internal/service"
	"hezzl_test_task/internal/utils"
	"net/http"
	"strconv"
)

type UpdateItemHandler struct {
	Service UpdateItemService
	logger  *zerolog.Logger
}

func NewUpdateItemHandler(s UpdateItemService, logger *zerolog.Logger) *UpdateItemHandler {
	return &UpdateItemHandler{
		Service: s,
		logger:  logger,
	}
}

type UpdateItemService interface {
	UpdateItem(ctx context.Context, id, campaignId int, values map[string]interface{}) (*models.Item, error)
}

func (h *UpdateItemHandler) Method() string {
	return http.MethodPatch
}

func (h *UpdateItemHandler) Path() string {
	return "/item/update/id/{id}/campaignId/{campaignId}"
}

func (h *UpdateItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	values := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	item, err := h.Service.UpdateItem(r.Context(), id, campaignId, values)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		if err.Error() == service.NotFoundError {
			writeResponse(w, http.StatusNotFound, utils.NotFoundMsg())
			return
		}
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, item)
}
