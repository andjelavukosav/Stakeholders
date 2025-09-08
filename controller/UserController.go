package controller

import (
	"database-example/model"
	"database-example/service"
	"database-example/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}

// Konstruktor
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: svc,
	}
}

// POST /users -> registracija i vraća JWT token
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	var reg model.Registration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username: reg.Username,
		Password: reg.Password,
		Email:    reg.Email,
		Role:     reg.Role,
	}

	// Kreiraj korisnika i hashiraj lozinku
	if err := h.UserService.CreateUser(user); err != nil {
		http.Error(w, `{"error":"could not create user"}`, http.StatusInternalServerError)
		return
	}

	// Generiši JWT token
	token, err := util.GenerateToken(user.ID.String(), user.Username, user.Role)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.AuthenticationResponse{AccessToken: token})
}

// POST /users/login -> login i vraća JWT token
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	var login model.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Authenticate(login.Username, login.Password)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if user.IsBlocked {
		http.Error(w, `{"error":"account blocked"}`, http.StatusForbidden)
		return
	}

	token, err := util.GenerateToken(user.ID.String(), user.Username, user.Role)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.AuthenticationResponse{AccessToken: token})
}

// GET /users/{username} -> dohvata korisnika po username
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	username := mux.Vars(r)["username"]
	user, err := h.UserService.GetByUsername(username)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// GET /users/all
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
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

// 
func (h *UserHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
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

// PATCH /users/unblock/{id} 
func (h *UserHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
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
