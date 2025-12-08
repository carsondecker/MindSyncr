package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var registerRequest sqlc.InsertUserParams
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		w.WriteHeader(400)
		return
	}

	newUser, err := h.cfg.Queries.InsertUser(ctx, registerRequest)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	res, err := json.Marshal(newUser)

	w.WriteHeader(201)
	w.Write(res)
}
