package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/permissions"
	"github.com/xyproto/siteengines"
	"github.com/xyproto/webhandle"

	//"github.com/xyproto/personplan"
)

// TODO: Norwegian everywhere
// TODO: Different Redis database than the other sites

const JQUERY_VERSION = "2.0.0"

func ServeEngines(userState *permissions.UserState, mainMenuEntries genericsite.MenuEntries) {
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
	// Create a Negroni and a ServeMux instance
	n := negroni.Classic()
	mux := http.NewServeMux()

	// UserState with a Redis Connection Pool, using database index 2
	userState := permissions.NewUserState(2, true, ":6379")
	defer userState.Close()

	// The archlinux.no webpage,
	mainMenuEntries := ServeFTLS(userState, "/js/jquery-"+JQUERY_VERSION+".min.js")

	ServeEngines(userState, mainMenuEntries)

	// Compilation errors, vim-compatible filename
	mux.HandleFunc("/error", webhandle.GenerateErrorHandler("errors.err"))
	mux.HandleFunc("/errors", webhandle.GenerateErrorHandler("errors.err"))

	// Various .php and .asp urls that showed up in the log
	genericsite.ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3002 for the Nginx instance to use
	n.Run(":3002")
}
