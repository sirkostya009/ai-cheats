package main

type Customer struct {
	Id        int
	Telegram  string
	Active    bool
	Hashes    []string
	MaxHashes int
	Model     string
}

func (c *Customer) HasHash(hash string) bool {
	for _, _hash := range c.Hashes {
		if _hash == hash {
			return true
		}
	}

	return false
}
