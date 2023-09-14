package main

//сделать запрос на сервер и вывести json в командную строку +
//подключиться к базе данных +
//запрос к бд +
//получать один json и переводить его в структуру +
//научиться делать insert запрос с исскуственными данными +
//вводить полученные данные в базу данных +
// Transaction сделать модулем и импортировать из models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "postgres"
)

type Transaction struct {
	UUID        uuid.UUID `json:"uuid"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	Amount      int       `json:"amount"`
	SenderID    int       `json:"sender id"`
	RecipientID int       `json:"recipient id"`
}

func getResponse(url string, j []Transaction) ([]Transaction, error) {
	client := http.Client{
		Timeout: 6 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("request error", err)
		return nil, err
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		log.Println("decode error", err)
		return nil, err
	}
	return j, nil
}

func insertRows(db *sql.DB, transactionsData []Transaction) error {
	stmt := `INSERT INTO transactionsdata (uuid, type, date, amount, senderid, recipientid)
	VALUES($1, $2, $3, $4, $5, $6)`
	for _, tr := range transactionsData {
		_, err := db.Exec(stmt, tr.UUID, tr.Type, tr.Date,
			tr.Amount, tr.SenderID, tr.RecipientID)
		if err != nil {
			log.Println("exec error")
			return err
		} else {
			log.Println("Row inserted successfully!")
		}
	}
	return nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("db open error", err)
	}
	defer db.Close()

	var j []Transaction
	url := "http://127.0.0.1:8080/generate/transactions?count=3"
	transactionsData, err := getResponse(url, j)
	if err != nil {
		log.Println("getRespons func error", err)
	}
	err = insertRows(db, transactionsData)
	if err != nil {
		log.Println("insert error", err)
	}
}
