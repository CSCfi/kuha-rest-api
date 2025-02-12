package utvapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/DeRuina/KUHA-REST-API/internal/store/utv"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type PolarDataHandler struct {
	store utv.PolarData
}

// NewPolarDataHandler initializes PolarData handler
func NewPolarDataHandler(store utv.PolarData) *PolarDataHandler {
	return &PolarDataHandler{store: store}
}

// Get available dates from Polar data (with optional filtering)
func (h *PolarDataHandler) GetDates(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	afterDate := r.URL.Query().Get("after_date")
	beforeDate := r.URL.Query().Get("before_date")

	if userID == "" {
		utils.BadRequestResponse(w, r, utils.ErrMissingUserID)
		return
	}

	err := utils.ValidateParams(r, []string{"user_id", "after_date", "before_date"})
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if afterDate != "" && beforeDate != "" && afterDate > beforeDate {
		utils.UnprocessableEntityResponse(w, r, utils.ErrInvalidDateRange)
		return
	}

	dates, err := h.store.GetDates(context.Background(), userID, &afterDate, &beforeDate)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidUUID):
			utils.BadRequestResponse(w, r, err)
		case errors.Is(err, utils.ErrInvalidDate):
			utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	if len(dates) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dates)
}

// Get all JSON keys (types) from Polar data on a specific date
func (h *PolarDataHandler) GetTypes(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	summaryDate := r.URL.Query().Get("date")

	if userID == "" {
		utils.BadRequestResponse(w, r, utils.ErrMissingUserID)
		return
	}
	if summaryDate == "" {
		utils.BadRequestResponse(w, r, utils.ErrMissingDate)
		return
	}

	err := utils.ValidateParams(r, []string{"user_id", "date"})
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	types, err := h.store.GetTypes(context.Background(), userID, summaryDate)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidUUID):
			utils.BadRequestResponse(w, r, err)
		case errors.Is(err, utils.ErrInvalidDate):
			utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
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
func (h *PolarDataHandler) GetData(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	summaryDate := r.URL.Query().Get("date")
	key := r.URL.Query().Get("key")

	if userID == "" {
		utils.BadRequestResponse(w, r, utils.ErrMissingUserID)
		return
	}
	if summaryDate == "" {
		utils.BadRequestResponse(w, r, utils.ErrMissingDate)
		return
	}

	err := utils.ValidateParams(r, []string{"user_id", "date", "key"})
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	data, err := h.store.GetData(context.Background(), userID, summaryDate, utils.NilIfEmpty(&key))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidUUID):
			utils.BadRequestResponse(w, r, err)
		case errors.Is(err, utils.ErrInvalidDate):
			utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	if len(data) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	utils.WriteJSON(w, http.StatusOK, data)
}
