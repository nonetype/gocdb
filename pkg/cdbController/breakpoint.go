package cdbcontroller

import (
	"fmt"
	"log"
	"strings"
)

type BreakpointType string

const (
	Normal     = BreakpointType("bp")
	Unresolved = BreakpointType("bu")
	Hardware   = BreakpointType("ba")
	Symbolic   = BreakpointType("bm")
)

type BreakpointEvent struct {
	cdbEvent CdbEvent
	bpNum    int
}

type breakpointHandler func(BreakpointEvent) bool

// func (ctrl *CdbController) InstallBreakpoint(installAddress int, handler breakpointHandler) (runError error) {
func (ctrl *CdbController) InstallBreakpoint(installAddress int, bpType BreakpointType, condition string) (runError error) {
	hexAddress := addressToHexString(installAddress)
	// Fixup condition
	if 0 < len(condition) {
		if !(strings.HasPrefix(condition, "j") || strings.HasPrefix(condition, ".if")) {
			condition = fmt.Sprintf("j %s ''; 'gc' ", condition)
		}
	}
	bpCommand := fmt.Sprintf("%s %s %s", bpType, hexAddress, condition)
	output, _ := ctrl.Execute(bpCommand)
	log.Printf("breakpoint result: %s\n", output)

	// get breakpoint list before install
	//     install breakpoint
	// get breakpoint list after  install
	// if after-before == 0, raise error
	// else, find installed breakpoint number, save it with handler
	return
}

func (ctrl *CdbController) ListBreakpoint() (breakpoints string, runError error) {
	breakpoints, err := ctrl.Execute("bl")
	if err != nil {
		runError = err
	}
	return
}
