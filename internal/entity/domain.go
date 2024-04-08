package entity

import (
	"time"
)

type Domain struct {
	ID        string    `json:"id"`
	AccountId string    `json:"account_id"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"created_at"`
}
