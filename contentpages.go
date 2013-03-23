package main

import (
	"time"

	"github.com/drbawb/mustache"
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

type ContentPage struct {
	generatedCSSurl          string
	extraCSSurls             []string
	jqueryJSurl              string
	faviconurl               string
	bgImageURL               string
	stretchBackground        bool
	title                    string
	subtitle                 string
	links                    []string
	contentTitle             string
	contentHTML              string
	headerJS                 string
	contentJS                string
	searchButtonText         string
	searchURL                string
	footerText               string
	backgroundTextureURL     string
	darkBackgroundTextureURL string
	footerTextColor          string
	footerColor              string
	userState                *UserState
	roundedLook              bool
	url                      string
}

type CPgen (func(userState *UserState) *ContentPage)

// A collection of ContentPages
type PageCollection []ContentPage

// The default settings
// Do not publish this page directly, but use it as a basis for the other pages
func DefaultCP(userState *UserState) *ContentPage {
	var cp ContentPage
	cp.generatedCSSurl = "/css/style.css"
	cp.extraCSSurls = []string{"/css/extra.css"}
	// TODO: fallback to local jquery.min.js, google how
	cp.jqueryJSurl = "//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js" // "/js/jquery-1.9.1.js"
	cp.faviconurl = "/favicon.ico"
	cp.links = []string{"Overview:/", "Login:/login", "Logout:/logout", "Register:/register", "Admin:/admin"}
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

	cp.footerText = "NOP"

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

	// This in combination with hiding the link in genericsite.go is cool, but the layout becomes weird :/
	//cp.headerJS += ShowAnimatedIf("/showmenu/admin", "#menuAdmin")

	// This keeps the layout but is less cool
	cp.headerJS += HideIfNot("/showmenu/admin", "#menuAdmin")

	cp.url = "/" // To be filled in when published

	return &cp
}

// TODO: Consider using Mustache for replacing elements after the page has been generated
// (for showing/hiding "login", "logout" or "register"
func genericPageBuilder(cp *ContentPage) *Page {
	// TODO: Record the time from one step out, because content may be generated and inserted into this generated conten
	startTime := time.Now()

	page := NewHTML5Page(cp.title + " " + cp.subtitle)

	page.LinkToCSS(cp.generatedCSSurl)
	for _, cssurl := range cp.extraCSSurls {
		page.LinkToCSS(cssurl)
	}
	page.LinkToJS(cp.jqueryJSurl)
	page.LinkToFavicon(cp.faviconurl)

	AddHeader(page, cp.headerJS)
	AddBodyStyle(page, cp.bgImageURL, cp.stretchBackground)
	AddTopBox(page, cp.title, cp.subtitle, cp.searchURL, cp.searchButtonText, cp.backgroundTextureURL, cp.roundedLook)

	// TODO:
	// Use something dynamic to add or remove /login and /register depending on the login status
	// The login status can be fetched over AJAX or REST or something.

	// TODO: Move the menubox into the TopBox

	// TODO: Do this with templates instead
	// Hide login/logout/register by default
	hidelist := []string{"/login", "/logout", "/register"} //, "/admin"}
	AddMenuBox(page, cp.links, hidelist, cp.darkBackgroundTextureURL)

	AddContent(page, cp.contentTitle, cp.contentHTML+BodyJS(cp.contentJS))

	elapsed := time.Since(startTime)
	AddFooter(page, cp.footerText, cp.footerTextColor, cp.footerColor, elapsed)

	return page
}

// Publish a list of ContentPaages, a colorscheme and template content
func PublishCPs(pc PageCollection, cs *ColorScheme, tp map[string]string, cssurl string) {
	// For each content page in the page collection
	for _, cp := range pc {
		// TODO: different css urls for all of these?
		cp.Pub(cp.url, cssurl, cs, tp)
	}
}

// Returns a BaseCP with the contentTitle set
func BaseTitleCP(contentTitle string, userState *UserState) *ContentPage {
	cp := BaseCP(userState)
	cp.contentTitle = contentTitle
	return cp
}

func ServeSite(userState *UserState, cps PageCollection, tp map[string]string, cs *ColorScheme) {
	// Add pages for login, logout and register
	cps = append(cps, *LoginCP(userState, "/login"))
	cps = append(cps, *RegisterCP(userState, "/register"))

	PublishCPs(cps, cs, tp, "/css/extra.css")

	ServeSearchPages(userState, cps, cs, tp)
	ServeAdminPages(userState, cps, cs, tp)

	// TODO: Add fallback to this local version
	Publish("/js/jquery-"+JQUERY_VERSION+".js", "static/js/jquery-"+JQUERY_VERSION+".js", true)

	// TODO: Generate these
	Publish("/robots.txt", "static/various/robots.txt", false)
	Publish("/sitemap_index.xml", "static/various/sitemap_index.xml", false)
}

// Make an html and css page available
func (cp *ContentPage) Pub(url, cssurl string, cs *ColorScheme, templateContent map[string]string) {
	genericpage := genericPageBuilder(cp)
	web.Get(url, GenerateHTMLwithTemplate(genericpage, templateContent))
	web.Get(cp.generatedCSSurl, GenerateCSS(genericpage))
	web.Get(cssurl, GenerateArchMenuCSS(cp.stretchBackground, cs))
}

// Wrap a lonely string in an entire webpage
func (cp *ContentPage) Surround(s string, tp map[string]string) (string, string) {
	cp.contentHTML = s
	archpage := genericPageBuilder(cp)
	// NOTE: Use GetXML(true) instead of .String() or .GetHTML() because some things are rendered
	// differently with different text layout!
	return mustache.Render(archpage.GetXML(true), tp), archpage.GetCSS()
}

// Uses a given SimpleWebHandle as the contents for the the ContentPage contents
func (cp *ContentPage) WrapSimpleWebHandle(swh SimpleWebHandle, tp map[string]string) SimpleWebHandle {
	return func(val string) string {
		html, css := cp.Surround(swh(val), tp)
		web.Get(cp.generatedCSSurl, css)
		return html
	}
}

// Uses a given WebHandle as the contents for the the ContentPage contents
func (cp *ContentPage) WrapWebHandle(wh WebHandle, tp map[string]string) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := cp.Surround(wh(ctx, val), tp)
		web.Get(cp.generatedCSSurl, css)
		return html
	}
}

// Uses a given SimpleContextHandle as the contents for the the ContentPage contents
func (cp *ContentPage) WrapSimpleContextHandle(sch SimpleContextHandle, tp map[string]string) SimpleContextHandle {
	return func(ctx *web.Context) string {
		html, css := cp.Surround(sch(ctx), tp)
		web.Get(cp.generatedCSSurl, css)
		return html
	}
}
