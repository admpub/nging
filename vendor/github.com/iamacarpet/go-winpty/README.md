# Go-WinPTY
## GoLang Wrapper for WinPTY.dll

Small wrapper around the [WinPTY](https://github.com/rprichard/winpty) DLL.

This, for example, should allow a Go app to serve up a Windows "cmd" or "powershell" prompt via WebSocket to [xterm.js](https://github.com/sourcelair/xterm.js)

Currently requires `winpty.dll` and `winpty-agent.exe` to be in the same directory as the compiled Go executable.

These are best obtained from the GitHub releases page, in the msvc2015 package.

The error handling needs a bit of work currently though!
