package generatorclient

import (
	"context"
	"encoding/json"
	"errors"
	"go/transaction/client/models"
	"log"
	"net/http"
	"time"
)

type errorJsonResponse struct {
	Text string `json:"-_-"`
}

func GetResp(ctx context.Context) ([]models.Transaction, bool, error) {
	client := http.Client{
		Timeout: 6 * time.Second,
	}
	url := "http://127.0.0.1:8080/generate/transactions?count=2"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	xUserHeader := ctx.Value("X-User").(string)
	req.Header.Set("X-User", xUserHeader)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request error", err)
		return nil, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		j := []models.Transaction{}
		err = json.NewDecoder(resp.Body).Decode(&j)
		if err != nil {
			log.Println("decode error", err)
			return nil, false, err
		}
		return j, false, nil
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		errJs := errorJsonResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errJs)
		if err != nil {
			log.Println("decode errorJsonResponse error", err)
			return nil, false, err
		}
		return nil, false, errors.New(errJs.Text)
	}
	return nil, true, nil
}

func GetTransaction(ctx context.Context) ([]models.Transaction, error) {
	for {
		transactions, is503, err := GetResp(ctx)
		if err != nil {
			return nil, err
		}
		if !is503 {
			return transactions, nil
		}
	}
}
