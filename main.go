package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"myGolangProject/Nastya/generator"
	"net/http"
	"strconv"
)

func changeSettings(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	var s generator.Settings
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(s.AmountMin, s.AmountMax)

	if s.AmountMin < 1 {
		http.NotFound(w, r)
		return
	}
	if s.AmountMax < 1 || s.AmountMax < s.AmountMin {
		http.NotFound(w, r)
		return
	}

	generator.ChangeSettingsAmount(s.AmountMin, s.AmountMax)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("setting was successfuly changed"))
	// w.Write([]byte(strconv.Itoa(s.AmountMin)))
}

func showTransactions(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if count < 1 || err != nil {
		http.NotFound(w, r)
		return
	}
	transactions, err := json.Marshal(generator.Generate(count))
	if err != nil {
		log.Println(err)
	}
	w.Write(transactions)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/generate/transactions", showTransactions)
	mux.HandleFunc("/change/settings", changeSettings)
	// mux.HandleFunc("/test", test)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
