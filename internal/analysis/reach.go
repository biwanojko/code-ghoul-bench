package analysis

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

// LoadConfig loads a config from a JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ComputeReachability performs BFS from entry points to determine reachability.
// The graph output carries entry_points and symbols so no separate extract output is needed.
// Statistics: total_symbols = reachable + conditionally_reachable + unreachable + test_scope_excluded
// (test_scope_excluded is a distinct category, not counted in unreachable)
func ComputeReachability(graph *GraphOutput, cfg *Config) *ReachOutput {
	// Build adjacency list
	adj := map[string][]Edge{}
	for _, e := range graph.Edges {
		adj[e.From] = append(adj[e.From], e)
	}

	// Build symbol map from graph.Symbols
	symByID := map[string]*Symbol{}
	for i := range graph.Symbols {
		symByID[graph.Symbols[i].ID] = &graph.Symbols[i]
	}

	// Determine which nodes are test-scope
	isTestScope := func(id string) bool {
		if sym, ok := symByID[id]; ok {
			return sym.Scope == "test"
		}
		return false
	}

	// BFS status: 1=reachable, 2=conditionally_reachable
	status := map[string]int{}
	reachedVia := map[string]string{}

	type qItem struct {
		id          string
		conditional bool
		via         string
	}
	var queue []qItem

	// Add entry points to queue
	for _, ep := range graph.EntryPoints {
		if cfg != nil && !cfg.ConsiderTestScope && isTestScope(ep.SymbolID) {
			continue
		}
		queue = append(queue, qItem{id: ep.SymbolID, conditional: false, via: ""})
	}

	// BFS
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		id := item.id
		if cfg != nil && !cfg.ConsiderTestScope && isTestScope(id) {
			continue
		}

		// Check conditions for this symbol
		isConditional := item.conditional
		if !isConditional {
			if sym, ok := symByID[id]; ok {
				if len(sym.Conditions) > 0 && cfg != nil {
					condActive := isConditionActive(sym.Conditions, sym.Language, cfg)
					if !condActive {
						if cfg.ConsiderAllConfigurations {
							isConditional = true
						} else {
							continue
						}
					}
				}
			}
		}

		targetStatus := 1
		if isConditional {
			targetStatus = 2
		}

		current, exists := status[id]
		if exists && current >= targetStatus {
			continue
		}

		// Update status
		status[id] = targetStatus
		if item.via != "" {
			reachedVia[id] = item.via
		} else {
			reachedVia[id] = "entry_point"
		}

		// Explore neighbors
		for _, edge := range adj[id] {
			neighbor := edge.To
			edgeConditional := isConditional
			if len(edge.Conditions) > 0 && cfg != nil {
				if !isConditionActive(edge.Conditions, "", cfg) {
					if cfg.ConsiderAllConfigurations {
						edgeConditional = true
					} else {
						continue
					}
				}
			}
			neighborStatus := status[neighbor]
			newStatus := 1
			if edgeConditional {
				newStatus = 2
			}
			if newStatus > neighborStatus {
				queue = append(queue, qItem{id: neighbor, conditional: edgeConditional, via: id})
			}
		}
	}

	// Build output
	var items []ReachabilityItem
	stats := ReachStats{}

	for _, node := range graph.Nodes {
		sym, hasSym := symByID[node]

		// Test scope check — separate category, not counted in unreachable
		if cfg != nil && !cfg.ConsiderTestScope && hasSym && sym.Scope == "test" {
			items = append(items, ReachabilityItem{
				ID:     node,
				Status: "UNREACHABLE",
			})
			stats.TestScopeExcluded++
			continue
		}

		s, exists := status[node]
		var ri ReachabilityItem
		ri.ID = node

		if !exists || s == 0 {
			// Unreachable — but check if conditionally reachable via build conditions
			if cfg != nil && cfg.ConsiderAllConfigurations && hasSym && len(sym.Conditions) > 0 {
				if !isConditionActive(sym.Conditions, sym.Language, cfg) {
					ri.Status = "CONDITIONALLY_REACHABLE"
					ri.ReachedVia = "build_condition"
					stats.ConditionallyReachable++
				} else {
					ri.Status = "UNREACHABLE"
					stats.Unreachable++
				}
			} else {
				ri.Status = "UNREACHABLE"
				stats.Unreachable++
			}
		} else if s == 1 {
			ri.Status = "REACHABLE"
			ri.ReachedVia = reachedVia[node]
			stats.Reachable++
		} else {
			ri.Status = "CONDITIONALLY_REACHABLE"
			ri.ReachedVia = reachedVia[node]
			stats.ConditionallyReachable++
		}
		items = append(items, ri)
	}

	// total_symbols = reachable + conditionally_reachable + unreachable + test_scope_excluded
	stats.TotalSymbols = stats.Reachable + stats.ConditionallyReachable + stats.Unreachable + stats.TestScopeExcluded

	// Sort by id
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	// Pass symbols through
	syms := make([]Symbol, len(graph.Symbols))
	copy(syms, graph.Symbols)

	return &ReachOutput{
		Reachability: items,
		Statistics:   stats,
		Symbols:      syms,
	}
}

// isConditionActive checks if all conditions in the list are active in the config.
// Handles negated conditions (!cgo) and compound conditions (linux && amd64).
func isConditionActive(conditions []string, lang string, cfg *Config) bool {
	if cfg == nil {
		return true
	}
	for _, cond := range conditions {
		if !evalCondition(cond, lang, cfg) {
			return false
		}
	}
	return true
}

// evalCondition evaluates a single condition string (may be compound like "linux && amd64")
func evalCondition(cond, lang string, cfg *Config) bool {
	cond = strings.TrimSpace(cond)
	// Handle compound AND conditions
	if strings.Contains(cond, "&&") {
		parts := strings.Split(cond, "&&")
		for _, part := range parts {
			if !evalCondition(strings.TrimSpace(part), lang, cfg) {
				return false
			}
		}
		return true
	}
	// Handle compound OR conditions
	if strings.Contains(cond, "||") {
		parts := strings.Split(cond, "||")
		for _, part := range parts {
			if evalCondition(strings.TrimSpace(part), lang, cfg) {
				return true
			}
		}
		return false
	}
	// Handle negated condition
	negated := false
	if strings.HasPrefix(cond, "!") {
		negated = true
		cond = strings.TrimPrefix(cond, "!")
	}
	active := lookupCondKey(cond, lang, cfg)
	if negated {
		return !active
	}
	return active
}

// lookupCondKey checks whether a single condition key is active in the config
func lookupCondKey(key, lang string, cfg *Config) bool {
	if langConds, ok := cfg.ActiveConditions[lang]; ok {
		if v, ok := langConds[key]; ok {
			return v
		}
	}
	// Fall back: search all language configs
	for _, lc := range cfg.ActiveConditions {
		if v, ok := lc[key]; ok {
			return v
		}
	}
	// Unknown condition keys are treated as inactive by default
	return false
}
