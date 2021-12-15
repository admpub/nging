package cr

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// ErrNotFound is returned when an XPATH is provided
// for a DOM element, but it can not be located.
var ErrNotFound = errors.New("element not found")

const minTimeout = time.Second

// Browser represents a Chrome browser controlled by chromedp.
type Browser struct {
	ctx       context.Context
	cancelCtx context.CancelFunc
	timeout   time.Duration
	taskCtx   context.Context
	logger    *log.Logger
}

// New instantiates a new Chrome browser and returns
// a *Browser used to control it.
func New(ctx context.Context, args ...chromedp.ExecAllocatorOption) (*Browser, error) {
	b := &Browser{
		timeout: time.Second * 30,
		logger:  log.GetLogger(`ChromeDP`),
	}
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Headless,
	)
	for _, option := range args {
		options = append(options, option)
	}

	allocCtx, _ := chromedp.NewExecAllocator(ctx, options...)

	// also set up a custom logger
	taskCtx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(b.logger.Errorf))

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		cancel()
		return b, err
	}
	b.ctx = taskCtx
	b.taskCtx = taskCtx
	b.cancelCtx = cancel

	return b, nil
}

// SetTimeout accepts a time.Duration. This duration will
// be used as the maximum timeout when waiting for a node to exist.
func (b *Browser) SetTimeout(d time.Duration) {
	if d < minTimeout {
		d = minTimeout
	}
	b.timeout = d
}

func (b *Browser) Context() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(b.taskCtx, b.timeout)
	b.ctx, _ = chromedp.NewContext(ctx)
	return b.ctx, cancel
}

// Close cleans up the *Browser; this should be called
// on every *Browser once its work is complete.
func (b *Browser) Close() error {
	b.cancelCtx()
	return nil
}

// RunAction run single action
func (b *Browser) RunAction(action chromedp.Action) error {
	return chromedp.Run(b.ctx, action)
}

// RunTasks run mutiple action
func (b *Browser) RunTasks(actions ...chromedp.Action) error {
	return chromedp.Run(b.ctx, chromedp.Tasks(actions))
}

// RunTaskWithOther run mutiple action
func (b *Browser) RunTaskWithOther(action chromedp.Action, otherActions ...chromedp.Action) error {
	actions := append([]chromedp.Action{action}, otherActions...)
	return chromedp.Run(b.ctx, chromedp.Tasks(actions))
}

// Navigate sends the browser to a URL.
func (b *Browser) Navigate(url string, otherActions ...chromedp.Action) error {
	return b.RunTaskWithOther(chromedp.Navigate(url), otherActions...)
}

// MustNavigate calls Navigate and ends execution on error.
func (b *Browser) MustNavigate(url string, otherActions ...chromedp.Action) {
	if err := b.Navigate(url, otherActions...); err != nil {
		log.Fatalf("Failed to navigate to %q: %s\n", url, err)
	}
}

// Location returns the current URL.
func (b *Browser) Location(otherActions ...chromedp.Action) (string, error) {
	var location string
	err := b.RunTaskWithOther(chromedp.Location(&location), otherActions...)
	return location, err
}

// SendKeys sends keystrokes to a DOM element.
func (b *Browser) SendKeys(xpath, value string) error {
	return chromedp.Run(b.ctx, chromedp.SendKeys(xpath, value))
}

// MustSendKeys sends keystrokes to a DOM element or halts execution.
func (b *Browser) MustSendKeys(xpath, value string) {
	if err := b.SendKeys(xpath, value); err != nil {
		log.Fatalf("Failed to send %q to %q: %s\n", value, xpath, err)
	}
}

// Click performs a mouse click on a DOM element.
func (b *Browser) Click(xpath string) error {
	return chromedp.Run(b.ctx, chromedp.Click(xpath))
}

// MustClick performs a mouse click or ends the program.
func (b *Browser) MustClick(xpath string) {
	if err := b.Click(xpath); err != nil {
		log.Fatalf("Failed to click %q: %s\n", xpath, err)
	}
}

// GetSource returns the HTML source from the browser tab.
func (b *Browser) GetSource() (string, error) {
	var html string
	err := chromedp.Run(b.ctx, chromedp.OuterHTML("html", &html))
	return html, err
}

// GetAttributes returns the HTML attributes of a DOM element.
func (b *Browser) GetAttributes(xpath string) (map[string]string, error) {
	attrs := make(map[string]string)
	err := chromedp.Run(b.ctx, chromedp.Attributes(xpath, &attrs))
	return attrs, err
}

// ClickByXY clicks the browser window in a specific location.
func (b *Browser) ClickByXY(xpath string) error {
	x, y, err := b.GetTopLeft(xpath)
	if err != nil {
		return err
	}
	return chromedp.Run(b.ctx, chromedp.MouseClickXY(x, y))
}

// GetTopLeft returns the x, y coordinates of a DOM element.
func (b *Browser) GetTopLeft(xpath string) (float64, float64, error) {
	var top, left float64
	js := fmt.Sprintf(topLeftJS, xpath)
	var result string
	err := chromedp.Run(b.ctx, chromedp.Evaluate(js, &result))
	parts := strings.Split(result, ":")
	if len(parts) == 2 {
		top, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			b.logger.Errorf("Failed to parse top coordinate: %s", err)
			return 0, 0, err
		}
		left, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			b.logger.Errorf("Failed to parse left coordinate: %s", err)
			return 0, 0, err
		}
	}
	return top + 1, left + 1, err
}

func (b *Browser) ElementScreenshot(urlStr string, selectionElem string, by ...func(s *chromedp.Selector)) ([]byte, error) {
	byType := chromedp.ByID
	if len(by) > 0 {
		byType = by[0]
	}
	var buf []byte
	err := chromedp.Run(b.ctx, chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(selectionElem, byType),
		chromedp.Screenshot(selectionElem, &buf, chromedp.NodeVisible, byType),
	})
	return buf, err
}

func (b *Browser) Screenshot(urlStr string, quality int64) ([]byte, error) {
	var buf []byte
	err := chromedp.Run(b.ctx, fullScreenshot(urlStr, quality, &buf))
	return buf, err
}

// FindElement attempts to locate a DOM element.
func (b *Browser) FindElement(xpath string) error {
	nodes, err := b.GetNodes(xpath)
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		return nil
	}
	return ErrNotFound
}

// GetNodes returns a slice of *chromedp.Node from the chromedp package.
func (b *Browser) GetNodes(xpath string) ([]*cdp.Node, error) {
	var nodes []*cdp.Node
	ctx, cancel := context.WithTimeout(b.ctx, b.timeout)
	defer cancel()
	err := chromedp.Run(ctx, chromedp.Nodes(xpath, &nodes))
	return nodes, err
}

var topLeftJS = `
	function getTopLeft() {
		var element = document.evaluate("%s",document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null ).singleNodeValue;
		var rect = element.getBoundingClientRect();
		return rect.top + ":" + rect.left;
	}
	(function main() {
	   return getTopLeft();
	})(); 
	`
