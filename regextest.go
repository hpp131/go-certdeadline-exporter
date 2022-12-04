package main

import (
	"fmt"
	"strings"
)

func main() {
	testString := "aaa,bbb,ccc,ddd"
	sliceExample := strings.Split(testString, ",")
	fmt.Println(sliceExample)
	for _, value := range sliceExample {
		fmt.Println(value)
	}
}
