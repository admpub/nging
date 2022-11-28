# conpty

Windows Pseudo Console (ConPTY) for Golang

See:

https://devblogs.microsoft.com/commandline/windows-command-line-introducing-the-windows-pseudo-console-conpty/

## Usage

```go
package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/UserExistsError/conpty"
)

func main() {
	commandLine := `c:\windows\system32\cmd.exe`
	cpty, err := conpty.Start(commandLine)
	if err != nil {
		log.Fatalf("Failed to spawn a pty:  %v", err)
	}
	defer cpty.Close()

	go func() {
		go io.Copy(os.Stdout, cpty)
		io.Copy(cpty, os.Stdin)
	}()

	exitCode, err := cpty.Wait(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ExitCode: %d", exitCode)
}
```