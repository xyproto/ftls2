package main

import (
	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
	"github.com/xyproto/web"
)

// An Engine is a specific piece of a website
// This part handles the "chat" pages

type ChatEngine struct {
	state *UserState
}

type ChatState struct {
	active *RedisSet       // A list of all users that are in the chat, must correspond to the users in UserState.users
	said   *RedisSet       // A list of everything that has been said so far
	pool   *ConnectionPool // A connection pool for Redis
}

func NewChatEngine(state *UserState) *ChatEngine {
	return &ChatEngine{state}
}

func (ce *ChatEngine) ServePages(basecp BaseCP, menuEntries MenuEntries) {
	state := ce.state

	chatCP := basecp(state)
	chatCP.ContentTitle = "Chat"
	chatCP.ExtraCSSurls = append(chatCP.ExtraCSSurls, "/css/chat.css")

	tpvf := DynamicMenuFactoryGenerator(menuEntries)

	web.Get("/chat", chatCP.WrapSimpleContextHandle(GenerateChatCurrentUser(state), tpvf(state)))
	web.Get("/css/chat.css", ce.GenerateCSS(chatCP.ColorScheme))
}

func GenerateChatCurrentUser(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return "No user logged in"
		}
		if !state.IsLoggedIn(username) {
			return "Not logged in"
		}

		// TODO: Use a list of users in the chat in the database
		//       add this user to that list and display the list.
		//       Also have a list for chat content and a form for saying stuff.
		//       Add a limit to how much can be said in the chat box.

		return username + " is ready for chatting"
	}
}

func (ce *ChatEngine) GenerateCSS(cs *ColorScheme) SimpleContextHandle {
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
