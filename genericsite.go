package main

import (
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

// TODO: a form
func LoginCP(userState *UserState, url string) *ContentPage {
	cp := BaseCP(userState)
	cp.contentTitle = "Login"
	// TODO: jquery get + ensure cookie is set
	// TODO: a form using jquery to post
	// TODO
	cp.contentHTML = LoginForm()
	cp.contentJS += OnClick("#loginButton", "$('#loginForm').get(0).setAttribute('action', '/login/' + $('#username').val());")
	cp.url = url
	return cp
}

// TODO: a form
func RegisterCP(userState *UserState, url string) *ContentPage {
	cp := BaseCP(userState)
	cp.contentTitle = "Register"
	cp.contentHTML = RegisterForm()
	cp.contentJS += OnClick("#registerButton", "$('#registerForm').get(0).setAttribute('action', '/register/' + $('#username').val());")
	cp.url = url
	return cp
}

// TODO: Consider using Mustache for replacing elements after the page has been generated
// (for showing/hiding "login", "logout" or "register"
func genericbuilder(cp *ContentPage) *Page {
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
	hidelist := []string{"/login", "/logout", "/register"}
	AddMenuBox(page, cp.links, hidelist, cp.darkBackgroundTextureURL)

	AddContent(page, cp.contentTitle, cp.contentHTML+BodyJS(cp.contentJS))
	AddFooter(page, cp.footerText, cp.footerTextColor, cp.footerColor)

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

func GenerateSearchCSS(cs *ColorScheme) SimpleContextHandle {
	return func(ctx *web.Context) string {
		ctx.ContentType("css")
		// TODO: Rename niceblue to something non-color specific
		return `
#searchresult {
	color: ` + cs.niceblue + `;
	text-decoration: underline;
}
`
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

	web.Get("/showmenu/loginlogoutregister", GenerateShowLoginLogoutRegister(userState))

	PublishCPs(cps, cs, tp, "/css/extra.css")

	ServeSearchPages(userState, cps, cs, tp)

	// TODO: Add fallback to this local version
	Publish("/js/jquery-"+JQUERY_VERSION+".js", "static/js/jquery-"+JQUERY_VERSION+".js", true)
	Publish("/robots.txt", "static/various/robots.txt", false)
	Publish("/sitemap_index.xml", "static/various/sitemap_index.xml", false)
}

func ServeSearchPages(userState *UserState, cps PageCollection, cs *ColorScheme, tp map[string]string) {
	searchCP := BaseTitleCP("Search results", userState)
	searchCP.extraCSSurls = append(searchCP.extraCSSurls, "/css/search.css")
	// Note, no slash between "search" and "(.*)". A typical search is "/search?q=blabla"
	web.Get("/search(.*)", searchCP.WrapWebHandle(GenerateSearchHandle(cps), tp))
	web.Get("/css/search.css", GenerateSearchCSS(cs))
}

// Make an html and css page available
func (cp *ContentPage) Pub(url, cssurl string, cs *ColorScheme, templateContent map[string]string) {
	genericpage := genericbuilder(cp)
	web.Get(url, GenerateHTMLwithTemplate(genericpage, templateContent))
	web.Get(cp.generatedCSSurl, GenerateCSS(genericpage))
	web.Get(cssurl, GenerateArchMenuCSS(cp.stretchBackground, cs))
}

// Wrap a lonely string in an entire webpage
func (cp *ContentPage) Surround(s string, tp map[string]string) (string, string) {
	cp.contentHTML = s
	archpage := genericbuilder(cp)
	return mustache.Render(archpage.GetHTML(), tp), archpage.GetCSS()
}

// Uses a given SimpleWebHandle as the contents for the the ArchPage contents
func (cp *ContentPage) WrapSimpleWebHandle(swh SimpleWebHandle, tp map[string]string) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := cp.Surround(swh(val), tp)
		web.Get(cp.generatedCSSurl, css)
		return html
	}
}

// Uses a given SimpleWebHandle as the contents for the the ArchPage contents
func (cp *ContentPage) WrapWebHandle(wh WebHandle, tp map[string]string) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := cp.Surround(wh(ctx, val), tp)
		web.Get(cp.generatedCSSurl, css)
		return html
	}
}
