package analysis

// Symbol represents a code symbol extracted from source files
type Symbol struct {
	ID          string   `json:"id"`
	Language    string   `json:"language"`
	Kind        string   `json:"kind"`
	Visibility  string   `json:"visibility"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	FFIExported bool     `json:"ffi_exported,omitempty"`
	FFIName     string   `json:"ffi_name,omitempty"`
	Scope       string   `json:"scope,omitempty"`
	Conditions  []string `json:"conditions,omitempty"`
	Attributes  []string `json:"attributes,omitempty"`
}

// Callsite represents a call from one symbol to another (possibly cross-language)
type Callsite struct {
	From      string `json:"from"`
	ToName    string `json:"to_name"`
	Mechanism string `json:"mechanism"`
	File      string `json:"file"`
	Line      int    `json:"line"`
}

// EntryPoint is a symbol that is a root of the call graph
type EntryPoint struct {
	SymbolID string `json:"symbol_id"`
	Reason   string `json:"reason"`
}

// ExtractOutput is the output of the extract command
type ExtractOutput struct {
	Symbols     []Symbol     `json:"symbols"`
	Callsites   []Callsite   `json:"callsites"`
	EntryPoints []EntryPoint `json:"entry_points"`
}

// Edge is a directed edge in the call graph
type Edge struct {
	From       string   `json:"from"`
	To         string   `json:"to"`
	Mechanism  string   `json:"mechanism"`
	Confidence float64  `json:"confidence"`
	Conditions []string `json:"conditions,omitempty"`
}

// Unresolved is a callsite that could not be resolved to a known symbol
type Unresolved struct {
	From      string `json:"from"`
	ToName    string `json:"to_name"`
	Mechanism string `json:"mechanism"`
	Reason    string `json:"reason"`
}

// GraphOutput is the output of the graph command
// It carries the full symbol table forward so downstream steps don't need --extract.
type GraphOutput struct {
	Nodes       []string     `json:"nodes"`
	Edges       []Edge       `json:"edges"`
	Unresolved  []Unresolved `json:"unresolved"`
	EntryPoints []EntryPoint `json:"entry_points,omitempty"`
	Symbols     []Symbol     `json:"symbols,omitempty"`
}

// ReachabilityItem is the reachability status of a symbol
type ReachabilityItem struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	ReachedVia string `json:"reached_via,omitempty"`
}

// ReachStats are statistics about reachability.
// Note: test_scope_excluded is a separate category, not included in unreachable.
// total_symbols = reachable + conditionally_reachable + unreachable + test_scope_excluded
type ReachStats struct {
	TotalSymbols           int `json:"total_symbols"`
	Reachable              int `json:"reachable"`
	ConditionallyReachable int `json:"conditionally_reachable"`
	Unreachable            int `json:"unreachable"`
	TestScopeExcluded      int `json:"test_scope_excluded"`
}

// ReachOutput is the output of the reach command.
// Carries symbols forward so plan can access metadata without --extract.
type ReachOutput struct {
	Reachability []ReachabilityItem `json:"reachability"`
	Statistics   ReachStats         `json:"statistics"`
	Symbols      []Symbol           `json:"symbols,omitempty"`
}

// Removal is a symbol that can be removed
type Removal struct {
	ID         string  `json:"id"`
	File       string  `json:"file"`
	Line       int     `json:"line"`
	Priority   string  `json:"priority"`   // "P0", "P1", "P2", "P3"
	Score      float64 `json:"score"`      // 1.0, 0.8, 0.5, 0.3
	LocSavings int     `json:"loc_savings"` // always positive int
	Reason     string  `json:"reason"`
}

// Retained is a symbol that should be kept
type Retained struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// FileLevelAction describes what to do with a file
type FileLevelAction struct {
	File                  string   `json:"file"`
	Action                string   `json:"action"`
	AllSymbolsUnreachable bool     `json:"all_symbols_unreachable"`
	Removals              []string `json:"removals"`
}

// PlanSummary summarizes the removal plan
type PlanSummary struct {
	TotalSymbols        int `json:"total_symbols"`
	Removable           int `json:"removable"`
	Retained            int `json:"retained"`
	FilesToDelete       int `json:"files_to_delete"`
	EstimatedLocSavings int `json:"estimated_loc_savings"`
}

// PlanOutput is the output of the plan command
type PlanOutput struct {
	Removals         []Removal         `json:"removals"`
	Retained         []Retained        `json:"retained"`
	FileLevelActions []FileLevelAction `json:"file_level_actions"`
	Summary          PlanSummary       `json:"summary"`
}

// Config is the configuration for the analysis
type Config struct {
	ActiveConditions          map[string]map[string]bool `json:"active_conditions"`
	ConsiderTestScope         bool                       `json:"consider_test_scope"`
	ConsiderAllConfigurations bool                       `json:"consider_all_configurations"`
}
