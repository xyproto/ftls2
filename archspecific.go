package main

import (
	"strings"
	"time"

	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
	//"github.com/hoisie/mustache"
)

type ArchPageContents struct {
	generatedCSSurl          string
	extraCSSurl              string
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

type APCgen (func(userState *UserState) *ArchPageContents)

// TODO: Consider using Mustache for replacing elements after the page has been generated
// (for showing/hiding "login", "logout" or "register"
func archbuilder(apc *ArchPageContents) *Page {
	page := NewHTML5Page(apc.title + " " + apc.subtitle)

	page.LinkToCSS(apc.generatedCSSurl)
	page.LinkToCSS(apc.extraCSSurl)
	page.LinkToJS(apc.jqueryJSurl)
	page.LinkToFavicon(apc.faviconurl)

	AddHeader(page, apc.headerJS)
	AddBodyStyle(page, apc.bgImageURL, apc.stretchBackground)
	AddTopBox(page, apc.title, apc.subtitle, apc.searchURL, apc.searchButtonText, apc.backgroundTextureURL, apc.roundedLook)

	// TODO:
	// Use something dynamic to add or remove /login and /register depending on the login status
	// The login status can be fetched over AJAX or REST or something.

	// TODO: Move the menubox into the TopBox

	// Hide login/logout/register by default
	hidelist := []string{"/login", "/logout", "/register"}
	AddMenuBox(page, apc.links, hidelist, apc.darkBackgroundTextureURL)

	AddContent(page, apc.contentTitle, apc.contentHTML+BodyJS(apc.contentJS))
	AddFooter(page, apc.footerText, apc.footerTextColor, apc.footerColor)

	return page
}

// Make an html and css page available
func (apc *ArchPageContents) Pub(url string) {
	archpage := archbuilder(apc)
	web.Get(url, HTML(archpage))
	web.Get(apc.generatedCSSurl, CSS(archpage))
	web.Get(apc.extraCSSurl, GenerateExtraCSS(apc.stretchBackground))
}

// Wrap a lonely string in an entire webpage
func (apc *ArchPageContents) Surround(s string) (string, string) {
	apc.contentHTML = s
	archpage := archbuilder(apc)
	return archpage.GetXML(true), archpage.GetCSS()
}

// Uses a given SimpleWebHandle as the contents for the the ArchPage contents
func (apc *ArchPageContents) WrapSimpleWebHandle(swh SimpleWebHandle) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := apc.Surround(swh(val))
		web.Get(apc.generatedCSSurl, css)
		return html
	}
}

// Uses a given SimpleWebHandle as the contents for the the ArchPage contents
func (apc *ArchPageContents) WrapWebHandle(wh WebHandle) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := apc.Surround(wh(ctx, val))
		web.Get(apc.generatedCSSurl, css)
		return html
	}
}



// Publish a list of ArchPageContents
func PublishAPCs(apcs []ArchPageContents) {
	for _, apc := range apcs {
		apc.Pub(apc.url)
	}
}

// Generate "true" or "false" depending on the cookie status
func GenerateShowLogin(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		return "1" // true
	}
}

// Generate "true" or "false" depending on the cookie status
func GenerateShowLogout(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		return "1" // true
	}
}

// Generate "true" or "false" depending on the cookie status
func GenerateShowRegister(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		return "1" // true
	}
}

