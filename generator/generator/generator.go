package generator

import (
	"math/rand"
	"time"

	"myGolangProject/Nastya/models"

	"github.com/gofrs/uuid"
)

type Settings struct {
	AmountMin int `json:"min"`
	AmountMax int `json:"max"`
}

var transactionTypes = []string{"cash", "card", "sbp"} //интерфейс для каждого типа оплаты с последующими ручками???

func ChangeSettingsAmount(min, max int) {
	currentSettings.AmountMin = min
	currentSettings.AmountMax = max
}

var currentSettings Settings = Settings{
	AmountMin: 2,
	AmountMax: 1000,
}

func Generate(count int) []models.Transaction {
	uuid, _ := uuid.NewV4()
	// rand.Seed(time.Now().UnixNano())
	var transactions = make([]models.Transaction, 0)
	for i := 0; i < count; i++ {
		transactions = append(transactions, models.Transaction{
			UUID: uuid,
			Type: transactionTypes[rand.Intn(3)],
			Date: time.Now(),
			Amount: currentSettings.AmountMin +
				rand.Intn(currentSettings.AmountMax-currentSettings.AmountMin+1), //сумма
			SenderID:    rand.Intn(1000),
			RecipientID: rand.Intn(1000),
		})
	}
	return transactions
}
