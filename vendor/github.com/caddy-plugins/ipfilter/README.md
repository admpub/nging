# ipfilter
[![Go Report Card](https://goreportcard.com/badge/pyed/ipfilter)](https://goreportcard.com/report/pyed/ipfilter)

This is middleware for the [Caddy](http://caddyserver.com)
web server that implements black and whitelisting based on
IP addresses (or CIDR ranges) or country of origin using a
[MaxMind](https://dev.maxmind.com/geoip/geoip2/geolite2/) database.

## Syntax

```
ipfilter <basepath> {
    rule       <block | allow>
    ip         <addresses or CIDR ranges to block>
    prefix_dir <IP addr directory prefix>
    database   </path/to/GeoLite2-Country.mmdb>
    country    <ISO two letter country codes>
    blockpage  <blockpage.html>
    strict
}
```

You can specify zero or more `ipfilter` blocks. Each `ipfilter` block has
to specify at least one `ip`, `prefix_dir` or `country` directive. If no
`ipfilter` blocks are defined this middleware will allow every request.

* **basepath**: A sequence of URI path prefixes to match for the filter
to be active. You have to specify at least one path prefix. Use `/` to
match every request. If the request doesn't match one of these prefixes
the filter is ignored for purposes of determining if the request is
blocked or allowed.

* **rule**: Should the filter `block` (blacklist) or `allow` (whitelist)
the addresses. This directive is mandatory. It is an error to use it more
than once per ipfilter block. The **rule** in effect for the last `ipfilter`
block to match a request determines if it is blocked or allowed.

  Note that if you only have `ipfilter` blocks that specify `rule allow`
  then any request which doesn't match those filters will be implicitly
  blocked.

* **ip**: A sequence of IP adddresses or CIDR ranges to match. For example,
`ip 1.2.3.4 192.168.0.0/24` This is optional. It can be used more than
once in each `ipfilter` block rather than enumerating all IPs after a single
`ip` directive.

* **prefix_dir**: Specifies a directory in which to search for file names
matching the IP address of the request. This is optional. It is an error
to use this more than once per `ipfilter` block.

  You can specify a relative pathname to place it relative to the Caddy
  server CWD (which should be the content root dir).  When putting the
  blacklisted directory in the web server document tree you should also add
  an `internal` directive to ensure those files are not visible via HTTP
  GET requests. For example, `internal /blacklist/`. You can also specify
  an absolute pathname to locate the blacklist directory outside the
  document tree. And the path can include environment vars. For example,
  `prefix_dir {$HOME}/etc/www/blacklist`.

  You can create the file in the root of the blacklist directory. This is
  known as using a "flat" namespace. For example, *blacklist/127.0.0.1*
  or *blacklist/2601:647:4601:fa93:1865:4b6c:d055:3f3*. However,
  putting thousands of files in a single directory may cause
  poor performance of the lookup function. So you can also,
  and should, use a "sharded" namespace. This involves creating
  the file in a subdirectory based on the first two components
  of the address. For example, *blacklist/127/0/127.0.0.1* or
  *blacklist/2601/647/2601:647:4601:fa93:1865:4b6c:d055:3f3*.

  **Note:** IPv6 addresses as file names can use
  colons or equal-signs to separate the components; e.g.,
  *blacklist/2601/647/2601=647=4601=fa93==3f3*. Using equal-signs in
  place of colons in the file name may be necessary on platforms like MS
  Windows which assign special meaning to colons in file names. You have
  to use one or the other; you cannot mix them in the same file name.

  Note that you can also whitelist IP addresses using this mechanism
  by specifying `rule allow`. This may be useful when it follows a more
  general blocking rule (e.g., by country) and you want to selectively
  allow some addresses through but don't want to hardcode the addresses
  in the Caddy config file.

  This mechanism is most useful when coupled with automated monitoring of
  your web server activity to detect signals that your server is under
  attack from malware. All your monitoring software has to do is create
  a file in the blacklist directory.

  At this time the content of the file is ignored. In the future the
  contents will probably be read and exposed as a placeholder variable
  for use in conjuction with a template to be filled in via the `markdown`
  directive. So you should consider putting some explanatory text in the
  file explaining why the address was blocked.

* **database**: Specifies the path to a
[MaxMind](https://dev.maxmind.com/geoip/geoip2/geolite2/) database. This
is required if using the **country** directive; otherwise it should
be omitted.

* **country**: A whitespace separated sequence of ISO two letter country
codes to filter. This is optional but if used also requires a **database**
directive. Note that if a country could not be found for the address it
will be the empty string. This can be specified more than once per block
rather than enumerating all countries on a single line.

* **blockpage**: Names the file to be returned if the ipfilter
matches. Note that a `http.StatusOK` (200) status is returned if the
page is successfully returned to the client. This is optional. If not
specified then a `http.StatusForbidden` (403) status is returned.

* **strict**: Use this to disallow use of the address in the
`X-Forwarded-For` request header if any. This is optional and defaults
to false. If true or there is no `X-Forwarded-For` header use the address
from the request remote address.

## Caddyfile examples

#### Filter clients based on a given IP or range of IPs

```
ipfilter / {
	rule block
	ip 70.1.128.0/19 2001:db8::/122 9.12.20.16
}
```
`caddy` will block any clients with IPs that fall into one of these two ranges `70.1.128.0/19` and `2001:db8::/122` , or a client that has an IP of `9.12.20.16` explicitly.

```
ipfilter / {
	rule allow
	blockpage default.html
	ip 55.3.4.20 2e80::20:f8ff:fe31:77cf
}
```
`caddy` will serve only these 2 IPs, eveyone else will get `default.html`

```
ipfilter / {
	rule block
	prefix_dir blacklisted
}
```
`caddy` will block any client IP that appears as a file name in the
*blacklisted* directory.

#### Filter clients based on their [Country ISO Code](https://en.wikipedia.org/wiki/ISO_3166-1#Current_codes)

filtering with country codes requires a local copy of the Geo database, can be downloaded for free from [MaxMind](https://dev.maxmind.com/geoip/geoip2/geolite2/)
```
ipfilter / {
	rule allow
	database /data/GeoLite.mmdb
	country US JP
}
```
with that in your `Caddyfile` caddy will only serve users from the `United States` or `Japan`

```
ipfilter /notglobal /secret {
	rule block
	database /data/GeoLite.mmdb
	blockpage default.html
	country US JP
}
```
having that in your `Caddyfile` caddy will ignore any requests from `United States` or `Japan` to `/notglobal` or `/secret` and it will show `default.html` instead, `blockpage` is optional.

#### Using mutiple `ipfilter` blocks

The `ipfilter` blocks are evaluated for each HTTP request in the order they
appear. The last rule which matches a request is used to decide if the request
is allowed. So in general you will want more general rules (e.g., blacklist an
entire country) to appear before more specific rules (e.g., to whitelist
specific address ranges).

```
ipfilter / {
	rule allow
	ip 32.55.3.10
}

ipfilter /webhook {
	rule allow
	ip 192.168.1.0/24
}
```
You can use as many `ipfilter` blocks as you please, the above says: block everyone but `32.55.3.10`, Unless it falls in `192.168.1.0/24` and requesting a path in `/webhook`. Note that this is slightly subtle. Any request doesn't match any of those filters is implicitly blocked. In other words, there is no need to explicitly block every  address followed by "allow" filters like those above.

## Backward compatibility

`ipfilter` supports [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing). This is the recommended way of specifiying ranges. The old formats of ranging over IPs will get converted to CIDR via [range2CIDRs](https://github.com/pyed/ipfilter/blob/master/range2CIDRs.go) for the purpose of backward compatibility.
