package cdbcontroller

import (
	"fmt"
	"regexp"
	"strconv"
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

var breakpointHandlerMap = make(map[int]breakpointHandler)

func (ctrl *CdbController) InstallBreakpoint(installAddress int, bpType BreakpointType, condition string, handler breakpointHandler) (runError error) {
	hexAddress := addressToHexString(installAddress)
	// Fixup condition
	if 0 < len(condition) {
		if !(strings.HasPrefix(condition, "j") || strings.HasPrefix(condition, ".if")) {
			condition = fmt.Sprintf("j %s ''; 'gc' ", condition)
		}
	}

	bpCountBefore, _ := ctrl.GetBreakpointCount()
	fmt.Printf("Before: %#v\n", bpCountBefore)

	bpCommand := fmt.Sprintf("%s %s %s", bpType, hexAddress, condition)
	ctrl.Execute(bpCommand)
	bpCountAfter, _ := ctrl.GetBreakpointCount()
	fmt.Printf("After: %#v\n", bpCountAfter)

	if 0 < len(bpCountAfter)-len(bpCountBefore) {
		// save new breakpoint with handler
		// find installed breakpoint number, save it with handler
		sub := subset(bpCountAfter, bpCountBefore)[0]
		breakpointHandlerMap[sub] = handler
	} else {
		// if after-before == 0, raise error
		// TODO: RAISE ERROR
	}
	return
}

var breakpointIdPattern = regexp.MustCompile(`[\s]*(\d)+ [\w]+`)

func (ctrl *CdbController) GetBreakpointCount() (bpCount []int, runError error) {
	breakpoints, err := ctrl.Execute("bl")
	if err != nil {
		runError = err
	}
	found := breakpointIdPattern.FindAllStringSubmatch(breakpoints, -1)
	for _, e := range found {
		num, _ := strconv.Atoi(e[1])
		bpCount = append(bpCount, num)
	}
	return
}
