package toukibo

import "fmt"

const DebugOn = false

func PrintSlice(s []string) {
	for i, v := range s {
		fmt.Println(i, v)
	}
}

func PrintBar() {
	fmt.Println("--------------------------------------------------------")
}
