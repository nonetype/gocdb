package main

import (
	cdbcontroller "github.com/nonetype/gocdb/cdbController"
)

func main() {
	controller := cdbcontroller.NewController(".\\test.exe")
	controller.Run()
	controller.Test()
	controller.Stop()
}
