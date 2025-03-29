package authapi

import "github.com/DeRuina/KUHA-REST-API/internal/store"

type AuthHandler struct {
	store store.Auth
}

func NewAuthHandler(store store.Auth) *AuthHandler {
	return &AuthHandler{store: store}
}
