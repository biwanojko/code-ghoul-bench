package analysis

import (
	"sort"
	"strings"
)

var safetyAttributes = map[string]bool{
	"used":             true,
	"keep":             true,
	"do_not_strip":     true,
	"panic_handler":    true,
	"global_allocator": true,
}

var safetyKinds = map[string]bool{
	"init_function":      true,
	"static_initializer": true,
}

func locSavings(kind string) int {
	switch kind {
	case "function":
		return 10
	case "method":
		return 8
	case "type", "sealed_class", "singleton", "companion":
		return 5
	case "constant", "variable":
		return 3
	default:
		return 5
	}
}

// priorityLabel maps a numeric score to a priority label
func priorityLabel(score float64) string {
	switch score {
	case 1.0:
		return "P0"
	case 0.8:
		return "P1"
	case 0.5:
		return "P2"
	case 0.3:
		return "P3"
	default:
		return "P1"
	}
}

// BuildPlan creates a removal plan from reachability output.
// The reach output carries symbol metadata (via reach.Symbols) so no separate extract is needed.
func BuildPlan(reach *ReachOutput) *PlanOutput {
	// Build symbol map from reach.Symbols
	symByID := map[string]*Symbol{}
	for i := range reach.Symbols {
		symByID[reach.Symbols[i].ID] = &reach.Symbols[i]
	}

	// Build reachability map
	reachMap := map[string]string{}
	for _, r := range reach.Reachability {
		reachMap[r.ID] = r.Status
	}

	// Compute per-file symbol sets
	fileSymbols := map[string][]string{}
	for _, sym := range reach.Symbols {
		fileSymbols[sym.File] = append(fileSymbols[sym.File], sym.ID)
	}

	var removals []Removal
	var retained []Retained
	removedSet := map[string]bool{}

	for _, sym := range reach.Symbols {
		status := reachMap[sym.ID]

		// Safety filter
		isSafe := false
		for _, attr := range sym.Attributes {
			if safetyAttributes[attr] {
				isSafe = true
				break
			}
		}
		if safetyKinds[sym.Kind] {
			isSafe = true
		}
		if sym.Scope == "test" {
			retained = append(retained, Retained{ID: sym.ID, Reason: "test_scope"})
			continue
		}

		switch status {
		case "UNREACHABLE":
			if isSafe {
				reason := "safety_filter"
				if safetyKinds[sym.Kind] {
					reason = sym.Kind
				}
				retained = append(retained, Retained{ID: sym.ID, Reason: reason})
			} else {
				// P2=0.5 for types that are unreachable but may be test-reachable
				score := 0.8
				if sym.Kind == "type" || sym.Kind == "sealed_class" || sym.Kind == "singleton" || sym.Kind == "companion" {
					score = 0.5
				}
				removals = append(removals, Removal{
					ID:         sym.ID,
					File:       sym.File,
					Line:       sym.Line,
					Priority:   priorityLabel(score),
					Score:      score,
					LocSavings: locSavings(sym.Kind),
					Reason:     "UNREACHABLE",
				})
				removedSet[sym.ID] = true
			}
		case "CONDITIONALLY_REACHABLE":
			if isSafe {
				retained = append(retained, Retained{ID: sym.ID, Reason: "safety_filter"})
			} else {
				removals = append(removals, Removal{
					ID:         sym.ID,
					File:       sym.File,
					Line:       sym.Line,
					Priority:   "P3",
					Score:      0.3,
					LocSavings: locSavings(sym.Kind),
					Reason:     "CONDITIONALLY_REACHABLE",
				})
				removedSet[sym.ID] = true
			}
		default: // REACHABLE or unknown
			retained = append(retained, Retained{ID: sym.ID, Reason: "REACHABLE"})
		}
	}

	// Compute file-level actions
	var fileLevelActions []FileLevelAction
	for file, syms := range fileSymbols {
		// Skip test files for file-level actions
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		var fileRemovals []string
		nonTestSyms := 0
		for _, symID := range syms {
			if sym, ok := symByID[symID]; ok && sym.Scope == "test" {
				continue
			}
			nonTestSyms++
			if removedSet[symID] {
				fileRemovals = append(fileRemovals, symID)
			}
		}

		if len(fileRemovals) == 0 {
			continue
		}

		sort.Strings(fileRemovals)

		allUnreachable := len(fileRemovals) == nonTestSyms && nonTestSyms > 0
		action := "PARTIAL_REMOVAL"
		if allUnreachable {
			action = "DELETE_ENTIRE_FILE"
			// Upgrade priority to P0 for all removals in this file
			for i := range removals {
				if removals[i].File == file {
					removals[i].Priority = "P0"
					removals[i].Score = 1.0
				}
			}
		}

		fileLevelActions = append(fileLevelActions, FileLevelAction{
			File:                  file,
			Action:                action,
			AllSymbolsUnreachable: allUnreachable,
			Removals:              fileRemovals,
		})
	}

	// Sort
	sort.Slice(removals, func(i, j int) bool {
		if removals[i].Score != removals[j].Score {
			return removals[i].Score > removals[j].Score
		}
		return removals[i].ID < removals[j].ID
	})
	sort.Slice(retained, func(i, j int) bool {
		return retained[i].ID < retained[j].ID
	})
	sort.Slice(fileLevelActions, func(i, j int) bool {
		return fileLevelActions[i].File < fileLevelActions[j].File
	})

	// Compute summary
	totalLoc := 0
	for _, r := range removals {
		totalLoc += r.LocSavings
	}
	filesToDelete := 0
	for _, fla := range fileLevelActions {
		if fla.Action == "DELETE_ENTIRE_FILE" {
			filesToDelete++
		}
	}

	return &PlanOutput{
		Removals:         removals,
		Retained:         retained,
		FileLevelActions: fileLevelActions,
		Summary: PlanSummary{
			TotalSymbols:        len(reach.Symbols),
			Removable:           len(removals),
			Retained:            len(retained),
			FilesToDelete:       filesToDelete,
			EstimatedLocSavings: totalLoc,
		},
	}
}
