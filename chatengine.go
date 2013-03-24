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
	url string
}

func NewChatEngine(state *UserState, url string) *ChatEngine {
	return &ChatEngine{state, url}
}

func (ce *ChatEngine) ShowMenu(url string, ctx *web.Context) bool {
	if url == ce.url {
		return false
	}
	if ce.state.UserRights(ctx) {
		return true
	}
	return false
}

// For other webpages, used internally
func (ce *ChatEngine) ServeSystem() {
}

func (ce *ChatEngine) ServePages(basecp BaseCP) {
	state := ce.state

	chatCP := basecp(state)
	chatCP.ContentTitle = "Chat"
	chatCP.ExtraCSSurls = append(chatCP.ExtraCSSurls, "/css/chat.css")

	// Hide the Chat menu if we're on the Chat page
	chatCP.HiddenMenuIDs = append(chatCP.HiddenMenuIDs, "menuChat")

	tp := Kake()
	web.Get("/chat", chatCP.WrapSimpleContextHandle(GenerateChatCurrentUser(state), tp))
	web.Get("/css/chat.css", GenerateChatCSS(chatCP.ColorScheme))
}

func GenerateChatCurrentUser(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return MessageOKback("Chat", "No user logged in")
		}
		if !state.IsLoggedIn(username) {
			return MessageOKback("Chat", "Not logged in")
		}
		return username + " is ready for chatting"
	}
}

func GenerateChatCSS(cs *ColorScheme) SimpleContextHandle {
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
