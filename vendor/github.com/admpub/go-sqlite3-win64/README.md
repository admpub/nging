# Go-SQLite3-Win64
## GoLang Wrapper for sqlite3.dll on Windows 64Bit

This package provides an alternative to [go-sqlite3](https://github.com/mattn/go-sqlite3) on Windows 64bit.

Although I don't think you can call it pure-Go since it requires the DLL, it doesn't require CGo to build, so the build is pure-Go.

This should allow really easy cross compile support when building from another OS (e.g. Linux).

Basic functionality is implemented, but it doesn't support things like user defined functions that call back to Go code.

You'll need [sqlite3.dll](https://sqlite.org/download.html) in either the same folder as you finished executable, or in a "support" folder in the same path as your exe.

Otherwise, usage should be the same as for go-sqlite3.
