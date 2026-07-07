package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code-ghoul/internal/analysis"
)

func main() {
	// Read ExtractOutput from stdin
	var extract analysis.ExtractOutput
	if err := json.NewDecoder(os.Stdin).Decode(&extract); err != nil {
		fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
		os.Exit(1)
	}

	// TODO: build the cross-language call graph from extract.Symbols and extract.Callsites.
	// Resolve cgo/JNI edges and output nodes + edges + unresolved callsites.
	out := map[string]interface{}{
		"nodes":      []string{},
		"edges":      []interface{}{},
		"unresolved": []interface{}{},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(1)
	}
}
