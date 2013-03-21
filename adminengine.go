package main

import (
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

// An Engine is a specific piece of a website
// This part handles the "admin" pages

type AdminEngine Engine

func NewAdminEngine(state *UserState) *AdminEngine {
	return &AdminEngine{state}
}

// The "/admin/x" page
func GenerateHi(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val != "" {
			return "Hi " + val + "!"
		}
		return "HI"
	}
}

// Checks if the given username is an administrator
func (state *UserState) Administrator(username string) bool {
	if !state.HasUser(username) {
		return false
	}
	status, err := state.users.Get(username, "admin")
	if err != nil {
		return false
	}
	return TruthValue(status)
}

func GenerateStatusCurrentUser(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return MessageOKback("Current user status", "No user logged in")
		}
		if !state.HasUser(username) {
			return MessageOKback("Current user status", username + " does not exist")
		}
		if !(state.LoggedIn(username)) {
			return MessageOKback("Current user status", "User " + username + " is not logged in")
		}
		return MessageOKback("Current user status", "User " + username + " is logged in")
	}
}

func GenerateStatusUser(state *UserState) SimpleWebHandle {
	return func(username string) string {
		if username == "" {
			return MessageOKback("Status", "No username given")
		}
		if !state.HasUser(username) {
			return MessageOKback("Status", username + " does not exist")
		}
		loggedinStatus := "not logged in"
		if IsLoggedIn(state, username) {
			loggedinStatus = "logged in"
		}
		confirmStatus := "email has not been confirmed"
		if IsConfirmed(state, username) {
			confirmStatus = "email has been confirmed"
		}
		return MessageOKback("Status", username + " is " + loggedinStatus + " and " + confirmStatus)
	}
}

// Create a user by adding the username to the list of usernames
func GenerateRemoveUser(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
		if username == "" {
			return "Can't remove blank user"
		}
		if !state.HasUser(username) {
			return username + " doesn't exists, could not remove"
		}

		// Remove the user
		state.usernames.Del(username)

		// Remove additional data as well
		state.users.Del(username, "loggedin")

		return "OK, " + username + " removed"
	}
}

func GenerateAllUsernames(state *UserState) SimpleWebHandle {
	return func(_ string) string {
		s := ""
		usernames, err := state.usernames.GetAll()
		if err == nil {
			for _, username := range usernames {
				s += username + "<br />"
			}
		}
		return MessageOKback("Usernames", s)
	}
}

func GenerateGetCookie(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		return "Cookie: username = " + username
	}
}

func GenerateSetCookie(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
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

// Needed to fullfill the Engine interface, serves the pages
func (ae *AdminEngine) ServeSystem() {
	state := ae.state

	web.Get("/admin/(.*)", GenerateHi(state))

	// TODO: admin pages should only be accessible as administrator

	web.Get("/status", GenerateStatusCurrentUser(state))
	web.Get("/status/(.*)", GenerateStatusUser(state))
	web.Get("/remove/(.*)", GenerateRemoveUser(state))
	web.Get("/users/(.*)", GenerateAllUsernames(state))
	web.Get("/cookie/get", GenerateGetCookie(state))
	web.Get("/cookie/set/(.*)", GenerateSetCookie(state))
}

