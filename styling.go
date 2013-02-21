package main

// Various functions that can be used to style a webpage

import (
	"github.com/hoisie/web"
	"github.com/xyproto/browserspeak"
)

// TODO: get style values from a file

func AddHeader(page *browserspeak.Page) {
	page.MetaCharset("UTF-8")
	page.LinkToGoogleFont("Armata")
	page.LinkToGoogleFont("Junge")
}

func AddBodyStyle(page *browserspeak.Page, bgimageurl string) {
	body, _ := page.SetMargin(1)
	body.SansSerif()
	body.AddStyle("background", "url('"+bgimageurl+"') no-repeat center center fixed")
	//body.AddStyle("background-size", "cover")
	//body.AddStyle("background-color", "#808080")
	//body.AddStyle("background-size", "100% 100%")
	////body.AddStyle("background-repeat", "no-repeat")
	////body.RepeatBackground(bgimageurl, "repeat-x")
	//page.SetColor("gray", "#a0e0e0") // gray text, turquise background color
	//page.SetColor("gray", "#202020") // gray text, turquise background color
	//page.SetColor("gray", "#202020") // gray text, turquise background color
}

// extra.css
// TODO: Rename this function
func hover(ctx *web.Context) string {
	menucolor := "#c0c0c0"   // light gray
	hovercolor := "#efefe0"  // very light gray, with some yellow
	activecolor := "#ffffff" // white
	ctx.ContentType("css")
	return `
#menulink:link {color:` + menucolor + `;}
#menulink:visited {color:` + menucolor + `;}
#menulink:hover {color:` + hovercolor + `;}
#menulink:active {color:` + activecolor + `;}
body {
	background-color: #808080;
	background-size: cover;
}
`
	// The load order of background-color, background-size and background-image is actually significant in Chrome! Do not reorder lightly!
}
