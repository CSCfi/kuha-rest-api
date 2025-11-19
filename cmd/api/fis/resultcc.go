package fisapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DeRuina/KUHA-REST-API/internal/auth/authz"
	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
	"github.com/DeRuina/KUHA-REST-API/internal/store/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type ResultCCHandler struct {
	store       fis.Resultcc
	competitors fis.Competitors
	cache       *cache.Storage
}

func NewResultCCHandler(store fis.Resultcc, competitors fis.Competitors, cache *cache.Storage) *ResultCCHandler {
	return &ResultCCHandler{store: store, competitors: competitors, cache: cache}
}

// helper
func parseListParam(r *http.Request, key string) []string {
	vals := r.URL.Query()[key]
	if len(vals) == 1 && strings.Contains(vals[0], ",") {
		return strings.Split(vals[0], ",")
	}
	return vals
}

// GetLastRowResultCC godoc
//
//	@Summary		Get last Cross-Country result record
//	@Description	Returns the last row in a_resultcc (by RecID DESC)
//	@Tags			FIS - Result Management – Cross-Country
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	swagger.FISLastResultCCResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/lastrow/resultcc [get]
func (h *ResultCCHandler) GetLastRowResultCC(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), fisResultCCLastRowPrefix); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	row, err := h.store.GetLastRowResultCC(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"result": FISResultCCFullFromSqlc(row)}
	cache.SetCacheJSON(r.Context(), h.cache, fisResultCCLastRowPrefix, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// InsertResultCC godoc
//
//	@Summary		Add new Cross-Country result
//	@Description	Inserts a new a_resultcc row
//	@Tags			FIS - Result Management – Cross-Country
//	@Accept			json
//	@Produce		json
//	@Param			resultcc	body	swagger.FISInsertResultCCExample	true	"Result payload"
//	@Success		201			"Created"
//	@Failure		400			{object}	swagger.ValidationErrorResponse
//	@Failure		401			{object}	swagger.UnauthorizedResponse
//	@Failure		403			{object}	swagger.ForbiddenResponse
//	@Failure		409			{object}	swagger.ConflictResponse
//	@Failure		500			{object}	swagger.InternalServerErrorResponse
//	@Failure		503			{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/resultcc [post]
func (h *ResultCCHandler) InsertResultCC(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in InsertResultCCInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapInsertResultCCInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.InsertResultCC(r.Context(), clean); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateResultCC(r.Context(), h.cache, clean.Recid)
	w.WriteHeader(http.StatusCreated)
}

// UpdateResultCC godoc
//
//	@Summary		Update Cross-Country result by RecID
//	@Description	Updates an existing a_resultcc row
//	@Tags			FIS - Result Management – Cross-Country
//	@Accept			json
//	@Produce		json
//	@Param			resultcc	body	swagger.FISUpdateResultCCExample	true	"Result payload"
//	@Success		200			"Updated"
//	@Failure		400			{object}	swagger.ValidationErrorResponse
//	@Failure		401			{object}	swagger.UnauthorizedResponse
//	@Failure		403			{object}	swagger.ForbiddenResponse
//	@Failure		404			{object}	swagger.NotFoundResponse
//	@Failure		500			{object}	swagger.InternalServerErrorResponse
//	@Failure		503			{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/resultcc [put]
func (h *ResultCCHandler) UpdateResultCC(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in UpdateResultCCInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapUpdateResultCCInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.UpdateResultCCByRecID(r.Context(), clean); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateResultCC(r.Context(), h.cache, clean.Recid)
	w.WriteHeader(http.StatusOK)
}

// DeleteResultCC godoc
//
//	@Summary		Delete Cross-Country result
//	@Description	Deletes a result by RecID
//	@Tags			FIS - Result Management – Cross-Country
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int32	true	"Result RecID"
//	@Success		200	"Deleted"
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/resultcc [delete]
func (h *ResultCCHandler) DeleteResultCC(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}
	if err := utils.ValidateParams(r, []string{"id"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := utils.ParsePositiveInt32(idStr)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.DeleteResultCCByRecID(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateResultCC(r.Context(), h.cache, id)
	w.WriteHeader(http.StatusOK)
}

// GetRaceResultsCC godoc
//
//	@Summary	Get results for a Cross-Country race
//	@Tags		FIS - Race Results
//	@Accept		json
//	@Produce	json
//	@Param		raceid	query		int32	true	"Race ID"
//	@Success	200		{object}	swagger.FISRaceResultsCCResponse
//	@Failure	400		{object}	swagger.ValidationErrorResponse
//	@Failure	401		{object}	swagger.UnauthorizedResponse
//	@Failure	403		{object}	swagger.ForbiddenResponse
//	@Failure	500		{object}	swagger.InternalServerErrorResponse
//	@Failure	503		{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/resultcc [get]
func (h *ResultCCHandler) GetRaceResultsCC(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}
	if err := utils.ValidateParams(r, []string{"raceid"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	raceStr := r.URL.Query().Get("raceid")
	raceID, err := utils.ParsePositiveInt32(raceStr)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	cacheKey := fmt.Sprintf("%s:race=%d", fisResultCCRacePrefix, raceID)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetRaceResultsCCByRaceID(r.Context(), raceID)
	if err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	if len(rows) == 0 {
		utils.NotFoundResponse(w, r, fmt.Errorf("no results found for raceid %d", raceID))
		return
	}

	out := make([]FISResultCCFullResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, FISResultCCFullFromSqlc(row))
	}

	body := map[string]any{"results": out}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetAthleteResultsCC godoc
//
//	@Summary	Get Cross-Country results for an athlete
//	@Tags		FIS - Athlete
//	@Accept		json
//	@Produce	json
//	@Param		fiscode			query		int32		true	"FIS Code"
//	@Param		seasoncode		query		[]int32		false	"Season code (repeat or comma-separated)"
//	@Param		disciplinecode	query		[]string	false	"Discipline code (repeat or comma-separated)"
//	@Param		catcode			query		[]string	false	"Category code (repeat or comma-separated)"
//	@Success	200				{object}	swagger.FISAthleteResultsCCResponse
//	@Failure	400				{object}	swagger.ValidationErrorResponse
//	@Failure	401				{object}	swagger.UnauthorizedResponse
//	@Failure	403				{object}	swagger.ForbiddenResponse
//	@Failure	404				{object}	swagger.NotFoundResponse
//	@Failure	500				{object}	swagger.InternalServerErrorResponse
//	@Failure	503				{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/resultathletecc [get]
func (h *ResultCCHandler) GetAthleteResultsCC(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}
	if err := utils.ValidateParams(r, []string{"fiscode", "seasoncode", "disciplinecode", "catcode"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	fiscodeStr := r.URL.Query().Get("fiscode")
	fiscode, err := utils.ParsePositiveInt32(fiscodeStr)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	// get CC competitor ID from fiscode
	competitorID, err := h.competitors.GetCompetitorIDByFiscodeCC(r.Context(), fiscode)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, fmt.Errorf("competitor with FIS code %d not found", fiscode))
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	// filters: seasoncode, disciplinecode, catcode
	seasonsS := parseListParam(r, "seasoncode")
	discs := parseListParam(r, "disciplinecode")
	cats := parseListParam(r, "catcode")

	var seasons []int32
	for _, s := range seasonsS {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			utils.BadRequestResponse(w, r, fmt.Errorf("invalid seasoncode: %s", s))
			return
		}
		seasons = append(seasons, int32(n))
	}

	cacheKey := fmt.Sprintf("%s:fis=%d:sc=%v:dc=%v:cc=%v", fisResultCCAthletePrefix, fiscode, seasons, discs, cats)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetAthleteResultsCC(r.Context(), competitorID, seasons, discs, cats)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	out := make([]FISAthleteResultCCRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, FISAthleteResultCCFromSqlc(row))
	}

	body := map[string]any{"results": out}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}
