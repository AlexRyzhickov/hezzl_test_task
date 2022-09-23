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

type UpdateItemHandler struct {
	Service UpdateItemService
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
	values := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Bad request"})
		return
	}
	item, err := h.Service.UpdateItem(r.Context(), id, campaignId, values)
	if err != nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, models.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, item)
}
