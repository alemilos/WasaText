package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type readMessagesRequest struct {
	Messages []int64 `json:"messages"`
}

func (rt *_router) readMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid conversation id"), http.StatusBadRequest)
		return
	}

	// Check user membership in conversation
	if err := rt.db.IsMember(convID, ctx.User.ID); err != nil {
		http.Error(w, ErrorMessage("User doesn't belong to the conversation"), http.StatusForbidden)
		return
	}

	// Parse body
	var req readMessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrorMessage("Invalid request body"), http.StatusBadRequest)
		return
	}

	if len(req.Messages) == 0 {
		http.Error(w, ErrorMessage("messages field cannot be empty"), http.StatusBadRequest)
		return
	}

	// Iterate and mark each message as read
	for _, messageID := range req.Messages {
		if err := rt.db.CreateMessageRead(messageID, ctx.User.ID); err != nil {
			ctx.Logger.WithError(err).Errorf("failed to mark message %d as read", messageID)
			// partial success is allowed, just continue to next message
			continue
		}
	}

	// Return success (204)
	w.WriteHeader(http.StatusNoContent)
}
