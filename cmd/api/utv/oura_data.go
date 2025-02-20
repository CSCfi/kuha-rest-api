package utvapi

import (
	"context"
	"net/http"

	"github.com/DeRuina/KUHA-REST-API/internal/store/utv"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

// Validation structs
type OuraGetDatesParams struct {
	UserID     string `form:"user_id" validate:"required,uuid4"`
	AfterDate  string `form:"after_date" validate:"omitempty,datetime=2006-01-02"`
	BeforeDate string `form:"before_date" validate:"omitempty,datetime=2006-01-02"`
}

type OuraGetTypesParams struct {
	UserID string `form:"user_id" validate:"required,uuid4"`
	Date   string `form:"date" validate:"required,datetime=2006-01-02"`
}

type OuraGetDataParams struct {
	UserID string `form:"user_id" validate:"required,uuid4"`
	Date   string `form:"date" validate:"required,datetime=2006-01-02"`
	Key    string `form:"key" validate:"omitempty,alphanum"`
}

type OuraDataHandler struct {
	store utv.OuraData
}

// NewOuraDataHandler initializes OuraData handler
func NewOuraDataHandler(store utv.OuraData) *OuraDataHandler {
	return &OuraDataHandler{store: store}
}

// Get available dates from Oura data (with optional filtering)
func (h *OuraDataHandler) GetDates(w http.ResponseWriter, r *http.Request) {
	err := utils.ValidateParams(r, []string{"user_id", "after_date", "before_date"})
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := OuraGetDatesParams{
		UserID:     r.URL.Query().Get("user_id"),
		AfterDate:  r.URL.Query().Get("after_date"),
		BeforeDate: r.URL.Query().Get("before_date"),
	}

	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if params.AfterDate != "" && params.BeforeDate != "" && params.AfterDate > params.BeforeDate {
		utils.UnprocessableEntityResponse(w, r, utils.ErrInvalidDateRange)
		return
	}

	dates, err := h.store.GetDates(context.Background(), params.UserID, &params.AfterDate, &params.BeforeDate)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if len(dates) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dates)
}

// Get all JSON keys (types) from Oura data on a specific date
func (h *OuraDataHandler) GetTypes(w http.ResponseWriter, r *http.Request) {
	err := utils.ValidateParams(r, []string{"user_id", "date"})
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := OuraGetTypesParams{
		UserID: r.URL.Query().Get("user_id"),
		Date:   r.URL.Query().Get("date"),
	}

	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	types, err := h.store.GetTypes(context.Background(), params.UserID, params.Date)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if len(types) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types)
}

// Get all data or a specific type for a given user and date
func (h *OuraDataHandler) GetData(w http.ResponseWriter, r *http.Request) {
	err := utils.ValidateParams(r, []string{"user_id", "date", "key"})
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := OuraGetDataParams{
		UserID: r.URL.Query().Get("user_id"),
		Date:   r.URL.Query().Get("date"),
		Key:    r.URL.Query().Get("key"),
	}

	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	data, err := h.store.GetData(context.Background(), params.UserID, params.Date, utils.NilIfEmpty(&params.Key))
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if len(data) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	utils.WriteJSON(w, http.StatusOK, data)
}
