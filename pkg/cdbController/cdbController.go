package cdbcontroller

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nonetype/gocdb/subprocess"
)

type CdbController struct {
	cdb Cdb
}

func NewController() *CdbController {
	cdb := NewCdb()
	cdb.Run()
	return &CdbController{
		cdb: *cdb,
	}
}

func (ctrl CdbController) Stop() {
	ctrl.cdb.process.Kill()
}

type Cdb struct {
	process   os.Process
	targetPID int
	stdin     io.WriteCloser
	stdout    io.ReadCloser
}

func NewCdb() *Cdb {
	return &Cdb{
		process:   os.Process{},
		targetPID: 0,
		stdin:     nil,
		stdout:    nil,
	}
}

func (c Cdb) Run() (pid int, runError error) {
	c.process, c.stdin, c.stdout, runError = subprocess.Run("./", "echo", "HELLO WORLD")
	if runError != nil {
		log.Fatal(runError)
	}
	stdoutValue, _ := ioutil.ReadAll(c.stdout)
	output := string(stdoutValue)
	fmt.Println("CDB OUTPUT:", output)
	fmt.Println("CDB PID:", c.process.Pid)

	return
}

func searchCdb() (cdbPath string, err error) {
	pathToProgramFiles := os.Getenv("PROGRAMFILES")
	if strings.Contains(pathToProgramFiles, "ProgramW6432") {
		// SET HOST BIT TO 64
	} // ELSE SET HOST BIT TO 32
	// CHECK ProgramFiles(x86) TOO
	// LIST ALL PROGRAMFILE PATH, SEARCH 'Windows Kits\\10', 'Windows Kits\\8.*', 'Debugging Tools for Windows*', ...
	// IF 'cdb.exe' EXISTS, RETURN PATH
	// ELSE, RAISE NOT FOUND ERROR
	return
}
