# server

A FTP server framework forked from [github.com/yob/graval](http://github.com/yob/graval) and changed a lot.

Full documentation for the package is available on [godoc](http://godoc.org/github.com/goftp/server)

## Version

    v0.2.1104

## Installation

    go get github.com/goftp/server

## Usage

To boot an FTP server you will need to provide a driver that speaks to your
persistence layer - the required driver contract is listed below.

There is a sample in-memory driver available as a demo. You can build it with
this command:

    go install github.com/goftp/ftpd

Then run it:

    ./bin/ftpd

And finally, connect to the server with any FTP client and the following
details:

    host: 127.0.0.1
    port: 2121
    username: test
    password: 1234

## Contributors

* Lunny Xiao <xiaolunwen@gmail.com>
* James Healy <james@yob.id.au> [http://www.yob.id.au](http://www.yob.id.au)

## Warning

FTP is an incredibly insecure protocol. Be careful about forcing users to authenticate
with a username or password that are important.

## License

This library is distributed under the terms of the MIT License. See the included file for
more detail.

## Contributing

All suggestions and patches welcome, preferably via a git repository I can pull from.
If this library proves useful to you, please let me know.

## Further Reading

There are a range of RFCs that together specify the FTP protocol. In chronological
order, the more useful ones are:

* [http://tools.ietf.org/rfc/rfc959.txt](http://tools.ietf.org/rfc/rfc959.txt)
* [http://tools.ietf.org/rfc/rfc1123.txt](http://tools.ietf.org/rfc/rfc1123.txt)
* [http://tools.ietf.org/rfc/rfc2228.txt](http://tools.ietf.org/rfc/rfc2228.txt)
* [http://tools.ietf.org/rfc/rfc2389.txt](http://tools.ietf.org/rfc/rfc2389.txt)
* [http://tools.ietf.org/rfc/rfc2428.txt](http://tools.ietf.org/rfc/rfc2428.txt)
* [http://tools.ietf.org/rfc/rfc3659.txt](http://tools.ietf.org/rfc/rfc3659.txt)
* [http://tools.ietf.org/rfc/rfc4217.txt](http://tools.ietf.org/rfc/rfc4217.txt)

For an english summary that's somewhat more legible than the RFCs, and provides
some commentary on what features are actually useful or relevant 24 years after
RFC959 was published:

* [http://cr.yp.to/ftp.html](http://cr.yp.to/ftp.html)

For a history lesson, check out Appendix III of RCF959. It lists the preceding
(obsolete) RFC documents that relate to file transfers, including the ye old
RFC114 from 1971, "A File Transfer Protocol"

This library is heavily based on [em-ftpd](https://github.com/yob/em-ftpd), an FTPd
framework with similar design goals within the ruby and EventMachine ecosystems. It
worked well enough, but you know, callbacks and event loops make me something
something.
