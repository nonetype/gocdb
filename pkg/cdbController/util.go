package cdbcontroller

import "fmt"

func addressToHexString(address int) (hexString string) {
	return fmt.Sprintf("0x%x", address)
}
