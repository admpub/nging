// Package rest is an extremely simple to use, lightweight, yet powerful REST Client.
//
// Motivation
//
// The Go http standard library is a great library, but it might sometimes be a bit too low level to use,
// and it doesn't offer features like fork-join requests for better performance, response caching based on headers,
// and the possibility to mockup responses.
//
// Features and Roadmap
//
// v0.1
//  * GET, POST, PUT, PATCH, DELETE, HEAD & OPTIONS HTTP verbs
//  * Dead simple, synchronous requests
//  * Automatic caching of transport object and hosts connections
//  * Response Caching, based on response headers (cache-control, last-modified, etag, expires)
//  * Local caching strategies: TTL, LRU & Max Byte Size.
//  * Mockups!
//  * Fork-Join request pattern, for sending many requests concurrently, getting better client performance.
//  * Async request pattern.
//  * Request Body can be `string`, `[]byte`, `struct` & `map`
//  * Automatic marshal and unmarshal for `JSON` and `XML` Content-Type. Default JSON.
//  * Full access to http.Response object.
//
// v0.2
//  * Retries
//  * BasicAuth
//  * UserAgent
//
// v0.3
//  * Connection usage metrics
//  * Response Time metrics
//  * Brotli Content-Encoding Responses support - http://google-opensource.blogspot.com.ar/2015/09/introducing-brotli-new-compression.html
//
// v0.4
//  * Custom Root Certificates and Client Certificates
//  * Testing +95%
//
// v0.5
//  * Plugable external caches like Memcached
//
// Caching
//
// Caching is done by two strategies working together: Time To Live (TTL) and
// Least Recently Used (LRU). Objects are inserted in the cache based on
// Response Headers. You can establish a maximum Memory Size for the cache
// and objects are flushed based on time expiration (TTL) or by hitting the maximum
// memory limit. In the last case, least accessed objects will be removed first.
//
// Examples
//
// Installation
//  go get github.com/go-loco/restful/rest
//
// Importing
//  import "github.com/go-loco/restful/rest"
//
// Simple GET
//  resp := rest.Get("https://api.restfulsite.com/resource")
//
// Simple POST
//  // Using a `string` as body
//  resp := rest.Post("https://api.restfulsite.com/resource", "Body")
//
//
// Simple POST, with Struct Body
//  type User struct {
//    Id   int    `json:"id"`
//    Name string `json:"name"`
//  }
//
//  body := new(User)
//  body.Id = 1
//  body.Name = "Hernan"
//
//  // body will be marshall as JSON
//  resp := rest.Post("https://api.restfulsite.com/resource/1", body)
//  fmt.Println(resp)
//
// Fork Join
//
// ForkJoin let you **fork** requests, and **wait** until all of them have return.
//
// Concurrent has methods for Get, Post, Put, Patch, Delete, Head & Options,
// with the almost the same API as the synchronous methods.
// The difference is that these methods return a FutureResponse, which holds a pointer to
// Response. Response inside FutureResponse is nil until request has finished.
//
//  var f [3]*rest.FutureResponse
//
//  // ForkJoin will send all requests concurrently
//  // and will wait until all requests have their correspondent responses
//  rest.ForkJoin(func(c *rest.Concurrent) {
//    f[0] = c.Get("https://api.restfulsite.com/resource/1")
//    f[1] = c.Get("https://api.restfulsite.com/resource/2")
//    f[2] = c.Get("https://api.restfulsite.com/resource/3")
//   })
//
//  for i := range f {
//    if f[i].Response().StatusCode == http.StatusOK {
//      fmt.Println(f[i].Response())
//    }
//  }
//
// Async
//
// Async let you make Restful requests in an **asynchronous** way, without blocking
// the go routine calling the Async function.
//
// Whenever the Response is ready, the **f** function will be called back.
//
//  // This won't be blocked.
//  rest.AsyncGet("https://api.restfulsite.com/user", func(r *rest.Response) {
//    if r.StatusCode == http.StatusOK {
//      fmt.Println(r)
//    }
//  })
//
//  // This will be printed first.
//  fmt.Println("print first")
//
// Defaults
// * Headers: keep-alive, Cache-Control: no-cache
// * Timeout: 2 seconds
// * ContentType: JSON (for body requests in POST, PUT and PATCH)
// * Cache: enable
// * Cache Size: 1GB
// * Idle Connections Per Host: 2 (the default of http.net package)
// * HTTP/2: automatic with Go 1.6
// * Gzip: automatic support for gzip responses
//
// RequestBuilder
//
// RequestBuilder gives you the power to go beyond defaults.
// It is possible to set up headers, timeout, baseURL, proxy, contentType, not to use
// cache, directly disabling timeout (in an explicit way), and setting max idle connections.
//  headers := make(http.Header)
//  headers.Add("myHeader", "myValue")
//
//  var rb = rest.RequestBuilder{
//    Headers:             headers,
//    Timeout:             200 * time.Millisecond,
//    BaseURL:             "https://baseURL",
//    Proxy:               "http://myproxy",
//    ContentType:         rest.JSON,
//    DisableCache:        false,
//    DisableTimeout:      false,
//    MaxIdleConnsPerHost: 10,
//  }
//
//  resp := rb.Get("/mypath")
//
// Mockups
//
// When using mockups, all requests will be sent to the mockup server.
// To activate the mockup *environment* you have two ways: using the flag -mock
//	go test -mock
//
// Or by programmatically starting the mockup server
// 	StartMockupServer()
// An example
//  myURL := "http://mytest.com/foo"
//
//  myHeaders := make(http.Header)
//  myHeaders.Add("Hello", "world")
//
//  mock := rest.Mock{
//    URL:          myURL,
//  	HTTPMethod:   http.MethodGet,
//  	ReqHeaders:   myHeaders,
//  	RespHTTPCode: http.StatusOK,
//  	RespBody:     "foo",
//  }
//
//  rest.AddMockups(&mock)
//
//  v := rest.Get(myURL)
package rest
