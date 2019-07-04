
cors gives you easy control over Cross Origin Resource Sharing for your site.

It allows you to whitelist particular domains per route, or to simply allow all domains `*` If desired you may customize nearly every aspect of the specification.

### Syntax

```
cors [path] [domains...] {
	origin            [origin]
	origin_regexp     [regexp]
	methods           [methods]
	allow_credentials [allowCredentials]
	max_age           [maxAge]
	allowed_headers   [allowedHeaders]
	exposed_headers   [exposedHeaders]
}
```

*   **path** is the file or directory this applies to (default is /).
*   **domains** is a space-seperated list of domains to allow. If ommitted, all domains will be granted access.
*   **origin** is a domain to grant access to. May be specified multiple times or ommitted.
*   **origin_regexp** is a regexp that will be matched to the `Origin` header. Access will be granted accordingly. It can be used in conjonction with the `origin` config (executed as a fallback to `origin`). May be specified multiple times or ommitted.
*   **methods** is set of http methods to allow. Default is these: POST,GET,OPTIONS,PUT,DELETE.
*   **allow_credentials** sets the value of the Access-Control-Allow-Credentials header. Can be true or false. By default, header will not be included.
*   **max_age** is the length of time in seconds to cache preflight info. Not set by default.
*   **allowed_headers** is a comma-seperated list of request headers a client may send.
*   **exposed_headers** is a comma-seperated list of response headers a client may access.

### Examples

Simply allow all domains to request any path:

<code class="block"><span class="hl-directive">cors</span></code>

Protect specific paths only, and only allow a few domains:

<code class="block"><span class="hl-directive">cors</span> <span class="hl-arg">/foo http://mysite.com http://anothertrustedsite.com</span></code>

Full configuration:

```
cors / {
  origin            http://allowedSite.com
  origin            http://anotherSite.org https://anotherSite.org
  origin_regexp     .+\.example\.com$
  methods           POST,PUT
  allow_credentials false
  max_age           3600
  allowed_headers   X-Custom-Header,X-Foobar
  exposed_headers   X-Something-Special,SomethingElse
}
```
