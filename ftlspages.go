package main

import (
	"strconv"
	"time"

	. "github.com/xyproto/genericsite"
	"github.com/xyproto/permissions2"
	"github.com/xyproto/pinterface"
	. "github.com/xyproto/siteengines"
)

// The default settings for FTLS content pages
func FTLSBaseCP(state pinterface.IUserState) *ContentPage {
	cp := DefaultCP(state)
	cp.Title = "Timeliste System"
	cp.Subtitle = "2"

	y := time.Now().Year()

	// TODO: Use templates for the footer, for more accurate measurment of the
	//       time made to generate the page.
	cp.FooterText = "Alexander F. Rødseth, " + strconv.Itoa(y)

	cp.Url = "/" // Is replaced when the contentpage is published

	cp.ColorScheme = NewFTLSColorScheme()

	// Behind the text
	cp.BgImageURL = "/img/rough_diagonal.png"
	cp.StretchBackground = false

	// Behind the menu
	cp.BackgroundTextureURL = "/img/hixs_pattern_evolution.png"

	cp.SearchBox = false

	return cp
}

func OverviewCP(userState *permissions.UserState, url string) *ContentPage {
	cp := FTLSBaseCP(userState)
	cp.ContentTitle = "Om"
	cp.ContentHTML = `FTLS 2 - under utvikling`
	cp.Url = url
	return cp
}

func TextCP(userState *permissions.UserState, url string) *ContentPage {
	apc := FTLSBaseCP(userState)
	apc.ContentTitle = "Text"
	apc.ContentHTML = `<p id='textparagraph'>Hei<br/>der<br/></p>`
	apc.Url = url
	return apc
}

// This is where the possibilities for the menu are listed
func Cps2MenuEntries(cps []ContentPage) MenuEntries {
	links := []string{
		"Om:/",
		"Logg inn:/login",
		"Registrer:/register",
		"Logg ut:/logout",
		"Admin:/admin",
		"Timer:/timetable",
		"Wiki:/wiki",
	}
	return Links2menuEntries(links)
}

// Routing for the FTLS2 webpage
// Admin, search and user management is already provided
func ServeFTLS(userState *permissions.UserState, jquerypath string) MenuEntries {
	cps := []ContentPage{
		*OverviewCP(userState, "/"),
		*TextCP(userState, "/text"),
		*LoginCP(FTLSBaseCP, userState, "/login"),
		*RegisterCP(FTLSBaseCP, userState, "/register"),
	}

	menuEntries := Cps2MenuEntries(cps)

	// template content generator
	tvgf := DynamicMenuFactoryGenerator(menuEntries)

	//ServeSearchPages(FTLSBaseCP, userState, cps, FTLSBaseCP(userState).ColorScheme, tvgf(userState))
	ServeSite(FTLSBaseCP, userState, cps, tvgf, jquerypath)

	return menuEntries
}

func NewFTLSColorScheme() *ColorScheme {
	var cs ColorScheme
	cs.Darkgray = "#202020"
	cs.Nicecolor = "#d80000"
	cs.Menu_link = "#c0c0c0"
	cs.Menu_hover = "#efefe0"
	cs.Menu_active = "#ffffff"
	cs.Default_background = "#000030"
	cs.TitleText = "#e0e0e0" // The first word of the title text
	return &cs
}
