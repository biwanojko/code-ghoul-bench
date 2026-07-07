package main

import "fmt"

func main() {
	fmt.Println("result: direct")
}

func helper() string {
	return "hello from helper"
}
