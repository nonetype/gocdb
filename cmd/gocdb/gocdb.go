package main

import (
	cdbcontroller "github.com/nonetype/gocdb/cdbController"
)

func main() {
	controller := cdbcontroller.NewController(".\\test_x64.exe")
	controller.Run()
	controller.Test()
	controller.Stop()
}
