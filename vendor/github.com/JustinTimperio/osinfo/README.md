# osinfo

![GitHub](https://img.shields.io/github/license/JustinTimperio/osinfo)
[![Go Reference](https://pkg.go.dev/badge/github.com/JustinTimperio/osinfo.svg)](https://pkg.go.dev/github.com/JustinTimperio/osinfo)
[![Go Report Card](https://goreportcard.com/badge/github.com/JustinTimperio/osinfo)](https://goreportcard.com/report/github.com/JustinTimperio/osinfo)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/180976560a9b46678d9d67053c0d7fc9)](https://www.codacy.com/gh/JustinTimperio/osinfo/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=JustinTimperio/osinfo&amp;utm_campaign=Badge_Grade)

## What is osinfo?
OSinfo is a cross-platform OS Version collection tool. It is designed to unify multiple OS detection methods in a single module that can easily be integrated into other projects. 

### Officially Supported:

| Windows             | MacOS                 | Linux               | BSD     |
|---------------------|-----------------------|---------------------|---------|
| Windows Server 2012 | 10.6  - Snow Leopard  | Ubuntu              | FreeBSD | 
| Windows Server 2016 | 10.7  - Lion          | Debian              | OpenBSD | 
| Windows Server 2019 | 10.8  - Mountain Lion | MXLinux             |         | 
| Windows 7           | 10.9  - Mavericks     | Mint                |         | 
| Windows 8           | 10.10 - Yosemite      | Kali                |         | 
| Windows 10          | 10.11 - El Capitan    | ParrotOS            |         |
|                     | 10.12 - Sierra        | OpenSUSE Leap       |         | 
|                     | 10.13 - High Sierra   | OpenSUSE TumbleWeed |         |
|                     | 10.14 - Mojave        | OpenSUSE SLES       |         |
|                     | 10.15 - Catalina      | Arch                |         |
|                     | 11.0  - Big Sur       | Manjaro             |         |
|                     |                       | Alpine              |         |
|                     |                       | Fedora              |         |
|                     |                       | RHEL                |         |
|                     |                       | CentOS              |         |
|                     |                       | Oracle              |         |


## Example Usage
 1. Create `fetchinfo.go`
```go
   package main

   import (
      "fmt"

      "github.com/JustinTimperio/osinfo"
   )

   func main() {
		release := osinfo.GetVersion()
		fmt.Printf(release.String())
	 }
```
 2. `go mod init`
 3. `go mod tidy`
 4. `go run fetchinfo.go`

## Example Outputs
```sh
--------------------

Runtime: linux
Architecture: amd64
OS Name: Arch Linux
Version: rolling
Kernel: 5.11.2-arch1-1
Distro: arch
Package Manager: pacman

--------------------

Runtime: linux
Architecture: amd64
OS Name: Debian GNU/Linux 10 (buster)
Version: 10
Kernel: 4.19.0-13-amd64
Distro: debian
Package Manager: apt

--------------------

Runtime: windows
Architecture: amd64
OS Name: Windows Server 2016 Standard Evaluation
Version: 10
Build: 14393

--------------------

Runtime: darwin
Architecture: amd64
OS Name: Mac OS X
Version: 11.0.1
Version Name: MacOS - Big Sur 

--------------------

Runtime: freebsd
Architecture: amd64
OS Name: FreeBSD 12.1-RELEASE
Version: 12.1-RELEASE
Kernel: 1201000
Package Manger: pkg

--------------------

Runtime: openbsd 
Architecture: amd64
OS Name: OpenBSD 6.7 
Version: 6.7 
Kernel: GENERIC.MP#182 
Package Manger: pkg_add

--------------------
```
