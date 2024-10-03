package domain

import (
	"time"
)

type Save struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	MapLink   string    `json:"mapLink"`
	Link      string    `json:"link"`
}
