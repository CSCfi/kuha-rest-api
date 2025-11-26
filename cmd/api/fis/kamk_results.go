package fisapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DeRuina/KUHA-REST-API/internal/auth/authz"
	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
	"github.com/DeRuina/KUHA-REST-API/internal/store/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

// Handles competitor/result-related KAMK endpoints.
type ResultKAMKHandler struct {
	resultCC fis.Resultcc
	resultJP fis.Resultjp
	resultNK fis.Resultnk
	cache    *cache.Storage
}

func NewResultKAMKHandler(
	resultCC fis.Resultcc,
	resultJP fis.Resultjp,
	resultNK fis.Resultnk,
	cache *cache.Storage,
) *ResultKAMKHandler {
	return &ResultKAMKHandler{
		resultCC: resultCC,
		resultJP: resultJP,
		resultNK: resultNK,
		cache:    cache,
	}
}

// GetCompetitorSeasonsCatcodes godoc
//
//	@Summary		Get competitor seasons and categories
//	@Description	Gets distinct (Seasoncode, Catcode) combinations for the races a competitor has actually competed in, for a given sector (CC, JP, NK). Uses both race and result tables.
//	@Tags			FIS - KAMK
//	@Accept			json
//	@Produce		json
//	@Param			competitorid	query		int32	true	"Competitor ID"
//	@Param			sector			query		string	true	"Sector code (CC,JP,NK)"
//	@Success		200				{object}	swagger.FISCompetitorSeasonsCatcodesResponse
//	@Failure		400				{object}	swagger.ValidationErrorResponse
//	@Failure		401				{object}	swagger.UnauthorizedResponse
//	@Failure		403				{object}	swagger.ForbiddenResponse
//	@Failure		500				{object}	swagger.InternalServerErrorResponse
//	@Failure		503				{object}	swagger.ServiceUnavailableResponse
//	@Security		BearerAuth
//	@Router			/fis/competitor/seasons-catcodes [get]
func (h *ResultKAMKHandler) GetCompetitorSeasonsCatcodes(w http.ResponseWriter, r *http.Request) {
	if !authz.Authorize(r) {
		utils.ForbiddenResponse(w, r, fmt.Errorf("access denied"))
		return
	}

	if err := utils.ValidateParams(r, []string{
		"competitorid", "sector",
	}); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	rawComp := strings.TrimSpace(r.URL.Query().Get("competitorid"))
	if rawComp == "" {
		utils.BadRequestResponse(w, r, fmt.Errorf("missing required query param: competitorid"))
		return
	}
	competitorID, err := utils.ParsePositiveInt32(rawComp)
	if err != nil {
		utils.BadRequestResponse(w, r, fmt.Errorf("invalid competitorid: %s", rawComp))
		return
	}

	sector := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("sector")))
	switch sector {
	case "CC", "JP", "NK":
	default:
		if sector == "" {
			utils.BadRequestResponse(w, r, fmt.Errorf("missing required query param: sector"))
		} else {
			utils.BadRequestResponse(w, r, fmt.Errorf("invalid sector: %s", sector))
		}
		return
	}

	type item struct {
		Seasoncode int32   `json:"seasoncode"`
		Catcode    *string `json:"catcode,omitempty"`
	}

	var items []item

	switch sector {
	case "CC":
		if h.resultCC == nil {
			utils.InternalServerError(w, r, fmt.Errorf("resultCC store not configured"))
			return
		}
		rows, err := h.resultCC.GetSeasonsCatcodesCCByCompetitor(r.Context(), competitorID)
		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}
		for _, row := range rows {
			if !row.Seasoncode.Valid {
				continue
			}
			var cat *string
			if row.Catcode.Valid && strings.TrimSpace(row.Catcode.String) != "" {
				c := strings.TrimSpace(row.Catcode.String)
				cat = &c
			}
			items = append(items, item{
				Seasoncode: row.Seasoncode.Int32,
				Catcode:    cat,
			})
		}

	case "JP":
		if h.resultJP == nil {
			utils.InternalServerError(w, r, fmt.Errorf("resultJP store not configured"))
			return
		}
		rows, err := h.resultJP.GetSeasonsCatcodesJPByCompetitor(r.Context(), competitorID)
		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}
		for _, row := range rows {
			if !row.Seasoncode.Valid {
				continue
			}
			var cat *string
			if row.Catcode.Valid && strings.TrimSpace(row.Catcode.String) != "" {
				c := strings.TrimSpace(row.Catcode.String)
				cat = &c
			}
			items = append(items, item{
				Seasoncode: row.Seasoncode.Int32,
				Catcode:    cat,
			})
		}

	case "NK":
		if h.resultNK == nil {
			utils.InternalServerError(w, r, fmt.Errorf("resultNK store not configured"))
			return
		}
		rows, err := h.resultNK.GetSeasonsCatcodesNKByCompetitor(r.Context(), competitorID)
		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}
		for _, row := range rows {
			if !row.Seasoncode.Valid {
				continue
			}
			var cat *string
			if row.Catcode.Valid && strings.TrimSpace(row.Catcode.String) != "" {
				c := strings.TrimSpace(row.Catcode.String)
				cat = &c
			}
			items = append(items, item{
				Seasoncode: row.Seasoncode.Int32,
				Catcode:    cat,
			})
		}
	}

	body := map[string]any{
		"competitorid": competitorID,
		"sector":       sector,
		"items":        items,
	}

	utils.WriteJSON(w, http.StatusOK, body)
}
