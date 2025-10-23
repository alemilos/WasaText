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

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type setGroupPhotoResponse struct {
	Message   string `json:"message"`
	PhotoPath string `json:"photoPath"`
}

// CleanupOldGroupPhotos deletes all existing group photos for a given conversationID,
// except the newly uploaded file.
func CleanupOldGroupPhotos(conversationID int64, uploadDir string, keepExt string) error {
	extensions := []string{".jpg", ".png", ".webp"}
	fileBaseName := fmt.Sprintf("%d", conversationID)

	var errAcc error

	for _, ext := range extensions {
		if ext == keepExt {
			continue // skip the new file
		}

		oldFile := filepath.Join(uploadDir, fileBaseName+ext)

		if _, err := os.Stat(oldFile); err == nil {
			if removeErr := os.Remove(oldFile); removeErr != nil {
				errAcc = fmt.Errorf("failed to remove %s: %w", oldFile, removeErr)
			}
		}
	}

	return errAcc
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	convIDStr := ps.ByName("id")
	conversationID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, ErrorMessage("Invalid conversation id"), http.StatusBadRequest)
		return
	}

	// Check conversation exists
	conversation, err := rt.db.GetConversationByID(conversationID)
	if err != nil {
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Ensure requester is member
	if err := rt.db.IsMember(conversationID, ctx.User.ID); err != nil {
		http.Error(w, ErrorMessage("User doesn't belong to this conversation"), http.StatusForbidden)
		return
	}

	// Ensure it's a group
	if conversation.Type != constants.CONV_GROUP {
		http.Error(w, ErrorMessage("Conversation is not a group"), http.StatusBadRequest)
		return
	}

	uploadDir := filepath.Join(rt.storagePath, constants.GROUP_PHOTO_PATH)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.Logger.WithError(err).Error("failed to create upload directory")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, ErrorMessage("Invalid form data"), http.StatusBadRequest)
		return
	}

	// Retrieve file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, ErrorMessage("No File Uploaded"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Limit file size
	limitedFile := http.MaxBytesReader(w, file, constants.MAX_UPLOAD_SIZE)
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

	// Save file as <conversationID>.<ext>
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%d%s", conversationID, ext))
	relPath := path.Join(constants.GROUP_PHOTO_PATH, fmt.Sprintf("%d%s", conversationID, ext))

	out, err := os.Create(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to create file")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, limitedFile); err != nil {
		ctx.Logger.WithError(err).Error("failed to save file")
		http.Error(w, ErrorMessage("File Too Big"), http.StatusBadRequest)
		return
	}

	// Cleanup old photos (jpg/png/webp), keep only the current one
	if err := CleanupOldGroupPhotos(conversationID, uploadDir, ext); err != nil {
		ctx.Logger.WithError(err).Warn("failed to cleanup old group photos")
	}

	// Update DB
	if err := rt.db.SetGroupPhoto(conversationID, relPath); err != nil {
		ctx.Logger.WithError(err).Error("failed to update group photo in DB")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(setGroupPhotoResponse{
		Message:   "Group photo updated successfully",
		PhotoPath: relPath,
	})
}
