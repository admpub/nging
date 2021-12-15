package minify

import (
	min "github.com/tdewolff/minify"
	minCSS "github.com/tdewolff/minify/css"
	minHTML "github.com/tdewolff/minify/html"
	minJS "github.com/tdewolff/minify/js"
)

var Minify = min.New()

func init() {
	Minify.AddFunc(`text/javascript`, minJS.Minify)
	Minify.AddFunc(`text/css`, minCSS.Minify)
	Minify.AddFunc(`text/html`, minHTML.Minify)
}

func MinifyCSS2(input []byte) ([]byte, error) {
	return Minify.Bytes(`text/css`, input)
}

func MinifyJS2(input []byte) ([]byte, error) {
	return Minify.Bytes(`text/javascript`, input)
}

func MinifyHTML2(input []byte) ([]byte, error) {
	return Minify.Bytes(`text/html`, input)
}
