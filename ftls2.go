package main

import (
	"fmt"

	"github.com/xyproto/browserspeak"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/instapage"
	"github.com/xyproto/siteengines"
	"github.com/xyproto/web"
)

const JQUERY_VERSION = "2.0.0"

func hello(val string) string {
	return instapage.Message("root page", "hello: "+val)
}

func helloHandle(ctx *web.Context, name string) string {
	return "Hello, " + name
}

func Hello() string {
	msg := "Hi"
	return instapage.Message("Hello", msg)
}

func ParamExample(ctx *web.Context) string {
	return fmt.Sprintf("%v\n", ctx.Params)
}

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(browserspeak.NotFound(ctx, val)))
}

func ServeEngines(userState *genericsite.UserState, mainMenuEntries genericsite.MenuEntries) {
	// The user engine
	userEngine := siteengines.NewUserEngine(userState)
	userEngine.ServePages()

	// The admin engine
	adminEngine := siteengines.NewAdminEngine(userState)
	adminEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	// The chat system (see also the menu entry in FTLSBaseCP)
	chatEngine := siteengines.NewChatEngine(userState)
	chatEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	// Wiki engine
	wikiEngine := siteengines.NewWikiEngine(userState)
	wikiEngine.ServePages(FTLSBaseCP, mainMenuEntries)
}

// TODO: Caching, login
func main() {

	// UserState with a Redis Connection Pool
	userState := genericsite.NewUserState()
	defer userState.Close()

	// The archlinux.no webpage,
	mainMenuEntries := ServeFTLS(userState, "/js/jquery-"+JQUERY_VERSION+".min.js")

	ServeEngines(userState, mainMenuEntries)

	// Compilation errors, vim-compatible filename
	web.Get("/error", browserspeak.GenerateErrorHandle("errors.err"))
	web.Get("/errors", browserspeak.GenerateErrorHandle("errors.err"))

	// Various .php and .asp urls that showed up in the log
	ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3000 for the Nginx instance to use
	web.Run("0.0.0.0:3000")
}
