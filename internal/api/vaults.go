package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vmpyr/afterlight/internal/core"
	"github.com/vmpyr/afterlight/internal/store"
)

type VaultHandler struct {
	store *store.Store
}

func NewVaultHandler(s *store.Store) *VaultHandler {
	return &VaultHandler{store: s}
}

func (h *VaultHandler) Routes(authMiddleware func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(authMiddleware)

	r.Post("/", h.CreateVault)
	r.Get("/", h.ListVaults)
	r.Post("/{id}/artifacts", h.CreateArtifact)
	r.Get("/{id}/artifacts", h.ListArtifacts)

	return r
}

// Vault Handlers
func (h *VaultHandler) CreateVault(w http.ResponseWriter, r *http.Request) {
	var req core.CreateVaultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(UserKey).(*store.User).ID

	vault, err := h.store.CreateVault(r.Context(), store.CreateVaultParams{
		ID:        uuid.New().String(),
		UserID:    userID,
		VaultName: req.VaultName,
		Hint:      sql.NullString{String: req.Hint, Valid: req.Hint != ""},
		KdfSalt:   req.KdfSalt,
	})
	if err != nil {
		http.Error(w, "Failed to create vault", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vault)
}

func (h *VaultHandler) ListVaults(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(*store.User).ID

	vaults, err := h.store.GetVaultsByUser(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			vaults = []store.Vault{}
		} else {
			http.Error(w, "Failed to retrieve vaults", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vaults)
}

func (h *VaultHandler) GetVaultByID(r *http.Request, vaultID, userID string) (*store.Vault, error) {
	vault, err := h.store.GetVaultByID(r.Context(), store.GetVaultByIDParams{
		ID:     vaultID,
		UserID: userID,
	})
	return &vault, err
}

// Artifact Handlers
func (h *VaultHandler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	var req core.CreateArtifactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(UserKey).(*store.User).ID

	vaultID := chi.URLParam(r, "id")
	if vaultID == "" {
		http.Error(w, "Missing vault ID", http.StatusBadRequest)
		return
	}

	_, err := h.GetVaultByID(r, vaultID, userID)
	if err != nil {
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}

	artifact, err := h.store.CreateArtifact(r.Context(), store.CreateArtifactParams{
		ID:            uuid.New().String(),
		VaultID:       vaultID,
		MessageType:   req.MessageType,
		EncryptedBlob: req.EncryptedBlob,
		Iv:            req.IV,
	})
	if err != nil {
		http.Error(w, "Failed to create artifact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artifact)
}

func (h *VaultHandler) ListArtifacts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(*store.User).ID

	vaultID := chi.URLParam(r, "id")
	if vaultID == "" {
		http.Error(w, "Missing vault ID", http.StatusBadRequest)
		return
	}

	vault, err := h.GetVaultByID(r, vaultID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Vault not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve vault", http.StatusInternalServerError)
		return
	}

	artifacts, err := h.store.GetArtifactsByVault(r.Context(), store.GetArtifactsByVaultParams{
		VaultID: vaultID,
		UserID:  userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			artifacts = []store.Artifact{}
		} else {
			http.Error(w, "Failed to retrieve artifacts", http.StatusInternalServerError)
			return
		}
	}

	if artifacts == nil {
		artifacts = []store.Artifact{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(store.ListArtifactsResponse{
		VaultName: vault.VaultName,
		Hint:      vault.Hint.String,
		Artifacts: artifacts,
		CreatedAt: vault.CreatedAt,
	})
}
