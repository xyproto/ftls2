package main

import (
	"github.com/hoisie/web"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/permissions2"
	"github.com/xyproto/pinterface"
	"github.com/xyproto/siteengines"
	"github.com/xyproto/webhandle"
	//"github.com/xyproto/personplan"
)

// TODO: Norwegian everywhere
// TODO: Different Redis database than the other sites

const JQUERY_VERSION = "2.0.0"

func ServeEngines(userState pinterface.IUserState, mainMenuEntries genericsite.MenuEntries) error {

	// The user engine
	userEngine, err := siteengines.NewUserEngine(userState)
	if err != nil {
		return err
	}
	userEngine.ServePages("ftls2.roboticoverlords.org")

	// The admin engine
	adminEngine, err := siteengines.NewAdminEngine(userState)
	if err != nil {
		return err
	}
	adminEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	// Wiki engine
	wikiEngine, err := siteengines.NewWikiEngine(userState)
	if err != nil {
		return err
	}
	wikiEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	// Timetable engine
	ftlsEngine, err := siteengines.NewTimeTableEngine(userState)
	if err != nil {
		return err
	}
	ftlsEngine.ServePages(FTLSBaseCP, mainMenuEntries)

	return nil
}

func main() {
	// UserState with a Redis Connection Pool, using database index 2
	userState := permissions.NewUserState(2, true, ":6379")
	defer userState.Close()

	// The FTLS2 webpage,
	mainMenuEntries := ServeFTLS(userState, "/js/jquery-"+JQUERY_VERSION+".min.js")

	ServeEngines(userState, mainMenuEntries)

	// Compilation errors, vim-compatible filename
	web.Get("/error", webhandle.GenerateErrorHandle("errors.err"))
	web.Get("/errors", webhandle.GenerateErrorHandle("errors.err"))

	// Various .php and .asp urls that showed up in the log
	genericsite.ServeForFun()

	// TODO: Consider adding support for the HEAD HTTP verb

	// Serve on port 3002 for the Nginx instance to use
	web.Run(":3002")
}