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
		return username + " is ready for chatting"
	}
}

func (ce *ChatEngine) GenerateCSS(cs *ColorScheme) SimpleContextHandle {
	return func(ctx *web.Context) string {
		ctx.ContentType("css")

		// TODO: Consider if hiding all the menus with CSS is the way to go or not
		//#menuChat {
		//	display: none;
		//}

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
