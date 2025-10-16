package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/validation"
	"github.com/julienschmidt/httprouter"
)

type commentMessageRequest struct {
	Emoji string `json:"emoji"`
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract message ID
	messageID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid Message Id"), http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var body commentMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, ErrorMessage("emoji is required"), http.StatusBadRequest)
		return
	}

	// Validate emoji
	if err := validation.ValidateEmoji(body.Emoji); err != nil {
		http.Error(w, ErrorMessage("Invalid emoji unicode"), http.StatusBadRequest)
		return
	}

	userID := ctx.User.ID

	// Check if message exists
	msg, err := rt.db.GetMessageByID(messageID)
	if err != nil || msg == nil {
		http.Error(w, ErrorMessage("Message Not Found"), http.StatusNotFound)
		return
	}

	// Check if comment already exists
	_, err = rt.db.GetCommentByUser(messageID, userID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Make sure user is in the conversation before adding the comment
	if err := rt.db.IsMember(msg.ConversationID, userID); err != nil {
		http.Error(w, ErrorMessage("User doesn't belong to this conversation"), http.StatusForbidden)
		return
	}

	// Add or update comment
	err = rt.db.AddOrUpdateComment(messageID, userID, body.Emoji)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Response message
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Reacted to message successfully"})
}
