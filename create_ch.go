package main

import (
	"fmt"
)

func create_chan() {
	var a chan int
	if a == nil {
		a = make(chan int)
	}
	test1(a)
}

func test1(ch chan int) {
	fmt.Printf("type is %T\n, value is %v\n", a, a)
}
