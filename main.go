package main

//
// TODO:
//
//     Refactor into:
//       * database related functions
//       * routing-related functions
//       * archlinux specific functions
//       * fronted specific functions
//
//     Add a <div> titlebox, make it fixed, let it be the title and search box
//
//     Make the color of the background beneath the titlebox a bit lighter (not black)
//

import (
	"fmt"
	"mime"
	"strconv"
	"strings"
	"time"
	"bytes"

	// For serving webpages and handling requests
	"github.com/hoisie/web"

	// For generating html/xml/css
	"github.com/xyproto/browserspeak"

	// For generating images
	"github.com/gosexy/canvas"

	// For connecting to Redis
	"github.com/garyburd/redigo/redis"
)

type ArchPageContents struct {
	generatedCSSurl  string
	extraCSSurl      string
	faviconurl       string
	bgImageFilename  string
	bgImageURL       string
	title            string
	subtitle         string
	links            []string
	contentTitle     string
	contentHTML      string
	searchButtonText string
	searchURL        string
	footerText       string
}

type RedisList struct {
	c redis.Conn
	id string
}

type State struct {
	ips *RedisList
}

const (
	NICEBLUE = "#5080D0"
)

func addHeader(page *browserspeak.Page) {
	page.MetaCharset("UTF-8")
	page.LinkToGoogleFont("Armata")
	page.LinkToGoogleFont("Junge")
	//page.LinkToGoogleFont("Geostar Fill")
}

func addBodyStyle(page *browserspeak.Page, bgimageurl string) {
	body, _ := page.SetMargin(1)
	SansSerif(body)
	//body.RepeatBackground(bgimageurl, "repeat-x")
	page.SetColor("gray", "black") // "#d9d9d9")
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
	CustomSansSerif(h1, "Armata")
	//body.RepeatBackground(bgimageurl, "repeat-x")

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
	CustomSansSerif(inputButton, "Armata")
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

func CustomSansSerif(tag *browserspeak.Tag, custom string) {
	tag.AddStyle("font-family", custom + ", Verdana, Geneva, sans-serif")
}

func SansSerif(tag *browserspeak.Tag) {
	tag.AddStyle("font-family", "Verdana, Geneva, sans-serif")
}

//func 

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
			CustomSansSerif(a, "Armata")
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
	div.AddStyle("color", "black") // content headline color
	div.AddStyle("min-height", "80%")
	div.AddStyle("min-width", "60%")
	div.AddStyle("float", "left")
	div.AddStyle("position", "relative")
	div.AddStyle("margin-left", "150px")
	div.AddStyle("margin-top", "1em")
	div.AddStyle("padding-left", "4em")
	div.AddStyle("padding-right", "5em")
	div.AddStyle("padding-top", "1em")
	div.AddStyle("padding-bottom", "2em")
	boxStyle(div)

	h2 := div.AddNewTag("h2")
	h2.AddAttr("id", "textheader")
	h2.AddContent(contentTitle)
	CustomSansSerif(h2, "Armata")

	p := div.AddNewTag("p")
	p.AddAttr("id", "textparagraph")
	p.AddStyle("margin-top", "0.5em")
	CustomSansSerif(p, "Junge")
	p.AddStyle("font-size", "1.0em")
	p.AddStyle("color", "black") // content text color
	p.AddContent(contentHTML)
}

// TODO: Use strings like {{.Title}} instead and run it through the html.template system
func archbuilder(apc *ArchPageContents) *browserspeak.Page {
	page := browserspeak.NewHTML5Page(apc.title + " " + apc.subtitle)

	page.LinkToCSS(apc.generatedCSSurl)
	page.LinkToCSS(apc.extraCSSurl)
	page.LinkToFavicon(apc.faviconurl)

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
	menucolor := NICEBLUE;
	hovercolor := "white";
	activecolor := "yellow";
	ctx.ContentType("css")
	return `
#menulink:link {color:` + menucolor + `;}
#menulink:visited {color:` + menucolor + `;}
#menulink:hover {color:` + hovercolor + `;}
#menulink:active {color:` + activecolor + `;}
`
}

func BaseAPC() *ArchPageContents {
	var apc ArchPageContents
	apc.generatedCSSurl = "/css/style.css"
	apc.extraCSSurl = "/css/extra.css"
	apc.faviconurl = "/favicon.ico"
	apc.bgImageFilename = "static/img/longbg.png"
	apc.bgImageURL = "/img/longbg.png"
	apc.title = "Arch Linux"
	apc.subtitle = "no"
	apc.links = []string{"Overview:/", "Hello:/hello/world", "Count:/counting"}
	apc.contentTitle = "Oh no"
	apc.contentHTML = "Hi there!"
	apc.searchButtonText = "Search"
	apc.searchURL = "/search"
	y := time.Now().Year()
	apc.footerText = "Alexander Rødseth, " + strconv.Itoa(y)
	return &apc
}

