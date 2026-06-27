package auth 

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/areyoush/surfspace/internal/models"
)

type Repository struct {
	db	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(email, passwordHash string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at`, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return u, err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *Repository) DenylistToken(token string, expires_at time.Time) error {
	query := `INSERT INTO denylist (token, expires_at) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING`
	_, err:= r.db.Exec(query, token, expires_at)

	return err
}

func (r *Repository) IsTokenDenylisted(token string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM denylist WHERE token = $1)"
	
	err := r.db.QueryRow(query, token).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking token denylist: %v", err)
		return false
	}
	return exists
}

func (r *Repository) CleanupDenylist() error {
	query := "DELETE FROM denylist WHERE expires_at < NOW()"
	
	_, err := r.db.Exec(query)
	if err != nil {
		log.Printf("Failed to cleanup token denylist: %v", err)
	}
	return err
	
}













