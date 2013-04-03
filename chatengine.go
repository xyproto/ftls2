package main

import (
	"time"

	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
	"github.com/xyproto/instapage"
	"github.com/xyproto/web"
)

// An Engine is a specific piece of a website
// This part handles the "chat" pages

type ChatEngine struct {
	userState *UserState
	chatState *ChatState
}

type ChatState struct {
	active   *RedisSet       // A list of all users that are in the chat, must correspond to the users in UserState.users
	said     *RedisList      // A list of everything that has been said so far
	lastSeen *RedisHashMap   // A list of everything that has been said so far
	pool     *ConnectionPool // A connection pool for Redis
}

func NewChatEngine(userState *UserState) *ChatEngine {
	pool := userState.GetPool()
	chatState := new(ChatState)
	chatState.active = NewRedisSet(pool, "active")
	chatState.said = NewRedisList(pool, "said")
	chatState.lastSeen = NewRedisHashMap(pool, "lastSeen") // lastSeen.time is an encoded timestamp for when the user was last seen chatting
	chatState.pool = pool
	return &ChatEngine{userState, chatState}
}

func (ce *ChatEngine) ServePages(basecp BaseCP, menuEntries MenuEntries) {
	chatCP := basecp(ce.userState)
	chatCP.ContentTitle = "Chat"
	chatCP.ExtraCSSurls = append(chatCP.ExtraCSSurls, "/css/chat.css")

	tvgf := DynamicMenuFactoryGenerator(menuEntries)
	tvg := tvgf(ce.userState)

	web.Get("/chat", chatCP.WrapSimpleContextHandle(ce.GenerateChatCurrentUser(), tvg))
	web.Post("/say", chatCP.WrapSimpleContextHandle(ce.GenerateSayCurrentUser(), tvg))
	web.Get("/css/chat.css", ce.GenerateCSS(chatCP.ColorScheme))
}

// Mark a user as seen
func (ce *ChatEngine) Seen(username string) {
	now := time.Now()
	encodedTime, err := now.GobEncode()
	if err != nil {
		panic("ERROR: Can't encode the time")
	}
	ce.chatState.lastSeen.Set(username, "time", string(encodedTime))
}

func (ce *ChatEngine) GetLastSeen(username string) string {
	encodedTime, err := ce.chatState.lastSeen.Get(username, "time")
	if err == nil {
		var then time.Time
		err = then.GobDecode([]byte(encodedTime))
		if err == nil {
			timestamp := then.String()
			return timestamp[11:19]
		}
	}
	return "never"
}

func (ce *ChatEngine) IsChatting(username string) bool {
	encodedTime, err := ce.chatState.lastSeen.Get(username, "time")
	if err == nil {
		var then time.Time
		err = then.GobDecode([]byte(encodedTime))
		if err == nil {
			elapsed := time.Since(then)
			if elapsed.Minutes() > 20 {
				// 20 minutes since last seen saying anything, set as not chatting
				ce.SetChatting(username, false)
				return false
			}
		}
	}
	// TODO: If the user was last seen more than N minutes ago, set as not chatting and return false
	return ce.userState.GetBooleanField(username, "chatting")
}

// Set "chatting" to "true" or "false" for a given user
func (ce *ChatEngine) SetChatting(username string, val bool) {
	ce.userState.SetBooleanField(username, "chatting", val)
}

func (ce *ChatEngine) JoinChat(username string) {
	// Join the chat
	ce.chatState.active.Add(username)
	// Change the chat status for the user
	ce.SetChatting(username, true)
	// Mark the user as seen
	ce.Seen(username)
}

func (ce *ChatEngine) Say(username, text string) {
	timestamp := time.Now().String()
	textline := timestamp[11:19] + "&nbsp;&nbsp;" + username + "> " + text
	ce.chatState.said.Add(textline)
	// Store the timestamp for when the user was last seen as well
	ce.Seen(username)
}

func LeaveChat(ce *ChatEngine, username string) {
	// Leave the chat
	ce.chatState.active.Del(username)
	// Change the chat status for the user
	ce.SetChatting(username, false)
}

func (ce *ChatEngine) GetChatUsers() []string {
	chatUsernames, err := ce.chatState.active.GetAll()
	if err != nil {
		return []string{}
	}
	return chatUsernames
}

func (ce *ChatEngine) GetChatText() []string {
	chatText, err := ce.chatState.said.GetAll()
	if err != nil {
		return []string{}
	}
	return chatText
}

// Get the last N entries
func (ce *ChatEngine) GetLastChatText(n int) []string {
	chatText, err := ce.chatState.said.GetLastN(n)
	if err != nil {
		return []string{}
	}
	return chatText
}

func (ce *ChatEngine) GenerateChatCurrentUser() SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return "No user logged in"
		}
		if !ce.userState.IsLoggedIn(username) {
			return "Not logged in"
		}

		ce.JoinChat(username)

		// TODO: Add a button for someone to see the entire chat
		// TODO: Add some protection against random monkeys that only fling poo

		retval := "Hi " + username + "<br />"
		retval += "<br />"
		retval += "Participants:" + "<br />"
		for _, otherUser := range ce.GetChatUsers() {
			if otherUser == username {
				continue
			}
			retval += "&nbsp;&nbsp;" + otherUser + ", last seen " + ce.GetLastSeen(otherUser) + "<br />"
		}
		retval += "<br />"
		retval += "<div style='background-color: white; margin: 1em;'>"
		for _, said := range ce.GetLastChatText(20) {
			retval += said + "<br />"
		}
		retval += "</div>"
		retval += "<br />"
		retval += "<form method='post' action='/say'><input name='said' type='text'><button>Say</button></form>"

		return retval
	}
}

func (ce *ChatEngine) GenerateSayCurrentUser() SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return "No user logged in"
		}
		if !ce.userState.IsLoggedIn(username) {
			return "Not logged in"
		}
		if !ce.IsChatting(username) {
			return "Not currently chatting"
		}
		said, found := ctx.Params["said"]
		if !found {
			return instapage.MessageOKback("Chat", "Can't chat without saying anything")
		}

		ce.Say(username, CleanUpUserInput(said))
		ctx.SetHeader("Refresh", "0; url=/chat", true)
		return ""
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
