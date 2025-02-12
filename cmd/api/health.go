package main

import (
	"net/http"

	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

// Health handler
func (app *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.config.env,
	}

	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		utils.InternalServerError(w, r, err)
	}
}
