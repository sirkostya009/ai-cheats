package main

import "net/netip"

type Customer struct {
	Id              int          `json:"id,omitempty"`
	Telegram        string       `json:"telegram,omitempty"`
	Active          bool         `json:"active,omitempty"`
	Ips             []netip.Addr `json:"ips,omitempty"`
	MaxIps          int          `json:"maxIps,omitempty"`
	Model           string       `json:"model,omitempty"`
	Requests        int          `json:"requests,omitempty"`
	RequestTokens   int          `json:"requestTokens,omitempty"`
	GeneratedTokens int          `json:"generatedTokens,omitempty"`
}

func (c *Customer) IpContains(ip netip.Addr) bool {
	for _, _ip := range c.Ips {
		if _ip == ip {
			return true
		}
	}

	return false
}
