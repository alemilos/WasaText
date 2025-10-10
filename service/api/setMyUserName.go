package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/validation"
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
		http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
		return
	}

	// === Validate the username ===
	if err := validation.ValidateUsername(req.Username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// === Check if username already exists ===
	existingUser, err := rt.db.GetUserByUsername(req.Username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// unexpected DB error
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// === Requester is not the same as the username's owner ===
	if existingUser != nil && existingUser.ID != ctx.User.ID {
		http.Error(w, ErrorMessage("Username already taken"), http.StatusConflict)
		return
	}

	// === Update the username ===
	err = rt.db.SetUsername(ctx.User.ID, req.Username)
	if err != nil {
		http.Error(w, ErrorMessage("Could not update username"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "username updated successfully",
	})
}
