package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strconv"

	"myGolangProject/Nastya/generator"
	"myGolangProject/Nastya/models"
)

type jsonResponse struct {
	Text string `json:"-_-"`
}

func serviceUnavailable(w http.ResponseWriter) bool {
	if rand.Int()%2 == 0 {
		sr := &jsonResponse{"Service Unavailable(((("}
		jr, _ := json.Marshal(sr)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(jr)
		return true
	}
	return false
}

func changeSettings(w http.ResponseWriter, r *http.Request) {
	if serviceUnavailable(w) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("reading the body error", err)
		return
	}
	log.Println(string(body))
	var s generator.Settings
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Println("unmarshal error", err)
		return
	}
	log.Println(s.AmountMin, s.AmountMax)

	if s.AmountMin < 1 {
		sr := &jsonResponse{"Bad request(((("}
		jr, _ := json.Marshal(sr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jr)
		return
	}
	if s.AmountMax < 1 || s.AmountMax < s.AmountMin {
		sr := &jsonResponse{"Bad request(((("}
		jr, _ := json.Marshal(sr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jr)
		return
	}

	generator.ChangeSettingsAmount(s.AmountMin, s.AmountMax)
	w.WriteHeader(http.StatusOK)
	sr := &jsonResponse{"settings was successfuly changed"}
	jr, _ := json.Marshal(sr)
	w.Write(jr)
}

func xUserAccess(w http.ResponseWriter, r *http.Request) bool {
	ua := r.Header.Get("X-User")
	if !slices.Contains(models.GlobalConfig.AvailableXUsers, ua) {
		sr := &jsonResponse{"Unauthorized"}
		jr, _ := json.Marshal(sr)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jr)
		return true
	}
	return false
}

func showTransactions(w http.ResponseWriter, r *http.Request) {
	if serviceUnavailable(w) {
		return
	}

	if xUserAccess(w, r) {
		return
	}
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if count < 1 || err != nil {
		sr := &jsonResponse{"Not found((("}
		jr, _ := json.Marshal(sr)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jr)
		return
	}
	transactions, err := json.Marshal(generator.Generate(count))
	if err != nil {
		log.Println("marshal error", err)
	}
	w.Write(transactions)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/generate/transactions", showTransactions)
	mux.HandleFunc("/change/settings", changeSettings)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
