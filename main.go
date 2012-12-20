package main

import (
	"strings"

	"github.com/hoisie/web"
	"github.com/xyproto/browserspeak"
)

type ArchPageContents struct {
	generatedCSSurl string
	extraCSSurl string
	bgImageFilename string
	bgImageURL string
	title string
	subtitle string
	links []string
	contentTitle string
	contentHTML string
	searchButtonText string
	searchURL string
	footerText string
}

const (
	NICEBLUE = "#5080D0"
)

func addHeader(page *browserspeak.Page) {
	page.MetaCharset("UTF-8")
	page.LinkToGoogleFont("Russo One")
	page.LinkToGoogleFont("Geostar Fill")
}

func addBodyStyle(page *browserspeak.Page, bgimageurl string) {
	body, _ := page.SetMargin(4)
	body.RepeatBackground(bgimageurl, "repeat-x")
	page.SetColor("gray", "#d9d9d9")
	page.SetFontFamily("sans serif")
}

func addTitle(page *browserspeak.Page, title, subtitle string) {
	body, err := page.GetTag("body")
	if err != nil {
		return
	}

	// TODO: Add in a div and take parameters for x and y pos

	word1 := title
	word2 := ""
	if strings.Contains(title, " ") {
		word1 = strings.SplitN(title, " ", 2)[0]
		word2 = strings.SplitN(title, " ", 2)[1]
	}
	h1 := body.AddNewTag("h1")
	h1.AddAttr("id", "titletext")
	h1.AddStyle("font-family", "Russo One")

	a := h1.AddNewTag("a")
	a.AddAttr("id", "homelink")
	a.AddAttr("href", "/")
	a.AddContent(word1)
	a.AddStyle("color", "white")
	a.AddStyle("text-decoration", "none")

	font := a.AddNewTag("font")
	font.AddAttr("id", "bluetitle")
	font.AddStyle("color", NICEBLUE)
	font.AddContent(word2)

	font = a.AddNewTag("font")
	font.AddAttr("id", "graytitle")
	font.AddStyle("font-size", "0.5em")
	font.AddStyle("color", "#808080")
	font.AddContent(subtitle)
}

// Add a search box to the page, actionURL is the url to use as a get action, buttonText is the text on the search button
func addSearchBox(page *browserspeak.Page, actionURL string, buttonText string) {
	body, err := page.GetTag("body")
	if err != nil {
		return
	}
	div := body.AddNewTag("div")
	div.AddAttr("id", "searchbox")
	div.AddStyle("width", "25%")
	div.AddStyle("margin-left", "75%")
	div.AddStyle("position", "relative")
	div.AddStyle("top", "-3.2em")
	div.AddStyle("text-align", "right")

	form := div.AddNewTag("form")
	form.AddAttr("id", "search")
	form.AddAttr("method", "get")
	form.AddAttr("action", actionURL)

	inputText := form.AddNewTag("input")
	inputText.AddAttr("name", "q")
	inputText.AddAttr("size", "22")

	inputButton := form.AddNewTag("input")
	inputButton.AddStyle("margin-left", "0.4em")
	inputButton.AddAttr("type", "submit")
	inputButton.AddAttr("value", buttonText)
}

// Split a string at the colon into two strings
// If there's no colon, return the string and an empty string
func colonsplit(s string) (string, string) {
	if strings.Contains(s, ":") {
		sl := strings.SplitN(s, ":", 2)
		return sl[0], sl[1]
	}
	return s, ""
}

func boxStyle(tag *browserspeak.Tag) {
	tag.AddStyle("border", "solid 1px #b4b4b4")
	tag.AddStyle("border-radius", "10px")
	tag.AddStyle("box-shadow", "0 1px 3px rgba(0,0,0, .3)")
	tag.AddStyle("background-color", "#eeeeee")
}

// Takes a page and a colon-separated string slice of text:url
func addMenu(page *browserspeak.Page, links []string) {
	body, err := page.GetTag("body")
	if err != nil {
		return
	}
	div := body.AddNewTag("div")
	div.AddAttr("id", "menubox")
	div.AddStyle("margin-left", "-1em")
	div.AddStyle("margin-top", "1.5em")
	div.AddStyle("padding-top", "1em")
	div.AddStyle("padding-left", "1em")
	div.AddStyle("padding-bottom", "0.5em")
	div.AddStyle("position", "absolute")
	div.AddStyle("top", "150px")
	div.AddStyle("float", "left")

	//ul := div.AddNewTag("ul")
	//ul.AddAttr("id", "nav")

	//ul.AddStyle("margin", "0")
	//ul.AddStyle("padding", "0")
	//ul.AddStyle("width", "185px")
	//ul.AddStyle("position", "absolute")
	//ul.AddStyle("top", "35px")
	//ul.AddStyle("left", "0")
	//boxStyle(ul)

	//var li
	var a *browserspeak.Tag
	var text, url string
	//styleadded1 := false
	styleadded2 := false
	for _, text_url := range links {
		text, url = colonsplit(text_url)
		//li = ul.AddNewTag("li")
		//li.AddAttr("id", "menuitem") //strings.ToLower(text))

		//if !styleadded1 {
		//	li.AddStyle("margin", "1px 5px")
		//	li.AddStyle("padding", "0 0 8px")
		//	li.AddStyle("float", "left")
		//	li.AddStyle("position", "relative")
		//	li.AddStyle("list-style", "none")
		//	li.AddStyle("display", "inline-block")
		//	styleadded1 = true
		//}

		a = div.AddNewTag("a")
		a.AddAttr("id", "menulink")
		a.AddAttr("href", url)
		if !styleadded2 {
			a.AddStyle("font-weight", "bold")
			a.AddStyle("color", "#303030")
			a.AddStyle("text-decoration", "none")
			a.AddStyle("padding", "8px 1.2em")
			a.AddStyle("margin", "0")
			a.AddStyle("font-family", "sans serif")
			a.AddStyle("display", "block")
			a.AddStyle("width", "60px")
			styleadded2 = true
		}
		//a.AddStyle("font-size", "0.5em")
		a.AddContent(text)

	}
}

