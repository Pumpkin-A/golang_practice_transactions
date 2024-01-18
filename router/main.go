// /hello - просто отдёт 200
// /api/v1/numbers - принимает джсон с массивом цифр, возвращает умноженный на 2
// /api/v1/text- принимает джсон с текстом, возвращает текст с уникальными буквами set
// /api/v2/numbers - принимает джсон с массивом цифр, возвращает умноженный на 3
// /api/v2/text- принимает джсон с текстом, возвращает текст с повторных букв (>1)

// Ручки v1 и ручки v2 должны работать с разными данными авторизации (принимать разные заголовки для работы),
// для hello авторизация не нужна

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nastya/router/changeStrings"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type NumbersArray struct {
	Numbers []int `json:"numbers"`
}

type StringArray struct {
	Text string `json:"text"`
}

func v1Num(w http.ResponseWriter, r *http.Request) {
	var nums NumbersArray
	// unmarshall json string to struct
	err := json.NewDecoder(r.Body).Decode(&nums)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i := range nums.Numbers {
		nums.Numbers[i] *= 2
	}
	newNums, err := json.Marshal(nums)
	if err != nil {
		log.Println("marshaling error")
		return
	}
	w.Write([]byte(newNums))
}

func v1Text(w http.ResponseWriter, r *http.Request) {
	var text StringArray
	// unmarshall json string to struct
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	changedStr := changeStrings.RemoveDuplicates(text.Text)
	newText, err := json.Marshal(changedStr)
	if err != nil {
		log.Println("marshaling error")
		return
	}
	w.Write([]byte(newText))

}

func v2Num(w http.ResponseWriter, r *http.Request) {
	var nums NumbersArray
	// unmarshall json string to struct
	err := json.NewDecoder(r.Body).Decode(&nums)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i := range nums.Numbers {
		nums.Numbers[i] *= 3
	}
	newNums, err := json.Marshal(nums)
	if err != nil {
		log.Println("marshaling error")
		return
	}
	w.Write([]byte(newNums))
}

func v2Text(w http.ResponseWriter, r *http.Request) {
	var text StringArray
	// unmarshall json string to struct
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	changedStr := changeStrings.ReturnDuplicates(text.Text)
	newText, err := json.Marshal(changedStr)
	if err != nil {
		log.Println("marshaling error")
		return
	}
	w.Write([]byte(newText))
}

func v1AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("User-Auth")

		if auth != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("invalid auth")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func v2AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("User-Auth")

		if auth != "surprise!" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("invalid auth")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	port := ":3000"
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(http.StatusText(200)))
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Use(v1AuthMiddleware)
			r.Post("/numbers", v1Num)
			r.Post("/text", v1Text)
		})
		r.Route("/v2", func(r chi.Router) {
			r.Use(v2AuthMiddleware)
			r.Post("/numbers", v2Num)
			r.Post("/text", v2Text)
		})
	})
	fmt.Printf("Запуск сервера на %s порте\n", port)
	http.ListenAndServe(port, r)
}
