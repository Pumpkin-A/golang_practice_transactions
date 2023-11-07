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
	SenderID    int       `json:"senderId"`
	RecipientID int       `json:"recipientId"`
}

type Config struct {
	AvailableXUsers []string
	ServerPort      string
}

var GlobalConfig = Config{
	AvailableXUsers: []string{"Nastya", "Maxim", "clown"},
	ServerPort:      ":8080",
}
