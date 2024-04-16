package entity

import (
	"time"
)

type Domain struct {
	ID        int       `db:"id"`
	AccountId int       `json:"account_id"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"created_at"`
}
