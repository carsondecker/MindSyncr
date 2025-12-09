package auth

import (
	"encoding/json"
	"fmt"
	"log"
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
	log.Println("Hi")

	var registerRequest sqlc.InsertUserParams
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		fmt.Printf("failed to decode data: %w\n", err)
		w.WriteHeader(400)
		return
	}

	newUser, err := h.cfg.Queries.InsertUser(r.Context(), registerRequest)

	if err != nil {
		fmt.Printf("failed to insert user: %w\n", err)
		w.WriteHeader(500)
		return
	}

	res, err := json.Marshal(newUser)

	if err != nil {
		fmt.Printf("failed to marshal data: %w\n", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Write(res)
}
