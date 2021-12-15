# cr

Simple wrapper for [chromedp](https://github.com/chromedp/chromedp).

Licence: BSD.

## Example usage:

```go
package main

import (
    "log"
    "time"

    "github.com/admpub/cr"
)

func main() {

    browser, err := cr.New()
    if err != nil {
        log.Fatalf("Unable to create browser: %s\n", err)
    }
    defer browser.Close()

    if err := browser.Navigate("https://google.com"); err != nil {
        log.Fatalf("Couldn't navigate to page: %s\n", err)
    }

    if err := browser.SendKeys("//input[@id='lst-ib']", "The Big Lebowski"); err != nil {
        log.Fatalf("Couldn't fill out form: %s\n", err)
    }

    if err = browser.Click("//input[@name='btnK']"); err != nil {
        log.Fatalf("Couldn't find the submit button: %s\n", err)
    }
    time.Sleep(time.Second * 15)

}
```

## Why XPath?

XPath seems to be the easiest way to find DOM elements. There is great documentation online, such as the [w3schools](https://www.w3schools.com/xml/xpath_intro.asp) tutorial, and Google searches for specific things tend to turn up good results.

It is also easy to find an XPath in Chrome by right-clicking on an element and selecting "Inspect."

Clicking on any element in the inspector will show three dots (...) that you can click on and select "Copy XPath" from the "Copy" menu.

Typing `CTRL+f` in the inspector will provide a text box which can be used to test XPath values before running your code.
