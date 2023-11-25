package main

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"os"
	"time"
)

type customError struct {
	message string
	status  int
}

func (e *customError) Error() string {
	return e.message
}

func extractHash(body string) (hash string) {
	for i := 1; i < len(body); i += 2 {
		if body[i] != '_' {
			hash += string(body[i])
		}
	}
	return
}

func clearBody(body string) (result string) {
	for i := 0; i < len(body); i += 2 {
		result += string(body[i])
	}
	return
}

func request(w http.ResponseWriter, r *http.Request) {
	var err error
	createdAt := time.Now()
	var response openai.ChatCompletionResponse

	customer := GetById(r.URL.Path[1:])

	defer func() {
		event := RequestEvent{
			CustomerId: customer.Id,
			CreatedAt:  createdAt,
			Status:     200,
			Model:      customer.Model,
		}

		if err == nil {
			fmt.Println("Successfully processed request", r.URL.Path[1:])

			event.CompletionTokens = response.Usage.CompletionTokens
			event.PromptTokens = response.Usage.PromptTokens
		} else {
			event.Status = (err.(*customError)).status
			event.Reason = &(err.(*customError)).message
			w.WriteHeader(event.Status)
			fmt.Println(err)
		}

		go AddEvent(event)
	}()

	if customer == nil {
		err = &customError{"Failed to fetch customer " + r.URL.Path[1:], http.StatusNotFound}
		return
	}

	if !customer.Active {
		err = &customError{"Customer is not active", http.StatusForbidden}
		return
	}

	raw, err := io.ReadAll(r.Body)
	body := string(raw)
	fmt.Println(body)
	hash := extractHash(body)
	if len(customer.Hashes) < customer.MaxIps && !customer.HasHash(hash) {
		customer.Hashes = append(customer.Hashes, hash)
		go Update(customer)
	}

	if !customer.HasHash(hash) {
		err = &customError{"Hash not allowed " + hash, http.StatusForbidden}
		return
	}

	interval, count := AverageRequestInterval(customer.Id)
	if interval < 1000 && count > 5 {
		err = &customError{
			fmt.Sprintf("Customer with id %v has been denied access to due to too many requests", customer.Id),
			http.StatusTooManyRequests,
		}
		return
	}

	response, err = CallAI(customer.Model, clearBody(body))
	if err != nil {
		err = &customError{"Failed to call AI " + err.Error(), http.StatusInternalServerError}
		return
	}

	_, err = w.Write([]byte(response.Choices[0].Message.Content))
	if err != nil {
		err = &customError{"Failed to write response " + err.Error(), http.StatusInternalServerError}
		return
	}
}

func main() {
	InitializeDB()
	InitializeAPI()
	defer CloseDB()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is not set")
	}

	http.HandleFunc("/", request)

	_ = http.ListenAndServe(":"+port, cors.New(cors.Options{AllowedOrigins: []string{"https://test.vntu.edu.ua"}}).Handler(http.DefaultServeMux))
}
