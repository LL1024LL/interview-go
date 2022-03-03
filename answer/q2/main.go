package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(help("1123"))
}

func help(str string) bool{
	for i, s := range str{
		if strings.Index(str, string(s)) != i{
			return true
		}
	}
	return false
}
