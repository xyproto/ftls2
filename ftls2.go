package main

import (
	"github.com/hoisie/web"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/siteengines"
	"github.com/xyproto/webhandle"
	//"github.com/xyproto/personplan"
)

// TODO: Norwegian everywhere
// TODO: Different Redis database than the other sites

const JQUERY_VERSION = "2.0.0"

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(webhandle.NotFound(ctx, val)))
}

func ServeEngines(userState *genericsite.UserState, mainMenuEntries genericsite.MenuEntries) {
	// The user engine
	userEngine := siteengines.NewUserEngine(userState)
	userEngine.ServePages("ftls2.roboticoverlords.org")

	// The admin engine
	adminEngine := siteengines.NewAdminEngine(userState)
	adminEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	// Wiki engine
	wikiEngine := siteengines.NewWikiEngine(userState)
	wikiEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	// Timetable engine
	ftlsEngine := siteengines.NewTimeTableEngine(userState)
	ftlsEngine.ServePages(FTLSBaseCP, mainMenuEntries)
}

func main() {

	// UserState with a Redis Connection Pool, using database index 2
	userState := genericsite.NewUserState(2)
	defer userState.Close()

	// The archlinux.no webpage,
	mainMenuEntries := ServeFTLS(userState, "/js/jquery-"+JQUERY_VERSION+".min.js")

	ServeEngines(userState, mainMenuEntries)

	// Compilation errors, vim-compatible filename
	web.Get("/error", webhandle.GenerateErrorHandle("errors.err"))
	web.Get("/errors", webhandle.GenerateErrorHandle("errors.err"))

	// Various .php and .asp urls that showed up in the log
	genericsite.ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3002 for the Nginx instance to use
	web.Run("0.0.0.0:3002")
}
