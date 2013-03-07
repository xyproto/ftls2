package main

// Various functions that can be used to style a webpage

import (
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

const (
	NICEGRAY  = "#202020"
	NICEBLUE  = "#5080D0"
)

type ColorScheme struct {
	darkgray string
	niceblue string
	menu_link string
	menu_hover string
	menu_active string
	default_background string
}

func NewArchColorScheme() *ColorScheme {
	var cs ColorScheme
	cs.darkgray  = "#202020"
	cs.niceblue  = "#5080D0"
	cs.menu_link = "#c0c0c0" // light gray
	cs.menu_hover = "#efefe0" // light gray, somewhat yellow
	cs.menu_active = "#ffffff" // white
	cs.default_background = "#000030"
	return &cs
}

// TODO: get style values from a file instead?
func AddHeader(page *Page, js string) {
	AddGoogleFonts(page, []string{"Armata"}) //, "Junge"})
	// TODO: Move to browserspeak
	page.MetaCharset("UTF-8")
	AddScriptToHeader(page, js)
}

func AddGoogleFonts(page *Page, googleFonts []string) {
	for _, fontname := range googleFonts {
		page.LinkToGoogleFont(fontname)
	}
}

func AddScriptToHeader(page *Page, js string) error {
	// Check if there's anything to add
	if js == "" {
		// Nope
		return nil
	}
	// Add a script tag
	head, err := page.GetTag("head")
	if err == nil {
		script := head.AddNewTag("script")
		script.AddAttr("type", "text/javascript")
		script.AddContent(js)
	}
	return err
}

func AddBodyStyle(page *Page, bgimageurl string, stretchBackground bool) {
	body, _ := page.SetMargin(1)
	body.SansSerif()
	if stretchBackground {
		body.AddStyle("background", "url('"+bgimageurl+"') no-repeat center center fixed")
	} else {
		body.AddStyle("background", "url('"+bgimageurl+"')")
	}
}

func GenerateArchMenuCSS(stretchBackground bool, cs *ColorScheme) SimpleContextHandle {
	return func(ctx *web.Context) string {
		ctx.ContentType("css")
		// one of the extra css files that are loaded after the main style
		retval := `
a {
  text-decoration: none;
  color: #303030;
  font-weight: regular;
}
a:link {color:` + cs.menu_link + `;}
a:visited {color:` + cs.menu_link + `;}
a:hover {color:` + cs.menu_hover + `;}
a:active {color:` + cs.menu_active + `;}
`
		// The load order of background-color, background-size and background-image
		// is actually significant in Chrome! Do not reorder lightly!
		if stretchBackground {
			retval = "body {\nbackground-color: " + cs.default_background + ";\nbackground-size: cover;\n}\n" + retval
		} else {
			retval = "body {\nbackground-color: " + cs.default_background + ";\n}\n" + retval
		}
		return retval
	}
}
