package analysis

import ("encoding/json"; "os"; "sort")

// LoadConfig loads the active_conditions config JSON.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil { return nil, err }
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil { return nil, err }
	return &cfg, nil
}

// ComputeReachability analyzes which symbols are reachable from entry points.
// TODO: implement BFS from graph.EntryPoints through graph.Edges.
// Statuses: REACHABLE, CONDITIONALLY_REACHABLE, UNREACHABLE.
func ComputeReachability(graph *GraphOutput, cfg *Config) *ReachOutput {
	var items []ReachabilityItem
	stats := ReachStats{}
	for _, node := range graph.Nodes {
		items = append(items, ReachabilityItem{ID: node, Status: "UNREACHABLE"})
		stats.Unreachable++
	}
	stats.TotalSymbols = stats.Unreachable
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	syms := make([]Symbol, len(graph.Symbols))
	copy(syms, graph.Symbols)
	return &ReachOutput{Reachability: items, Statistics: stats, Symbols: syms}
}
