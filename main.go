package main

//
// TODO:
//
//     Make the title text <a href>
//     Refactor out getip/setip part to a separate project which is an alternative to dyndns
//

import (
	"fmt"
	"mime"

	"github.com/gosexy/canvas"          // For generating images
	. "github.com/xyproto/browserspeak" // For generating html/xml/css
	"github.com/xyproto/web"            // For serving webpages and handling requests
)

const (
	NICEBLUE = "#5080D0"
	NICEGRAY = "#202020"
)

type (
	// Every input from the user must be intitially stored in a UserInput variable, not in a string!
	// This is for security and to keep it clean.
	UserInput      string
	StringFunction SimpleWebHandle
)

func Publish(url, filename string, cache bool) {
	if cache {
		web.Get(url, CacheWrapper(url, File(filename)))
	} else {
		web.Get(url, File(filename))
	}
}

var globalStringCache map[string]string

// Wrap a SimpleContextHandle so that the output is cached (with an id)
// Do not cache functions with side-effects! (that sets the mimetype for instance)
// The safest thing for now is to only cache images.
func CacheWrapper(id string, f SimpleContextHandle) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if _, ok := globalStringCache[id]; !ok {
			globalStringCache[id] = f(ctx)
		}
		return globalStringCache[id]
	}
}

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

// TODO: Don't write to a file, but return the image data
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

	// Connect to Redis
	connection, err := NewRedisConnection()
	if err != nil {
		panic("ERROR: Can't connect to redis")
	}
	defer connection.Close()

	// The dynamic IP webpage (returns an *IPState)
	ServeIPs(connection)

	// The login system, returns a *UserState
	userState := ServeUserSystem(connection)

	// The archlinux.no webpage
	ServeArchlinuxNo(userState)

	// Compilation errors
	web.Get("/error", Errorlog)
	web.Get("/errors", Errorlog)

	// Various .php and .asp urls that showed up in the log
	ServeForFun()

	// Not found
	web.Get("/(.*)", NotFound)

	// Serve on port 3000 for the Nginx instance to use
	web.Run("0.0.0.0:3000")
}
