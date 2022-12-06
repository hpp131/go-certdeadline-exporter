package main

import (
	"flag"
	"fmt"
)

func main() {
	//fmt.Println("this is a addtest branch code")
	testString := flag.String("domains", "", "help")
	flag.Parse()
	fmt.Println(*testString)
}
