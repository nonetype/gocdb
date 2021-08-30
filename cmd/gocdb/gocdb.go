package main

import (
	"fmt"

	cdbcontroller "github.com/nonetype/gocdb/cdbController"
)

func main() {
	controller := cdbcontroller.NewController()
	fmt.Printf("controller: %v\n", controller)
}
