package analytics

import (	

	"log"
	
	"database/sql"
	"github.com/areyoush/surfspace/internal/models"
	
)

type Repository struct {
	db	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RecordClick(click *models.Click) error {
	_, err := r.db.Exec(
		`INSERT INTO clicks (link_id, referrer, user_agent, country, city)
		 VALUES ($1, $2, $3, $4, $5)`,
		click.LinkID, click.Referrer, click.UserAgent, click.Country, click.City,
	)
	return err
}

func (r *Repository) GetClicksByLinkID(linkID string) ([]*models.Click, error) {
	rows, err := r.db.Query(
		`SELECT id, link_id, clicked_at, referrer, user_agent, country, city
		 FROM clicks WHERE link_id = $1 ORDER BY clicked_at DESC`,
		linkID,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()


	var clicks []*models.Click
	for rows.Next() {
		c := &models.Click{}
		if err := rows.Scan(&c.ID, &c.LinkID, &c.ClickedAt, &c.Referrer, &c.UserAgent, &c.Country, &c.City); err != nil {
			return nil, err
		}
		clicks = append(clicks, c)
	}
	return clicks, nil
}