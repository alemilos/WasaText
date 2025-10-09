package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type setMyUserNameRequest struct {
	Username string `json:"username"`
}

// login endpoint handler
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decode request body
	var req setMyUserNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// === Check if username already exists ===
	existingUser, err := rt.db.GetUserByUsername(req.Username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// unexpected DB error
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// === Requester is not the same as the username's owner ===
	if existingUser != nil && existingUser.ID != ctx.User.ID {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	// === Update the username ===
	err = rt.db.SetUsername(ctx.User.ID, req.Username)
	if err != nil {
		http.Error(w, "could not update username: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "username updated successfully",
	})
}
