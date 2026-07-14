package main

import ("encoding/json"; "flag"; "fmt"; "os"; "code-ghoul/internal/analysis")

func main() {
	configPath := flag.String("config", "", "config JSON")
	flag.Parse()
	var graph analysis.GraphOutput
	if err := json.NewDecoder(os.Stdin).Decode(&graph); err != nil {
		fmt.Fprintf(os.Stderr, "reach: %v\n", err); os.Exit(2)
	}
	var cfg *analysis.Config
	if *configPath != "" {
		var err error
		if cfg, err = analysis.LoadConfig(*configPath); err != nil {
			fmt.Fprintf(os.Stderr, "reach: %v\n", err); os.Exit(2)
		}
	}
	out := analysis.ComputeReachability(&graph, cfg)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}
