package main

// Functions for generating images

import (
	"github.com/gosexy/canvas"
)

// TODO: Don't write to a file, but return the image data
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


