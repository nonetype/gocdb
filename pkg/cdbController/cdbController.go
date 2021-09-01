package cdbcontroller

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"

	"github.com/nonetype/gocdb/subprocess"
)

type CdbController struct {
	targetProgram string
	cdb           *Cdb
	hostBits      int
}

func NewController(targetProgram string) *CdbController {
	cdbPath, hostBits, _ := searchCdb()
	cdb := NewCdb(cdbPath)
	return &CdbController{
		targetProgram: targetProgram,
		cdb:           cdb,
		hostBits:      hostBits,
	}
}

func (ctrl *CdbController) Run() {
	ctrl.cdb.Run(ctrl.targetProgram)
}

func (ctrl *CdbController) Stop() {
	log.Println("Stopping CDB")
	ctrl.cdb.process.Kill()
}

func (ctrl *CdbController) readCdb() (output string, runError error) {
	output, runError = ctrl.cdb.Read()
	return
}

func (ctrl *CdbController) writeCdb(input string) (runError error) {
	runError = ctrl.cdb.Write(input)
	return
}

func (ctrl *CdbController) Test() (runError error) {
	cdbConsoleRegex, _ := regexp.Compile("([0-9]{1}):([0-9]{3})> ")
	for {
		output, _ := ctrl.readCdb()
		fmt.Print(output)
		if cdbConsoleRegex.MatchString(output) {
			ctrl.writeCdb("?\n") // |.
			break
		}
	}
	output, _ := ctrl.readCdb()
	fmt.Print(output)
	return
}

type Cdb struct {
	cdbPath   string
	process   *os.Process
	targetPID int
	stdin     *io.WriteCloser
	stdout    *io.ReadCloser
}

func NewCdb(cdbPath string) *Cdb {
	cdb := &Cdb{
		cdbPath:   cdbPath,
		process:   nil,
		targetPID: 0,
		stdin:     nil,
		stdout:    nil,
	}
	return cdb
}

func (c *Cdb) Run(programArgs ...string) (runError error) {
	c.process, c.stdin, c.stdout, runError = subprocess.Run(".", c.cdbPath, programArgs...)
	if runError != nil {
		log.Fatal(runError)
	}

	return
}

func (c *Cdb) Read() (output string, runError error) {
	if c.stdout != nil {
		reader := bufio.NewReader(*c.stdout)
		buf := make([]byte, 1024)
		_, err := reader.Read(buf)
		if err != nil {
			runError = err
		}
		output = string(buf)
	}
	return
}

func (c *Cdb) Write(input string) (runError error) {
	if c.stdout != nil {
		writer := bufio.NewWriter(*c.stdin)
		writer.WriteString(input)
		writer.Flush()
	}
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
