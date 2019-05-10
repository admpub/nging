package rest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var readVerbs = []string{http.MethodGet, http.MethodHead, http.MethodOptions}
var contentVerbs = []string{http.MethodPost, http.MethodPut, http.MethodPatch}

var maxAge = regexp.MustCompile(`(?:max-age|s-maxage)=(\d+)`)
var httpDateFormat = "Mon, 01 Jan 2006 15:04:05 GMT"

func (rb *RequestBuilder) doRequest(verb string, reqURL string, reqBody interface{}) (response *Response) {

	var cacheURL string
	var cacheResp *Response

	response = new(Response)
	reqURL = rb.BaseURL + reqURL

	//If Cache enable && operation is read: Cache GET
	if !rb.DisableCache && match(verb, readVerbs) {
		if cacheResp = resourceCache.get(reqURL); cacheResp != nil {
			cacheResp.cacheHit.Store(true)
			if !cacheResp.revalidate {
				return cacheResp
			}
		}
	}

	//Marshal request to JSON or XML
	body, err := rb.marshalReqBody(reqBody)
	if err != nil {
		response.Err = err
		return
	}

	// Change URL to point to Mockup server
	reqURL, cacheURL, err = checkMockup(reqURL)
	if err != nil {
		response.Err = err
		return
	}

	//Get Client (client + transport)
	client := rb.getClient()

	//Create request
	request, err := http.NewRequest(verb, reqURL, bytes.NewBuffer(body))
	if err != nil {
		response.Err = err
		return
	}

	// Set extra parameters
	rb.setParams(client, request, cacheResp, cacheURL)

	// Make the request
	httpResp, err := client.Do(request)
	if err != nil {
		response.Err = err
		return
	}

	// Read response
	defer httpResp.Body.Close()
	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		response.Err = err
		return
	}

	// If we get a 304, return response from cache
	if httpResp.StatusCode == http.StatusNotModified {
		response = cacheResp
		return
	}

	response.Response = httpResp
	response.byteBody = respBody

	ttl := setTTL(response)
	lastModified := setLastModified(response)
	etag := setETag(response)

	if !ttl && (lastModified || etag) {
		response.revalidate = true
	}

	//If Cache enable: Cache SETNX
	if !rb.DisableCache && match(verb, readVerbs) && (ttl || lastModified || etag) {
		resourceCache.setNX(cacheURL, response)
	}

	return
}

func checkMockup(reqURL string) (string, string, error) {

	cacheURL := reqURL

	if mockUpEnv {

		rURL, err := url.Parse(reqURL)
		if err != nil {
			return reqURL, cacheURL, err
		}

		rURL.Scheme = mockServerURL.Scheme
		rURL.Host = mockServerURL.Host

		return rURL.String(), cacheURL, nil
	}

	return reqURL, cacheURL, nil
}

func (rb *RequestBuilder) marshalReqBody(body interface{}) (b []byte, err error) {

	if body != nil {
		switch rb.ContentType {
		case JSON:
			b, err = json.Marshal(body)
		case XML:
			b, err = xml.Marshal(body)
		}
	}

	return
}

func (rb *RequestBuilder) getClient() *http.Client {

	// This will be executed only once
	// per request builder
	rb.clientMtxOnce.Do(func() {

		dTransportMtxOnce.Do(func() {

			if defaultTransport == nil {
				defaultTransport = &http.Transport{
					MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
					Proxy:               http.ProxyFromEnvironment}
			}
		})

		tr := defaultTransport

		if cp := rb.CustomPool; cp != nil {
			tr = &http.Transport{MaxIdleConnsPerHost: rb.CustomPool.MaxIdleConnsPerHost}

			//Set Proxy
			if cp.Proxy != "" {
				if proxy, err := url.Parse(cp.Proxy); err == nil {
					tr.Proxy = http.ProxyURL(proxy)
				}
			}

		}

		rb.client = &http.Client{Transport: tr}

		//Timeout
		rb.client.Timeout = rb.getTimeout()

	})
	//

	return rb.client
}

func (rb *RequestBuilder) getTimeout() time.Duration {

	switch {
	case rb.DisableTimeout:
		return 0
	case rb.Timeout > 0:
		return rb.Timeout
	default:
		return DefaultTimeout
	}

}

func (rb *RequestBuilder) setParams(client *http.Client, req *http.Request, cacheResp *Response, cacheURL string) {

	//Custom Headers
	if rb.Headers != nil {
		req.Header = rb.Headers
	}

	//Default headers
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")

	//If mockup
	if mockUpEnv {
		req.Header.Set("X-Original-URL", cacheURL)
	}

	// Basic Auth
	if rb.BasicAuth != nil {
		req.SetBasicAuth(rb.BasicAuth.UserName, rb.BasicAuth.Password)
	}

	// User Agent
	req.Header.Set("User-Agent", func() string {
		if rb.UserAgent != "" {
			return rb.UserAgent
		}
		return "github.com/admpub/restful"
	}())

	//Encoding
	var cType string

	switch rb.ContentType {
	case JSON:
		cType = "json"
	case XML:
		cType = "xml"
	}

	req.Header.Set("Accept", "application/"+cType)

	if match(req.Method, contentVerbs) {
		req.Header.Set("Content-Type", "application/"+cType)
	}

	if cacheResp != nil && cacheResp.revalidate {
		switch {
		case cacheResp.etag != "":
			req.Header.Set("If-None-Match", cacheResp.etag)
		case cacheResp.lastModified != nil:
			req.Header.Set("If-Modified-Since", cacheResp.lastModified.Format(httpDateFormat))
		}
	}

}

func match(s string, sarray []string) bool {

	for _, v := range sarray {
		if v == s {
			return true
		}
	}

	return false
}

func setTTL(resp *Response) (set bool) {

	now := time.Now()

	//Cache-Control Header
	cacheControl := maxAge.FindStringSubmatch(resp.Header.Get("Cache-Control"))

	if len(cacheControl) > 1 {

		ttl, err := strconv.Atoi(cacheControl[1])
		if err != nil {
			return
		}

		if ttl > 0 {
			t := now.Add(time.Duration(ttl) * time.Second)
			resp.ttl = &t
			set = true
		}

		return
	}

	//Expires Header
	//Date format from RFC-2616, Section 14.21
	expires, err := time.Parse(httpDateFormat, resp.Header.Get("Expires"))
	if err != nil {
		return
	}

	if expires.Sub(now) > 0 {
		resp.ttl = &expires
		set = true
	}

	return
}

func setLastModified(resp *Response) bool {
	lastModified, err := time.Parse(httpDateFormat, resp.Header.Get("Last-Modified"))
	if err != nil {
		return false
	}

	resp.lastModified = &lastModified
	return true
}

func setETag(resp *Response) bool {

	resp.etag = resp.Header.Get("ETag")

	if resp.etag != "" {
		return true
	}

	return false
}
