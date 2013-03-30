package main

import (
	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
	"github.com/xyproto/web"
)

// Ok, "archlinuxno"

//type TemplateValues map[string]string
//type TemplateValueGenerator func(*web.Context) TemplateValues
//type TemplateValueGeneratorFactory func(*UserState) TemplateValueGenerator

// Functions that generate functions that generate content that can be used in templates

func DynamicMenuFactory(state *UserState) TemplateValueGenerator {
	return func(*web.Context) TemplateValues {
		// Can generate the menu based on both the user state and the web context here
		return TemplateValues{"menu": "<div style='margin-left: 3em;'><a href='/login'>Login</a> | <a href='/register'>Register</a></div>"}
	}
}

