package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// === Parse message ID from path ===
	messageIDStr := ps.ByName("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid message ID"), http.StatusBadRequest)
		return
	}

	// === Retrieve the message ===
	msg, err := rt.db.GetMessageByID(messageID)
	if err != nil || msg == nil {
		http.Error(w, ErrorMessage("Message not found"), http.StatusNotFound)
		return
	}

	// === Check author ===
	userID := ctx.User.ID
	if msg.AuthorID != userID {
		http.Error(w, ErrorMessage("Unauthorized"), http.StatusUnauthorized)
		return
	}

	// === Check if created less than 30 minutes ago ===
	if time.Since(msg.CreatedAt) > 30*time.Minute {
		http.Error(w, ErrorMessage("Cannot delete messages older than 30 minutes"), http.StatusForbidden)
		return
	}

	// === Delete the message ===
	if err := rt.db.DeleteMessage(messageID); err != nil {
		http.Error(w, ErrorMessage("Failed to delete message"), http.StatusInternalServerError)
		return
	}

	// === Return success JSON ===
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Message deleted successfully",
	})
}
