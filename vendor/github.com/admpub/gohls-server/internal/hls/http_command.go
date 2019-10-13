package hls

import (
	"sync"
	//"net/http"
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func calculateCommandHash(cmd string, args []string) string {
	h := sha1.New()
	h.Write([]byte(cmd))
	for _, v := range args {
		h.Write([]byte(v))
	}
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

type Empty struct{}

type HttpCommandHandler struct {
	tokenChannel    chan Empty
	cacheDir        string
	inProgress      map[string]string
	inProgressMutex *sync.RWMutex
	// path string
}

func NewHttpCommandHandler(workerCount int, cacheDir string) *HttpCommandHandler {
	ch := &HttpCommandHandler{make(chan Empty, workerCount), cacheDir, make(map[string]string), new(sync.RWMutex)}
	for i := workerCount; i > 0; i-- {
		ch.tokenChannel <- Empty{}
	}
	go ch.start()
	return ch
}

func (s *HttpCommandHandler) start() {

}

func (s *HttpCommandHandler) ServeCommand(cmdPath string, args []string, key string, w io.Writer) error {
	cachePath := filepath.Join(HomeDir, cacheDirName, s.cacheDir, key)
	mkerr := os.MkdirAll(filepath.Join(HomeDir, cacheDirName, s.cacheDir), 0777)
	if mkerr != nil {
		log.Errorf("Could not create cache dir %v: %v", filepath.Join(cacheDirName, s.cacheDir), mkerr)
		return mkerr
	}
	if file, err := os.Open(cachePath); err == nil {
		defer file.Close()
		_, err = io.Copy(w, file)
		if err != nil {
			log.Errorf("Error copying file to client: %v", err)
			return err
		}
		return nil
	}
	token := <-s.tokenChannel
	//log.Printf("Token: %v",key)
	defer func() {
		s.tokenChannel <- token
		//log.Printf("Released token")
	}()
	cacheFile, ferr := os.Create(cachePath)
	if ferr != nil {
		log.Errorf("Could not create cache file %v: %v", cacheFile, ferr)
		return ferr
	}
	defer cacheFile.Close()
	log.Debugf("Executing %v %v", cmdPath, args)
	cmd := exec.Command(cmdPath, args...)
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		log.Errorf("Error opening stdout of command: %v", err)
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.Errorf("Error starting command: %v", err)
		return err
	}
	filew := bufio.NewWriter(cacheFile)
	multiw := io.MultiWriter(filew, w)
	_, err = io.Copy(multiw, stdout)
	if err != nil {
		log.Errorf("Error copying data to client: %v", err)
		cacheFile.Close()
		os.Remove(cachePath)
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		return err
	}
	perr := cmd.Wait()
	if perr != nil {
		log.Errorf("Error waiting for process: %v", perr)
		cacheFile.Close()
		os.Remove(cachePath)
		return perr
	}
	/*
		if !state.Success {
			log.Errorf("Process didn't end successfully")
			cacheFile.Close()
			os.Remove(cachePath)
			return fmt.Errorf("Process didn't end successfully")
		}
	*/
	filew.Flush()
	log.Debugf("HTTP command success")
	return nil
	//s.inProgressMutex.Lock()
	//s.inProgressMutex.Unlock()
}
