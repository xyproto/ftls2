package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

type IPState struct {
	data       *RedisList
	pool *redis.Pool
}

func InitIPs(pool *redis.Pool) *IPState {

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

// TODO: RESTful services
func ServeIPs(pool *redis.Pool) *IPState {
	state := InitIPs(pool)

	web.Get("/setip/(.*)", GenerateSetIP(state))
	web.Get("/getip/(.*)", GenerateGetLastIP(state))
	web.Get("/getallips/(.*)", GenerateGetAllIPs(state))

	return state
}
