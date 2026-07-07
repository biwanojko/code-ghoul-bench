package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("code-ghoul testdata main")
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("running")
	return nil
}

// unusedTopLevel is never called - dead code
func unusedTopLevel() string {
	return "unused"
}
