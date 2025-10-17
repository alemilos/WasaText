package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/validation"
	"github.com/julienschmidt/httprouter"
)

type setGroupNameRequest struct {
	Name string `json:"name"`
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse conversation ID
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid conversation id"), http.StatusBadRequest)
		return
	}

	// Ensure requester is a member
	if err := rt.db.IsMember(convID, ctx.User.ID); err != nil {
		http.Error(w, ErrorMessage("You are not a member of this conversation"), http.StatusForbidden)
		return
	}

	// Ensure conversation is a group
	conv, err := rt.db.GetConversationByID(convID)
	if err != nil || conv == nil {
		http.Error(w, ErrorMessage("Conversation not found"), http.StatusNotFound)
		return
	}
	if conv.Type != constants.CONV_GROUP {
		http.Error(w, ErrorMessage("Conversation is not a group"), http.StatusBadRequest)
		return
	}

	// Parse and validate request body
	var req setGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateGroupName(req.Name); err != nil {
		http.Error(w, ErrorMessage(err.Error()), http.StatusBadRequest)
		return
	}

	// Update the group name
	if err := rt.db.SetGroupName(convID, req.Name); err != nil {
		http.Error(w, ErrorMessage("Failed to update group name"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Group name changed successfully."})
}
