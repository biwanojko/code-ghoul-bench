//go:build linux && amd64

package server

import "fmt"

func init() {
	fmt.Println("linux/amd64 init")
}

// LinuxAMD64Optimize is only available on linux/amd64
func LinuxAMD64Optimize(data []byte) []byte {
	fmt.Println("using linux/amd64 optimized path")
	return data
}

// platformSpecificCleanup is dead code on all platforms
func platformSpecificCleanup() {
	fmt.Println("cleanup")
}
