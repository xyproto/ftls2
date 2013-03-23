package main

/*
 * This is the actual static content for the Arch Linux website.
 *
 * There is also login/logout/register, searching and dynamic pages, in other source files
 */

import (
	"strconv"
	"time"

	"github.com/xyproto/web"
)

// The default settings for Arch Linux content pages
func ArchBaseCP(state *UserState) *ContentPage {
	cp := DefaultCP(state)
	cp.bgImageURL = "/img/norway4.jpg"
	cp.stretchBackground = true
	cp.title = "Arch Linux"
	cp.subtitle = "no"

	//cp.links = []string{"Overview:/", "Mirrors:/mirrors", "Login:/login", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//cp.links = []string{"Overview:/", "Text:/text", "Bob:/bob", "JQuery:/jquery", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//News, Norwegian AUR
	cp.links = append(cp.links, "JQuery:/jquery"}

	y := time.Now().Year()

	//cp.footerText = "Alexander Rødseth &lt;rodseth@gmail.com&gt;, " + strconv.Itoa(y)
	cp.footerText = "Alexander Rødseth, " + strconv.Itoa(y)


	// Javascript that applies everywhere
	//cp.contentJS += HideIfNot("/showmenu/login", "#menuLogin")
	//cp.contentJS += HideIfNot("/showmenu/logout", "#menuLogout")
	//cp.contentJS += HideIfNot("/showmenu/register", "#menuRegister")
	//cp.contentJS += HideIfNotLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")
	//cp.contentJS += ShowIfLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")

	// This only works at first page load in Internet Explorer 8. Fun times. Oh well, why bother.
	cp.headerJS += ShowIfLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")

	// This in combination with hiding the link in genericsite.go is cool, but the layout becomes weird :/
	//cp.headerJS += ShowAnimatedIf("/showmenu/admin", "#menuAdmin")

	// This keeps the layout but is less cool
	cp.headerJS += HideIfNot("/showmenu/admin", "#menuAdmin")

	cp.url = "/" // To be filled in when published

	return &cp
}

func OverviewCP(userState *UserState, url string) *ContentPage {
	cp := BaseCP(userState)
	cp.contentTitle = "Overview"
	cp.contentHTML = `This site is currently under construction.<br />Visit the <a href="https://bbs.archlinux.org/viewtopic.php?id=4998">Arch Linux Forum</a> in the meantime.<br /><br /><i>- Alexander Rødseth &lt;rodseth / gmail&gt;</i>`
	cp.url = url
	return cp
}

func MirrorsCP(userState *UserState, url string) *ContentPage {
	cp := BaseCP(userState)
	cp.contentTitle = "Mirrors"
	cp.contentHTML = "List over Norwegian Arch Linux mirrors:"
	cp.url = url
	return cp
}

func PublishArchImages() {
	//faviconFilename := "/static/generated/img/favicon.ico"
	//genFavicon(faviconFilename)
	//Publish("/favicon.ico", faviconFilename, false)
	Publish("/favicon.ico", "static/img/favicon.ico", false)

	// Tried previously:
	// "rough.png", "longbg.png", "donutbg.png", "donutbg_light.jpg",
	// "felix_predator2.jpg", "centerimage.png", "underwater.png",
	// "norway.jpg", "norway2.jpg", "underwater.jpg"

	// Publish and cache images
	imgs := []string{"norway4.jpg", "norway3.jpg", "gray.jpg", "darkgray.jpg"}
	for _, img := range imgs {
		Publish("/img/"+img, "static/img/"+img, true)
	}
}

// Routing for the archlinux.no webpage
// Admin, search and user management is already provided
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

	// "dynamic" pages
	// Makes helloSF handle the content for /hello/(.*) urls, but wrapped in a BaseCP with the title "Hello"
	web.Get("/hello/(.*)", BaseTitleCP("Hello", userState).WrapSimpleWebHandle(helloSF, Kake()))

   // static images
	PublishArchImages()
}
