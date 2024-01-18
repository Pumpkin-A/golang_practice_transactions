package main

import (
	"context"
	"database/sql"
	"fmt"
	generatorClient "go/transaction/client/generatorClient"
	"go/transaction/client/models"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "postgres"
)

func customHandler(w http.ResponseWriter, r *http.Request) error {
	xUserHeader := r.Header.Get("X-User")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("db open error", err)
		return err
	}
	defer db.Close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "X-User", xUserHeader)
	transactionsData, err := generatorClient.GetTransaction(ctx)
	if err != nil {
		log.Println("getRespons func error", err)
		return err
	}
	// fmt.Printf("%v\n", transactionsData)
	err = insertRows(db, transactionsData)
	if err != nil {
		log.Println("insert error", err)
		return err
	}

	w.Write([]byte("Data insertion was successfully completed"))
	return nil

}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// handle returned error here.
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func insertRows(db *sql.DB, transactionsData []models.Transaction) error {
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
	port := ":3333"
	log.Printf("сервер запущен на %s порту", port)
	r := chi.NewRouter()
	r.Method("GET", "/", Handler(customHandler))
	http.ListenAndServe(port, r)
}

// err := headersCheck("http://127.0.0.1:8080/generate/transactions?count=3")
// if err != nil {
// 	log.Println("headersCheck func error", err)
// }
