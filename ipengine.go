package main

// OK, only IP-specific stuff, 23-03-13

// TODO: Split out as it's own application, then add DNS functionality

import (
	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
	"github.com/xyproto/web"
	"github.com/xyproto/instapage"
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
		return instapage.Message("IPs", s)
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

func ServeIPs(userState *UserState) {
	ipState := InitIPs(userState.GetPool())

	web.Get("/setip/(.*)", GenerateSetIP(ipState))
	web.Get("/getip/(.*)", GenerateGetLastIP(ipState))
	web.Get("/getallips/(.*)", GenerateGetAllIPs(ipState))
}
