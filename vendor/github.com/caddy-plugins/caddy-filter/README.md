[![Go Report Card](https://goreportcard.com/badge/github.com/echocat/caddy-filter)](https://goreportcard.com/report/github.com/echocat/caddy-filter)
[![Build Status](https://travis-ci.org/echocat/caddy-filter.svg?branch=master)](https://travis-ci.org/echocat/caddy-filter)
[![Coverage Status](https://img.shields.io/coveralls/echocat/caddy-filter/master.svg?style=flat-square)](https://coveralls.io/github/echocat/caddy-filter?branch=master)
[![License](https://img.shields.io/github/license/echocat/caddy-filter.svg?style=flat-square)](LICENSE)

# caddy-filter

filter allows you to modify the responses.

This could be useful to modify static HTML files to add (for example) Google Analytics source code to it.

* [Syntax](#syntax)
* [Examples](#examples)
* [Run tests](#run-tests)
* [Contributing](#contributing)
* [License](#license)

## Syntax

```
filter rule {
    path                          <regexp pattern>
    content_type                  <regexp pattern>
    path_content_type_combination <and|or>
    search_pattern                <regexp pattern>
    replacement                   <replacement pattern>
}
filter rule ...
filter max_buffer_size    <maximum buffer size in bytes>
```

* **rule**: Defines a new filter rule for a file to respond.
    > **Important:** Define ``path`` and/or ``content_type`` not to open. Slack rules could dramatically impact the system performance because every response is recorded to memory before returning it.

    * **path**: Regular expression that matches the requested path.
    * **content_type**: Regular expression that matches the requested content type that results after the evaluation of the whole request.
    * **path_content_type_combination**: _(Since 0.8)_ Could be `and` or `or`. (Default: `and` - before this parameter existed it was `or`)
    * **search_pattern**: Regular expression to find in the response body to replace it.
    * **replacement**: Pattern to replace the ``search_pattern`` with. 
        <br>You can use parameters. Each parameter must be formatted like: ``{name}``.
        * Regex group: Every group of the ``search_pattern`` could be addressed with ``{index}``.
          <br>Example: ``"My name is (.*?) (.*?)." => "Name: {2}, {1}."``
        
        * Request context: Parameters like URL ... could be accessed.
          <br>Example: ``Host: {request_host}``
            * ``request_header_<header name>``: Contains a header value of the request, if provided or empty.
            * ``request_url``: Full requested url
            * ``request_path``: Requested path
            * ``request_method``: Used method
            * ``request_host``: Target host
            * ``request_proto``: Used proto
            * ``request_proto``: Used proto
            * ``request_remoteAddress``: Remote address of the calling client
            * ``response_header_<header name>``: Contains a header value of the response, if provided or empty.
            * ``env_<environment variable name>``: Contains an environment variable value, if provided or empty.
            * ``now[:<pattern>]``: Current timestamp. If pattern not provided, `RFC` or `RFC3339` [RFC3339](https://tools.ietf.org/html/rfc3339) is used. Other values: [`unix`](https://en.wikipedia.org/wiki/Unix_time), [`timestamp`](https://developer.mozilla.org/en/docs/Web/JavaScript/Reference/Global_Objects/Date/now) or free format following [Golang time formatting rules](https://golang.org/pkg/time/#pkg-constants).
            * ``response_header_last_modified[:<pattern>]``: Same like `now` for last modification time of current resource - see above. If not send by server current time will be used.
        * Replacements in files: If the replacement is prefixed with a ``@`` character it will be tried
           to find a file with this name and load the replacement from there. This will help you to also
           add replacements with larger payloads which will be ugly direct within the Caddyfile.
           <br>Example: ``@myfile.html``
* **max_buffer_size**: Limit the buffer size to the specified maximum number of bytes. If a rules matches the whole body will be recorded at first to memory before delivery to HTTP client. If this limit is reached no filtering will executed and the content is directly forwarded to the client to prevent memory overload. Default is: ``10485760`` (=10 MB)

## Examples

Replace in every text file ``Foo`` with ``Bar``.

```
filter rule {
    path .*\.txt
    search_pattern Foo
    replacement Bar
}
```

Add Google Analytics to every HTML page from a file.

**``Caddyfile``**:
```
filter rule {
    path .*\.html
    search_pattern </title>
    replacement @header.html
}
```

**``header.html``**:
```html
<script>(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)})(window,document,'script','//www.google-analytics.com/analytics.js','ga');ga('create', 'UA-12345678-9', 'auto');ga('send', 'pageview');</script>
</title>
```

Insert server name in every HTML page

```
filter rule {
    content_type text/html.*
    search_pattern Server
    replacement "This site was provided by {response_header_Server}"
}
```

## Run tests

### Full

This includes download of all dependencies and also creation and upload of coverage reports.

> No working golang installation is required but Java 8+ (in ``PATH`` or ``JAVA_HOME`` set.). 

```bash
# On Linux/macOS
$ ./gradlew test

# On Windows
$ gradlew test
```
### Native

> Requires a working golang installation in ``PATH`` and ``GOPATH`` set.

```bash
$ go test
```

## Contributing

caddy-filter is an open source project by [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/caddy-filter)
by [fork us](https://github.com/echocat/caddy-filter/fork).

If you commit code to this project, you have to accept that this code will be released under the [license](#license) of this project.

## License

See the [LICENSE](LICENSE) file.
