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
		"update customers set ips = $1 where id = $2",
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
		"insert into requests (customer_id, created_at, completion_tokens, prompt_tokens, status, reason, model) values ($1, $2, $3, $4, $5, $6, $7)",
		e.CustomerId, e.CreatedAt, e.CompletionTokens, e.PromptTokens, e.Status, e.Reason, e.Model,
	)
	if err != nil {
		fmt.Println("Failed to save RequestEvent", e, err)
	}
}

func AverageRequestInterval(customerId int) (interval float64, count int) {
	err := db.QueryRow(
		context.TODO(),
		`
		select avg(time_diff), count(time_diff)
		from (
			select extract(epoch from created_at - lag(created_at) over (order by created_at)) as time_diff
			from requests
			where customer_id = 1
			  and status = 200
			  and created_at > NOW() - interval '1 hour'
		) time_diffs`,
		customerId,
	).Scan(&interval, &count)
	if err != nil {
		fmt.Println("Failed to calculate average request interval", err)
	}
	return
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
