//go:build cgo

package server

/*
#include <stdlib.h>

extern int rust_compute(int x);
*/
import "C"

import "fmt"

// CallRustCompute calls the Rust compute function via CGO
//
//export GoCallback
func GoCallback(x C.int) C.int {
	result := C.rust_compute(x)
	fmt.Println("CGO round-trip:", int(result))
	return result
}

// unusedCGOHelper is never called - dead code
func unusedCGOHelper(data []byte) []byte {
	return data
}
