package api

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type createConversationRequest struct {
	RecipientID int64 `json:"recipientId"`
}

type createConversationResponse struct {
	ConversationID     int64                `json:"conversationId"`
	Type               string               `json:"type"`
	CreatedBy          int64                `json:"createdBy"`
	CreatedAt          time.Time            `json:"createdAt"`
	Members            []ConversationMember `json:"members,omitempty"`            // shown only for group chats
	OtherParticipantID *int64               `json:"otherParticipantId,omitempty"` // if the conversation is private
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decode request body
	var req createConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
		return
	}

	if req.RecipientID == 0 {
		http.Error(w, ErrorMessage("Missing recipient id"), http.StatusBadRequest)
		return
	}

	userID := ctx.User.ID

	// Recipient and requester must be different users
	if req.RecipientID == userID {
		http.Error(w, ErrorMessage("The recipient must be another user"), http.StatusBadRequest)
		return
	}

	// Check if user exists
	user, err := rt.db.GetUserById(req.RecipientID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, ErrorMessage("User not found"), http.StatusNotFound)
		return
	}

	// Check if private conversation already exists
	conversation, err := rt.db.GetPrivateConversation(userID, req.RecipientID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	var resp createConversationResponse

	if conversation != nil {
		// Existing conversation found -> return 200
		resp = createConversationResponse{
			ConversationID: conversation.ID,
			Type:           conversation.Type,
			CreatedBy:      conversation.CreatedBy,
			CreatedAt:      conversation.CreatedAt,
		}

		if conversation.Type == constants.CONV_GROUP {
			members, err := rt.db.GetMembersByConversation(conversation.ID)
			if err == nil {
				apiMembers, err := rt.databaseToApiMembers(conversation.ID, members)
				if err == nil {
					resp.Members = apiMembers
				}
			}
		} else {
			resp.OtherParticipantID = &req.RecipientID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		return

	}

	// Create a new private conversation
	newConv, err := rt.db.CreatePrivateConversation(userID, req.RecipientID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	resp = createConversationResponse{
		ConversationID:     newConv.ID,
		Type:               newConv.Type,
		CreatedBy:          newConv.CreatedBy,
		CreatedAt:          newConv.CreatedAt,
		OtherParticipantID: &req.RecipientID,
	}

	// newly created conversation -> 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)

}
