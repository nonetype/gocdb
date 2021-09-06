package cdbcontroller

import "fmt"

func addressToHexString(address int) (hexString string) {
	return fmt.Sprintf("0x%x", address)
}

func subset(a, b []int) (out []int) {
	set := make(map[int]int)
	for _, value := range b {
		set[value] += 1
	}

	for _, value := range a {
		if count, found := set[value]; !found {
			out = append(out, value)
		} else {
			set[value] = count - 1
		}
	}
	return
}
