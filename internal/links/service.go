package links

import (

	"math/rand"
	
	"github.com/areyoush/surfspace/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ShortenURL(userID, originalURL string) (*models.Link, error) {
	shortCode := generateShortCode()
	return s.repo.CreateLink(userID, originalURL, shortCode)
}

func (s *Service) GetLink(shortCode string) (*models.Link, error) {
	return s.repo.GetLinkByShortCode(shortCode)
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}