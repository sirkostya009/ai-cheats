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
		"update customers set ips = $1, requests = $2, request_tokens = $3, generated_tokens = $4 where id = $5",
		c.Ips, c.Requests, c.RequestTokens, c.GeneratedTokens, c.Id,
	)
	if err != nil {
		fmt.Println("Failed to save Customer", c.Id, err)
	}
}

func GetById(id string) (c *Customer) {
	c = mapCustomer(db.QueryRow(context.TODO(), "select * from customers where id = $1", id))
	return
}

func GetAll() (result []Customer) {
	rows, err := db.Query(context.TODO(), "select * from customers")
	if err != nil {
		fmt.Println("Failed to get all Customers", err)
		return nil
	}
	for rows.Next() {
		c := mapCustomer(rows)
		if c == nil {
			continue
		}
		result = append(result, *c)
	}
	return
}

func mapCustomer(row interface {
	Scan(dest ...any) error
}) (c *Customer) {
	c = &Customer{}
	err := row.Scan(&c.Id, &c.Telegram, &c.Active, &c.Ips, &c.MaxIps, &c.Model, &c.Requests, &c.RequestTokens, &c.GeneratedTokens)
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
