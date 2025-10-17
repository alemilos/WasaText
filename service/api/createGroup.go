package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/validation"
	"github.com/julienschmidt/httprouter"
)

type createGroupRequest struct {
	Name    string  `json:"name"`
	Members []int64 `json:"members"`
}

type createGroupResponse struct {
	ConversationID int64     `json:"conversationId"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateGroupName(req.Name); err != nil {
		http.Error(w, ErrorMessage(err.Error()), http.StatusBadRequest)
		return
	}

	if len(req.Members) == 0 {
		http.Error(w, ErrorMessage("At least one member is required"), http.StatusBadRequest)
		return
	}

	userID := ctx.User.ID

	// Make sure the member ids are real users
	validMembers := make([]int64, 0)
	for _, id := range req.Members {
		if id == userID {
			continue // skip the creator; they'll be added as admin automatically
		}

		_, err := rt.db.GetUserById(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue // skip invalid users
			}
			http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
			return
		}

		validMembers = append(validMembers, id)
	}

	if len(validMembers) == 0 {
		http.Error(w, ErrorMessage("No valid members found in request"), http.StatusBadRequest)
		return
	}

	// Create a new group conversation
	newGroup, err := rt.db.CreateGroupConversation(req.Name, userID, validMembers)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createGroupResponse{
		ConversationID: newGroup.ID,
		Type:           newGroup.Type,
		Name:           *newGroup.Name,
		CreatedBy:      newGroup.CreatedBy,
		CreatedAt:      newGroup.CreatedAt,
	})
}
