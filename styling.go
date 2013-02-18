package main

// Various functions that can be used to style a webpage

import (
	"github.com/xyproto/browserspeak"
	"github.com/hoisie/web"
)

// TODO: get style values from a file

// Boxes for content
func RoundedBox(tag *browserspeak.Tag) {
	tag.AddStyle("border", "solid 1px #b4b4b4")
	tag.AddStyle("border-radius", "10px")
	tag.AddStyle("box-shadow", "0 1px 3px rgba(0,0,0, .3)")
	tag.AddStyle("background-color", "#e0e0e0") // light gray
}

// Set the tag font to some sort of sans-serif
func SansSerif(tag *browserspeak.Tag) {
	tag.AddStyle("font-family", "Verdana, Geneva, sans-serif")
}

// Set the tag font to the given font or just some sort of sans-serif
func CustomSansSerif(tag *browserspeak.Tag, custom string) {
	tag.AddStyle("font-family", custom+", Verdana, Geneva, sans-serif")
}

func AddHeader(page *browserspeak.Page) {
	page.MetaCharset("UTF-8")
	page.LinkToGoogleFont("Armata")
	page.LinkToGoogleFont("Junge")
}

func AddBodyStyle(page *browserspeak.Page, bgimageurl string) {
	body, _ := page.SetMargin(1)
	SansSerif(body)
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
	menucolor := NICEBLUE
	hovercolor := "#c0c0a0" // light gray, with some yellow
	activecolor := "#d0d0b0" // very light gray, with some yellow
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

