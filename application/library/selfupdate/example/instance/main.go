package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/webx-top/com"
)

func main() {
	fmt.Println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~1`)
	shutdown := make(chan os.Signal, 1)
	// ctrl+c信号os.Interrupt
	// pkill信号syscall.SIGTERM
	signal.Notify(
		shutdown,
		os.Interrupt, syscall.SIGTERM,
	)
	fmt.Println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~2`)
	for i := 0; true; i++ {
		sig := <-shutdown
		fmt.Printf(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>%s`+com.StrLF, sig)
		err := os.WriteFile(`./sig.txt`, []byte(sig.String()), os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~3`)
}
