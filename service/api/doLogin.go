package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type doLoginRequest struct {
	Username string `json:"username"`
}

type doLoginResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	PhotoPath *string `json:"photo_path"`
}

// login endpoint handler
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request body
	var req doLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, ErrorMessage("Username is required"), http.StatusBadRequest)
		return
	}

	// Check if user exists
	user, err := rt.db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// User exists, return the user with 200 (performs as a login)
	if user != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(doLoginResponse{
			ID:        user.ID,
			Username:  user.Username,
			PhotoPath: user.PhotoPath,
		})
		return
	}

	// User doesn’t exist, attempt the creation
	newUser, err := rt.db.CreateUser(req.Username, nil)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Return new user with 201 Created (performs as a register)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(doLoginResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		PhotoPath: newUser.PhotoPath,
	})
}
