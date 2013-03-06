package main

import (
	"strings"
	"time"

	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
	//"github.com/hoisie/mustache"
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

// TODO: Consider using Mustache for replacing elements after the page has been generated
// (for showing/hiding "login", "logout" or "register"
func archbuilder(apc *ContentPage) *Page {
	page := NewHTML5Page(apc.title + " " + apc.subtitle)

	page.LinkToCSS(apc.generatedCSSurl)
	for _, cssurl := range apc.extraCSSurls {
		page.LinkToCSS(cssurl)
	}
	page.LinkToJS(apc.jqueryJSurl)
	page.LinkToFavicon(apc.faviconurl)

	AddHeader(page, apc.headerJS)
	AddBodyStyle(page, apc.bgImageURL, apc.stretchBackground)
	AddTopBox(page, apc.title, apc.subtitle, apc.searchURL, apc.searchButtonText, apc.backgroundTextureURL, apc.roundedLook)

	// TODO:
	// Use something dynamic to add or remove /login and /register depending on the login status
	// The login status can be fetched over AJAX or REST or something.

	// TODO: Move the menubox into the TopBox

	// TODO: Do this with templates instead
	// Hide login/logout/register by default
	hidelist := []string{"/login", "/logout", "/register"}
	AddMenuBox(page, apc.links, hidelist, apc.darkBackgroundTextureURL)

	AddContent(page, apc.contentTitle, apc.contentHTML+BodyJS(apc.contentJS))
	AddFooter(page, apc.footerText, apc.footerTextColor, apc.footerColor)

	return page
}

// Make an html and css page available
func (apc *ContentPage) Pub(url, cssurl string) {
	archpage := archbuilder(apc)
	web.Get(url, GenerateHTMLwithTemplate(archpage, Kake()))
	web.Get(apc.generatedCSSurl, GenerateCSS(archpage))
	cs := NewArchColorScheme()
	web.Get(cssurl, GenerateArchMenuCSS(apc.stretchBackground, cs))
}

// Wrap a lonely string in an entire webpage
func (apc *ContentPage) Surround(s string) (string, string) {
	apc.contentHTML = s
	archpage := archbuilder(apc)
	return archpage.GetXML(true), archpage.GetCSS()
}

// Uses a given SimpleWebHandle as the contents for the the ArchPage contents
func (apc *ContentPage) WrapSimpleWebHandle(swh SimpleWebHandle) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := apc.Surround(swh(val))
		web.Get(apc.generatedCSSurl, css)
		return html
	}
}

// Uses a given SimpleWebHandle as the contents for the the ArchPage contents
func (apc *ContentPage) WrapWebHandle(wh WebHandle) WebHandle {
	return func(ctx *web.Context, val string) string {
		html, css := apc.Surround(wh(ctx, val))
		web.Get(apc.generatedCSSurl, css)
		return html
	}
}

