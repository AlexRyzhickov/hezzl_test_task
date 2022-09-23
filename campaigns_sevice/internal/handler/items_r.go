package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"hezzl_test_task/internal/models"
	"hezzl_test_task/internal/service"
	"hezzl_test_task/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type ReadItemsHandler struct {
	Service ReadItemsService
}

type ReadItemsService interface {
	ReadItem(ctx context.Context, id, campaignId int) (*models.Item, error)
}

func (h *ReadItemsHandler) Method() string {
	return http.MethodGet
}

func (h *ReadItemsHandler) Path() string {
	return "/item/id/{id}/campaignId/{campaignId}"
}

func (h *ReadItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	item, err := h.Service.ReadItem(r.Context(), id, campaignId)
	if err != nil {
		log.Println(err)
		if err.Error() == service.NotFoundError {
			writeResponse(w, http.StatusNotFound, utils.NotFoundMsg())
			return
		}
		writeResponse(w, http.StatusInternalServerError, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, item)
}
