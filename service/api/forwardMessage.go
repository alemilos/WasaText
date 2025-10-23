package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type forwardMessageRequest struct {
	ConversationID int64 `json:"conversationId"`
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// === Parse message ID from path ===
	messageIDStr := ps.ByName("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// === Parse body for destination conversation ID ===
	var body forwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	destConvID := body.ConversationID

	// === Get requester’s user ID ===
	userID := ctx.User.ID
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// === Validate that the message exists ===
	originalMsg, err := rt.db.GetMessageByID(messageID)
	if err != nil || originalMsg == nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// === Validate that the destination conversation exists ===
	destConv, err := rt.db.GetConversationByID(destConvID)
	if err != nil || destConv == nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// === Ensure the requester is a member of the destination conversation ===
	if err := rt.db.IsMember(destConvID, userID); err != nil {
		http.Error(w, "User not in conversation", http.StatusUnauthorized)
		return
	}

	// === Create the forwarded message ===
	newMsg, err := rt.db.CreateMessage(
		destConvID,
		userID,
		originalMsg.Type,
		originalMsg.Content,
		originalMsg.SecondaryContent,
		true, // isForwarded
	)
	if err != nil {
		http.Error(w, "Failed to forward message", http.StatusInternalServerError)
		return
	}

	// Sender will mark as read its own message
	if err := rt.db.CreateMessageRead(newMsg.ID, ctx.User.ID); err != nil {
		ctx.Logger.WithError(err).Error("failed to mark message as read for sender")
	}

	// --- 8. Return the new message as JSON ---
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": newMsg,
	})
}
