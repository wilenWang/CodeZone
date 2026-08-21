package users

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"codezone/backend/internal/httpx"
	"codezone/backend/internal/storage"
)

const maxAvatarBytes = 5 << 20

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	var req struct {
		DisplayName string `json:"displayName"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}
	name := strings.TrimSpace(req.DisplayName)
	if name == "" {
		httpx.WriteError(w, http.StatusBadRequest, "display_name_invalid", "Display name is required")
		return
	}
	user, err := h.service.UpdateProfile(r.Context(), userID, name, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "profile_update_failed", "Could not update profile")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) UploadAvatar(store storage.PublicStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := httpx.UserID(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
			return
		}
		if store == nil {
			httpx.WriteError(w, http.StatusServiceUnavailable, "avatar_storage_unavailable", "Avatar storage is unavailable")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxAvatarBytes+1024)
		if err := r.ParseMultipartForm(maxAvatarBytes); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "avatar_too_large", "Avatar must be 5 MB or smaller")
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "avatar_required", "Avatar file is required")
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil || len(content) == 0 || len(content) > maxAvatarBytes {
			httpx.WriteError(w, http.StatusBadRequest, "avatar_invalid", "Invalid avatar file")
			return
		}
		mime, ext, ok := avatarType(content)
		if !ok || (header.Header.Get("Content-Type") != "" && header.Header.Get("Content-Type") != mime) {
			httpx.WriteError(w, http.StatusBadRequest, "avatar_type_invalid", "Avatar must be JPEG, PNG, or WebP")
			return
		}
		var id [16]byte
		if _, err := rand.Read(id[:]); err != nil {
			httpx.WriteError(w, 500, "avatar_upload_failed", "Could not upload avatar")
			return
		}
		key := fmt.Sprintf("%s/avatars/%d/%s.%s", strings.Trim(h.pathPrefix, "/"), userID, hex.EncodeToString(id[:]), ext)
		if err := store.Upload(r.Context(), key, bytes.NewReader(content), mime); err != nil {
			httpx.WriteError(w, 502, "avatar_upload_failed", "Could not upload avatar")
			return
		}
		url := store.URL(key)
		user, err := h.service.UpdateProfile(r.Context(), userID, "", &url)
		if err != nil {
			httpx.WriteError(w, 500, "profile_update_failed", "Could not update profile")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, user)
	}
}

func avatarType(data []byte) (string, string, bool) {
	switch {
	case len(data) >= 3 && bytes.Equal(data[:3], []byte{0xff, 0xd8, 0xff}):
		return "image/jpeg", "jpg", true
	case len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}):
		return "image/png", "png", true
	case len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")):
		return "image/webp", "webp", true
	default:
		return "", "", false
	}
}
