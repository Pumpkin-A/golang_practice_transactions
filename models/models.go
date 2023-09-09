package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Transaction struct {
	UUID        uuid.UUID `json:"uuid"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	Amount      int       `json:"amount"`
	SenderID    int       `json:"sender id"`
	RecipientID int       `json:"recipient id"`
}
