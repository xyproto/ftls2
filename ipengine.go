package main

// OK, only IP-specific stuff, 23-03-13

// TODO: Split out as it's own application, with a DNS server too
//       (or configure a DNS server)

import (
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
	// "github.com/xyproto/simpleredis"
)

type IPState struct {
	data *RedisList
	pool *ConnectionPool
}

func InitIPs(pool *ConnectionPool) *IPState {

	// Create a RedisList for storing IP adresses
	ips := NewRedisList(pool, "IPs")

	state := new(IPState)
	state.data = ips
	state.pool = pool

	return state
}

// Set an IP adress and generate a confirmation page for it
func GenerateSetIP(state *IPState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Empty value, IP not set"
		}
		state.data.Store(val)
		return "OK, set IP to " + val
	}
}

// Get all the stored IP adresses and generate a page for it
func GenerateGetAllIPs(state *IPState) SimpleWebHandle {
	return func(val string) string {
		s := ""
		iplist, err := state.data.GetAll()
		if err == nil {
			for _, val := range iplist {
				s += "IP: " + val + "<br />"
			}
		}
		return Message("IPs", s)
	}
}

// Get the last stored IP adress and generate a page for it
func GenerateGetLastIP(state *IPState) SimpleWebHandle {
	return func(val string) string {
		s := ""
		ip, err := state.data.GetLast()
		if err == nil {
			s = "IP: " + ip
		}
		return s
	}
}

func ServeIPs(pool *ConnectionPool) *IPState {
	state := InitIPs(pool)

	web.Get("/setip/(.*)", GenerateSetIP(state))
	web.Get("/getip/(.*)", GenerateGetLastIP(state))
	web.Get("/getallips/(.*)", GenerateGetAllIPs(state))

	return state
}
