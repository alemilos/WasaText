package api

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type Conversation struct {
	ConversationID     int64     `json:"conversationId"`
	Type               string    `json:"type"` // "private" | "group"
	Name               *string   `json:"name"`
	PhotoPath          *string   `json:"photoPath"`
	CreatedBy          int64     `json:"createdBy"`
	CreatedAt          time.Time `json:"createdAt"`
	LastMessage        *Message  `json:"lastMessage"`
	OtherParticipantID *int64    `json:"otherParticipantId,omitempty"` // if the conversation is private
}

type Message struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"` // "text" | "image"
	IsForwarded bool      `json:"isForwarded"`
	AuthorID    int64     `json:"authorId"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}

type getMyConversationsResponse struct {
	Conversations []Conversation `json:"conversations"`
}

// a serializer, converter from database user to api response user.
func databaseToApiConversations(dbConversations []database.Conversation) []Conversation {
	apiConversations := make([]Conversation, len(dbConversations))
	for i, c := range dbConversations {
		apiConversations[i] = Conversation{
			ConversationID: c.ID,
			Type:           c.Type,
			Name:           c.Name,
			PhotoPath:      c.PhotoPath,
			CreatedBy:      c.CreatedBy,
			CreatedAt:      c.CreatedAt,
		}
	}
	return apiConversations
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := ctx.User.ID

	conversations, err := rt.db.GetConversationsByUserID(userID)
	if err != nil {
		http.Error(w, "failed to fetch conversations", http.StatusInternalServerError)
		return
	}

	apiConversations := databaseToApiConversations(conversations)

	for i, conversation := range apiConversations {
		if conversation.Type == constants.CONV_PRIVATE {
			members, err := rt.db.GetMembersByConversation(conversation.ConversationID)
			if err == nil {
				for _, m := range members {
					if m != userID {
						apiConversations[i].OtherParticipantID = &m
						break
					}
				}
			}
		}

		lastMsg, err := rt.db.GetLastMessageByConversation(conversation.ConversationID)
		if err == nil && lastMsg != nil {
			apiConversations[i].LastMessage = &Message{
				ID:          lastMsg.ID,
				Type:        lastMsg.Type,
				IsForwarded: lastMsg.IsForwarded,
				AuthorID:    lastMsg.AuthorID,
				Content:     lastMsg.Content,
				CreatedAt:   lastMsg.CreatedAt,
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		getMyConversationsResponse{Conversations: apiConversations})

}
