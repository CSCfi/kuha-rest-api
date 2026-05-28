package archapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DeRuina/KUHA-REST-API/internal/auth/authz"
	archsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/archinisis"
	"github.com/DeRuina/KUHA-REST-API/internal/store/archinisis"
	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type DataHandler struct {
	store archinisis.Data
	cache *cache.Storage
}

func NewDataHandler(store archinisis.Data, cache *cache.Storage) *DataHandler {
	return &DataHandler{store: store, cache: cache}
}

type RaceReportSessionsQuery struct {
	ID string `validate:"required,numeric"`
}

type RaceReportHTMLQuery struct {
	ID        string `validate:"required,numeric"`
	SessionID string `validate:"required,numeric"`
}

// GetRaceReportSessions godoc
//
//	@Summary		List race-report session IDs for a Sportti ID
//	@Description	Returns all session_id values that have race reports for the given sportti_id
//	@Tags			ARCHINISIS - Data
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"Sportti ID"
//	@Success		200	{object}	swagger.RaceReportSessionsResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/race-report/sessions [get]
func (h *DataHandler) GetRaceReportSessions(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	q := RaceReportSessionsQuery{
		ID: r.URL.Query().Get("id"),
	}
	if err := utils.GetValidator().Struct(q); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(q.ID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	cacheKey := fmt.Sprintf("arch:race-report:sessions:%s", sid)

	if h.cache != nil {
		if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil && cached != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(cached))
			return
		}
	}

	sessionIDs, err := h.store.GetRaceReportSessions(r.Context(), sid)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if sessionIDs == nil {
		sessionIDs = []int32{}
	}

	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, map[string]any{"race_report": sessionIDs}, ARCHCacheTTL)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"race_report": sessionIDs,
	})
}

// GetRaceReportHTML godoc
//
//	@Summary		Get a specific race report (HTML)
//	@Description	Returns the raw HTML race report for a (sportti_id, session_id). Content-Type is text/html.
//	@Tags			ARCHINISIS - Data
//	@Accept			json
//	@Produce		html
//	@Param			id			query		string	true	"Sportti ID"
//	@Param			session_id	query		string	true	"Session ID"
//	@Success		200			{string}	string	"<!DOCTYPE html><html><head><title>Race Report</title></head><body><h1>HTML RACE REPORT</h1><p>full report returned in html DOCTYPE</p></body></html>"
//	@Failure		400			{object}	swagger.ValidationErrorResponse
//	@Failure		401			{object}	swagger.UnauthorizedResponse
//	@Failure		403			{object}	swagger.ForbiddenResponse
//	@Failure		404			{object}	swagger.NotFoundResponse
//	@Failure		500			{object}	swagger.InternalServerErrorResponse
//	@Failure		503			{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/race-report [get]
func (h *DataHandler) GetRaceReportHTML(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"id", "session_id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	q := RaceReportHTMLQuery{
		ID:        r.URL.Query().Get("id"),
		SessionID: r.URL.Query().Get("session_id"),
	}
	if err := utils.GetValidator().Struct(q); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(q.ID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sessionID, err := utils.ParsePositiveInt32(q.SessionID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	cacheKey := fmt.Sprintf("arch:race-report:html:%s:%d", sid, sessionID)

	if h.cache != nil {
		if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil && cached != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(cached))
			return
		}
	}

	html, err := h.store.GetRaceReport(r.Context(), sid, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	if h.cache != nil {
		_ = h.cache.Set(r.Context(), cacheKey, html, ARCHCacheTTL)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

type RaceReportUpsertInput struct {
	SporttiID  string `json:"sportti_id" validate:"required,numeric"`
	SessionID  int32  `json:"session_id" validate:"required,gt=0"`
	RaceReport string `json:"race_report" validate:"required"`
}

// PostRaceReport godoc
//
//	@Summary		Upsert a race report (HTML)
//	@Description	Inserts or updates the shared race report for session_id and links the given sportti_id to that session.
//	@Tags			ARCHINISIS - Data
//	@Accept			json
//	@Produce		json
//	@Param			data	body	swagger.ArchRaceReportUpsertRequest	true	"race report"
//	@Success		201		"Data processed successfully"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/race-report [post]
func (h *DataHandler) PostRaceReport(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in RaceReportUpsertInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(in.SporttiID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.UpsertRaceReport(r.Context(), sid, in.SessionID, in.RaceReport); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	sporttiIDs, err := h.store.GetSporttiIDsBySessionID(r.Context(), in.SessionID)
	if err == nil {
		for _, linkedID := range sporttiIDs {
			invalidateArchRaceReport(r.Context(), h.cache, linkedID, &in.SessionID)
		}
	} else {
		invalidateArchRaceReport(r.Context(), h.cache, sid, &in.SessionID)
	}

	w.WriteHeader(http.StatusCreated)
}

type archIDParam struct {
	ID string `validate:"required,numeric"`
}

// GetAthlete godoc
//
//	@Summary		Get athlete profile by ID
//	@Description	Returns athlete profile fields (no measurements) for the given ID.
//	@Tags			ARCHINISIS - Athlete
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"National ID (Sportti ID)"
//	@Success		200	{object}	swagger.ArchAthleteResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/athlete [get]
func (h *DataHandler) GetAthlete(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := archIDParam{ID: r.URL.Query().Get("id")}
	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(params.ID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	cacheKey := fmt.Sprintf("arch:athlete:%s", sid)
	if h.cache != nil {
		if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil && cached != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(cached))
			return
		}
	}

	res, err := h.store.GetAthleteByID(r.Context(), sid)
	if err == sql.ErrNoRows {
		utils.NotFoundResponse(w, r, err)
		return
	}
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, res, ARCHCacheTTL)
	utils.WriteJSON(w, http.StatusOK, res)
}

// PostAthlete godoc
//
//	@Summary		Upsert athlete profile
//	@Description	Inserts or updates athlete profile fields. New data for the same national_id overwrites existing.
//	@Tags			ARCHINISIS - Athlete
//	@Accept			json
//	@Produce		json
//	@Param			data	body	swagger.ArchAthleteUpsertRequest	true	"athlete profile"
//	@Success		201		"Data processed successfully"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/athlete [post]
func (h *DataHandler) PostAthlete(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in ArchAthleteInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(in.NationalID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	ath, err := mapAthleteToParams(in, sid)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.UpsertAthleteOnly(r.Context(), ath); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateArchAthlete(r.Context(), h.cache, sid)
	w.WriteHeader(http.StatusCreated)
}

// GetMeasurements godoc
//
//	@Summary		Get all measurements for an athlete
//	@Description	Returns all measurements for the athlete identified by the given ID.
//	@Tags			ARCHINISIS - Measurements
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"National ID (Sportti ID)"
//	@Success		200	{object}	swagger.ArchMeasurementsResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/measurements [get]
func (h *DataHandler) GetMeasurements(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := archIDParam{ID: r.URL.Query().Get("id")}
	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(params.ID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	cacheKey := fmt.Sprintf("arch:measurements:%s", sid)
	if h.cache != nil {
		if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil && cached != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(cached))
			return
		}
	}

	ms, err := h.store.GetMeasurementsByID(r.Context(), sid)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	if ms == nil {
		ms = []archinisis.ArchMeasurementResponse{}
	}

	resp := map[string]any{"measurements": ms}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, resp, ARCHCacheTTL)
	utils.WriteJSON(w, http.StatusOK, resp)
}

