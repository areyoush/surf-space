package analytics

import (

	"github.com/areyoush/surfspace/internal/models"

)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordClick(linkID, referrer, userAgent string) error {
	click := &models.Click{
		LinkID:    linkID,
		Referrer:  referrer,
		UserAgent: userAgent,
		Country:   "",
		City:      "",
	}
	return s.repo.RecordClick(click)
}

func (s *Service) GetAnalytics(linkID string) ([]*models.Click, error) {
	return s.repo.GetClicksByLinkID(linkID)
}