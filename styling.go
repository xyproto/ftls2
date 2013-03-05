package main

// Various functions that can be used to style a webpage

import (
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

// TODO: get style values from a file

func AddHeader(page *Page, js string) {
	page.MetaCharset("UTF-8")
	page.LinkToGoogleFont("Armata")
	page.LinkToGoogleFont("Junge")
	// TODO: Move to browserspeak
	AddScriptToHeader(page, js)
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
	//body.AddStyle("background-size", "cover")
	//body.AddStyle("background-color", "#808080")
	//body.AddStyle("background-size", "100% 100%")
	////body.AddStyle("background-repeat", "no-repeat")
	////body.RepeatBackground(bgimageurl, "repeat-x")
	//page.SetColor("gray", "#a0e0e0") // gray text, turquise background color
	//page.SetColor("gray", "#202020") // gray text, turquise background color
	//page.SetColor("gray", "#202020") // gray text, turquise background color
}

func GenerateExtraCSS(stretchBackground bool) SimpleContextHandle {
	return func(ctx *web.Context) string {
		ctx.ContentType("css")
		// extra.css, loaded after the other CSS
		menucolor := "#c0c0c0"   // light gray
		hovercolor := "#efefe0"  // very light gray, with some yellow
		activecolor := "#ffffff" // white
		retval := `
a {
  text-decoration: none;
  color: #303030;
  font-weight: regular;
}
a:link {color:` + menucolor + `;}
a:visited {color:` + menucolor + `;}
a:hover {color:` + hovercolor + `;}
a:active {color:` + activecolor + `;}
`
		// The load order of background-color, background-size and background-image
		// is actually significant in Chrome! Do not reorder lightly!
		if stretchBackground {
			retval = "body {\nbackground-color: #808080;\nbackground-size: cover;\n}\n" + retval
		} else {
			retval = "body {\nbackground-color: #808080;\n}\n" + retval
		}
		return retval
	}
}
