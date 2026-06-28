package models

import (
	"time"
)

type Link struct {
	ID				string
	UserID			string
	OriginalURL		string
	ShortCode		string
	ClickCount		int
	CreatedAt		time.Time
	ExpiresAt		*time.Time		
}	