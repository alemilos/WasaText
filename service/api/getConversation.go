package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type messageStatus struct {
	Status  string  `json:"status"`  // "read" or "received"
	Members []int64 `json:"members"` // user IDs who read the message
}

type conversationMessage struct {
	MessageID     int64         `json:"message_id"`
	Type          string        `json:"type"`
	IsForwarded   bool          `json:"is_forwarded"`
	Content       string        `json:"content"`
	AuthorID      int64         `json:"author_id"`
	CreatedAt     string        `json:"created_at"`
	MessageStatus messageStatus `json:"message_status"`
}

type getConversationResponse struct {
	Messages []conversationMessage `json:"messages"`
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid conversation id"), http.StatusBadRequest)
		return
	}

	if err := rt.db.IsMember(convID, ctx.User.ID); err != nil {
		http.Error(w, ErrorMessage("User doesn't belong to this conversation"), http.StatusForbidden)
		return
	}

	messages, err := rt.db.GetMessagesByConversation(convID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	members, err := rt.db.GetMembersByConversation(convID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}
	totalMembers := len(members)

	var resp getConversationResponse
	for _, msg := range messages {
		readMembers, err := rt.db.GetReadMembersByMessage(msg.ID)
		if err != nil {
			http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
			return
		}

		status := "received"
		if len(readMembers) == totalMembers {
			status = "read"
		}

		resp.Messages = append(resp.Messages, conversationMessage{
			MessageID:   msg.ID,
			Type:        msg.Type,
			IsForwarded: msg.IsForwarded,
			Content:     msg.Content,
			AuthorID:    msg.AuthorID,
			CreatedAt:   msg.CreatedAt.UTC().Format(time.RFC3339),
			MessageStatus: messageStatus{
				Status:  status,
				Members: readMembers,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
