package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"code-ghoul/internal/analysis"
)

func main() {
	// --extract accepted for backward compatibility but no longer needed
	_ = flag.String("extract", "", "(unused) path to extract output JSON")
	flag.Parse()

	// Read reach output from stdin
	var reach analysis.ReachOutput
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&reach); err != nil {
		fmt.Fprintf(os.Stderr, "plan: decode reach: %v\n", err)
		os.Exit(2)
	}

	out := analysis.BuildPlan(&reach)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "plan: encode: %v\n", err)
		os.Exit(2)
	}
}
