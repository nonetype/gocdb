package cdbcontroller

import (
	"container/list"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/nonetype/gocdb/subprocess"
)

type CdbController struct {
	cdb      Cdb
	hostBits int
}

func NewController() *CdbController {
	cdbPath, hostBits, _ := searchCdb()
	fmt.Printf("cdbPath: %v\n", cdbPath) // DUMMY
	cdb := NewCdb()
	cdb.Run()
	return &CdbController{
		cdb:      *cdb,
		hostBits: hostBits,
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

func searchCdb() (cdbPath string, hostBits int, runError error) {
	// Get program files paths (check "PROGRAMFILES", "ProgramW6432", "ProgramFiles(x86)")
	programFilePaths := list.New()
	programFilePaths.PushBack(os.Getenv("PROGRAMFILES"))
	if 0 < len(os.Getenv("ProgramW6432")) {
		hostBits = 64
		programFilePaths.PushBack(os.Getenv("ProgramW6432"))
	} else {
		hostBits = 32
	}
	if 0 < len(os.Getenv("ProgramFiles(x86)")) {
		programFilePaths.PushBack(os.Getenv("ProgramFiles(x86)"))
	}

	debuggerPaths := list.New()
	for _, bits := range []string{"x64", "x86"} {
		debuggerPaths.PushBack(fmt.Sprintf("Windows Kits\\10\\Debuggers\\%s", bits))
		debuggerPaths.PushBack(fmt.Sprintf("Windows Kits\\8.1\\Debuggers\\%s", bits))
		debuggerPaths.PushBack(fmt.Sprintf("Windows Kits\\8.0\\Debuggers\\%s", bits))
		debuggerPaths.PushBack(fmt.Sprintf("Debugging Tools for Windows (%s)", bits))
		debuggerPaths.PushBack("Debugging Tools for Windows")
	}

	for progPath := programFilePaths.Front(); progPath != nil; progPath = progPath.Next() {
		for dbgPath := debuggerPaths.Front(); dbgPath != nil; dbgPath = dbgPath.Next() {
			fullPath := fmt.Sprintf("%s\\%s\\cdb.exe", progPath.Value, dbgPath.Value)
			if _, err := os.Stat(fullPath); err == nil {
				cdbPath = fullPath
				return
			}
		}
	}

	runError = fs.ErrNotExist
	return
}
