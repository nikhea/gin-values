package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Upload limits and locations for user-supplied files.
const (
	// MaxUploadSize caps any single uploaded file (imports + avatars).
	MaxUploadSize = 5 << 20 // 5 MB
	// UploadBaseDir is the on-disk root served publicly at /uploads.
	UploadBaseDir = "uploads"
	// AvatarUsersDir / AvatarContactsDir hold avatar images per owner type.
	AvatarUsersDir    = "avatars/users"
	AvatarContactsDir = "avatars/contacts"
)

var (
	ErrFileTooLarge   = errors.New("file too large (max 5 MB)")
	ErrInvalidFileExt = errors.New("invalid file extension")
)

// allowedAvatarExts are the only image types accepted for avatars.
var allowedAvatarExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

// ValidateUploadSize rejects files over MaxUploadSize.
func ValidateUploadSize(size int64) error {
	if size > MaxUploadSize {
		return ErrFileTooLarge
	}
	return nil
}

// ValidateAvatarExt rejects non-image extensions (case-insensitive).
func ValidateAvatarExt(filename string) error {
	if !allowedAvatarExts[strings.ToLower(filepath.Ext(filename))] {
		return fmt.Errorf("%w: want one of .jpg .jpeg .png .gif .webp", ErrInvalidFileExt)
	}
	return nil
}

// ValidateImportExt rejects files whose extension isn't in allowed
// (e.g. ".csv" for the CSV endpoint, ".json" for the JSON endpoint).
func ValidateImportExt(filename, allowed string) error {
	if !strings.EqualFold(filepath.Ext(filename), allowed) {
		return fmt.Errorf("%w: want %s", ErrInvalidFileExt, allowed)
	}
	return nil
}

// AvatarTarget returns the public URL path and absolute disk path for a
// new avatar upload, creating the directory. The filename is a UUID so
// uploads can never collide or traverse directories.
func AvatarTarget(subdir, originalName string) (urlPath, diskPath string, err error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	rel := filepath.Join(subdir, uuid.New().String()+ext)
	abs := filepath.Join(UploadBaseDir, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return "", "", fmt.Errorf("create upload dir: %w", err)
	}
	return "/" + filepath.ToSlash(filepath.Join(UploadBaseDir, rel)), abs, nil
}
