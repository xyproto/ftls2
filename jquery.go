package main

// Creates an anonymous function
func fn(source string) string {
	return "function() { " + source + " }"
}

func quote(src string) string {
	return "\"" + src + "\""
}

// Creates an event
func event(tagname, event, source string) string {
	return "$(" + quote(tagname) + ")." + event + "(" + fn(source) + ");"
}

// Call a method on a tag
func method(tagname, methodname, value string) string {
	return "$(" + quote(tagname) + ")." + methodname + "(" + value + ");"
}

// Call a method on a tag that takes a string as a parameter
func methodString(tagname, methodname, value string) string {
	return method(tagname, methodname, quote(value))
}

// Run code when the document is ready
func OnDocumentReady(source string) string {
	return "$(document).ready(" + fn(source) + ");"
}

// Display an intruding message
func Alert(msg string) string {
	return "alert(" + quote(msg) + ");"
}

// When a tag is clicked at
func OnClick(tagname, source string) string {
	return event(tagname, "click", source)
}

func SetText(tagname, text string) string {
	return methodString(tagname, "text", text)
}

func SetHTML(tagname, html string) string {
	return method(tagname, "html", html)
}

func SetValue(tagname, val string) string {
	return methodString(tagname, "val", val)
}

func SetRawValue(tagname, val string) string {
	return method(tagname, "val", val)
}

func Hide(tagname string) string {
	return "$(" + quote(tagname) + ").hide();"
}

func Show(tagname string) string {
	return "$(" + quote(tagname) + ").show();"
}

func Load(tagname, url string) string {
	return methodString(tagname, "load", url)
}

// This never worked?
//func SetTextFromURL(tagname, url string) string {
//	//return "$.ajax({url: " + quote(url) + "}).done( function(data) {" + SetRawValue(tagname, "data") + "});"
//	return "$.get(" + quote(url) + ", function(data) {" + SetHTML(tagname, "data") + "});"
//}

// This works
// TODO: Use this to check the result of get and toggle menu items depending on the situation
func GetTest() string {
	return "$.get(\"/status/bob\", function(data) { alert(\"Data Loaded: \" + data); });"
}

// Returns html to run javascript once the document is ready
// Returns an empty string if there is no javascript to run.
func BodyJS(source string) string {
	if source != "" {
		return "<script>" + OnDocumentReady(source) + "</script>"
	}
	return ""
}
