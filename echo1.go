package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println("1:" + s)
	fmt.Println("2:" + strings.Join(os.Args[1:], " ")) //类似php 的implode?

	s, sep = "", ""
	for _, value := range os.Args[1:len(os.Args)] { //for配合range
		s += sep + value
		sep = " "
	}
	fmt.Println("3111:" + s)

	fmt.Println(os.Args[1:])
}
