package main

import (
	"errors"
	"math/rand"
	"time"

	. "github.com/xyproto/browserspeak"
	"github.com/garyburd/redigo/redis"
	"github.com/xyproto/web"
)

const (
	ONLY_LOGIN = "100"
	ONLY_LOGOUT = "010"
	ONLY_REGISTER = "001"
	EXCEPT_LOGIN = "011"
	EXCEPT_LOGOUT = "101"
	EXCEPT_REGISTER = "110"
	NOTHING = "000"
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

// TODO: Rethink this. Use templates for Login/Logout button?
// Generate "1" or "0" values for showing the login, logout or register menus,
// depending on the cookie status and UserState
func GenerateShowLoginLogoutRegister(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if username := GetBrowserUsername(ctx); username != "" {
			//print("USERNAME", username)
			// Has a username stored in the browser
			if state.LoggedIn(username) {
				// Ok, logged in to the system + login cookie in the browser
				// Only present the "Logout" menu
				return ONLY_LOGOUT
			} else {
				// Has a login cookie, but is not logged in.
				// Keep the browser cookie (could be tempting to remove it)
				// Present only the "Login" menu
				//return "100"
				// Present both "Login" and "Register", just in case it's a new user
				// in the same browser.
				return EXCEPT_LOGOUT
			}
		} else {
			// Does not have a username stored in the browser
			// Present the "Register" and "Login" menu
			return EXCEPT_LOGOUT
		}
		// Everything went wrong, should never reach this point
		return NOTHING
	}
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
