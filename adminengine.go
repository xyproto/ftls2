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

const (
	ADMIN = "1"
	USER  = "0"
)

// Checks if the current user is logged in as administrator right now
func (state *UserState) AdminNow(ctx *web.Context) bool {
	if username := GetBrowserUsername(ctx); username != "" {
		return state.IsLoggedIn(username) && state.IsAdministrator(username)
	}
	return false
}

func GenerateShowAdmin(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if state.AdminNow(ctx) {
			return ADMIN
		}
		return USER
	}
}

func GenerateAdminCSS(cs *ColorScheme) SimpleContextHandle {
	return func(ctx *web.Context) string {
		ctx.ContentType("css")
		return `
.yes {
	background-color: #90ff90;
	color: black;
}
.no {
	background-color: #ff9090;
	color: black;
}
table {
	border-collapse: collapse;
	padding: 1em;
}
table, th, tr, td {
	border: 1px solid black;
	padding: 1em;
}
.username:link { color: green; }
.username:visited { color: green; }
.username:hover { color: green; }
.username:active { color: green; }
.whitebg {
	background-color: white;
}
.darkgrey:link { color: #404040; }
.darkgrey:visited { color: #404040; }
.darkgrey:hover { color: #404040; }
.darkgrey:active { color: #404040; }
`
		//
	}
}

func ServeAdminPages(state *UserState, cps PageCollection, cs *ColorScheme, tp map[string]string) {
	adminCP := BaseTitleCP("Admin", state)
	adminCP.extraCSSurls = append(adminCP.extraCSSurls, "/css/admin.css")

	// Hide the Admin menu if we're on the Admin page
	adminCP.contentJS = Hide("#menuAdmin")

	web.Get("/admin", adminCP.WrapSimpleContextHandle(GenerateAdminStatus(state), tp))
	web.Get("/css/admin.css", GenerateAdminCSS(cs))

	web.Get("/showmenu/admin", GenerateShowAdmin(state))
}

// TODO: Log and graph when people visit pages and when people contribute content
// This one is wrapped by ServeAdminPages
func GenerateAdminStatus(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if !state.AdminNow(ctx) {
			return "<div class=\"no\">Not administrator</div>"
		}

		// TODO: List all sorts of info, edit users, etc
		s := "<h2>Welcome chief</h2>"

		s += "<strong>User table</strong><br />"
		s += "<table class=\"whitebg\">"
		s += "<tr>"
		s += "<th>Username</th><th>Confirmed</th><th>Logged in</th><th>Administrator</th><th>Admin toggle</th><th>Remove user</th>"
		s += "</tr>"
		usernames, err := state.usernames.GetAll()
		if err == nil {
			for _, username := range usernames {
				s += "<tr>"
				s += "<td>" + "<a class=\"username\" href=\"/status/" + username + "\">" + username + "</a></td>"
				s += bool2td(state.IsConfirmed(username))
				s += bool2td(state.IsLoggedIn(username))
				s += bool2td(state.IsAdministrator(username))
				s += "<td>" + "<a class=\"darkgrey\" href=\"/admintoggle/" + username + "\">admin toggle</a></td>"
				// TODO: Ask for confirmation first with a MessageOKurl("blabla", "blabla", "/actually/remove/stuff")
				s += "<td>" + "<a class=\"darkgrey\" href=\"/remove/" + username + "\">remove</a></td>"
				s += "</tr>"
			}
		}
		s += "</table>"
		s += "<br />"
		s += "<strong>Unconfirmed users</strong><br />"
		s += "<table>"
		s += "<tr>"
		s += "<th>Username</th><th>Confirmation link</th>"
		s += "</tr>"
		usernames, err = state.unconfirmed.GetAll()
		if err == nil {
			for _, username := range usernames {
				s += "<tr>"
				s += "<td>" + "<a class=\"username\" href=\"/status/" + username + "\">" + username + "</a></td>"
				s += "<td>" + state.GetConfirmationSecret(username) + "</td>"
				s += "<td>" + "<a class=\"username\" href=\"/removeunconfirmed/" + username + "\">remove</a></td>"
				s += "</tr>"
			}
		}
		s += "</table>"
		return s
	}
}

// Checks if the given username is an administrator
func (state *UserState) IsAdministrator(username string) bool {
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
		if !state.AdminNow(ctx) {
			return MessageOKback("Status", "Not administrator")
		}
		username := GetBrowserUsername(ctx)
		if username == "" {
			return MessageOKback("Current user status", "No user logged in")
		}
		hasUser := state.HasUser(username)
		if !hasUser {
			return MessageOKback("Current user status", username+" does not exist")
		}
		if !(state.IsLoggedIn(username)) {
			return MessageOKback("Current user status", "User "+username+" is not logged in")
		}
		return MessageOKback("Current user status", "User "+username+" is logged in")
	}
}

