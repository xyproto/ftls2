package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

type UserState struct {
	// see: http://redis.io/topics/data-types
	users      *RedisHashMap // "users:"username "loggedin" "true"/"false"
	usernames  *RedisSet     // "usernames" username,username,username
	connection redis.Conn
}

func InitUserSystem(connection redis.Conn) *UserState {

	// For the secure cookies
	rand.Seed(time.Now().UnixNano())
	// TODO: Move this somewhere else?
	web.Config.CookieSecret = "3a19QRmwf7mHZ9CPAaPQ0hsWezfKz"

	// For the database
	state := new(UserState)
	state.users = NewRedisHashMap(connection, "users")
	state.usernames = NewRedisSet(connection, "usernames")
	state.connection = connection

	return state
}

// TODO: Don't return false if there is an error, the user may exist
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
		username := val
		if username == "" {
			return "Can't log in a blank user"
		}
		if !state.HasUser(username) {
			return "user " + username + " does not exist, could not log in"
		}
		state.users.Set(username, "loggedin", "true")
		state.SetBrowserUsername(ctx, username)
		return "OK, user " + username + " logged in"
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
		return Message("Usernames", s)
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

// Checks if the given username is logged in or not
func (state *UserState) LoggedIn(username string) bool {
	if !state.HasUser(username) {
		return false
	}
	status, err := state.users.Get(username, "loggedin")
	if err != nil {
		return false
	}
	return truthValue(status)
}

func GenerateGetCookie(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		//username, _ := ctx.GetSecureCookie("user")
		return "Cookie: username = " + username // + " err: " + fmt.Sprintf("%v", exists) + " val: " + val
	}
}

// Gets the username that is stored in a cookie in the browser, if available
func GetBrowserUsername(ctx *web.Context) string {
	username, _ := ctx.GetSecureCookie("user")
	return username
}

func (state *UserState) SetBrowserUsername(ctx *web.Context, username string) error {
	if username == "" {
		return errors.New("Can't set cookie for empty username")
	}
	if !state.HasUser(username) {
		return errors.New("Can't store cookie for non-existsing user")
	}
	// Create a cookie that lasts for one hour,
	// this is the equivivalent of a session for a given username
	ctx.SetSecureCookiePath("user", username, 3600, "/")
	//"Cookie stored: user = " + username + "."
	return nil
}

// NB! Set the cookie at / for it to work in the paths underneath!
func GenerateSetCookie(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		username := val
		if username == "" {
			return "Can't set cookie for empty username"
		}
		if !state.HasUser(username) {
			return "Can't store cookie for non-existsing user"
		}
		// Create a cookie that lasts for one hour,
		// this is the equivivalent of a session for a given username
		ctx.SetSecureCookiePath("user", username, 3600, "/")
		return "Cookie stored: user = " + username + "."
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
	web.Get("/cookie/get", GenerateGetCookie(state))
	web.Get("/cookie/set/(.*)", GenerateSetCookie(state))

	return state
}
