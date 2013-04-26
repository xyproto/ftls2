package main

import (
	"github.com/xyproto/browserspeak"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/siteengines"
	"github.com/xyproto/web"
)

// TODO: Norwegian everywhere
// TODO: Different Redis database than the other sites

const JQUERY_VERSION = "2.0.0"

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(browserspeak.NotFound(ctx, val)))
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

	// FTLS engine
	ftlsEngine := siteengines.NewFTLSEngine(userState)
	ftlsEngine.ServePages(FTLSBaseCP, mainMenuEntries)
}

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
	genericsite.ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3002 for the Nginx instance to use
	web.Run("0.0.0.0:3002")
}
