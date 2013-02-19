package main

import (
	"strconv"

	"github.com/xyproto/browserspeak"
)

// Create an empty page only containing the given tag
// Returns both the page and the tag
func cowboyTag(tagname string) (*browserspeak.Page, *browserspeak.Tag) {
	page := browserspeak.NewPage("blank", tagname)
	tag, _ := page.GetTag(tagname)
	return page, tag
}

func tagString(tagname string) string {
	page := browserspeak.NewPage("blank", tagname)
	return page.String()
}

func SetPixelPosition(tag *browserspeak.Tag, xpx, ypx int) {
	tag.AddStyle("position", "absolute")
	xpxs := strconv.Itoa(xpx) + "px"
	ypxs := strconv.Itoa(ypx) + "px"
	tag.AddStyle("top", xpxs)
	tag.AddStyle("left", ypxs)
}

func SetRelativePosition(tag *browserspeak.Tag, x, y string) {
	tag.AddStyle("position", "relative")
	tag.AddStyle("top", x)
	tag.AddStyle("left", y)
}

func SetWidthAndSide(tag *browserspeak.Tag, width string, leftside bool) {
	side := "right"
	if leftside {
		side = "left"
	}
	tag.AddStyle("float", side)
	tag.AddStyle("width", width)
}
