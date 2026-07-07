package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code-ghoul/internal/analysis"
)

func main() {
	// Read extract output from stdin
	var extract analysis.ExtractOutput
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&extract); err != nil {
		fmt.Fprintf(os.Stderr, "graph: decode error: %v\n", err)
		os.Exit(2)
	}

	graph := analysis.BuildGraph(&extract)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(graph); err != nil {
		fmt.Fprintf(os.Stderr, "graph: encode error: %v\n", err)
		os.Exit(2)
	}
}
