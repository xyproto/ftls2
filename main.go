package main

import (
	"strings"
	"time"
	"strconv"

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

	var a *browserspeak.Tag
	var text, url string

	// TODO: Fix the problem with several #id's in the CSS from browserspeak
	styleadded := false
	for _, text_url := range links {
		text, url = colonsplit(text_url)

		a = div.AddNewTag("a")
		a.AddAttr("id", "menulink")
		a.AddAttr("href", url)
		if !styleadded {
			a.AddStyle("font-weight", "bold")
			a.AddStyle("color", "#303030")
			a.AddStyle("text-decoration", "none")
			a.AddStyle("padding", "8px 1.2em")
			a.AddStyle("margin", "0")
			a.AddStyle("font-family", "sans serif")
			a.AddStyle("display", "block")
			a.AddStyle("width", "60px")
			styleadded = true
		}
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

func BaseAPC() *ArchPageContents {
	var apc ArchPageContents
	apc.generatedCSSurl = "/css/style.css"
	apc.extraCSSurl = "/css/extra.css"
	apc.bgImageFilename = "img/longbg4.png"
	apc.bgImageURL = "/img/longbg4.png"
	apc.title = "Arch Linux"
	apc.subtitle = "Norge"
	apc.links = []string{"Overview:/", "Hello:/hello/world", "Count:/counting"}
	apc.contentTitle = "Welcome"
	apc.contentHTML = "Hi there!"
	apc.searchButtonText = "Search"
	apc.searchURL = "/search"
	y := time.Now().Year()
	apc.footerText = "Alexander Rødseth, " + strconv.Itoa(y)
	return &apc
}

func HiAPC() *ArchPageContents {
	apc := BaseAPC()
	apc.contentHTML = `Hi!</br></br>This page is under construction, you might want to visit the <a href="https://bbs.archlinux.org/">Arch Forum</a> in the mean while.</br></br>Alexander Rødseth &lt;rodseth _at gmail.com&gt;`
	return apc
}

func HelloAPC() *ArchPageContents {
	apc := BaseAPC()
	apc.contentTitle = "This is it"
	return apc
}

func helloSF(name string) string {
	return "Hello, " + name
}

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

// Generate a search handle. This is done in order to be able to modify the apc
func (apc *ArchPageContents) GenerateSearchHandle() WebHandle {
	return func (ctx *web.Context, val string) string {
	    q, found := ctx.Params["q"]
		var html, css string
		if found {
			apc.contentTitle = "Search results"
			content := "Search: " + q
			nl := tagString("br")
			content += nl + nl
			content += "No results found"
			html, css = apc.Surround(content)
		} else {
			apc.contentTitle = "Error"
			html, css = apc.Surround("Invalid parameters")
		}
		web.Get(apc.generatedCSSurl, css)
		return html
	}
}

// Make an html and css page available, plus a search handler (they are special)
func (apc *ArchPageContents) Pub(url string, search WebHandle) {
	archpage := archbuilder(apc)
	web.Get(url, browserspeak.HTML(archpage))
	web.Get(apc.generatedCSSurl, browserspeak.CSS(archpage))
	web.Get(apc.extraCSSurl, hover)
	web.Get(apc.bgImageURL, browserspeak.FILE(apc.bgImageFilename))
	web.Get(apc.searchURL + "(.*)", search)
}

// Wrap a lonely string in an entire webpage
func (apc *ArchPageContents) Surround(s string) (string, string) {
	apc.contentHTML = s
	archpage := archbuilder(apc)
	return archpage.GetXML(true), archpage.GetCSS()
}

type WebHandle (func (ctx *web.Context, val string) string)
type StringFunction (func (string) string)
type APCgen (func () *ArchPageContents)

// Creates a handle from s string function
func (apc *ArchPageContents) GetHandle(fn StringFunction) WebHandle {
	return func (ctx *web.Context, val string) string {
		html, css := apc.Surround(fn(val))
		web.Get(apc.generatedCSSurl, css)
		return html
	}
}

// Wraps the looks of an ArchPageContents + a string function together to a web.go handle
func wrapHandle(apcgen APCgen, sf StringFunction) WebHandle {
	apc := apcgen()
	return apc.GetHandle(sf)
}

// Creeates a page based on ArchPageContens generated by apcgen
// Publishes the HTML and CSS at the given URL
func pub(url string, apcgen APCgen) {
	apc := apcgen()
	apc.Pub(url, apc.GenerateSearchHandle())
}

func CountAPC() *ArchPageContents {
	apc := BaseAPC()
	apc.contentTitle = "Counting"
	apc.contentHTML = "1 2 3"
	return apc
}

// TODO: Caching
func main() {
	pub("/", HiAPC)

	web.Get("/hello/(.*)", wrapHandle(HelloAPC, helloSF))
	pub("/counting", CountAPC)

	web.Get("/error", browserspeak.Errorlog)
	web.Get("/(.*)", browserspeak.NotFound)
	web.Run("0.0.0.0:3000")
}