// PostMeasurements godoc
//
//	@Summary		Upsert measurements for an athlete
//	@Description	Inserts or updates one or more measurements linked to the given national_id.
//	@Tags			ARCHINISIS - Measurements
//	@Accept			json
//	@Produce		json
//	@Param			data	body	swagger.ArchMeasurementsUpsertRequest	true	"measurements"
//	@Success		201		"Data processed successfully"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/measurements [post]
func (h *DataHandler) PostMeasurements(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in ArchMeasurementsInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	sid, err := utils.ParseSporttiID(in.NationalID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	measParams := make([]archsqlc.UpsertMeasurementParams, 0, len(in.Measurements))
	for _, m := range in.Measurements {
		mp, err := mapMeasurementToParams(m, sid)
		if err != nil {
			utils.BadRequestResponse(w, r, err)
			return
		}
		measParams = append(measParams, mp)
	}

	if err := h.store.UpsertMeasurements(r.Context(), sid, measParams); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateArchMeasurements(r.Context(), h.cache, sid)
	w.WriteHeader(http.StatusCreated)
}

type archMeasurementIDParam struct {
	ID string `validate:"required,numeric"`
}

// GetMeasurement godoc
//
//	@Summary		Get a single measurement by measurement ID
//	@Description	Returns the measurement identified by the given measurement_id.
//	@Tags			ARCHINISIS - Measurements
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"Measurement ID"
//	@Success		200	{object}	swagger.ArchMeasurementResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/measurement [get]
func (h *DataHandler) GetMeasurement(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := archMeasurementIDParam{ID: r.URL.Query().Get("id")}
	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	mid, err := utils.ParsePositiveInt32(params.ID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	res, err := h.store.GetMeasurementByMeasurementID(r.Context(), mid)
	if err == sql.ErrNoRows {
		utils.NotFoundResponse(w, r, err)
		return
	}
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, res)
}

// DeleteMeasurement godoc
//
//	@Summary		Delete a single measurement by measurement ID
//	@Description	Removes the measurement identified by the given measurement_id.
//	@Tags			ARCHINISIS - Measurements
//	@Accept			json
//	@Produce		json
//	@Param			id	query	string	true	"Measurement ID"
//	@Success		200
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/archinisis/measurement [delete]
func (h *DataHandler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := archMeasurementIDParam{ID: r.URL.Query().Get("id")}
	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	mid, err := utils.ParsePositiveInt32(params.ID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	nationalID, err := h.store.DeleteMeasurementByMeasurementID(r.Context(), mid)
	if err == sql.ErrNoRows {
		utils.NotFoundResponse(w, r, err)
		return
	}
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if nationalID != "" {
		invalidateArchMeasurements(r.Context(), h.cache, nationalID)
	}

	w.WriteHeader(http.StatusOK)
}
