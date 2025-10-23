package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type addToGroupRequest struct {
	UserId int64 `json:"userId"`
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// get the conversation id (grp id) from the params
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid group id"), http.StatusBadRequest)
		return
	}

	// get the user id from the request
	var req addToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
		return
	}

	if req.UserId == 0 {
		http.Error(w, ErrorMessage("User id is required"), http.StatusBadRequest)
		return
	}

	requesterID := ctx.User.ID

	// Make sure conversation type is 'group'
	conversation, err := rt.db.GetConversationByID(convID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	if conversation.Type != constants.CONV_GROUP {
		http.Error(w, ErrorMessage("Cannot add members to a non-group conversation"), http.StatusBadRequest)
		return
	}

	// make sure the requester is member of conversation
	if err := rt.db.IsMember(convID, requesterID); err != nil {
		http.Error(w, "User not in conversation", http.StatusUnauthorized)
		return
	}

	// make sure the requester role is 'admin'
	role, err := rt.db.GetRoleByConversation(convID, requesterID)
	if err != nil {
		// unexpected DB error
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}
	if role != constants.ROLE_ADMIN {
		http.Error(w, ErrorMessage("Forbidden"), http.StatusInternalServerError)
		return
	}

	// Make sure user is not adding himself
	if req.UserId == requesterID {
		http.Error(w, ErrorMessage("User is already a member of the group"), http.StatusBadRequest)
		return
	}

	// make sure the user is not already in the group
	if err := rt.db.IsMember(convID, req.UserId); err == nil {
		http.Error(w, ErrorMessage("User is already a member of the group"), http.StatusBadRequest)
		return
	}

	// add the user to the group conversation
	if err := rt.db.AddMember(convID, req.UserId, constants.ROLE_MEMBER); err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// === Return success JSON ===
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to group",
	})
}
