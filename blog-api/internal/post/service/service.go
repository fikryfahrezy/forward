package service

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/fikryfahrezy/forward/blog-api/internal/post/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func generateSlug(title string) (string, error) {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Add random suffix
	suffix := make([]byte, 4)
	_, err := rand.Read(suffix)
	if err != nil {
		return "", err
	}
	slug = slug + "-" + hex.EncodeToString(suffix)

	return slug, nil
}
