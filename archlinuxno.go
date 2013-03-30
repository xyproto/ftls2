package main

// OK, "archlinuxno", 23-03-13

import (
	"fmt"

	"github.com/xyproto/browserspeak"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/web"
)

func hello(val string) string {
	return browserspeak.Message("root page", "hello: "+val)
}

func helloHandle(ctx *web.Context, name string) string {
	return "Hello, " + name
}

func Hello() string {
	msg := "Hi"
	return browserspeak.Message("Hello", msg)
}

func ParamExample(ctx *web.Context) string {
	return fmt.Sprintf("%v\n", ctx.Params)
}

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(browserspeak.NotFound(ctx, val)))
}

// TODO: Caching, login
func main() {

	// Redis Connection Pool
	pool := genericsite.InitSystem()
	defer pool.Close()

	userEngine := genericsite.NewUserEngine(pool)
	userEngine.ServeSystem()
	userState := userEngine.GetState()

	// The dynamic IP webpage (returns an *IPState)
	ServeIPs(userState)

	adminEngine := genericsite.NewAdminEngine(userState)
	adminEngine.ServeSystem()

	adminEngine.ServePages(ArchBaseCP, DynamicMenuFactory(userState))

	// The archlinux.no webpage
	ServeArchlinuxNo(userState)

	// The chat system (see also the menu entry in ArchBaseCP)
	chatEngine := NewChatEngine(userState)
	chatEngine.ServeSystem()
	chatEngine.ServePages(ArchBaseCP)

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
