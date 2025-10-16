package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract message ID
	messageID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid Message Id"), http.StatusBadRequest)
		return
	}

	userID := ctx.User.ID

	// Get the message
	msg, err := rt.db.GetMessageByID(messageID)
	if err != nil || msg == nil {
		http.Error(w, ErrorMessage("Message Not Found"), http.StatusNotFound)
		return
	}

	// Make sure user is in the conversation before deleting the comment
	if err := rt.db.IsMember(msg.ConversationID, userID); err != nil {
		http.Error(w, ErrorMessage("User doesn't belong to this conversation"), http.StatusForbidden)
		return
	}

	// Check if comment exists
	comment, err := rt.db.GetCommentByUser(messageID, userID)
	if err != nil {
		http.Error(w, ErrorMessage("Database Error"), http.StatusInternalServerError)
		return
	}
	if comment == nil {
		http.Error(w, ErrorMessage("Comment Not Found"), http.StatusNotFound)
		return
	}

	// Delete the comment
	err = rt.db.DeleteComment(messageID, userID)
	if err != nil {
		http.Error(w, ErrorMessage("Failed to remove comment"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).
		Encode(map[string]string{"message": "Message reaction deleted successfully"})
}