func HiAPC() *ArchPageContents {
	apc := BaseAPC()
	apc.contentHTML = `This place is currently under construction.<br />You may want to visit the <a href="https://bbs.archlinux.org/">Arch Linux Forum</a> in the mean time.<br /><br /><i>- Alexander Rødseth &lt;rodseth / gmail&gt;</i>`
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
	return func(ctx *web.Context, val string) string {
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
	web.Get(apc.searchURL+"(.*)", search)
}

// Wrap a lonely string in an entire webpage
func (apc *ArchPageContents) Surround(s string) (string, string) {
	apc.contentHTML = s
	archpage := archbuilder(apc)
	return archpage.GetXML(true), archpage.GetCSS()
}

type WebHandle (func(ctx *web.Context, val string) string)
type StringFunction (func(string) string)
type SimpleWebHandle StringFunction
type APCgen (func() *ArchPageContents)

// Creates a handle from s string function
func (apc *ArchPageContents) GetHandle(fn StringFunction) WebHandle {
	return func(ctx *web.Context, val string) string {
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

// Set an IP adress and generate a confirmation page for it
func GenerateSetIP(state *State) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Empty value, IP not set"
		}
		state.ips.Store(val)
		return "OK, set IP to " + val
	}
}

// Get all the stored IP adresses and generate a page for it
func GenerateGetAllIPs(state *State) SimpleWebHandle {
	return func(val string) string {
		s := ""
		iplist, err := state.ips.GetAll()
		if err == nil {
			for _, val := range iplist {
				s += "IP: " + val + "<br />"
			}
		}
		return browserspeak.Message("IPs", s)
	}
}

// Get the last stored IP adress and generate a page for it
func GenerateGetLastIP(state *State) SimpleWebHandle {
	return func(val string) string {
		s := ""
		ip, err := state.ips.GetLast()
		if err == nil {
			s = "IP: " + ip
		}
		return s
	}
}

func Hello() string {
	msg := "Hi"
	return browserspeak.Message("Hello", msg)
}

func Publish(url, filename string) {
	web.Get(url, browserspeak.FILE(filename))
}

func ParamExample(ctx *web.Context) string {
	return fmt.Sprintf("%v\n", ctx.Params)
}

func genFavicon(filename string) {
	img := canvas.New()
	img.Blank(16, 16)
	img.SetStrokeColor("#005090")

	// All the lines and translations use relative coordinates

	// "\"
	img.SetStrokeWidth(2)
	img.Translate(8, 2)
	img.Line(3, 11)
	img.Translate(-8, -2)

	// "/"
	img.SetStrokeWidth(2)
	img.Translate(8, 2)
	img.Line(-6, 12)
	img.Translate(-8, -2)

	// "-"
	img.SetStrokeWidth(2)
	img.Translate(2, 10)
	img.Line(12, -2)

	img.Write(filename)
}

func NewRedisList(c redis.Conn, id string) *RedisList {
	var rl RedisList
	rl.c = c
	rl.id = id
	return &rl
}

func (rl *RedisList) Store(value string) error {
	_, err := rl.c.Do("RPUSH", rl.id, value)
	return err
}

func bytes2string(b []uint8) string {
	return bytes.NewBuffer(b).String()
}

func getString(bi []interface{}, i int) string {
	return bytes2string(bi[i].([]uint8))
}

func (rl *RedisList) GetAll() ([]string, error) {
	result, err := redis.Values(rl.c.Do("LRANGE", rl.id, "0", "-1"))
	strs := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		strs[i] = getString(result, i)
	}
	return strs, err
}

func (rl *RedisList) GetLast() (string, error) {
	result, err := redis.Values(rl.c.Do("LRANGE", rl.id, "-1", "-1"))
	if len(result) == 1 {
		return getString(result, 0), err
	}
	return "", err
}

func (rl *RedisList) DelAll() error {
	_, err := rl.c.Do("DEL", rl.id)
	return err
}

// TODO: Caching, login
func main() {

	// Connect to Redis
	client, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("ERROR: Can't connect to redis")
	}
	defer client.Close()

	// Create a RedisList for storing IP adresses
	ips := NewRedisList(client, "IPs")

	faviconFilename := "generated/img/favicon.ico"

	genFavicon(faviconFilename)

	// These common ones are missing!
	mime.AddExtensionType(".txt", "text/plain; charset=utf-8")
	mime.AddExtensionType(".ico", "image/x-icon")

	state := new(State)
	state.ips = ips

	pub("/", HiAPC)

	web.Get("/hello/(.*)", wrapHandle(HelloAPC, helloSF))
	pub("/counting", CountAPC)

	web.Get("/setip/(.*)", GenerateSetIP(state))
	web.Get("/getip/(.*)", GenerateGetLastIP(state))
	web.Get("/getallips/(.*)", GenerateGetAllIPs(state))

	Publish("/robots.txt", "static/various/robots.txt")
	Publish("/sitemap_index.xml", "static/various/sitemap_index.xml")
	Publish("/favicon.ico", faviconFilename)

	//pubTemplate("/yes", templateGenerator(), templateContent{a:2})

	web.Get("/error", browserspeak.Errorlog)
	web.Get("/errors", browserspeak.Errorlog)

	// honeypot?
	web.Get("/index.php", Hello)
	web.Get("/viewtopic.php", ParamExample)

	web.Get("/(.*)", browserspeak.NotFound)
	web.Run("0.0.0.0:3000")
}
