package main

import "fmt"

func divide() int {
	i := 1
	return 10 / i
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			panic("this is a panic by person")
		}
	}()
	divide()
	fmt.Println("this is normal context")
}
