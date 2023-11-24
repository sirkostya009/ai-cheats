package backend

import "time"

type RequestEvent struct {
	CustomerId       int
	CreatedAt        time.Time
	CompletionTokens int
	PromptTokens     int
	Status           int
	Reason           *string
	Model            string
}
