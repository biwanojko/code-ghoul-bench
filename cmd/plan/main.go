package main

import ("encoding/json"; "fmt"; "os"; "code-ghoul/internal/analysis")

func main() {
	var reach analysis.ReachOutput
	if err := json.NewDecoder(os.Stdin).Decode(&reach); err != nil {
		fmt.Fprintf(os.Stderr, "plan: %v\n", err); os.Exit(2)
	}
	out := analysis.BuildPlan(&reach)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}
