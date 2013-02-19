package main


import (
	"time"
	"strconv"

	"github.com/xyproto/browserspeak"
	"github.com/hoisie/web"
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

type APCgen (func() *ArchPageContents)

// TODO: Use strings like {{.Title}} instead and run it through the html.template system
func archbuilder(apc *ArchPageContents) *browserspeak.Page {
	page := browserspeak.NewHTML5Page(apc.title + " " + apc.subtitle)

	page.LinkToCSS(apc.generatedCSSurl)
	page.LinkToCSS(apc.extraCSSurl)
	page.LinkToFavicon(apc.faviconurl)

	AddHeader(page)
	AddBodyStyle(page, apc.bgImageURL)
	AddTitleSearchMenu(page, apc.title, apc.subtitle, apc.links, apc.searchButtonText, apc.searchURL)
	AddContent(page, apc.contentTitle, apc.contentHTML)
	AddFooter(page, apc.footerText)

	return page
}

func BaseAPC() *ArchPageContents {
	var apc ArchPageContents
	apc.generatedCSSurl = "/css/style.css"
	apc.extraCSSurl = "/css/extra.css"
	apc.faviconurl = "/favicon.ico"
	apc.bgImageFilename = "static/img/longbg.png"
	apc.bgImageURL = "/img/longbg.png"
	//apc.bgImageURL = "http://home.online.no/~hakrist/Elg%20i%20solnedgang%201024.JPG"
	//apc.bgImageURL = "http://upload.wikimedia.org/wikipedia/commons/1/1f/Theodor_Kittelsen,_Soria_Moria.jpg"
	apc.bgImageURL = "http://www.dks-hordaland.no/media/bilde/elgisolnedgang.jpg"
	apc.title = "Arch Linux"
	apc.subtitle = "no"
	apc.links = []string{"Overview:/", "Hello:/hello/world", "Count:/counting"}
	apc.contentTitle = "Hi"
	apc.contentHTML = "Hi there!"
	apc.searchButtonText = "Search"
	apc.searchURL = "/search"
	y := time.Now().Year()
	apc.footerText = "Alexander Rødseth &lt;rodseth@gmail.com&gt;, " + strconv.Itoa(y)
	return &apc
}

func HiAPC() *ArchPageContents {
	apc := BaseAPC()
	//apc.contentHTML = `This site is currently under construction.<br />You may wish to visit the <a href="https://bbs.archlinux.org/">Arch Linux Forum</a> in the mean time.<br /><br /><i>- Alexander Rødseth &lt;rodseth / gmail&gt;</i>`
	apc.contentTitle = "YOLO narwhal"
	apc.contentHTML = `
<p>Locavore Austin fanny pack pickled.  Marfa hoodie pitchfork american apparel, flexitarian YOLO pickled keytar twee cred craft beer seitan authentic raw denim kogi.  Selvage mixtape blog, pickled cosby sweater williamsburg skateboard brooklyn lo-fi twee.  Blue bottle echo park kale chips, selvage fap skateboard swag chambray tousled.  Street art etsy four loko fap, iphone carles cliche banh mi fashion axe PBR authentic leggings.  Narwhal mumblecore street art tumblr.  Messenger bag vice art party, next level aesthetic church-key tumblr direct trade  typewriter street art.</p><p>Messenger bag blue bottle VHS before they sold out.  Artisan pickled swag, VHS meggings jean shorts blog tonx salvia cosby sweater mumblecore aesthetic literally narwhal.  Brunch tofu gluten-free disrupt blog occupy.  Austin bicycle rights sartorial narwhal, butcher trust fund cred.  Neutra kale chips letterpress literally, williamsburg kogi brunch bicycle rights.  Williamsburg craft beer brunch quinoa, forage YOLO swag put a bird on it four loko mixtape banksy.  Tumblr semiotics yr fixie.</p><p>Iphone banksy wolf squid wayfarers, VHS photo booth banh mi fap.  Tonx flexitarian vinyl scenester terry richardson squid synth deep v.  VHS tousled godard, cardigan american apparel lo-fi flannel.  Vice church-key cliche, hashtag banh mi direct trade  skateboard.  Sriracha meh pitchfork, wayfarers helvetica leggings try-hard viral YOLO lo-fi fingerstache synth ennui next level ugh.  Wayfarers organic american apparel fingerstache craft beer bicycle rights, beard keffiyeh banksy four loko butcher hashtag mumblecore banjo wes anderson.  Williamsburg next level deep v pickled typewriter kogi.</p><p>Meggings gastropub flexitarian, before they sold out DIY wes anderson cred authentic artisan dreamcatcher aesthetic ennui food truck.  Fanny pack selvage synth vegan pug.  YOLO shoreditch pitchfork, letterpress whatever put a bird on it truffaut mumblecore flannel terry richardson irony cray master cleanse ethnic gluten-free.  Fap banksy blog pickled meh ethnic food truck +1, vice leggings retro quinoa.  Small batch vice pop-up mustache.  +1 ethnic echo park semiotics letterpress raw denim.  Keytar photo booth wes anderson, freegan before they sold out skateboard seitan brooklyn.</p><p>Wes anderson high life banksy messenger bag art party plaid disrupt tattooed, next level swag viral raw denim.  Cliche meggings terry richardson cray.  Next level 3 wolf moon retro marfa.  Pork belly authentic banjo, iphone lomo williamsburg letterpress cosby sweater Austin typewriter quinoa skateboard hoodie.  Plaid kale chips godard farm-to-table.  Fashion axe mixtape freegan, pop-up chambray ugh etsy YOLO jean shorts dreamcatcher meggings.  Banh mi letterpress tousled, skateboard stumptown high life vegan fap typewriter shoreditch 8-bit lo-fi master cleanse selfies bespoke.</p>
`
	return apc
}

func HelloAPC() *ArchPageContents {
	apc := BaseAPC()
	apc.contentTitle = "This is it"
	return apc
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

// Creates a handle from a string function
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

// Routing for the archlinux.no webpage
func ServeArchlinuxNo() {
	faviconFilename := "generated/img/favicon.ico"
	genFavicon(faviconFilename)

	pub("/", HiAPC)
	web.Get("/hello/(.*)", wrapHandle(HelloAPC, helloSF))
	pub("/counting", CountAPC)

	Publish("/robots.txt", "static/various/robots.txt")
	Publish("/sitemap_index.xml", "static/various/sitemap_index.xml")
	Publish("/favicon.ico", faviconFilename)

	//pubTemplate("/yes", templateGenerator(), templateContent{a:2})
}
