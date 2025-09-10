package controller

import (
	"database-example/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
    UserService *service.UserService
}

func NewAdminHandler(svc *service.UserService) *AdminHandler {
    return &AdminHandler{UserService: svc}
}


// GET /users/all
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Omogući CORS za Angular frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/json")

	// Opcija za preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Dohvati sve korisnike
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, `{"error":"could not fetch users"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// PATCH admin/users/block/{id} 
func (h *AdminHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(w,r)
	if r.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	err := h.UserService.BlockUser(userID)
    if err != nil {
        http.Error(w, `{"error":"failed to block user"}`, http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message":"user blocked"}`))
}

// PATCH admin/users/unblock/{id} 
func (h *AdminHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	err := h.UserService.UnblockUser(userID)
	if err != nil {
		http.Error(w, `{"error":"failed to unblock user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"user unblocked"}`))
}
