package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	_ = flag.String("config", "", "build config JSON file")
	flag.Parse()

	// Read graph output from stdin
	var graph map[string]interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&graph); err != nil {
		fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
		os.Exit(1)
	}

	// TODO: compute reachability from entry points via BFS/DFS.
	// Each symbol gets status: REACHABLE, CONDITIONALLY_REACHABLE, or UNREACHABLE.
	out := map[string]interface{}{
		"reachability": []interface{}{},
		"statistics": map[string]int{
			"total_symbols":           0,
			"reachable":               0,
			"conditionally_reachable": 0,
			"unreachable":             0,
			"test_scope_excluded":     0,
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(1)
	}
}
