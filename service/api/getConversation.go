package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// ---------- API Models ----------

type MessageStatus struct {
	Status  string  `json:"status"`  // "read" or "received"
	Members []int64 `json:"members"` // user IDs who read the message
}

type Comment struct {
	AuthorID  int64  `json:"author_id"`
	Emoji     string `json:"emoji"`
	CreatedAt string `json:"created_at"`
}

type ConversationMember struct {
	UserID    int64  `json:"user_id"`
	PhotoPath string `json:"photo_path"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

type ConversationMessage struct {
	MessageID     int64         `json:"message_id"`
	Type          string        `json:"type"`
	IsForwarded   bool          `json:"is_forwarded"`
	Content       string        `json:"content"`
	AuthorID      int64         `json:"author_id"`
	CreatedAt     string        `json:"created_at"`
	MessageStatus MessageStatus `json:"message_status"`
	Comments      []Comment     `json:"comments"`
}

type GetConversationResponse struct {
	ConversationID int64                 `json:"conversation_id"`
	Name           string                `json:"name,omitempty"` // shown only for group chats
	Type           string                `json:"type"`
	PhotoPath      string                `json:"photo_path,omitempty"` // shown only for group chats
	CreatedBy      int64                 `json:"created_by"`
	CreatedAt      time.Time             `json:"created_at"`
	Members        []ConversationMember  `json:"members,omitempty"` // shown only for group chats
	Messages       []ConversationMessage `json:"messages"`
}

// ---------- Converters (Database → API) ----------

// Converts DB comments to API comments
func databaseToApiComments(dbComments []database.Comment) []Comment {
	apiComments := make([]Comment, len(dbComments))
	for i, c := range dbComments {
		apiComments[i] = Comment{
			AuthorID:  c.AuthorID,
			Emoji:     c.Emoji,
			CreatedAt: c.CreatedAt, // already a string in the DB model
		}
	}
	return apiComments
}

// Converts DB messages to API messages (with message status and comments)
func (rt *_router) databaseToApiConversationMessages(
	dbMessages []database.Message,
	totalMembers int,
) ([]ConversationMessage, error) {

	apiMessages := make([]ConversationMessage, 0, len(dbMessages))

	for _, msg := range dbMessages {
		readMembers, err := rt.db.GetReadMembersByMessage(msg.ID)
		if err != nil {
			return nil, err
		}

		status := "received"
		if len(readMembers) == totalMembers {
			status = "read"
		}

		dbComments, err := rt.db.GetCommentsByMessage(msg.ID)
		if err != nil {
			return nil, err
		}

		apiMessages = append(apiMessages, ConversationMessage{
			MessageID:   msg.ID,
			Type:        msg.Type,
			IsForwarded: msg.IsForwarded,
			Content:     msg.Content,
			AuthorID:    msg.AuthorID,
			CreatedAt:   msg.CreatedAt.UTC().Format(time.RFC3339),
			MessageStatus: MessageStatus{
				Status:  status,
				Members: readMembers,
			},
			Comments: databaseToApiComments(dbComments),
		})
	}

	return apiMessages, nil
}

func (rt *_router) databaseToApiMembers(conversationID int64, memberIDs []int64) ([]ConversationMember, error) {
	apiMembers := make([]ConversationMember, 0, len(memberIDs))

	for _, id := range memberIDs {
		user, err := rt.db.GetUserById(id)
		if err != nil {
			rt.baseLogger.Errorf("failed GetUserById(%d): %v", id, err)
			return nil, err
		}

		role := constants.ROLE_MEMBER
		role, err = rt.db.GetRoleByConversation(conversationID, id)
		if err != nil {
			rt.baseLogger.Errorf("failed GetRoleByConversation(%d, %d): %v", conversationID, id, err)
			return nil, err
		}

		photo := ""
		if user.PhotoPath != nil {
			photo = *user.PhotoPath
		}

		apiMembers = append(apiMembers, ConversationMember{
			UserID:    user.ID,
			PhotoPath: photo,
			Username:  user.Username,
			Role:      role,
		})
	}

	return apiMembers, nil
}

// ---------- Handler ----------

func (rt *_router) getConversation(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid conversation id"), http.StatusBadRequest)
		return
	}

	conversation, err := rt.db.GetConversationByID(convID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	if err := rt.db.IsMember(convID, ctx.User.ID); err != nil {
		http.Error(w, ErrorMessage("User doesn't belong to this conversation"), http.StatusForbidden)
		return
	}

	dbMessages, err := rt.db.GetMessagesByConversation(convID)
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

	apiMessages, err := rt.databaseToApiConversationMessages(dbMessages, totalMembers)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	resp := GetConversationResponse{
		ConversationID: conversation.ID,
		CreatedBy:      conversation.CreatedBy,
		CreatedAt:      conversation.CreatedAt,
		Type:           conversation.Type,
		Messages:       apiMessages,
	}

	if conversation.Type == constants.CONV_GROUP {
		apiMembers, err := rt.databaseToApiMembers(convID, members)
		if err != nil {
			http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
			return
		}
		resp.Members = apiMembers
		if conversation.Name != nil {
			resp.Name = *conversation.Name
		}
		if conversation.PhotoPath != nil {
			resp.PhotoPath = *conversation.PhotoPath
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
