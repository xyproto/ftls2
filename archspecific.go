package main

import (
	"github.com/xyproto/web"
)

func ServeDynamicPages(userState *UserState) {
	// Makes helloSF handle the content for /hello/(.*) urls,
	// but wrapped in a BaseCP with the title "Hello"
	web.Get("/hello/(.*)", BaseTitleCP("Hello", userState).WrapSimpleWebHandle(helloSF, Kake()))
}

// Routing for the archlinux.no webpage
func ServeArchlinuxNo(userState *UserState) {
	// Pages that only depends on the user state
	cps := []ContentPage{
		*OverviewCP(userState, "/"),
		*TextCP(userState, "/text"),
		*JQueryCP(userState, "/jquery"),
		*BobCP(userState, "/bob"),
		*CountCP(userState, "/counting"),
		*MirrorsCP(userState, "/mirrors"),
		*HelloCP(userState, "/feedback"),
	}

	// template content
	tp := Kake()

	// color scheme
	cs := NewArchColorScheme()

	ServeSite(userState, cps, tp, cs)
	ServeDynamicPages(userState)

	PublishArchImages()
}
