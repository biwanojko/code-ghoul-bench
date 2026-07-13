package analysis

import (
	"path/filepath"
	"sort"
	"strings"
)

// BuildGraph resolves callsites and builds the call graph.
// Used internally by the deadcode diff subcommand to build graphs per commit.
func BuildGraph(extract *ExtractOutput) *GraphOutput {
	// Build maps
	symByID := map[string]*Symbol{}
	for i := range extract.Symbols {
		symByID[extract.Symbols[i].ID] = &extract.Symbols[i]
	}

	// Map ffi_name -> symbol IDs (for CGO/JNI resolution)
	ffiNameToIDs := map[string][]string{}
	for _, sym := range extract.Symbols {
		if sym.FFIName != "" {
			ffiNameToIDs[sym.FFIName] = append(ffiNameToIDs[sym.FFIName], sym.ID)
		}
		// Also index by simple name for JNI resolution
		if sym.Language == "rust" && sym.FFIExported {
			parts := strings.Split(sym.ID, "#")
			if len(parts) > 1 {
				simpleName := parts[len(parts)-1]
				ffiNameToIDs[simpleName] = append(ffiNameToIDs[simpleName], sym.ID)
			}
		}
	}

	// Nodes = all symbol IDs
	nodes := make([]string, 0, len(extract.Symbols))
	for _, sym := range extract.Symbols {
		nodes = append(nodes, sym.ID)
	}
	sort.Strings(nodes)

	var edges []Edge
	var unresolved []Unresolved

	for _, cs := range extract.Callsites {
		switch cs.Mechanism {
		case "cgo":
			// Find Rust symbol with matching ffi_name or function name
			targets := ffiNameToIDs[cs.ToName]
			if len(targets) == 0 {
				// Try searching all rust symbols
				for _, sym := range extract.Symbols {
					if sym.Language == "rust" && sym.FFIExported {
						if sym.FFIName == cs.ToName || strings.HasSuffix(sym.ID, "#"+cs.ToName) {
							targets = append(targets, sym.ID)
						}
					}
				}
			}
			if len(targets) > 0 {
				for _, to := range targets {
					edges = append(edges, Edge{
						From:       cs.From,
						To:         to,
						Mechanism:  "cgo",
						Confidence: 1.0,
					})
				}
			} else {
				unresolved = append(unresolved, Unresolved{
					From:      cs.From,
					ToName:    cs.ToName,
					Mechanism: cs.Mechanism,
					Reason:    "no_match",
				})
			}

		case "jni":
			// Find Rust symbol with matching JNI export name
			targets := ffiNameToIDs[cs.ToName]
			if len(targets) == 0 {
				for _, sym := range extract.Symbols {
					if sym.Language == "rust" && sym.FFIExported {
						if sym.FFIName == cs.ToName || strings.HasSuffix(sym.ID, "#"+cs.ToName) {
							targets = append(targets, sym.ID)
						}
					}
				}
			}
			if len(targets) > 0 {
				for _, to := range targets {
					edges = append(edges, Edge{
						From:       cs.From,
						To:         to,
						Mechanism:  "jni",
						Confidence: 1.0,
					})
				}
			} else {
				unresolved = append(unresolved, Unresolved{
					From:      cs.From,
					ToName:    cs.ToName,
					Mechanism: cs.Mechanism,
					Reason:    "no_match",
				})
			}

		case "reflection":
			// Find all public methods on the named class (lower confidence)
			className := cs.ToName
			prefix := "java://" + className
			found := false
			for _, sym := range extract.Symbols {
				if strings.HasPrefix(sym.ID, prefix) && (sym.Visibility == "public" || sym.Visibility == "protected") {
					edges = append(edges, Edge{
						From:       cs.From,
						To:         sym.ID,
						Mechanism:  "reflection",
						Confidence: 0.7,
					})
					found = true
				}
			}
			if !found {
				unresolved = append(unresolved, Unresolved{
					From:      cs.From,
					ToName:    cs.ToName,
					Mechanism: cs.Mechanism,
					Reason:    "no_match",
				})
			}

		case "jni_load":
			// Library load - not a direct symbol reference
			unresolved = append(unresolved, Unresolved{
				From:      cs.From,
				ToName:    cs.ToName,
				Mechanism: cs.Mechanism,
				Reason:    "library_load_not_symbol",
			})

		case "ffi":
			// extern "C" declaration in Rust - points to Go side
			for _, sym := range extract.Symbols {
				if sym.Language == "go" && strings.HasSuffix(sym.ID, "#"+cs.ToName) {
					edges = append(edges, Edge{
						From:       cs.From,
						To:         sym.ID,
						Mechanism:  "ffi",
						Confidence: 0.9,
					})
				}
			}

		case "call":
			// Same-language call - resolve by full symbol ID
			if _, ok := symByID[cs.ToName]; ok {
				edges = append(edges, Edge{
					From:       cs.From,
					To:         cs.ToName,
					Mechanism:  "call",
					Confidence: 1.0,
				})
			}

		case "go_call":
			// Regular Go function call - resolve by function name within Go symbols
			callerPrefix := "go://"
			// Extract relative path from caller ID to find matching symbol in same or related package
			// Format: go://path/to/file.go#FuncName
			// We search all Go symbols with the matching name
			callerParts := strings.SplitN(cs.From, "#", 2)
			callerFileBase := ""
			if len(callerParts) > 0 {
				// Get directory of caller
				callerFileBase = strings.TrimPrefix(callerParts[0], callerPrefix)
				callerFileBase = filepath.Dir(callerFileBase)
			}
			var matched []string
			for _, sym := range extract.Symbols {
				if sym.Language != "go" {
					continue
				}
				// Check if name matches
				parts := strings.SplitN(sym.ID, "#", 2)
				if len(parts) < 2 || parts[1] != cs.ToName {
					continue
				}
				// Prefer same directory (same package)
				symFile := strings.TrimPrefix(parts[0], callerPrefix)
				symDir := filepath.Dir(symFile)
				if symDir == callerFileBase {
					matched = []string{sym.ID}
					break
				}
				matched = append(matched, sym.ID)
			}
			for _, to := range matched {
				edges = append(edges, Edge{
					From:       cs.From,
					To:         to,
					Mechanism:  "go_call",
					Confidence: 1.0,
				})
			}
		}
	}

	// Sort edges
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From != edges[j].From {
			return edges[i].From < edges[j].From
		}
		return edges[i].To < edges[j].To
	})

	// Sort unresolved
	sort.Slice(unresolved, func(i, j int) bool {
		if unresolved[i].From != unresolved[j].From {
			return unresolved[i].From < unresolved[j].From
		}
		return unresolved[i].ToName < unresolved[j].ToName
	})

	// Include entry points in graph output
	eps := make([]EntryPoint, len(extract.EntryPoints))
	copy(eps, extract.EntryPoints)

	// Include full symbol table so downstream steps don't need --extract
	syms := make([]Symbol, len(extract.Symbols))
	copy(syms, extract.Symbols)

	return &GraphOutput{
		Nodes:       nodes,
		Edges:       edges,
		Unresolved:  unresolved,
		EntryPoints: eps,
		Symbols:     syms,
	}
}

