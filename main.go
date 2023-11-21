package main

import (
	"github.com/rs/cors"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
)

func checkSalt(salt string) bool {
	return salt == ""
}

func request(w http.ResponseWriter, r *http.Request) {
	customer := GetById(r.URL.Path[1:])

	if !checkSalt(r.Header.Get("X-Salt")) {
		log.Println("Failed to verify salt", r.Header.Get("X-Salt"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if customer == nil {
		log.Println("Failed to fetch customer", r.URL.Path[1:])
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !customer.Active {
		log.Println("Customer deactivated", customer.Id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	remoteHost, _, _ := net.SplitHostPort(r.RemoteAddr)
	addr, err := netip.ParseAddr(remoteHost)
	if err != nil {
		log.Println("Failed to parse remote address", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(customer.Ips) < customer.MaxIps && !customer.IpContains(addr) {
		customer.Ips = append(customer.Ips, addr)
	}

	if !customer.IpContains(addr) {
		log.Println("IP not in customer's IP list", addr, customer.Ips)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prompt, err := io.ReadAll(r.Body)
	response, err := CallAI(customer.Model, string(prompt))
	if err != nil {
		log.Println("Failed to call OpenAI API", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	customer.Requests++
	customer.GeneratedTokens += response.Usage.CompletionTokens
	customer.RequestTokens += response.Usage.PromptTokens

	go Update(customer)

	_, err = w.Write([]byte(response.Choices[0].Message.Content))
	if err != nil {
		log.Println("Failed to write response", err)
		w.WriteHeader(http.StatusInternalServerError)
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
