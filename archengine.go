package main

// OK, only archlinux.no stuff, 23-03-13

// Move to "archlinuxno" once it has settled

import (
	"strconv"
	"time"

	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
	"github.com/xyproto/web"
)

// The default settings for Arch Linux content pages
func ArchBaseCP(state *UserState) *ContentPage {
	cp := DefaultCP(state)
	cp.BgImageURL = "/img/norway4.jpg"
	cp.StretchBackground = true
	cp.Title = "Arch Linux"
	cp.Subtitle = "no"

	//cp.links = []string{"Overview:/", "Mirrors:/mirrors", "Login:/login", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//cp.links = []string{"Overview:/", "Text:/text", "Bob:/bob", "JQuery:/jquery", "Register:/register", "Hello:/hello/world", "Count:/counting", "Feedback:/feedback"}
	//cp.Links = append(cp.Links, "Sample text:/text", "Chat:/chat")

	// The default links contains things like Login and Logout
	cp.Links = append(cp.Links, "Chat:/chat")

	y := time.Now().Year()

	//cp.footerText = "Alexander Rødseth &lt;rodseth@gmail.com&gt;, " + strconv.Itoa(y)
	cp.FooterText = "Alexander Rødseth, " + strconv.Itoa(y)

	// Hide and show the correct menus depending on the /showmenu/ pages for the various parts
	//cp.HeaderJS += UserMenuJS()
	//cp.HeaderJS += AdminMenuJS()
	//cp.HeaderJS += ChatMenuJS()

	// Additional hidden menus (by default)
	//cp.HiddenMenuIDs = append(cp.HiddenMenuIDs, "menuChat")
	//cp.HiddenMenuIDs = append(cp.HiddenMenuIDs, "menuAdmin")

	cp.Url = "/" // Is replaced when the contentpage is published

	cp.ColorScheme = NewArchColorScheme()

	return cp
}

//// Returns a ArchBaseCP with the contentTitle set
func ArchBaseTitleCP(contentTitle string, userState *UserState) *ContentPage {
	cp := ArchBaseCP(userState)
	cp.ContentTitle = contentTitle
	return cp
}

func OverviewCP(userState *UserState, url string) *ContentPage {
	cp := ArchBaseCP(userState)
	cp.ContentTitle = "Overview"
	cp.ContentHTML = `This site is currently under construction.<br />Visit the <a href="https://bbs.archlinux.org/viewtopic.php?id=4998">Arch Linux Forum</a> in the meantime.<br /><br /><i>- Alexander Rødseth &lt;rodseth / gmail&gt;</i>`
	cp.Url = url

	// TODO: !!!!!!!!!!!!!!!!!!!!!
	// TODO: Add CSS to hide Logout, see userengine
	// TODO: !!!!!!!!!!!!!!!!!!!!!

	return cp
}

func MirrorsCP(userState *UserState, url string) *ContentPage {
	cp := ArchBaseCP(userState)
	cp.ContentTitle = "Mirrors"
	cp.ContentHTML = "List over Norwegian Arch Linux mirrors:"
	cp.Url = url
	return cp
}

func PublishArchImages() {
	//faviconFilename := "/static/generated/img/favicon.ico"
	//genFavicon(faviconFilename)
	//Publish("/favicon.ico", faviconFilename, false)
	//Publish("/favicon.ico", "static/img/favicon.ico", false)

	// Tried previously:
	// "rough.png", "longbg.png", "donutbg.png", "donutbg_light.jpg",
	// "felix_predator2.jpg", "centerimage.png", "underwater.png",
	// "norway.jpg", "norway2.jpg", "underwater.jpg"

	// Publish and cache images
	imgs := []string{"norway4.jpg", "norway3.jpg", "gray.jpg", "darkgray.jpg"}
	for _, img := range imgs {
		Publish("/img/"+img, "static/img/"+img, true)
	}
}

func CountCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "Counting"
	apc.ContentHTML = "1 2 3"
	apc.Url = url
	return apc
}

// TODO: Find out why this only happens once the server starts
// and not every time the page reloads. Probably have to use
// more functions in functions. Try to use the model from sitespecific and ipspecific!
// That works fairly well.
func BobCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "Bob"
	if userState.HasUser("bob") {
		apc.ContentHTML = "has bob, l "
	} else {
		apc.ContentHTML = "no bob, l "
	}
	if userState.IsLoggedIn("bob") {
		apc.ContentHTML += "yes"
	} else {
		apc.ContentHTML += "no"
	}
	apc.Url = url
	return apc
}

func JQueryCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "JQuery"

	apc.ContentHTML = "<button id=clickme>bob</button><br />"
	apc.ContentHTML += "<div id=status>status</div>"

	//apc.contentJS = OnClick("#clickme", GetTest())
	//apc.contentJS += OnClick("#clickme", SetText("#clickme", "ost"))
	//apc.contentJS += OnClick("#clickme", SetTextFromURL("#clickme", "http://archlinux.no/status/bob"))
	//apc.contentJS += OnClick("#clickme", GetTest())

	apc.ContentJS += Load("#status", "/status/elg")
	apc.ContentJS += OnClick("#clickme", Load("#status", "/status/bob"))
	apc.ContentJS += SetText("#menuJQuery", "Heppa")

	apc.Url = url

	return apc
}

func TextCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "YOLO narwhal"
	apc.ContentHTML = `<p>Locavore Austin fanny pack pickled.  Marfa hoodie pitchfork american apparel, flexitarian YOLO pickled keytar twee cred craft beer seitan authentic raw denim kogi.  Selvage mixtape blog, pickled cosby sweater williamsburg skateboard brooklyn lo-fi twee.  Blue bottle echo park kale chips, selvage fap skateboard swag chambray tousled.  Street art etsy four loko fap, iphone carles cliche banh mi fashion axe PBR authentic leggings.  Narwhal mumblecore street art tumblr.  Messenger bag vice art party, next level aesthetic church-key tumblr direct trade  typewriter street art.</p><p>Messenger bag blue bottle VHS before they sold out.  Artisan pickled swag, VHS meggings jean shorts blog tonx salvia cosby sweater mumblecore aesthetic literally narwhal.  Brunch tofu gluten-free disrupt blog occupy.  Austin bicycle rights sartorial narwhal, butcher trust fund cred.  Neutra kale chips letterpress literally, williamsburg kogi brunch bicycle rights.  Williamsburg craft beer brunch quinoa, forage YOLO swag put a bird on it four loko mixtape banksy.  Tumblr semiotics yr fixie.</p><p>Iphone banksy wolf squid wayfarers, VHS photo booth banh mi fap.  Tonx flexitarian vinyl scenester terry richardson squid synth deep v.  VHS tousled godard, cardigan american apparel lo-fi flannel.  Vice church-key cliche, hashtag banh mi direct trade  skateboard.  Sriracha meh pitchfork, wayfarers helvetica leggings try-hard viral YOLO lo-fi fingerstache synth ennui next level ugh.  Wayfarers organic american apparel fingerstache craft beer bicycle rights, beard keffiyeh banksy four loko butcher hashtag mumblecore banjo wes anderson.  Williamsburg next level deep v pickled typewriter kogi.</p><p>Meggings gastropub flexitarian, before they sold out DIY wes anderson cred authentic artisan dreamcatcher aesthetic ennui food truck.  Fanny pack selvage synth vegan pug.  YOLO shoreditch pitchfork, letterpress whatever put a bird on it truffaut mumblecore flannel terry richardson irony cray master cleanse ethnic gluten-free.  Fap banksy blog pickled meh ethnic food truck +1, vice leggings retro quinoa.  Small batch vice pop-up mustache.  +1 ethnic echo park semiotics letterpress raw denim.  Keytar photo booth wes anderson, freegan before they sold out skateboard seitan brooklyn.</p><p>Wes anderson high life banksy messenger bag art party plaid disrupt tattooed, next level swag viral raw denim.  Cliche meggings terry richardson cray.  Next level 3 wolf moon retro marfa.  Pork belly authentic banjo, iphone lomo williamsburg letterpress cosby sweater Austin typewriter quinoa skateboard hoodie.  Plaid kale chips godard farm-to-table.  Fashion axe mixtape freegan, pop-up chambray ugh etsy YOLO jean shorts dreamcatcher meggings.  Banh mi letterpress tousled, skateboard stumptown high life vegan fap typewriter shoreditch 8-bit lo-fi master cleanse selfies bespoke.</p>`
	apc.Url = url
	return apc
}

func HelloCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "This is it"
	apc.Url = url
	return apc
}

// Routing for the archlinux.no webpage
// Admin, search and user management is already provided
func ServeArchlinuxNo(userState *UserState) {
	cps := []ContentPage{
		*OverviewCP(userState, "/"),
		*TextCP(userState, "/text"),
		*JQueryCP(userState, "/jquery"),
		*BobCP(userState, "/bob"),
		*CountCP(userState, "/counting"),
		*MirrorsCP(userState, "/mirrors"),
		*HelloCP(userState, "/feedback"),
	}

	// template content
	tp := Kake()

	ServeSite(ArchBaseCP, userState, cps, tp)

	// "dynamic" pages
	// Makes helloSF handle the content for /hello/(.*) urls, but wrapped in a BaseCP with the title "Hello"
	web.Get("/hello/(.*)", ArchBaseTitleCP("Hello", userState).WrapSimpleWebHandle(helloSF, Kake()))

	// static images
	PublishArchImages()
}

func NewArchColorScheme() *ColorScheme {
	var cs ColorScheme
	cs.Darkgray = "#202020"
	cs.Nicecolor = "#5080D0"   // nice blue
	cs.Menu_link = "#c0c0c0"   // light gray
	cs.Menu_hover = "#efefe0"  // light gray, somewhat yellow
	cs.Menu_active = "#ffffff" // white
	cs.Default_background = "#000030"
	return &cs
}
