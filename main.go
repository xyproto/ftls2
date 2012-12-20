package main

import (
	"io/ioutil"
	"bytes"
	"path/filepath"
	"strings"

	"github.com/hoisie/web"
	"github.com/xyproto/browserspeak"
)

func addHeader(page *browserspeak.Page) {
	page.metaCharset(page, "UTF-8")
	page.linkToGoogleFont(page, "Russo One")
	page.linkToGoogleFont(page, "Geostar Fill")
}

func addBodyStart(page *browserspeak.Page, bgimageurl string) {
	body, _ := page.SetMargin(4)
	repeatBackground(body, bgimageurl, "repeat-x")
	page.SetColor("gray", "#d9d9d9")
	page.SetFont("sans-serif")
}

func addLogoSearchMenu(page *browserspeak.Page, title, subtitle string) {
	body, err := page.GetTag("body")
	if err != nil {
		return
	}
	word1 := title
	word2 := ""
	if strings.Contains(title, " ") {
		word1 = strings.SplitN(title, " ", 2)[0]
		word2 = strings.SplitN(title, " ", 2)[1]
	}
	h1 := body.AddNewTag("h1")
	h1.AddStyle("font-family", "Russo One")
	a := h1.AddNewTag("a")
	a.AddAttr("href", "/")
	a.AddContent(word1)
	font := a.AddNewTag("font")
	font.AddAttr("id", "blue")
	font.AddStyle("color", "#5080D0")
	font.AddContent(word2)
	font = a.AddNewTag("font")
	font.AddAttr("id", "gray")
	font.AddStyle("font-size", "0.5em")
	font.AddStyle("color", "#808080")
	font.AddContent(subtitle)
}

func archbuilder(cssurl string) *browserspeak.Page {
	title := "Arch Linux Norway"

	page := browserspeak.NewHTML5Page(title)

	page.linkToCSS(page, cssurl)

	addHeader(page)
	addBodyStart(page, "/img/longbg4.png")
	addLogoSearchMenu(page, "Arch Linux", "Norway")

	if body, err := page.GetTag("body"); err == nil {
		h2 := body.AddNewTag("h2")
		h2.AddContent(title)
		h2.AddStyle("margin-left", "2%")
		h2.AddStyle("margin-top", "0.5em")
		h2.AddStyle("color", "black")
		h2.AddStyle("font-family", "tahoma, arial, sans-serif")
		p := body.AddNewTag("p")
		p.AddStyle("margin-left", "2%")
		p.AddStyle("margin-top", "0.5em")
		p.AddStyle("font-family", "sans-serif")
		p.AddStyle("font-size", "1.0em")
		p.AddStyle("color", "black")
	}

	return page
}

// This is a test function
func testbuilder(cssurl string) *browserspeak.Page {
	page := browserspeak.NewHTML5Page("Hello")
	body, _ := page.SetMargin(3)

	h1 := body.AddNewTag("h1")
	h1.SetMargin(1)
	h1.AddContent("Browser")

	h1, err := page.GetTag("h1")
	if err == nil {
		h1.AddContent("Spe")
	}

	if err = linkToCSS(page, cssurl); err == nil {
		h1.AddContent("ak")
	} else {
		h1.AddContent("akkkkkkkk")
	}

	page.SetColor("#202020", "#A0A0A0")
	page.SetFont("sans-serif")

	box, _ := page.AddBox("box0", true)
	box.AddStyle("margin-top", "-2em")
	box.AddStyle("margin-bottom", "3em")

	image := body.AddImage("http://www.shoutmeloud.com/wp-content/uploads/2010/01/successful-Blogger.jpeg", "50%")
	image.AddStyle("margin-top", "2em")
	image.AddStyle("margin-left", "3em")

	return page
}



func hello(val string) string {
	return browserspeak.Message("root page", "hello: "+val)
}



func main() {
	browserspeak.Publish("/test", "/main.css", testbuilder)
	browserspeak.Publish("/", "/css/style.css", archbuilder)

	// used by archbuilder
	browserspeak.PublishRootFile("img/longbg4.png")

	web.Get("/error", errorlog)

	web.Get("/hi", hi)

	web.Get("/hello/(.*)", hello)

	web.Get("/(.*)", notFound)
	web.Run("0.0.0.0:3000")
}
