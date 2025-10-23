package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type sendTextMessageRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	SecondaryContent string `json:"secondaryContent"`
}

type sendMessageResponse struct {
	MessageID int64     `json:"messageId"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	SecondaryContent string `json:"secondaryContent"`
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	convIDStr := ps.ByName("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid conversation id"), http.StatusBadRequest)
		return
	}

	if err := rt.db.IsMember(convID, ctx.User.ID); err != nil {
		ctx.Logger.WithError(err).Error("failed to check conversation membership")
		http.Error(w, ErrorMessage("User doesn't belong to the conversation"), http.StatusForbidden)
		return
	}

	ct := r.Header.Get("Content-Type")

	// === CASE 1: application/json (text message) ===
	if ct == "application/json" || ct == "application/json; charset=utf-8" {
		var req sendTextMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, ErrorMessage(InvalidRequestBody), http.StatusBadRequest)
			return
		}

		if req.Type != "text" {
			http.Error(w, ErrorMessage("type can either be 'text' or 'image'"), http.StatusBadRequest)
			return
		}
		if req.Content == "" {
			http.Error(w, ErrorMessage("content is required"), http.StatusBadRequest)
			return
		}

		msg, err := rt.db.CreateMessage(convID, ctx.User.ID, req.Type, req.Content, "", false)
		if err != nil {
			http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
			return
		}

		// Mark sender as having read their own message
		if err := rt.db.CreateMessageRead(msg.ID, ctx.User.ID); err != nil {
			ctx.Logger.WithError(err).Error("failed to mark message as read for sender")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sendMessageResponse{
			MessageID: msg.ID,
			Type:      msg.Type,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
		return
	}

	// === CASE 2: multipart/form-data (image message) ===
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, ErrorMessage("Invalid Form Data"), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	secondaryContent := r.FormValue("secondaryContent")
	if err != nil {
		http.Error(w, ErrorMessage("No File Uploaded"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	validTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}

	ext, ok := validTypes[contentType]
	if !ok {
		http.Error(w, ErrorMessage("Wrong File Format"), http.StatusBadRequest)
		return
	}

	// Limit file size and prepare destination
	limitedFile := http.MaxBytesReader(w, file, constants.MAX_UPLOAD_SIZE)
	uploadDir := filepath.Join(rt.storagePath, "messages", fmt.Sprintf("%d", convID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.Logger.WithError(err).Error("failed to create upload dir")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Make sure the file can be stored in the upload directory
	tempFile, err := os.CreateTemp(uploadDir, "upload-*"+ext)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to create temp file")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tempFile, limitedFile)
	tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		ctx.Logger.WithError(err).Error("failed to write temp image")
		http.Error(w, ErrorMessage("File Too Big or Invalid"), http.StatusBadRequest)
		return
	}

	// create DB message (only after successful local write)
	msg, err := rt.db.CreateMessage(convID, ctx.User.ID, "image", "", secondaryContent, false)
	if err != nil {
		os.Remove(tempFile.Name())
		ctx.Logger.WithError(err).Error("failed to create image message")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// rename the file and use the msg id just created
	finalName := fmt.Sprintf("%d%s", msg.ID, ext)
	finalPath := filepath.Join(uploadDir, finalName)
	if err := os.Rename(tempFile.Name(), finalPath); err != nil {
		os.Remove(tempFile.Name())
		ctx.Logger.WithError(err).Error("failed to move temp image to final location")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// --- Step 4: update message content in DB ---
	relPath := path.Join(constants.MESSAGE_PHOTO_PATH, fmt.Sprintf("%d", convID), finalName)
	if err := rt.db.UpdateMessageContent(msg.ID, relPath); err != nil {
		ctx.Logger.WithError(err).Error("failed to update message content")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Mark sender as having read their own message
	if err := rt.db.CreateMessageRead(msg.ID, ctx.User.ID); err != nil {
		ctx.Logger.WithError(err).Error("failed to mark message as read for sender")
	}

	// Update msg.Content for response
	msg.Content = relPath

	// --- Step 5: return response ---
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sendMessageResponse{
		MessageID: msg.ID,
		Type:      msg.Type,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	})
}
