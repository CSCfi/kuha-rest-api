package main

import (
	"net/http"
)

func (app *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.config.env,
	}

	if err := WriteJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
