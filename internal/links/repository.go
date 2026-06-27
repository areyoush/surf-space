package links

import (
	"database/sql"

	"github.com/areyoush/surfspace/internal/models"
)


type Repository struct {
	db	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLink(userID, originalURL, shortCode string)	(*models.Link, error) {
	l := &models.Link{}
	err := r.db.QueryRow(`INSERT INTO links (user_id, original_url, short_code) VALUES ($1, $2, $3) RETURNING id, user_id, original_url, short_code, click_count, created_at`, userID, originalURL, shortCode,).Scan(&l.ID, &l.UserID, &l.OriginalURL, &l.ShortCode, &l.ClickCount, &l.CreatedAt)
	return l, err 
} 

func (r *Repository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	l := &models.Link{}
	err := r.db.QueryRow(
		`SELECT id, user_id, original_url, short_code, click_count, created_at 
		 FROM links WHERE short_code = $1`,
		shortCode,
	).Scan(&l.ID, &l.UserID, &l.OriginalURL, &l.ShortCode, &l.ClickCount, &l.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return l, err
}

func (r *Repository) IncrementClickCount(shortCode string) error {
	_, err := r.db.Exec(
		`UPDATE links SET click_count = click_count + 1 WHERE short_code = $1`,
		shortCode,
	)
	return err
}