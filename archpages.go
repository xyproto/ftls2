package main

/*
 * This is the actual static content for the Arch Linux website.
 *
 * There is also login/logout/register, searching and dynamic pages, in other source files
 */

 import (
	"strconv"
	"time"
)

// The default settings
// Do not publish this page directly, but use it as a basis for the other pages
// TODO: Remove anything Arch-specific, rename the default background images
// TODO: Make it easy to replace logo and footer text
func BaseCP(userState *UserState) *ContentPage {
	var cp ContentPage
	cp.generatedCSSurl = "/css/style.css"
	cp.extraCSSurls = []string{"/css/extra.css"}
	// TODO: fallback to local jquery.min.js, google how
	cp.jqueryJSurl = "//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js" // "/js/jquery-1.9.1.js"
	cp.faviconurl = "/favicon.ico"
	cp.bgImageURL = "/img/norway4.jpg"
	cp.stretchBackground = true
	cp.title = "Arch Linux"
	cp.subtitle = "no"
	//cp.links = []string{"Overview:/", "Mirrors:/mirrors", "Login:/login", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//cp.links = []string{"Overview:/", "Text:/text", "Bob:/bob", "JQuery:/jquery", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//News, Norwegian AUR

	cp.links = []string{"Overview:/", "Login:/login", "Logout:/logout", "Register:/register"}
	cp.contentTitle = "NOP"
	cp.contentHTML = "NOP NOP NOP"
	cp.contentJS = ""
	cp.headerJS = ""
	cp.searchButtonText = "Search"
	cp.searchURL = "/search"
	// http://wptheming.wpengine.netdna-cdn.com/wp-content/uploads/2010/04/gray-texture.jpg
	// TODO: Draw these two backgroundimages with a canvas instead
	cp.backgroundTextureURL = "/img/gray.jpg"
	// http://turbo.designwoop.com/uploads/2012/03/16_free_subtle_textures_subtle_dark_vertical.jpg
	cp.darkBackgroundTextureURL = "/img/darkgray.jpg"
	cp.footerColor = "black"
	cp.footerTextColor = "#303040"
	y := time.Now().Year()
	//cp.footerText = "Alexander Rødseth &lt;rodseth@gmail.com&gt;, " + strconv.Itoa(y)
	cp.footerText = "Alexander Rødseth, " + strconv.Itoa(y)
	cp.userState = userState
	cp.roundedLook = false

	// Javascript that applies everywhere
	//cp.contentJS += HideIfNot("/showmenu/login", "#menuLogin")
	//cp.contentJS += HideIfNot("/showmenu/logout", "#menuLogout")
	//cp.contentJS += HideIfNot("/showmenu/register", "#menuRegister")
	//cp.contentJS += HideIfNotLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")
	//cp.contentJS += ShowIfLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")

	// This only works at first page load in Internet Explorer 8. Fun times. Oh well, why bother.
	cp.headerJS += ShowIfLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")

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

