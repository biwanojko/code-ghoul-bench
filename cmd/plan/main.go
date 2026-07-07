package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read reach output from stdin
	var reach map[string]interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&reach); err != nil {
		fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
		os.Exit(1)
	}

	// TODO: generate safe removal plan from UNREACHABLE symbols.
	// Apply safety filters (#[used], @Keep, init/static_initializer, etc.).
	// Assign priority: P0=1.0 (whole file), P1=0.8 (function), P2=0.5 (type), P3=0.3 (conditional).
	out := map[string]interface{}{
		"removals":          []interface{}{},
		"retained":          []interface{}{},
		"file_level_actions": []interface{}{},
		"summary": map[string]interface{}{
			"total_removable_symbols": 0,
			"total_retained_symbols":  0,
			"estimated_loc_savings":   0,
			"files_to_delete":         0,
			"files_to_edit":           0,
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(1)
	}
}
