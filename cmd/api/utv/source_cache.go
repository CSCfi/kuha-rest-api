package utvapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DeRuina/KUHA-REST-API/internal/auth/authz"
	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
	"github.com/DeRuina/KUHA-REST-API/internal/store/utv"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type SourceCacheHandler struct {
	store utv.SourceCache
	cache *cache.Storage
}

func NewSourceCacheHandler(store utv.SourceCache, cache *cache.Storage) *SourceCacheHandler {
	return &SourceCacheHandler{store: store, cache: cache}
}

type SourceParam struct {
	Source string `form:"source"`
}

type UpsertSourceCacheBody struct {
	Source string   `json:"source" validate:"required"`
	Data   []string `json:"data"   validate:"required,min=1,dive,required"`
}

// GetAllDataTypes godoc
//
//	@Summary		Get data types (source_cache)
//	@Description	Returns data types. If source is provided, returns only that source's keys.
//	@Tags			UTV - Source_Cache
//	@Accept			json
//	@Produce		json
//	@Param			source	query		string	false	"Source name (e.g., oura, garmin, polar, suunto)"
//	@Success		200		{object}	swagger.SourceCacheSingleResponse
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		404		{object}	swagger.NotFoundResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Failure		503		{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/utv/source_cache/data-types [get]
func (h *SourceCacheHandler) GetAllDataTypes(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{"source"}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	params := SourceParam{
		Source: r.URL.Query().Get("source"),
	}

	if err := utils.GetValidator().Struct(params); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	source := params.Source
	if source != "" {
		data, err := h.store.GetBySource(r.Context(), source)
		if err == sql.ErrNoRows {
			utils.NotFoundResponse(w, r, fmt.Errorf("source not found"))
			return
		}
		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}

		_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
			"source": source,
			"data":   json.RawMessage(data),
		})
		return
	}

	rows, err := h.store.GetAll(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"sources": rows,
	})
}

// UpsertDataTypes godoc
//
//	@Summary		Replace keys for a source (upsert)
//	@Description	Replaces the entire key list for the given source. One source per request.
//	@Tags			UTV - Source_Cache
//	@Accept			json
//	@Produce		json
//	@Param			body	body	swagger.SourceCacheUpsertInput	true	"source + data array"
//	@Success		201		"Created"
//	@Failure		400		{object}	swagger.ValidationErrorResponse
//	@Failure		401		{object}	swagger.UnauthorizedResponse
//	@Failure		403		{object}	swagger.ForbiddenResponse
//	@Failure		500		{object}	swagger.InternalServerErrorResponse
//	@Security		BearerAuth
//	@Router			/utv/source_cache/data-types [post]
func (h *SourceCacheHandler) UpsertDataTypes(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	var input UpsertSourceCacheBody
	if err := utils.ReadJSON(w, r, &input); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	if err := utils.GetValidator().Struct(input); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	raw, err := json.Marshal(input.Data)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := h.store.Upsert(r.Context(), input.Source, raw); err != nil {
		utils.HandleDatabaseError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
