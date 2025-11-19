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

type ResultNKHandler struct {
	store       fis.Resultnk
	competitors fis.Competitors
	cache       *cache.Storage
}

func NewResultNKHandler(store fis.Resultnk, competitors fis.Competitors, cache *cache.Storage) *ResultNKHandler {
	return &ResultNKHandler{store: store, competitors: competitors, cache: cache}
}

// GetLastRowResultNK godoc
//
//	@Summary		Get last Nordic Combined result record
//	@Description	Returns the last row in a_resultnk (by RecID DESC)
//	@Tags			FIS - Result Management – Nordic Combined
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	swagger.FISLastResultNKResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/lastrow/resultnk [get]
func (h *ResultNKHandler) GetLastRowResultNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), fisResultNKLastRowPrefix); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	row, err := h.store.GetLastRowResultNK(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"result": FISResultNKFullFromSqlc(row)}
	cache.SetCacheJSON(r.Context(), h.cache, fisResultNKLastRowPrefix, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// InsertResultNK godoc
//
//	@Summary		Add new Nordic Combined result
//	@Description	Inserts a new a_resultnk row
//	@Tags			FIS - Result Management – Nordic Combined
//	@Accept			json
//	@Produce		json
//	@Param			resultnk	body	swagger.FISInsertResultNKExample	true	"Result payload"
//	@Success		201			"Created"
//	@Failure		400			{object}	swagger.ValidationErrorResponse
//	@Failure		401			{object}	swagger.UnauthorizedResponse
//	@Failure		403			{object}	swagger.ForbiddenResponse
//	@Failure		409			{object}	swagger.ConflictResponse
//	@Failure		500			{object}	swagger.InternalServerErrorResponse
//	@Failure		503			{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/resultnk [post]
func (h *ResultNKHandler) InsertResultNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in InsertResultNKInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapInsertResultNKInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.InsertResultNK(r.Context(), clean); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateResultNK(r.Context(), h.cache, clean.Recid)
	w.WriteHeader(http.StatusCreated)
}

// UpdateResultNK godoc
//
//	@Summary		Update Nordic Combined result by RecID
//	@Description	Updates an existing a_resultnk row
//	@Tags			FIS - Result Management – Nordic Combined
//	@Accept			json
//	@Produce		json
//	@Param			resultnk	body	swagger.FISUpdateResultNKExample	true	"Result payload"
//	@Success		200			"Updated"
//	@Failure		400			{object}	swagger.ValidationErrorResponse
//	@Failure		401			{object}	swagger.UnauthorizedResponse
//	@Failure		403			{object}	swagger.ForbiddenResponse
//	@Failure		404			{object}	swagger.NotFoundResponse
//	@Failure		500			{object}	swagger.InternalServerErrorResponse
//	@Failure		503			{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/resultnk [put]
func (h *ResultNKHandler) UpdateResultNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in UpdateResultNKInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapUpdateResultNKInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.UpdateResultNKByRecID(r.Context(), clean); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateResultNK(r.Context(), h.cache, clean.Recid)
	w.WriteHeader(http.StatusOK)
}

// DeleteResultNK godoc
//
//	@Summary		Delete Nordic Combined result
//	@Description	Deletes a result by RecID
//	@Tags			FIS - Result Management – Nordic Combined
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
//	@Router			/fis/resultnk [delete]
func (h *ResultNKHandler) DeleteResultNK(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.DeleteResultNKByRecID(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateResultNK(r.Context(), h.cache, id)
	w.WriteHeader(http.StatusOK)
}

// GetRaceResultsNK godoc
//
//	@Summary	Get results for a Nordic Combined race
//	@Tags		FIS - Race Results
//	@Accept		json
//	@Produce	json
//	@Param		raceid	query		int32	true	"Race ID"
//	@Success	200		{object}	swagger.FISRaceResultsNKResponse
//	@Failure	400		{object}	swagger.ValidationErrorResponse
//	@Failure	401		{object}	swagger.UnauthorizedResponse
//	@Failure	403		{object}	swagger.ForbiddenResponse
//	@Failure	500		{object}	swagger.InternalServerErrorResponse
//	@Failure	503		{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/resultnk [get]
func (h *ResultNKHandler) GetRaceResultsNK(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("%s:race=%d", fisResultNKRacePrefix, raceID)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetRaceResultsNKByRaceID(r.Context(), raceID)
	if err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	if len(rows) == 0 {
		utils.NotFoundResponse(w, r, fmt.Errorf("no results found for raceid %d", raceID))
		return
	}

	out := make([]FISResultNKFullResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, FISResultNKFullFromSqlc(row))
	}

	body := map[string]any{"results": out}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetAthleteResultsNK godoc
//
//	@Summary	Get Nordic Combined results for an athlete
//	@Tags		FIS - Athlete
//	@Accept		json
//	@Produce	json
//	@Param		fiscode			query		int32		true	"FIS Code"
//	@Param		seasoncode		query		[]int32		false	"Season code (repeat or comma-separated)"
//	@Param		disciplinecode	query		[]string	false	"Discipline code (repeat or comma-separated)"
//	@Param		catcode			query		[]string	false	"Category code (repeat or comma-separated)"
//	@Success	200				{object}	swagger.FISAthleteResultsNKResponse
//	@Failure	400				{object}	swagger.ValidationErrorResponse
//	@Failure	401				{object}	swagger.UnauthorizedResponse
//	@Failure	403				{object}	swagger.ForbiddenResponse
//	@Failure	404				{object}	swagger.NotFoundResponse
//	@Failure	500				{object}	swagger.InternalServerErrorResponse
//	@Failure	503				{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/resultathletenk [get]
func (h *ResultNKHandler) GetAthleteResultsNK(w http.ResponseWriter, r *http.Request) {
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

	// get NK competitor ID from fiscode
	competitorID, err := h.competitors.GetCompetitorIDByFiscodeNK(r.Context(), fiscode)
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

	cacheKey := fmt.Sprintf("%s:fis=%d:sc=%v:dc=%v:cc=%v", fisResultNKAthletePrefix, fiscode, seasons, discs, cats)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetAthleteResultsNK(r.Context(), competitorID, seasons, discs, cats)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	out := make([]FISAthleteResultNKRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, FISAthleteResultNKFromSqlc(row))
	}

	body := map[string]any{"results": out}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}
