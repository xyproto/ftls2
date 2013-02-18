package main

//
// TODO:
//
//     Refactor into:
//       * database related functions
//       * routing-related functions
//       * archlinux specific functions
//       * fronted specific functions
//
//     Add a <div> titlebox, make it fixed, let it be the title and search box
//
//     Make the color of the background beneath the titlebox a bit lighter (not black)
//
//     Refactor out getip/setip and the redis parts to a separate project which is an alternative to dyndns
//

import (
	"fmt"
	"mime"

	"github.com/gosexy/canvas"        // For generating images
	"github.com/hoisie/web"           // For serving webpages and handling requests
	"github.com/xyproto/browserspeak" // For generating html/xml/css
)

const (
	NICEBLUE = "#5080D0"
	NICEGRAY = "#202020"
)


type (
	WebHandle (func(ctx *web.Context, val string) string)
	StringFunction (func(string) string)
	SimpleWebHandle StringFunction
)

func hello(val string) string {
	return browserspeak.Message("root page", "hello: "+val)
}

func helloSF(name string) string {
	return "Hello, " + name
}

func Hello() string {
	msg := "Hi"
	return browserspeak.Message("Hello", msg)
}

func Publish(url, filename string) {
	web.Get(url, browserspeak.FILE(filename))
}

func ParamExample(ctx *web.Context) string {
	return fmt.Sprintf("%v\n", ctx.Params)
}

func genFavicon(filename string) {
	img := canvas.New()
	img.Blank(16, 16)
	img.SetStrokeColor("#005090")

	// All the lines and translations use relative coordinates

	// "\"
	img.SetStrokeWidth(2)
	img.Translate(8, 2)
	img.Line(3, 11)
	img.Translate(-8, -2)

	// "/"
	img.SetStrokeWidth(2)
	img.Translate(8, 2)
	img.Line(-6, 12)
	img.Translate(-8, -2)

	// "-"
	img.SetStrokeWidth(2)
	img.Translate(2, 10)
	img.Line(12, -2)

	img.Write(filename)
}

// TODO: Caching, login
func main() {

	// These common ones are missing!
	mime.AddExtensionType(".txt", "text/plain; charset=utf-8")
	mime.AddExtensionType(".ico", "image/x-icon")

	// The archlinux.no webpage
	ServeArchlinuxNo()

	// The dynamic IP webpage
	state := ServeIPs()
	defer state.Close()

	// Compilation errors
	web.Get("/error", browserspeak.Errorlog)
	web.Get("/errors", browserspeak.Errorlog)

	// Honeypot? Found these in the logs
	web.Get("/index.php", Hello)
	web.Get("/viewtopic.php", ParamExample)

	// Not found
	web.Get("/(.*)", browserspeak.NotFound)

	// Serve on port 3000 for the Nginx instance to use
	web.Run("0.0.0.0:3000")
}