func addLogoSearchMenu(page *browserspeak.Page, title, subtitle string, links []string, searchButtonText, searchURL string) {
	addTitle(page, title, subtitle)
	addSearchBox(page, searchURL, searchButtonText)
	addMenu(page, links)
}

func addFooter(page *browserspeak.Page, footerText string) {
	body, err := page.GetTag("body")
	if err != nil {
		return
	}
	div := body.AddNewTag("div")
	div.AddAttr("id", "notice")
	div.AddStyle("position", "absolute")
	div.AddStyle("top", "90%")
	div.AddStyle("left", "90%")
	div.AddStyle("font-size", "0.7em")
	div.AddContent(footerText)
	//div.AddNewTag("br")
}

func addContent(page *browserspeak.Page, contentTitle, contentHTML string) {
	body, err := page.GetTag("body")
	if err != nil {
		return
	}
	//script := body.AddNewTag("script")
	//script.AddAttr("language", "javascript")
	//script.AddContent(`document.getElementById("overview").setAttribute("class", "current");`)

	div := body.AddNewTag("div")
	div.AddAttr("id", "content")
	div.AddStyle("z-index", "-1")
	div.AddStyle("color", "black")
	div.AddStyle("width", "85%")
	div.AddStyle("min-height", "600px")
	div.AddStyle("float", "left")
	div.AddStyle("position", "relative")
	div.AddStyle("margin-left", "150px")
	div.AddStyle("padding-left", "5em")
	div.AddStyle("padding-bottom", "0.5em")
	div.AddStyle("margin-top", "1.5em")
	div.AddStyle("padding-top", "1em")
	boxStyle(div)

	h2 := div.AddNewTag("h2")
	h2.AddAttr("id", "textheader")
	h2.AddContent(contentTitle)

	p := div.AddNewTag("p")
	p.AddAttr("id", "textparagraph")
	p.AddStyle("margin-top", "0.5em")
	p.AddStyle("font-family", "sans serif")
	p.AddStyle("font-size", "1.0em")
	p.AddStyle("color", "black")
	p.AddContent(contentHTML)
}

// TODO: Use strings like {{.Title}} instead and run it through the html.template system
func archbuilder(apc *ArchPageContents) *browserspeak.Page {
	page := browserspeak.NewHTML5Page(apc.title + " " + apc.subtitle)

	page.LinkToCSS(apc.generatedCSSurl)
	page.LinkToCSS(apc.extraCSSurl)

	addHeader(page)
	addBodyStyle(page, apc.bgImageURL)
	addLogoSearchMenu(page, apc.title, apc.subtitle, apc.links, apc.searchButtonText, apc.searchURL)
	addContent(page, apc.contentTitle, apc.contentHTML)
	addFooter(page, apc.footerText)

	return page
}

func hello(val string) string {
	return browserspeak.Message("root page", "hello: "+val)
}

func hover(ctx *web.Context) string {
	ctx.ContentType("css")
	return `
#menulink:link {color:black;}
#menulink:visited {color:black;}
#menulink:hover {color:` + NICEBLUE + `;}
#menulink:active {color:red;}
`
}

func search(ctx *web.Context, val string) string {
	q, found := ctx.Params["q"]
	if found == false {
		return browserspeak.Message("Error", "invalid params")
	}
	return browserspeak.Message("Search", q)
}

func BaseAPC() *ArchPageContents {
	var apc ArchPageContents
	apc.generatedCSSurl = "/css/style.css"
	apc.extraCSSurl = "/css/extra.css"
	apc.bgImageFilename = "img/longbg4.png"
	apc.bgImageURL = "/img/longbg4.png"
	apc.title = "Arch Linux"
	apc.subtitle = "Norway"
	apc.links = []string{"Overview:/", "List:/list", "Submit:/submit", "Ais:http://github.com/xyproto/ais", "Setconf:http://setconf.roboticoverlords.org/", "ArchFriend:https://play.google.com/store/apps/details?id=com.xyproto.archfriend"}
	apc.contentTitle = "Welcome"
	apc.contentHTML = "Hi there!"
	apc.searchButtonText = "Search"
	apc.searchURL = "/search"
	apc.footerText = "Alexander Rødseth, 2012"
	return &apc
}

func HiAPC() *ArchPageContents {
	apc := BaseAPC()
	apc.contentHTML = "Hi!"
	return apc
}

func main() {
	apc := HiAPC()
	archpage := archbuilder(apc)
	web.Get("/", browserspeak.HTML(archpage))
	web.Get(apc.generatedCSSurl, browserspeak.CSS(archpage))
	web.Get(apc.extraCSSurl, hover)
	web.Get(apc.bgImageURL, browserspeak.FILE(apc.bgImageFilename))

	web.Get(apc.searchURL + "(.*)", search)
	web.Get("/error", browserspeak.Errorlog)

	web.Get("/hello/(.*)", hello)

	web.Get("/(.*)", browserspeak.NotFound)
	web.Run("0.0.0.0:3000")
}
