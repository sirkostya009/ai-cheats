package main

type Customer struct {
	Id       int
	Telegram string
	Active   bool
	Hashes   []string
	MaxIps   int
	Model    string
}

func (c *Customer) HasHash(hash string) bool {
	for _, _ip := range c.Hashes {
		if _ip == hash {
			return true
		}
	}

	return false
}
