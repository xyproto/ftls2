package main

// OK, "archlinuxno", 23-03-13

import (
	"fmt"

	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

func hello(val string) string {
	return Message("root page", "hello: "+val)
}

func helloSF(name string) string {
	return "Hello, " + name
}

func Hello() string {
	msg := "Hi"
	return Message("Hello", msg)
}

func ParamExample(ctx *web.Context) string {
	return fmt.Sprintf("%v\n", ctx.Params)
}

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(NotFound(ctx, val)))
}

// TODO: Caching, login
func main() {

	userState := InitSystem()

	// The dynamic IP webpage (returns an *IPState)
	ServeIPs(userState)

	userEngine := NewUserEngine(userState)
	userEngine.ServeSystem()

	adminEngine := NewAdminEngine(userState)
	adminEngine.ServeSystem()

	// The archlinux.no webpage
	ServeArchlinuxNo(userState)

	// Compilation errors
	web.Get("/error", Errorlog)
	web.Get("/errors", Errorlog)

	// Various .php and .asp urls that showed up in the log
	ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3000 for the Nginx instance to use
	web.Run("0.0.0.0:3000")
}
