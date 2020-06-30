phantomjs [![godoc](https://godoc.org/github.com/admpub/phantomjs?status.svg)](https://godoc.org/github.com/admpub/phantomjs) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)
=========

This is a Go wrapper for the [`phantomjs`][phantomjs] command line program. It
provides the full `webpage` API and has a strongly typed API. The wrapper
provides an idiomatic Go interface while allowing you to communicate with the
underlying WebKit and JavaScript engine in a seamless way.

[phantomjs]: http://phantomjs.org/


## Installing

First, install `phantomjs` on your machine. This can be done using your package
manager (such as `apt-get` or `brew`). Then install this package using the Go
toolchain:

```sh
$ go get -u github.com/admpub/phantomjs
```


## Usage

### Starting the process

This wrapper works by communicating with a separate `phantomjs` process over
HTTP. The process can take several seconds to start up and shut down so you
should do that once and then share the process. There is a package-level
variable called `phantomjs.DefaultProcess` that exists for this purpose.

```go
package main

import (
	"github.com/admpub/phantomjs"
)

func main() {
	// Start the process once.
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer phantomjs.DefaultProcess.Close()

	// Do other stuff in your program.
	doStuff()
}
```


### Starting the process using proxy and additional arguments

This wrapper works by communicating with a separate `phantomjs` process over
HTTP using proxy and additional arguments. The process can take several seconds to start up and shut down so you
should do that once and then share the process. There is a package-level
variable called `phantomjs.DefaultProcess` that exists for this purpose.

```go
package main

import (
	"github.com/admpub/phantomjs"
)

func main() {
	// Start the process once.
	phantomjs.DefaultProcess = phantomjs.NewProcessUsingProxy("196.18.xxx.xxx:44","http","user:pass")
	phantomjs.DefaultProcess.AddArgs("--debug=true")
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer phantomjs.DefaultProcess.Close()

	// Do other stuff in your program.
	doStuff()
}
```


You can have multiple processes, however, you will need to change the port used
for each one so they do not conflict. This library uses port `20202` by default.


### Working with WebPage

The `WebPage` will be the primary object you work with in `phantomjs`. Typically
you will create a web page from a `Process` and then either open a URL or you
can set the content directly:

```go
p := phantomjs.DefaultProcess
// Create a web page.
// IMPORTANT: Always make sure you close your pages!
page, err := p.CreateWebPage()
if err != nil {
	return err
}
defer page.Close()

// Open a URL.
if err := page.Open("https://google.com"); err != nil {
	return err
}
```

The HTTP API uses a reference map to track references between the Go library
and the `phantomjs` process. Because of this, it is important to always
`Close()` your web pages or else you can experience memory leaks.



### Executing JavaScript

You can synchronously execute JavaScript within the context of a web page by
by using the `Evaluate()` function. This example below opens Hacker News,
retrieves the text and URL from the first link, and prints it to the terminal.

```go
// Open a URL.
if err := page.Open("https://news.ycombinator.com"); err != nil {
	return err
}

// Read first link.
info, err := page.Evaluate(`function() {
	var link = document.body.querySelector('.itemlist .title a');
	return { title: link.innerText, url: link.href };
}`)
if err != nil {
	return err
}

// Print title and URL.
link := info.(map[string]interface{})
fmt.Println("Hacker News Top Link:")
fmt.Println(link["title"])
fmt.Println(link["url"])
fmt.Println()
```

You can pass back any object from `Evaluate()` that can be marshaled over JSON.



### Rendering web pages

Another common task with PhantomJS is to render a web page to an image. Once
you have opened your web page, simply set the viewport size and call the
`Render()` method:

```go
// Open a URL.
if err := page.Open("https://news.ycombinator.com"); err != nil {
	return err
}

// Setup the viewport and render the results view.
if err := page.SetViewportSize(1024, 800); err != nil {
	return err
}
if err := page.Render("hackernews.png", "png", 100); err != nil {
	return err
}
```

You can also use the `RenderBase64()` to return a base64 encoded image to your
program instead of writing the file to disk.

