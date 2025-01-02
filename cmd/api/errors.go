package main

import (
	"log"
	"net/http"
)

func (app *api) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// will implement proper logging
	log.Printf("Internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *api) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	// will implement proper logging
	log.Printf("Bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *api) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	// will implement proper logging
	log.Printf("Not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusNotFound, "Not found")
}
