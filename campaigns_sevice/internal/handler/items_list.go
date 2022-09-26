package handler

import (
	"context"
	"github.com/rs/zerolog"
	"hezzl_test_task/campaigns_sevice/internal/models"
	"net/http"
)

type ListItemsHandler struct {
	Service ListItemsService
	logger  *zerolog.Logger
}

func NewListItemsHandler(s ListItemsService, logger *zerolog.Logger) *ListItemsHandler {
	return &ListItemsHandler{
		Service: s,
		logger:  logger,
	}
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
		h.logger.Error().Err(err).Msg("")
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, list)
}
