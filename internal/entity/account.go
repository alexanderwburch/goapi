package entity

import (
	"time"
)

type Account struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirebaseId string    `json:"firebase_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"created_at"`
}
