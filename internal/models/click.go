package models

import (

	"time"
	
)

type Click struct {
	ID        string
	LinkID    string
	ClickedAt time.Time
	Referrer  string
	UserAgent string
	Country   string
	City      string
}