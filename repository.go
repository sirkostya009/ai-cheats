package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var db *pgxpool.Pool

func Update(c *Customer) {
	_, err := db.Exec(
		context.TODO(),
		"update customers set ips = $1 where id = $5",
		c.Ips, c.Id,
	)
	if err != nil {
		fmt.Println("Failed to save Customer", c.Id, err)
	}
}

func GetById(id string) (c *Customer) {
	c = mapCustomer(db.QueryRow(context.TODO(), "select * from customers where id = $1", id))
	return
}

func AddEvent(e RequestEvent) {
	_, err := db.Exec(
		context.TODO(),
		"insert into requests (customer_id, created_at, completion_tokens, prompt_tokens, status, reason) values ($1, $2, $3, $4, $5, $6)",
		e.CustomerId, e.CreatedAt, e.CompletionTokens, e.PromptTokens, e.Status, e.Reason,
	)
	if err != nil {
		fmt.Println("Failed to save RequestEvent", e, err)
	}
}

func mapCustomer(row interface {
	Scan(dest ...any) error
}) (c *Customer) {
	c = &Customer{}
	err := row.Scan(&c.Id, &c.Telegram, &c.Active, &c.Ips, &c.MaxIps, &c.Model)
	if err != nil {
		fmt.Println("Failed to map Customer ", err)
		return nil
	}
	return
}

func InitializeDB() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("DATABASE_URL is not set")
	}

	var err error
	db, err = pgxpool.New(context.TODO(), url)
	if err != nil {
		panic(err)
	}
}

var CloseDB = db.Close
