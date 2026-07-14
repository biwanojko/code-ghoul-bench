package main

import ("fmt"; "os")

func main() {
	if len(os.Args) < 2 { fmt.Fprintf(os.Stderr, "usage: deadcode <analyze|diff> [flags]\n"); os.Exit(2) }
	switch os.Args[1] {
	case "analyze":
		fmt.Println(`{"removals":[],"retained":[],"file_level_actions":[],"summary":{"total_symbols":0,"removable":0,"retained":0,"files_to_delete":0,"estimated_loc_savings":0}}`)
	case "diff":
		fmt.Println(`{"delta":{"new_dead":[],"resolved_dead":[],"persistent_dead":[],"new_symbol_dead":[]}}`)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1]); os.Exit(2)
	}
}
