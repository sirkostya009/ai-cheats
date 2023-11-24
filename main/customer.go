package main

import "net/netip"

type Customer struct {
	Id       int
	Telegram string
	Active   bool
	Ips      []netip.Addr
	MaxIps   int
	Model    string
}

func (c *Customer) ContainsIp(ip netip.Addr) bool {
	for _, _ip := range c.Ips {
		if _ip == ip {
			return true
		}
	}

	return false
}
