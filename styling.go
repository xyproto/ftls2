package main

// Move to "browserspeak"

// Various functions that can be used to style a webpage

import (
	. "github.com/xyproto/browserspeak"
)

type ColorScheme struct {
	darkgray           string
	nicecolor          string
	menu_link          string
	menu_hover         string
	menu_active        string
	default_background string
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
