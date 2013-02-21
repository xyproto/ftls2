package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/hoisie/web"
	"github.com/xyproto/browserspeak"
)

type UserState struct {
	// see: http://redis.io/topics/data-types
	users      *RedisHashMap // "users:"username "loggedin" "true"/"false"
	usernames  *RedisSet     // "usernames" username,username,username
	connection redis.Conn
}

func InitUserSystem(connection redis.Conn) *UserState {

	state := new(UserState)
	state.users = NewRedisHashMap(connection, "users")
	state.usernames = NewRedisSet(connection, "usernames")
	state.connection = connection

	return state
}

func (state *UserState) HasUser(username string) bool {
	val, err := state.usernames.Has(username)
	if err != nil {
		return false
	}
	return val
}

// Create a user by adding the username to the list of usernames
func GenerateCreateUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Can't create blank user"
		}
		if state.HasUser(val) {
			return "user " + val + " already exists, could not create"
		}

		// Add he user
		state.usernames.Add(val)

		// Add additional data as well
		state.users.Set(val, "loggedin", "false")

		return "OK, user " + val + " created"
	}
}

// Create a user by adding the username to the list of usernames
func GenerateRemoveUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Can't remove blank user"
		}
		if !state.HasUser(val) {
			return "user " + val + " doesn't exists, could not remove"
		}

		// Remove the user
		state.usernames.Del(val)

		// Remove additional data as well
		state.users.Del(val, "loggedin")

		return "OK, user " + val + " removed"
	}
}

// Log in a user by changing the loggedin value
func GenerateLoginUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Can't log in a blank user"
		}
		if !state.HasUser(val) {
			return "user " + val + " does not exist, could not log in"
		}
		state.users.Set(val, "loggedin", "true")
		return "OK, user " + val + " logged in"
	}
}

// Log out a user by changing the loggedin value
func GenerateLogoutUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "No user logged out"
		}
		if !state.HasUser(val) {
			return "user " + val + " does not exist, could not log out"
		}
		state.users.Set(val, "loggedin", "false")
		return "OK, user " + val + " logged out"
	}
}

func GenerateGetAllUsernames(state *UserState) SimpleWebHandle {
	return func(val string) string {
		s := ""
		usernames, err := state.usernames.GetAll()
		if err == nil {
			for _, val := range usernames {
				s += "USERNAME: " + val + "<br />"
			}
		}
		return browserspeak.Message("Usernames", s)
	}
}

// Converts "true" or "false" to a bool
func truthValue(val string) bool {
	return "true" == val
}

func GenerateUserStatus(state *UserState) SimpleWebHandle {
	return func(val string) string {
		if !state.HasUser(val) {
			return val + " does not exist"
		}
		status, err := state.users.Get(val, "loggedin")
		if err != nil {
			return "No login/logout status for user " + val
		}
		if truthValue(status) {
			return val + " is logged in"
		}
		return val + " is not logged in"
	}
}

func GenerateSetCookie(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		ctx.SetSecureCookie("ost", "kake", 123)
		return "Cookiezyzz!"
	}
}

// TODO: RESTful services
func ServeUserSystem(connection redis.Conn) *UserState {
	state := InitUserSystem(connection)

	web.Get("/login/(.*)", GenerateLoginUser(state))
	web.Get("/logout/(.*)", GenerateLogoutUser(state))
	web.Get("/status/(.*)", GenerateUserStatus(state))
	web.Get("/create/(.*)", GenerateCreateUser(state))
	web.Get("/remove/(.*)", GenerateRemoveUser(state))
	web.Get("/users/(.*)", GenerateGetAllUsernames(state))
	web.Get("/cookie/(.*)", GenerateSetCookie(state))

	return state
}
