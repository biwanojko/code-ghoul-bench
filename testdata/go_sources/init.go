package server

import "fmt"

func init() {
	fmt.Println("server package initialized")
	_ = BuildInfo
}

// Version is the server version constant
const Version = "1.0.0"

// BuildInfo holds build metadata
var BuildInfo = map[string]string{
	"version": Version,
	"os":      "linux",
}

// debugMode is an unexported flag - dead code
var debugMode = false

// enableDebug enables debug mode - dead code
func enableDebug() {
	debugMode = true
}
