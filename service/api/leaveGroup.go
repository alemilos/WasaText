package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type leaveGroupRequest struct {
	UserId int64 `json:"userId"`
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// get the conversation id (grp id) from the params
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid group id"), http.StatusBadRequest)
		return
	}

	// get the user id from the request
	var req leaveGroupRequest
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

	// make sure the target user is in the group
	if err := rt.db.IsMember(convID, req.UserId); err != nil {
		http.Error(w, "User not in conversation", http.StatusUnauthorized)
		return
	}

	// Get roles
	requesterRole, err := rt.db.GetRoleByConversation(convID, requesterID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	targetRole, err := rt.db.GetRoleByConversation(convID, req.UserId)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// if user id === ctx.User.ID && role is 'admin' -> the first other user who joined the group will become admin
	// if user id === ctx.User.ID -> user is removed from group
	// if user id !== ctx.User.ID && role is 'admin' -> user is removed from group
	if requesterID == req.UserId {
		if targetRole == constants.ROLE_ADMIN {
			// Get all members
			members, err := rt.db.GetMembersByConversation(convID)
			if err != nil {
				http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
				return
			}

			// Remove self
			if err := rt.db.RemoveMember(convID, requesterID); err != nil {
				http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
				return
			}

			// Check if another admin already exists
			alreadyHasAdmin := false
			for _, m := range members {
				if m == requesterID {
					continue
				}
				role, err := rt.db.GetRoleByConversation(convID, m)
				if err != nil {
					http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
					return
				}
				if role == constants.ROLE_ADMIN {
					alreadyHasAdmin = true
					break
				}
			}

			// If no admin left, promote first available member
			if !alreadyHasAdmin {
				for _, m := range members {
					if m != requesterID {
						if err := rt.db.SetRole(convID, m, constants.ROLE_ADMIN); err != nil {
							http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
							return
						}
						break
					}
				}
			}
		} else {
			// Just remove the member
			if err := rt.db.RemoveMember(convID, requesterID); err != nil {
				http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Admin removes another user
	if requesterRole == constants.ROLE_ADMIN {
		if err := rt.db.RemoveMember(convID, req.UserId); err != nil {
			http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Non-admin tries to remove another member
	http.Error(w, ErrorMessage("Only admins can remove other members"), http.StatusForbidden)

}
