package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: deadcode <analyze|diff> [flags]\n")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "analyze":
		// TODO: run full pipeline (extract→graph→reach→plan) on --repo dir
		// Output removal plan JSON; exit 0 if no dead code, 1 if dead code found
		fmt.Println(`{"removals":[],"retained":[],"file_level_actions":[],"summary":{"total_removable_symbols":0,"total_retained_symbols":0,"estimated_loc_savings":0,"files_to_delete":0,"files_to_edit":0}}`)
	case "diff":
		// TODO: run pipeline on --base and --head commits, report delta
		fmt.Println(`{"base_commit":"","head_commit":"","delta":{"new_dead":[],"resolved_dead":[],"persistent_dead":[],"new_symbol_dead":[]},"summary":{}}`)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		os.Exit(2)
	}
}
