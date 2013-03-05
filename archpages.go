package main

import (
	"strconv"
	"time"
)

// The main page / default settings
func BaseCP(userState *UserState) *ContentPage {
	var apc ContentPage
	apc.generatedCSSurl = "/css/style.css"
	apc.extraCSSurls = []string{"/css/extra.css"}
	apc.jqueryJSurl = "//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js" // "/js/jquery-1.9.1.js"
	apc.faviconurl = "/favicon.ico"
	apc.bgImageURL = "/img/norway3.jpg"
	apc.stretchBackground = true
	apc.title = "Arch Linux"
	apc.subtitle = "no"
	//apc.links = []string{"Overview:/", "Mirrors:/mirrors", "Login:/login", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//apc.links = []string{"Overview:/", "Text:/text", "Bob:/bob", "JQuery:/jquery", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	apc.links = []string{"Overview:/", "Login:/login", "Logout:/logout", "Register:/register"}
	apc.contentTitle = "Hi"
	apc.contentHTML = "Hi there!"
	apc.contentJS = ""
	apc.headerJS = ""
	apc.searchButtonText = "Search"
	apc.searchURL = "/search"
	// http://wptheming.wpengine.netdna-cdn.com/wp-content/uploads/2010/04/gray-texture.jpg
	apc.backgroundTextureURL = "/img/gray.jpg"
	// http://turbo.designwoop.com/uploads/2012/03/16_free_subtle_textures_subtle_dark_vertical.jpg
	apc.darkBackgroundTextureURL = "/img/darkgray.jpg"
	apc.footerColor = "black"
	apc.footerTextColor = "#303040"
	y := time.Now().Year()
	//apc.footerText = "Alexander Rødseth &lt;rodseth@gmail.com&gt;, " + strconv.Itoa(y)
	apc.footerText = "Alexander Rødseth, " + strconv.Itoa(y)
	apc.userState = userState
	apc.roundedLook = false

	// Javascript that applies everywhere
	//apc.contentJS += HideIfNot("/showmenu/login", "#menuLogin")
	//apc.contentJS += HideIfNot("/showmenu/logout", "#menuLogout")
	//apc.contentJS += HideIfNot("/showmenu/register", "#menuRegister")
	//apc.contentJS += HideIfNotLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")
	//apc.contentJS += ShowIfLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")
	apc.headerJS += ShowIfLoginLogoutRegister("/showmenu/loginlogoutregister", "#menuLogin", "#menuLogout", "#menuRegister")

	apc.url = "/" // To be filled in when published

	return &apc
}

func LoginCP(userState *UserState, url string) *ContentPage {
	apc := BaseCP(userState)
	apc.contentTitle = "Login"
	// TODO: jquery get + ensure cookie is set
	apc.contentHTML = "<a href=\"/login/bob\">login</a>"
	apc.url = url
	return apc
}

func LogoutCP(userState *UserState, url string) *ContentPage {
	apc := BaseCP(userState)
	apc.contentTitle = "Logout"
	apc.contentHTML = "<a href=\"/logout/bob\">logout</a>"
	apc.url = url
	return apc
}

func RegisterCP(userState *UserState, url string) *ContentPage {
	apc := BaseCP(userState)
	apc.contentTitle = "Register"
	apc.contentHTML = "<a href=\"/create/bob\">register</a>"
	apc.url = url
	return apc
}

func OverviewCP(userState *UserState, url string) *ContentPage {
	apc := BaseCP(userState)
	apc.contentTitle = "Overview"
	apc.contentHTML = `This site is currently under construction.<br />You may wish to visit the <a href="https://bbs.archlinux.org/viewtopic.php?id=4998">Arch Linux Forum</a> in the meantime.<br /><br /><i>- Alexander Rødseth &lt;rodseth / gmail&gt;</i>`
	apc.url = url
	return apc
}