// TODO: Rethink this. Use templates for Login/Logout button?
// Generate "1" or "0" values for showing the login, logout or register menus,
// depending on the cookie status and UserState
func GenerateShowLoginLogoutRegister(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if username := GetBrowserUsername(ctx); username != "" {
			//print("USERNAME", username)
			// Has a username stored in the browser
			if state.LoggedIn(username) {
				// Ok, logged in to the system + login cookie in the browser
				// Only present the "Logout" menu
				return "010"
			} else {
				// Has a login cookie, but is not logged in.
				// Keep the browser cookie (could be tempting to remove it)
				// Present only the "Login" menu
				//return "100"
				// Present both "Login" and "Register", just in case it's a new user
				// in the same browser.
				return "101"
			}
		} else {
			// Does not have a username stored in the browser
			// Present the "Register" and "Login" menu
			return "101"
		}
		// Everything went wrong, should never reach this point
		return "000"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Search a list of ArchPageContents for a given searchText
// Returns a list of urls or an empty list, a list of page titles and the string that was actually searched for
func searchResults(userSearchText UserInput, apcs []ArchPageContents) ([]string, []string, string) {
	// Search for maximum 100 letters, lowercase and trimmed
	searchText := strings.ToLower(strings.TrimSpace(string(userSearchText)[:min(100, len(string(userSearchText)))]))

	if searchText == "" {
		// No search results for the empty string
		return []string{}, []string{}, ""
	}

	var matches, titles []string
	for _, apc := range apcs {
		if strings.Contains(strings.ToLower(apc.contentTitle), searchText) || strings.Contains(strings.ToLower(apc.contentHTML), searchText) {
			// Check if the url is already in the matches slices
			found := false
			for _, url := range matches {
				if url == apc.url {
					found = true
					break
				}
			}
			// If not, add it
			if !found {
				matches = append(matches, apc.url)
				titles = append(titles, apc.contentTitle)
			}
		}
	}
	return matches, titles, searchText
}

// Generate a search handle. This is done in order to be able to modify the apc
// Searches a list of ArchPageContents structs
func GenerateSearchHandle(apcs []ArchPageContents) WebHandle {
	return func(ctx *web.Context, val string) string {
		q, found := ctx.Params["q"]
		searchText := UserInput(q)
		if found {
			content := "Search: " + string(searchText)
			nl := tagString("br")
			content += nl + nl
			startTime := time.Now()
			urls, titles, searchedFor := searchResults(searchText, apcs)
			elapsed := time.Since(startTime)
			page, p := CowboyTag("p")
			if len(urls) == 0 {
				p.AddContent("No results found")
				p.AddNewTag("br")
			} else {
				for i, url := range urls {
					a := p.AddNewTag("a")
					a.AddAttr("id", "searchresult")
					a.AddAttr("href", url)
					a.AddContent(titles[i])
					font := p.AddNewTag("font")
					font.AddContent(" - contains \"" + searchedFor + "\"")
					p.AddNewTag("br")
				}
			}
			p.AddNewTag("br")
			p.AddContent("Search took: " + elapsed.String())
			return page.String()
		}
		return "Invalid parameters"
	}
}

func ServeImages() {
	faviconFilename := "generated/img/favicon.ico"
	genFavicon(faviconFilename)
	Publish("/favicon.ico", faviconFilename, false)

	// Images
	//Publish("/img/rough.png", "static/img/rough.png")
	//Publish("/img/longbg.png", "static/img/longbg.png")
	//Publish("/img/donutbg.png", "static/img/donutbg.png")
	//Publish("/img/donutbg_light.jpg", "static/img/donutbg_light.jpg")
	//Publish("/img/boxes_cartoon2.png", "static/img/boxes_cartoon2.png")
	//Publish("/img/boxes_softglow.png", "static/img/boxes_softglow.png")
	//Publish("/img/felix_predator2.jpg", "static/img/felix_predator2.jpg")
	//Publish("/img/space_predator.png", "static/img/space_predator.png")
	//Publish("/img/centerimage.png", "static/img/centerimage.png")
	//Publish("/img/underwater.png", "static/img/underwater.png")
	//Publish("/img/norway.jpg", "static/img/norway.jpg")
	//Publish("/img/norway2.jpg", "static/img/norway2.jpg")
	//Publish("/img/underwater.jpg", "static/img/underwater.jpg")
	Publish("/img/norway3.jpg", "static/img/norway3.jpg", true)
	Publish("/img/gray.jpg", "static/img/gray.jpg", true)
	Publish("/img/darkgray.jpg", "static/img/darkgray.jpg", true)
}

// Returns a BaseAPC with the contentTitle set
func BaseTitleAPC(contentTitle string, userState *UserState) *ArchPageContents {
	apc := BaseAPC(userState)
	apc.contentTitle = contentTitle
	return apc
}

// Routing for the archlinux.no webpage
func ServeArchlinuxNo(userState *UserState) {

	helloAPC := BaseAPC(userState)
	helloAPC.contentTitle = "Hello"
	web.Get("/hello/(.*)", helloAPC.WrapSimpleWebHandle(helloSF))

	//web.Get("/showmenu/login", GenerateShowLogin(userState))
	//web.Get("/showmenu/logout", GenerateShowLogout(userState))
	//web.Get("/showmenu/register", GenerateShowRegister(userState))
	web.Get("/showmenu/loginlogoutregister", GenerateShowLoginLogoutRegister(userState))

	// Pages that only depends on the user state
	apcs := []ArchPageContents{
		*OverviewAPC(userState, "/"),
		*LoginAPC(userState, "/login"),
		*LogoutAPC(userState, "/logout"),
		*RegisterAPC(userState, "/register"),
		*TextAPC(userState, "/text"),
		*JQueryAPC(userState, "/jquery"),
		*BobAPC(userState, "/bob"),
		*CountAPC(userState, "/counting"),
		*HelloAPC(userState, "/mirrors"),
		*HelloAPC(userState, "/feedback"),
	}
	PublishAPCs(apcs)

	// Dynamic pages

	// Note, no slash between "search" and "(.*)". A typical search is "/search?q=blabla"
	web.Get("/search(.*)", BaseTitleAPC("Search results", userState).WrapWebHandle(GenerateSearchHandle(apcs)))

	// TODO: Add fallback to this version
	Publish("/js/jquery-1.9.1.js", "static/js/jquery-1.9.1.js", true)
	Publish("/robots.txt", "static/various/robots.txt", false)
	Publish("/sitemap_index.xml", "static/various/sitemap_index.xml", false)

	//web.Get(apc.searchURL+"(.*)", search)

	ServeImages()
}
