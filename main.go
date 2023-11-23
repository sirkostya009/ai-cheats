package main

import (
	"github.com/rs/cors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
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

func checkSalt(salt string) bool {
	return salt == ""
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
			log.Println("Successfully processed request", r.URL.Path[1:])

			event.CompletionTokens = response.Usage.CompletionTokens
			event.PromptTokens = response.Usage.PromptTokens
		} else {
			event.Status = (err.(*customError)).status
			event.Reason = &(err.(*customError)).message
			w.WriteHeader(event.Status)
			log.Println(err)
		}

		go AddEvent(event)
	}()

	if !checkSalt(r.Header.Get("X-Salt")) {
		err = &customError{"Failed to verify salt " + r.Header.Get("X-Salt"), http.StatusBadRequest}
		return
	}

	if customer == nil {
		err = &customError{"Failed to fetch customer " + r.URL.Path[1:], http.StatusNotFound}
		return
	}

	if !customer.Active {
		err = &customError{"Customer is not active", http.StatusForbidden}
		return
	}

	remoteHost, _, _ := net.SplitHostPort(r.RemoteAddr)
	addr, err := netip.ParseAddr(remoteHost)
	if err != nil {
		err = &customError{"Failed to parse IP address " + err.Error(), http.StatusBadRequest}
		return
	}

	if len(customer.Ips) < customer.MaxIps && !customer.ContainsIp(addr) {
		customer.Ips = append(customer.Ips, addr)
		go Update(customer)
	}

	if !customer.ContainsIp(addr) {
		err = &customError{"IP address is not allowed " + addr.String(), http.StatusForbidden}
		return
	}

	prompt, err := io.ReadAll(r.Body)
	response, err = CallAI(customer.Model, string(prompt))
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
