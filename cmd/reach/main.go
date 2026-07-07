package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"code-ghoul/internal/analysis"
)

func main() {
	configPath := flag.String("config", "", "path to config JSON")
	// --extract is accepted for compatibility but no longer required
	_ = flag.String("extract", "", "(unused) path to extract output")
	flag.Parse()

	// Read graph from stdin
	var graph analysis.GraphOutput
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&graph); err != nil {
		fmt.Fprintf(os.Stderr, "reach: decode graph: %v\n", err)
		os.Exit(2)
	}

	var cfg *analysis.Config
	if *configPath != "" {
		var err error
		cfg, err = analysis.LoadConfig(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reach: load config: %v\n", err)
			os.Exit(2)
		}
	}

	out := analysis.ComputeReachability(&graph, cfg)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "reach: encode: %v\n", err)
		os.Exit(2)
	}
}
