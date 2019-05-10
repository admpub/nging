package rest

var dfltBuilder = RequestBuilder{}

// Get issues a GET HTTP verb to the specified URL.
//
// In Restful, GET is used for "reading" or retrieving a resource.
// Client should expect a response status code of 200(OK) if resource exists,
// 404(Not Found) if it doesn't, or 400(Bad Request).
//
// Get uses the DefaultBuilder.
func Get(url string) *Response {
	return dfltBuilder.Get(url)
}

// Post issues a POST HTTP verb to the specified URL.
//
// In Restful, POST is used for "creating" a resource.
// Client should expect a response status code of 201(Created), 400(Bad Request),
// 404(Not Found), or 409(Conflict) if resource already exist.
//
// Body could be any of the form: string, []byte, struct & map.
//
// Post uses the DefaultBuilder.
func Post(url string, body interface{}) *Response {
	return dfltBuilder.Post(url, body)
}

// Put issues a PUT HTTP verb to the specified URL.
//
// In Restful, PUT is used for "updating" a resource.
// Client should expect a response status code of of 200(OK), 404(Not Found),
// or 400(Bad Request). 200(OK) could be also 204(No Content)
//
// Body could be any of the form: string, []byte, struct & map.
//
// Put uses the DefaultBuilder.
func Put(url string, body interface{}) *Response {
	return dfltBuilder.Put(url, body)
}

// Patch issues a PATCH HTTP verb to the specified URL
//
// In Restful, PATCH is used for "partially updating" a resource.
// Client should expect a response status code of of 200(OK), 404(Not Found),
// or 400(Bad Request). 200(OK) could be also 204(No Content)
//
// Body could be any of the form: string, []byte, struct & map.
//
// Patch uses the DefaultBuilder.
func Patch(url string, body interface{}) *Response {
	return dfltBuilder.Patch(url, body)
}

// Delete issues a DELETE HTTP verb to the specified URL
//
// In Restful, DELETE is used to "delete" a resource.
// Client should expect a response status code of of 200(OK), 404(Not Found),
// or 400(Bad Request).
//
// Delete uses the DefaultBuilder.
func Delete(url string) *Response {
	return dfltBuilder.Delete(url)
}

// Head issues a HEAD HTTP verb to the specified URL
//
// In Restful, HEAD is used to "read" a resource headers only.
// Client should expect a response status code of 200(OK) if resource exists,
// 404(Not Found) if it doesn't, or 400(Bad Request).
//
// Head uses the DefaultBuilder.
func Head(url string) *Response {
	return dfltBuilder.Head(url)
}

// Options issues a OPTIONS HTTP verb to the specified URL
//
// In Restful, OPTIONS is used to get information about the resource
// and supported HTTP verbs.
// Client should expect a response status code of 200(OK) if resource exists,
// 404(Not Found) if it doesn't, or 400(Bad Request).
//
// Options uses the DefaultBuilder.
func Options(url string) *Response {
	return dfltBuilder.Options(url)
}

// AsyncGet is the *asynchronous* option for GET.
// The go routine calling AsyncGet(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncGet uses the DefaultBuilder
func AsyncGet(url string, f func(*Response)) {
	dfltBuilder.AsyncGet(url, f)
}

// AsyncPost is the *asynchronous* option for POST.
// The go routine calling AsyncPost(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncPost uses the DefaultBuilder
func AsyncPost(url string, body interface{}, f func(*Response)) {
	dfltBuilder.AsyncPost(url, body, f)
}

// AsyncPut is the *asynchronous* option for PUT.
// The go routine calling AsyncPut(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncPut uses the DefaultBuilder
func AsyncPut(url string, body interface{}, f func(*Response)) {
	dfltBuilder.AsyncPut(url, body, f)
}

// AsyncPatch is the *asynchronous* option for PATCH.
// The go routine calling AsyncPatch(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncPatch uses the DefaultBuilder
func AsyncPatch(url string, body interface{}, f func(*Response)) {
	dfltBuilder.AsyncPatch(url, body, f)
}

// AsyncDelete is the *asynchronous* option for DELETE.
// The go routine calling AsyncDelete(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncDelete uses the DefaultBuilder
func AsyncDelete(url string, f func(*Response)) {
	dfltBuilder.AsyncDelete(url, f)
}

// AsyncHead is the *asynchronous* option for HEAD.
// The go routine calling AsyncHead(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncHead uses the DefaultBuilder
func AsyncHead(url string, f func(*Response)) {
	dfltBuilder.AsyncHead(url, f)
}

// AsyncOptions is the *asynchronous* option for OPTIONS.
// The go routine calling AsyncOptions(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
//
// AsyncOptions uses the DefaultBuilder
func AsyncOptions(url string, f func(*Response)) {
	dfltBuilder.AsyncOptions(url, f)
}

// ForkJoin let you *fork* requests, and *wait* until all of them have return.
//
// Concurrent has methods for Get, Post, Put, Patch, Delete, Head & Options,
// with the almost the same API as the synchronous methods.
// The difference is that these methods return a FutureResponse, which holds a pointer to
// Response. Response inside FutureResponse is nil until request has finished.
//
// 	var futureA, futureB *rest.FutureResponse
//
// 	rest.ForkJoin(func(c *rest.Concurrent){
//		futureA = c.Get("/url/1")
//		futureB = c.Get("/url/2")
//	})
//
//	fmt.Println(futureA.Response())
//	fmt.Println(futureB.Response())
//
// AsyncOptions uses the DefaultBuilder
func ForkJoin(f func(*Concurrent)) {
	dfltBuilder.ForkJoin(f)
}
