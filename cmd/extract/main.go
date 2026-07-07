package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"code-ghoul/internal/analysis"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan for source files")
	flag.Parse()

	// TODO: implement ExtractDir to walk *dir and populate all three slices.
	// Walk files with extensions .go, .rs, .java, .kt.
	// See testdata/ for the fixture files and their expected output format.
	out := &analysis.ExtractOutput{
		Symbols:     []analysis.Symbol{},
		Callsites:   []analysis.Callsite{},
		EntryPoints: []analysis.EntryPoint{},
	}

	_ = *dir // remove this line once you implement the extraction

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(1)
	}
}
