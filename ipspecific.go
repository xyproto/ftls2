package main

import (
	"github.com/hoisie/web"
	"github.com/xyproto/browserspeak"
)

type State struct {
	ips *RedisList
}

func setupIPs() *State {

	// Connect to Redis
	client, err := NewRedisConnection()
	if err != nil {
		panic("ERROR: Can't connect to redis")
	}

	// Create a RedisList for storing IP adresses
	ips := NewRedisList(client, "IPs")

	state := new(State)
	state.ips = ips

	return state
}

// Set an IP adress and generate a confirmation page for it
func GenerateSetIP(state *State) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Empty value, IP not set"
		}
		state.ips.Store(val)
		return "OK, set IP to " + val
	}
}

// Get all the stored IP adresses and generate a page for it
func GenerateGetAllIPs(state *State) SimpleWebHandle {
	return func(val string) string {
		s := ""
		iplist, err := state.ips.GetAll()
		if err == nil {
			for _, val := range iplist {
				s += "IP: " + val + "<br />"
			}
		}
		return browserspeak.Message("IPs", s)
	}
}

// Get the last stored IP adress and generate a page for it
func GenerateGetLastIP(state *State) SimpleWebHandle {
	return func(val string) string {
		s := ""
		ip, err := state.ips.GetLast()
		if err == nil {
			s = "IP: " + ip
		}
		return s
	}
}

// Close the connection in the RedisList in the state
// TODO: Put the connection in the State instead
func (state *State) Close() {
	state.ips.c.Close()
}

// TODO: RESTful services
func ServeIPs() *State {
	state := setupIPs()

	web.Get("/setip/(.*)", GenerateSetIP(state))
	web.Get("/getip/(.*)", GenerateGetLastIP(state))
	web.Get("/getallips/(.*)", GenerateGetAllIPs(state))

	return state
}