func GenerateStatusUser(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
		if username == "" {
			return MessageOKback("Status", "No username given")
		}
		if !state.HasUser(username) {
			return MessageOKback("Status", username+" does not exist")
		}
		loggedinStatus := "not logged in"
		if state.IsLoggedIn(username) {
			loggedinStatus = "logged in"
		}
		confirmStatus := "email has not been confirmed"
		if state.IsConfirmed(username) {
			confirmStatus = "email has been confirmed"
		}
		return MessageOKback("Status", username+" is "+loggedinStatus+" and "+confirmStatus)
	}
}

// Remove an unconfirmed user
func GenerateRemoveUnconfirmedUser(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
		if !state.AdminNow(ctx) {
			return MessageOKback("Remove unconfirmed user", "Not administrator")
		}

		if username == "" {
			return MessageOKback("Remove unconfirmed user", "Can't remove blank user.")
		}

		found := false
		usernames, err := state.unconfirmed.GetAll()
		if err == nil {
			for _, unconfirmedUsername := range usernames {
				if username == unconfirmedUsername {
					found = true
					break
				}
			}
		}

		if !found {
			return MessageOKback("Remove unconfirmed user", "Can't find "+username+" in the list of unconfirmed users.")
		}

		// Remove the user
		state.unconfirmed.Del(username)

		// Remove additional data as well
		state.users.Del(username, "secret")

		return MessageOKurl("Remove unconfirmed user", "OK, removed "+username+" from the list of unconfirmed users.", "/admin")
	}
}

// TODO: Add possibility for Admin to restart the webserver

// TODO: Undo for removing users
// Remove a user
func GenerateRemoveUser(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
		if !state.AdminNow(ctx) {
			return MessageOKback("Remove user", "Not administrator")
		}

		if username == "" {
			return MessageOKback("Remove user", "Can't remove blank user")
		}
		if !state.HasUser(username) {
			return MessageOKback("Remove user", username+" doesn't exists, could not remove")
		}

		// Remove the user
		state.usernames.Del(username)

		// Remove additional data as well
		state.users.Del(username, "loggedin")

		return MessageOKurl("Remove user", "OK, removed "+username, "/admin")
	}
}

func GenerateAllUsernames(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if !state.AdminNow(ctx) {
			return MessageOKback("List usernames", "Not administrator")
		}
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
		if !state.AdminNow(ctx) {
			return MessageOKback("Get cookie", "Not administrator")
		}
		username := GetBrowserUsername(ctx)
		return MessageOKback("Get cookie", "Cookie: username = "+username)
	}
}

func GenerateSetCookie(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
		if !state.AdminNow(ctx) {
			return MessageOKback("Set cookie", "Not administrator")
		}
		if username == "" {
			return MessageOKback("Set cookie", "Can't set cookie for empty username")
		}
		if !state.HasUser(username) {
			return MessageOKback("Set cookie", "Can't store cookie for non-existsing user")
		}
		// Create a cookie that lasts for one hour,
		// this is the equivivalent of a session for a given username
		ctx.SetSecureCookiePath("user", username, 3600, "/")
		return MessageOKback("Set cookie", "Cookie stored: user = "+username+".")
	}
}

func GenerateToggleAdmin(state *UserState) WebHandle {
	return func(ctx *web.Context, username string) string {
		if !state.AdminNow(ctx) {
			return MessageOKback("Admin toggle", "Not administrator")
		}
		if username == "" {
			return MessageOKback("Admin toggle", "Can't set toggle empty username")
		}
		if !state.HasUser(username) {
			return MessageOKback("Admin toggle", "Can't toggle non-existing user")
		}
		// A special case
		if username == "admin" {
			return MessageOKback("Admin toggle", "Can't remove admin rights from the admin user")
		}
		if !state.IsAdministrator(username) {
			state.users.Set(username, "admin", "true")
			return MessageOKurl("Admin toggle", "OK, "+username+" is now an admin", "/admin")
		}
		state.users.Set(username, "admin", "false")
		return MessageOKurl("Admin toggle", "OK, "+username+" is now a regular user", "/admin")

	}
}

// Needed to fullfill the Engine interface, serves the pages
func (ae *AdminEngine) ServeSystem() {
	state := ae.state

	// These are available for everyone
	web.Get("/status/(.*)", GenerateStatusUser(state))

	// These are only available as administrator, all have checks
	web.Get("/status", GenerateStatusCurrentUser(state))
	web.Get("/remove/(.*)", GenerateRemoveUser(state))
	web.Get("/removeunconfirmed/(.*)", GenerateRemoveUnconfirmedUser(state))
	web.Get("/users/(.*)", GenerateAllUsernames(state))
	web.Get("/admintoggle/(.*)", GenerateToggleAdmin(state))
	//web.Get("/cookie/get", GenerateGetCookie(state))
	//web.Get("/cookie/set/(.*)", GenerateSetCookie(state))
}
