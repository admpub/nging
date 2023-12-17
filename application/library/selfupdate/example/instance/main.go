package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~1`)
	// shutdown := make(chan os.Signal, 1)
	// // ctrl+c信号os.Interrupt
	// // pkill信号syscall.SIGTERM
	// signal.Notify(
	// 	shutdown,
	// 	os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP,
	// )
	// fmt.Println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~2`)
	// for i := 0; true; i++ {
	// 	fmt.Printf(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>%s`+com.StrLF, `A`)
	// 	sig := <-shutdown
	// 	fmt.Printf(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>%s`+com.StrLF, sig)
	// 	err := os.WriteFile(`./sig.txt`, []byte(sig.String()), os.ModePerm)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	fmt.Println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~3`)
	time.Sleep(time.Hour)
}
