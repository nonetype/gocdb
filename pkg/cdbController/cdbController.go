package cdbcontroller

import (
	"container/list"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
)

var cdbConsoleRegex = regexp.MustCompile("([0-9]{1}):([0-9]{3})> ")

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

func (ctrl *CdbController) Execute(command string) (output string, runError error) {
	log.Printf("Executing '%s'\n", command)
	runError = ctrl.cdb.Write(command)
	output, _ = ctrl.cdb.ReadAll()
	log.Printf("result of '%s': %s\n", command, output)
	return
}

// func (ctrl *CdbController) RegisterExceptionHandler(callback func) (runError error)
// func (ctrl *CdbController) RegisterModuleLoadHandler(callback func) (runError error)
// func (ctrl *CdbController) RegisterMainHandler(callback func) (runError error)

func (ctrl *CdbController) Test() (runError error) {
	ctrl.cdb.ReadAll()
	ctrl.Execute("|.") // |.
	ctrl.InstallBreakpoint(0x1234, Normal, "")
	ctrl.ListBreakpoint()
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
