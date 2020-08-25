<p align="center">
  <a href="https://github.com/admpub/xgo/releases/latest"><img src="https://img.shields.io/github/release/admpub/xgo.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/admpub/xgo/releases/latest"><img src="https://img.shields.io/github/downloads/admpub/xgo/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://hub.docker.com/r/admpub/xgo/"><img src="https://img.shields.io/docker/stars/admpub/xgo.svg?style=flat-square" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/admpub/xgo/"><img src="https://img.shields.io/docker/pulls/admpub/xgo.svg?style=flat-square" alt="Docker Pulls"></a>
</p>

## Fork

This repository is a fork of [karalabe/xgo](https://github.com/karalabe/xgo) to push images and [tags to an unique docker repository](https://hub.docker.com/r/admpub/xgo/tags/?page=1&ordering=last_updated) to make things more consistent for users.

I use [GitHub Actions](https://github.com/admpub/xgo/actions) and his [matrix strategy](https://help.github.com/en/articles/workflow-syntax-for-github-actions#jobsjob_idstrategymatrix) to build the images instead of using automated builds of Docker Hub (see [workflows](.github/workflows) folder).

This also creates a [standalone xgo executable](https://github.com/admpub/xgo/releases) that can be used on many platforms.

## About

Although Go strives to be a cross platform language, cross compilation from one
platform to another is not as simple as it could be, as you need the Go sources
bootstrapped to each platform and architecture.

The first step towards cross compiling was Dave Cheney's [golang-crosscompile](https://github.com/davecheney/golang-crosscompile)
package, which automatically bootstrapped the necessary sources based on your
existing Go installation. Although this was enough for a lot of cases, certain
drawbacks became apparent where the official libraries used CGO internally: any
dependency to third party platform code is unavailable, hence those parts don't
cross compile nicely (native DNS resolution, system certificate access, etc).

A step forward in enabling cross compilation was Alan Shreve's [gonative](https://github.com/inconshreveable/gonative)
package, which instead of bootstrapping the different platforms based on the
existing Go installation, downloaded the official pre-compiled binaries from the
golang website and injected those into the local toolchain. Since the pre-built
binaries already contained the necessary platform specific code, the few missing
dependencies were resolved, and true cross compilation could commence... of pure
Go code.

However, there was still one feature missing: cross compiling Go code that used
CGO itself, which isn't trivial since you need access to OS specific headers and
libraries. This becomes very annoying when you need access only to some trivial
OS specific functionality (e.g. query the CPU load), but need to configure and
maintain separate build environments to do it.

## Documentation

* [Enter xgo](doc/enter-xgo.md)
* [Installation](doc/installation.md)
* [Usage](doc/usage.md)
  * [Build flags](doc/usage/build-flags.md)
  * [Go releases](doc/usage/go-releases.md)
  * [Output prefixing](doc/usage/output-prefixing.md)
  * [Branch selection](doc/usage/branch-selection.md)
  * [Remote selection](doc/usage/remote-selection.md)
  * [Package selection](doc/usage/package-selection.md)
  * [Limit build targets](doc/usage/limit-build-targets.md)
  * [Platform versions](doc/usage/platform-versions.md)
  * [CGO dependencies](doc/usage/cgo-dependencies.md)

## How can I help ?

All kinds of contributions are welcome :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Support me on Patreon](.res/patreon.png)](https://www.patreon.com/crazymax) 
[![Paypal Donate](.res/paypal.png)](https://www.paypal.me/crazyws)

## License

MIT. See `LICENSE` for more details.
