package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/constants"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// CleanupOldProfilePhotos deletes all existing profile photos for a given userID,
// except the newly uploaded file.
func CleanupOldProfilePhotos(userID int64, uploadDir string, keepExt string) error {
	extensions := []string{".jpg", ".png", ".webp"}
	fileBaseName := fmt.Sprintf("%d", userID)

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

type setMyPhotoResponse struct {
	Message   string `json:"message"`
	PhotoPath string `json:"photo_path"`
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uploadDir := filepath.Join(rt.storagePath, constants.PROFILE_PHOTO_PATH)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.Logger.WithError(err).Error("failed to create upload directory")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	// Parse multipart form without strict limit first
	if err := r.ParseMultipartForm(10 << 20); err != nil { // allow up to 10MB temp form parsing
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

	// Now limit the file size
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

	userID := ctx.User.ID
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%d%s", userID, ext))
	relPath := fmt.Sprintf("%s%d%s", constants.PROFILE_PHOTO_PATH, userID, ext)

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

	// Now cleanup old photos (avoid cases where you have all jpg/png/webp)
	if err := CleanupOldProfilePhotos(userID, uploadDir, ext); err != nil {
		ctx.Logger.WithError(err).Warn("failed to cleanup old profile photos")
	}

	if err := rt.db.SetProfilePhoto(userID, relPath); err != nil {
		ctx.Logger.WithError(err).Error("failed to update user photo in DB")
		http.Error(w, ErrorMessage(InternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(setMyPhotoResponse{
		Message:   "Profile photo updated",
		PhotoPath: relPath,
	})
}
