package analysis

// BuildPlan creates a removal plan from reachability output.
// TODO: implement safety filtering, priority scoring, file-level actions, LOC estimates.
func BuildPlan(reach *ReachOutput) *PlanOutput {
	return &PlanOutput{Removals: []Removal{}, Retained: []Retained{}, FileLevelActions: []FileLevelAction{}, Summary: PlanSummary{TotalSymbols: len(reach.Symbols)}}
}
