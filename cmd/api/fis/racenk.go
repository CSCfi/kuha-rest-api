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

type RaceNKHandler struct {
	store fis.Racenk
	cache *cache.Storage
}

func NewRaceNKHandler(store fis.Racenk, cache *cache.Storage) *RaceNKHandler {
	return &RaceNKHandler{store: store, cache: cache}
}

// GetSeasonCodesNK godoc
//
//	@Summary	Get Nordic Combined season codes
//	@Tags		FIS - Season Discipline & Category Codes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	swagger.FISSeasonsNKResponse
//	@Failure	400	{object}	swagger.ValidationErrorResponse
//	@Failure	401	{object}	swagger.UnauthorizedResponse
//	@Failure	403	{object}	swagger.ForbiddenResponse
//	@Failure	500	{object}	swagger.InternalServerErrorResponse
//	@Failure	503	{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/seasoncodeNK [get]
func (h *RaceNKHandler) GetSeasonCodesNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	cacheKey := fmt.Sprintf("%s:seasons", fisRaceNKCodesPrefix)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetNordicCombinedSeasons(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"seasons": rows}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetDisciplineCodesNK godoc
//
//	@Summary	Get Nordic Combined discipline codes
//	@Tags		FIS - Season Discipline & Category Codes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	swagger.FISDisciplinesNKResponse
//	@Failure	400	{object}	swagger.ValidationErrorResponse
//	@Failure	401	{object}	swagger.UnauthorizedResponse
//	@Failure	403	{object}	swagger.ForbiddenResponse
//	@Failure	500	{object}	swagger.InternalServerErrorResponse
//	@Failure	503	{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/disciplinecodeNK [get]
func (h *RaceNKHandler) GetDisciplineCodesNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	cacheKey := fmt.Sprintf("%s:disciplines", fisRaceNKCodesPrefix)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetNordicCombinedDisciplines(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"disciplines": rows}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetCategoryCodesNK godoc
//
//	@Summary	Get Nordic Combined category codes
//	@Tags		FIS - Season Discipline & Category Codes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	swagger.FISCategoriesNKResponse
//	@Failure	400	{object}	swagger.ValidationErrorResponse
//	@Failure	401	{object}	swagger.UnauthorizedResponse
//	@Failure	403	{object}	swagger.ForbiddenResponse
//	@Failure	500	{object}	swagger.InternalServerErrorResponse
//	@Failure	503	{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/catcodeNK [get]
func (h *RaceNKHandler) GetCategoryCodesNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	cacheKey := fmt.Sprintf("%s:categories", fisRaceNKCodesPrefix)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetNordicCombinedCategories(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"categories": rows}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetRacesNK godoc
//
//	@Summary	Get list of Nordic Combined races
//	@Tags		FIS - Race Data
//	@Accept		json
//	@Produce	json
//	@Param		seasoncode		query		[]int32		false	"Season code (repeat or comma-separated)"
//	@Param		disciplinecode	query		[]string	false	"Discipline code (repeat or comma-separated)"
//	@Param		catcode			query		[]string	false	"Category code (repeat or comma-separated)"
//	@Success	200				{object}	swagger.FISRacesNKResponse
//	@Failure	400				{object}	swagger.ValidationErrorResponse
//	@Failure	401				{object}	swagger.UnauthorizedResponse
//	@Failure	403				{object}	swagger.ForbiddenResponse
//	@Failure	500				{object}	swagger.InternalServerErrorResponse
//	@Failure	503				{object}	swagger.ServiceUnavailableResponse
//	@Security	BearerAuth
//	@Router		/fis/racenk [get]
func (h *RaceNKHandler) GetRacesNK(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("%s:sc=%v:dc=%v:cc=%v", fisRaceNKListPrefix, seasons, discs, cats)
	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), cacheKey); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	rows, err := h.store.GetRacesNK(r.Context(), seasons, discs, cats)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	out := make([]FISRaceNKFullResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, FISRaceNKFullFromSqlc(row))
	}

	body := map[string]any{"races": out}
	cache.SetCacheJSON(r.Context(), h.cache, cacheKey, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// GetLastRowRaceNK godoc
//
//	@Summary		Get last Nordic Combined race record
//	@Description	Returns the last row in a_racenk (by RaceID DESC)
//	@Tags			FIS - Race Management – Nordic Combined
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	swagger.FISLastRaceNKResponse
//	@Failure		400	{object}	swagger.ValidationErrorResponse
//	@Failure		401	{object}	swagger.UnauthorizedResponse
//	@Failure		403	{object}	swagger.ForbiddenResponse
//	@Failure		404	{object}	swagger.NotFoundResponse
//	@Failure		500	{object}	swagger.InternalServerErrorResponse
//	@Failure		503	{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/lastrow/racenk [get]
func (h *RaceNKHandler) GetLastRowRaceNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if h.cache != nil {
		if raw, err := h.cache.Get(r.Context(), fisRaceNKLastRowPrefix); err == nil && raw != "" {
			utils.WriteJSON(w, http.StatusOK, json.RawMessage(raw))
			return
		}
	}

	row, err := h.store.GetLastRowRaceNK(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	body := map[string]any{"race": FISRaceNKFullFromSqlc(row)}
	cache.SetCacheJSON(r.Context(), h.cache, fisRaceNKLastRowPrefix, body, FISCacheTTL)
	utils.WriteJSON(w, http.StatusOK, body)
}

// InsertRaceNK godoc
//
//	@Summary		Add new Nordic Combined race
//	@Description	Inserts a new a_racenk row
//	@Tags			FIS - Race Management – Nordic Combined
//	@Accept			json
//	@Produce		json
//	@Param			racenk	body	swagger.FISInsertRaceNKExample	true	"Race payload"
//	@Success		201		"Created"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		409		{object}	swagger.ConflictResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/racenk [post]
func (h *RaceNKHandler) InsertRaceNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in InsertRaceNKInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapInsertRaceNKInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.InsertRaceNK(r.Context(), clean); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateRaceNK(r.Context(), h.cache, clean.Raceid)
	w.WriteHeader(http.StatusCreated)
}

// UpdateRaceNK godoc
//
//	@Summary		Update Nordic Combined race by ID
//	@Description	Updates an existing a_racenk row
//	@Tags			FIS - Race Management – Nordic Combined
//	@Accept			json
//	@Produce		json
//	@Param			racenk	body	swagger.FISUpdateRaceNKExample	true	"Race payload"
//	@Success		200		"Updated"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		404		{object}	swagger.NotFoundResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/racenk [put]
func (h *RaceNKHandler) UpdateRaceNK(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var in UpdateRaceNKInput
	if err := utils.ReadJSON(w, r, &in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(in); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	clean, err := mapUpdateRaceNKInput(in)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := h.store.UpdateRaceNKByID(r.Context(), clean); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateRaceNK(r.Context(), h.cache, clean.Raceid)
	w.WriteHeader(http.StatusOK)
}

// DeleteRaceNK godoc
//
//	@Summary		Delete Nordic Combined race
//	@Description	Deletes a race by RaceID
//	@Tags			FIS - Race Management – Nordic Combined
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
//	@Router			/fis/racenk [delete]
func (h *RaceNKHandler) DeleteRaceNK(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.DeleteRaceNKByID(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, err)
			return
		}
		utils.HandleDatabaseError(w, r, err)
		return
	}

	invalidateRaceNK(r.Context(), h.cache, id)
	w.WriteHeader(http.StatusOK)
}