// Publish a list of ContentPage
func PublishCPs(pc PageCollection) {
	// For each content page in the page collection
	for _, cp := range pc {
		// TODO: different css urls for all of these?
		cp.Pub(cp.url, "/css/extra.css")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Search a list of ContentPage for a given searchText
// Returns a list of urls or an empty list, a list of page titles and the string that was actually searched for
func searchResults(userSearchText UserInput, pc PageCollection) ([]string, []string, string) {
	// Search for maximum 100 letters, lowercase and trimmed
	searchText := strings.ToLower(strings.TrimSpace(string(userSearchText)[:min(100, len(string(userSearchText)))]))

	if searchText == "" {
		// No search results for the empty string
		return []string{}, []string{}, ""
	}

	var matches, titles []string
	for _, apc := range pc {
		if strings.Contains(strings.ToLower(apc.url), searchText) ||
		   strings.Contains(strings.ToLower(apc.contentTitle), searchText) ||
		   strings.Contains(strings.ToLower(apc.contentHTML), searchText) {
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
// Searches a list of ContentPage structs
func GenerateSearchHandle(pc PageCollection) WebHandle {
	return func(ctx *web.Context, val string) string {
		q, found := ctx.Params["q"]
		searchText := UserInput(q)
		if found {
			content := "Search: " + string(searchText)
			nl := tagString("br")
			content += nl + nl
			startTime := time.Now()
			urls, titles, searchedFor := searchResults(searchText, pc)
			elapsed := time.Since(startTime)
			page, p := CowboyTag("p")
			if len(urls) == 0 {
				p.AddContent("No results found")
				p.AddNewTag("br")
			} else {
				for i, url := range urls {
					a := p.AddNewTag("a")
					a.AddAttr("id", "searchresult")
					a.AddStyle("color", "red")
					a.AddAttr("href", url)
					a.AddContent(titles[i])
					font := p.AddNewTag("font")
					font.AddContent(" - contains \"" + searchedFor + "\" in the url, title or content")
					p.AddNewTag("br")
				}
			}
			p.AddNewTag("br")
			p.AddLastContent("Search took: " + elapsed.String())
			return page.String()
		}
		return "Invalid parameters"
	}
}

func PublishImages() {
	faviconFilename := "generated/img/favicon.ico"
	genFavicon(faviconFilename)
	Publish("/favicon.ico", faviconFilename, false)

	// Tried previously:
	// "rough.png", "longbg.png", "donutbg.png", "donutbg_light.jpg",
	// "felix_predator2.jpg", "centerimage.png", "underwater.png",
	// "norway.jpg", "norway2.jpg", "underwater.jpg"

	// Publish and cache images
	imgs := []string{"norway3.jpg", "gray.jpg", "darkgray.jpg"}
	for _, img := range imgs {
		Publish("/img/"+img, "static/img/"+img, true)
	}
}

// Returns a BaseCP with the contentTitle set
func BaseTitleCP(contentTitle string, userState *UserState) *ContentPage {
	apc := BaseCP(userState)
	apc.contentTitle = contentTitle
	return apc
}

func GenerateSearchCSS() SimpleContextHandle {
	return func(ctx *web.Context) string {
		ctx.ContentType("css")
		return `
#searchresult {
	color: ` + NICEBLUE + `;
	text-decoration: underline;
}
`
	}
}

func ServeDynamicPages(userState *UserState, cps []ContentPage) {
	web.Get("/hello/(.*)", BaseTitleCP("Hello", userState).WrapSimpleWebHandle(helloSF))
	// Note, no slash between "search" and "(.*)". A typical search is "/search?q=blabla"
	searchCP := BaseTitleCP("Search results", userState)
	searchCP.extraCSSurls = append(searchCP.extraCSSurls, "/css/search.css")
	web.Get("/search(.*)", searchCP.WrapWebHandle(GenerateSearchHandle(cps)))
	web.Get("/css/search.css", GenerateSearchCSS())
}

func ServeSite(userState *UserState, cps []ContentPage) {
	// Add pages for login, logout and register
	cps = append(cps, *LoginCP(userState, "/login"))
	cps = append(cps, *LogoutCP(userState, "/logout"))
	cps = append(cps, *RegisterCP(userState, "/register"))

	web.Get("/showmenu/loginlogoutregister", GenerateShowLoginLogoutRegister(userState))

	PublishCPs(cps)

	ServeDynamicPages(userState, cps)

	// TODO: Add fallback to this local version
	Publish("/js/jquery-"+JQUERY_VERSION+".js", "static/js/jquery-"+JQUERY_VERSION+".js", true)
	Publish("/robots.txt", "static/various/robots.txt", false)
	Publish("/sitemap_index.xml", "static/various/sitemap_index.xml", false)

	PublishImages()
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
		*HelloCP(userState, "/mirrors"),
		*HelloCP(userState, "/feedback"),
	}
	ServeSite(userState, cps)
}
