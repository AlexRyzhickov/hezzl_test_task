package handler

import (
	"context"
	"hezzl_test_task/internal/models"
	"net/http"
)

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
