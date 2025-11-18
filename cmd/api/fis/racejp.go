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

type RaceJPHandler struct {
	store fis.Racejp
	cache *cache.Storage
}

func NewRaceJPHandler(store fis.Racejp, cache *cache.Storage) *RaceJPHandler {
	return &RaceJPHandler{store: store, cache: cache}
}

// GetSeasonCodesJP godoc
//
//	@Summary	Get Ski Jumping season codes
//	@Tags		FIS - Season Discipline & Category Codes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	swagger.FISSeasonsJPResponse
//	@Failure	400	{object}	swagger.ValidationErrorResponse
//	@Failure	401	{object}	swagger.UnauthorizedResponse
//	@Failure	403	{object}	swagger.ForbiddenResponse
//	@Failure	500	{object}	swagger.InternalServerErrorResponse
//	@Failure	503	{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/seasoncodeJP [get]
func (h *RaceJPHandler) GetSeasonCodesJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	cacheKey := fmt.Sprintf("%s:seasons", fisRaceJPCodesPrefix)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetSkiJumpingSeasons(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"seasons": rows}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetDisciplineCodesJP godoc
//
//	@Summary	Get Ski Jumping discipline codes
//	@Tags		FIS - Season Discipline & Category Codes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	swagger.FISDisciplinesJPResponse
//	@Failure	400	{object}	swagger.ValidationErrorResponse
//	@Failure	401	{object}	swagger.UnauthorizedResponse
//	@Failure	403	{object}	swagger.ForbiddenResponse
//	@Failure	500	{object}	swagger.InternalServerErrorResponse
//	@Failure	503	{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/disciplinecodeJP [get]
func (h *RaceJPHandler) GetDisciplineCodesJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	cacheKey := fmt.Sprintf("%s:disciplines", fisRaceJPCodesPrefix)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetSkiJumpingDisciplines(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"disciplines": rows}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetCategoryCodesJP godoc
//
//	@Summary	Get Ski Jumping category codes
//	@Tags		FIS - Season Discipline & Category Codes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	swagger.FISCategoriesJPResponse
//	@Failure	400	{object}	swagger.ValidationErrorResponse
//	@Failure	401	{object}	swagger.UnauthorizedResponse
//	@Failure	403	{object}	swagger.ForbiddenResponse
//	@Failure	500	{object}	swagger.InternalServerErrorResponse
//	@Failure	503	{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/catcodeJP [get]
func (h *RaceJPHandler) GetCategoryCodesJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	cacheKey := fmt.Sprintf("%s:categories", fisRaceJPCodesPrefix)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetSkiJumpingCategories(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"categories": rows}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetRacesJP godoc
//
//	@Summary	Get list of Ski Jumping races
//	@Tags		FIS - Race Data
//	@Accept		json
//	@Produce	json
//	@Param		seasoncode		query		[]int32		false	"Season code (repeat or comma-separated)"
//	@Param		disciplinecode	query		[]string	false	"Discipline code (repeat or comma-separated)"
//	@Param		catcode			query		[]string	false	"Category code (repeat or comma-separated)"
//	@Success	200				{object}	swagger.FISRacesJPResponse
//	@Failure	400				{object}	swagger.ValidationErrorResponse
//	@Failure	401				{object}	swagger.UnauthorizedResponse
//	@Failure	403				{object}	swagger.ForbiddenResponse
//	@Failure	500				{object}	swagger.InternalServerErrorResponse
//	@Failure	503				{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/racejp [get]
func (h *RaceJPHandler) GetRacesJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	parseList := func(key string) []string {
		vals := r.URL.Query()[key]
		if len(vals) == 1 && strings.Contains(vals[0], ",") {
			return strings.Split(vals[0], ",")
		}
		return vals
	}
	seasonsS := parseList("seasoncode")
	discs := parseList("disciplinecode")
	cats := parseList("catcode")

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

	cacheKey := fmt.Sprintf("%s:sc=%v:dc=%v:cc=%v", fisRaceJPListPrefix, seasons, discs, cats)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetRacesJP(r.Context(), seasons, discs, cats)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	out := make([]FISRaceJPFullResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, FISRaceJPFullFromSqlc(row))
	}

	body := map[string]any{"races": out}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetLastRowRaceJP godoc
//
//	@Summary		Get last Ski Jumping race record
//	@Description	Returns the last row in a_racejp (by RaceID DESC)
//	@Tags			FIS - Race Management – Ski Jumping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	swagger.FISLastRaceJPResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/lastrow/racejp [get]
func (h *RaceJPHandler) GetLastRowRaceJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), fisRaceJPLastRowPrefix); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	row, err := h.store.GetLastRowRaceJP(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"race": FISRaceJPFullFromSqlc(row)}
	cache.SetCacheJSON(r.Context(), h.cache, fisRaceJPLastRowPrefix, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// InsertRaceJP godoc
//
//	@Summary		Add new Ski Jumping race
//	@Description	Inserts a new a_racejp row
//	@Tags			FIS - Race Management – Ski Jumping
//	@Accept			json
//	@Produce		json
//	@Param			racejp	body	swagger.FISInsertRaceJPExample	true	"Race payload"
//	@Success		201		"Created"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		409		{object}	swagger.ConflictResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/racejp [post]
func (h *RaceJPHandler) InsertRaceJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in InsertRaceJPInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapInsertRaceJPInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.InsertRaceJP(r.Context(), clean); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateRaceJP(r.Context(), h.cache, clean.Raceid)
	w.WriteHeader(http.StatusCreated)
}

// UpdateRaceJP godoc
//
//	@Summary		Update Ski Jumping race by ID
//	@Description	Updates an existing a_racejp row
//	@Tags			FIS - Race Management – Ski Jumping
//	@Accept			json
//	@Produce		json
//	@Param			racejp	body	swagger.FISUpdateRaceJPExample	true	"Race payload"
//	@Success		200		"Updated"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		404		{object}	swagger.NotFoundResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/racejp [put]
func (h *RaceJPHandler) UpdateRaceJP(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in UpdateRaceJPInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapUpdateRaceJPInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.UpdateRaceJPByID(r.Context(), clean); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateRaceJP(r.Context(), h.cache, clean.Raceid)
	w.WriteHeader(http.StatusOK)
}

// DeleteRaceJP godoc
//
//	@Summary		Delete Ski Jumping race
//	@Description	Deletes a race by RaceID
//	@Tags			FIS - Race Management – Ski Jumping
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int32	true	"Race ID"
//	@Success		200	"Deleted"
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/racejp [delete]
func (h *RaceJPHandler) DeleteRaceJP(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.DeleteRaceJPByID(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateRaceJP(r.Context(), h.cache, id)
	w.WriteHeader(http.StatusOK)
}
