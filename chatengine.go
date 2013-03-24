package main

import (
	"github.com/xyproto/web"
	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
)

// An Engine is a specific piece of a website
// This part handles the "chat" pages

type ChatEngine Engine

func NewChatEngine(state *UserState) *ChatEngine {
	var engine Engine
	engine.SetState(state)
	chatengine := ChatEngine(engine)
	return &chatengine
}

func ChatMenuJS() string {
	// This in combination with hiding the link in genericsite.go is cool, but the layout becomes weird :/
	return ShowInlineAnimatedIf("/showmenu/chat", "#menuChat")

	// This keeps the layout but is less cool
	//return HideIfNot("/showmenu/chat", "#menuChat")
}

func GenerateShowChat(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if state.UserRights(ctx) {
			return "1"
		}
		return "0"
	}
}

func (ce *ChatEngine) ServePages(basecp BaseCP) {
	engine := Engine(*ce)
	state := engine.GetState()
	chatCP := basecp(state)
	chatCP.ContentTitle = "Chat"
	chatCP.ExtraCSSurls = append(chatCP.ExtraCSSurls, "/css/chat.css")

	// Hide the Chat menu if we're on the Chat page
	chatCP.ContentJS = Hide("#menuChat")

	tp := Kake()
	web.Get("/chat", chatCP.WrapSimpleContextHandle(GenerateChatCurrentUser(state), tp))
	web.Get("/css/chat.css", GenerateChatCSS(chatCP.ColorScheme))
	web.Get("/showmenu/chat", GenerateShowChat(state))
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


