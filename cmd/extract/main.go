package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"code-ghoul/internal/analysis"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan")
	flag.Parse()

	out, err := analysis.ExtractDir(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "extract error: %v\n", err)
		os.Exit(2)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(2)
	}
}
