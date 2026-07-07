package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"code-ghoul/internal/analysis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: deadcode <analyze|diff> [flags]\n")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "analyze":
		runAnalyze(os.Args[2:])
	case "diff":
		runDiff(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		os.Exit(2)
	}
}

func runAnalyze(args []string) {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	repo := fs.String("repo", ".", "repository directory")
	configFile := fs.String("config", "", "config file path")
	jsonOut := fs.Bool("json", false, "output JSON")
	_ = fs.Bool("human", false, "output human-readable")
	fs.Parse(args)

	plan, err := runPipeline(*repo, *configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "analyze error: %v\n", err)
		os.Exit(2)
	}

	if *jsonOut || true { // default to JSON
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(plan)
	}

	if plan.Summary.Removable > 0 {
		os.Exit(1)
	}
}

func runPipeline(repo, configFile string) (*analysis.PlanOutput, error) {
	// Step 1: Extract
	extract, err := analysis.ExtractDir(repo)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	// Step 2: Graph
	graph := analysis.BuildGraph(extract)

	// Step 3: Reach
	var cfg *analysis.Config
	if configFile != "" {
		cfg, err = analysis.LoadConfig(configFile)
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
	}
	reach := analysis.ComputeReachability(graph, cfg)

	// Step 4: Plan
	plan := analysis.BuildPlan(reach)

	return plan, nil
}

func runDiff(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	repo := fs.String("repo", ".", "repository directory")
	base := fs.String("base", "", "base commit")
	head := fs.String("head", "HEAD", "head commit")
	configFile := fs.String("config", "", "config file path")
	_ = fs.Bool("json", true, "output JSON (always true)")
	fs.Parse(args)

	if *base == "" {
		fmt.Fprintf(os.Stderr, "diff: --base is required\n")
		os.Exit(2)
	}

	// Analyze base commit
	basePlan, err := analyzeCommit(*repo, *base, *configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "diff: analyze base: %v\n", err)
		os.Exit(2)
	}

	// Analyze head commit
	headPlan, err := analyzeCommit(*repo, *head, *configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "diff: analyze head: %v\n", err)
		os.Exit(2)
	}

	// Compute delta
	baseSet := map[string]bool{}
	for _, r := range basePlan.Removals {
		baseSet[r.ID] = true
	}
	headSet := map[string]bool{}
	for _, r := range headPlan.Removals {
		headSet[r.ID] = true
	}

	newDead := []string{}
	resolvedDead := []string{}
	for id := range headSet {
		if !baseSet[id] {
			newDead = append(newDead, id)
		}
	}
	for id := range baseSet {
		if !headSet[id] {
			resolvedDead = append(resolvedDead, id)
		}
	}
	sort.Strings(newDead)
	sort.Strings(resolvedDead)

	result := map[string]interface{}{
		"base_removals": basePlan.Removals,
		"head_removals": headPlan.Removals,
		"delta": map[string]interface{}{
			"new_dead":      newDead,
			"resolved_dead": resolvedDead,
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result)
}

func analyzeCommit(repo, commit, configFile string) (*analysis.PlanOutput, error) {
	// List files in commit
	cmd := exec.Command("git", "-C", repo, "ls-tree", "-r", "--name-only", commit)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ls-tree: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "deadcode-commit-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Extract files
	for _, file := range files {
		if file == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file))
		switch ext {
		case ".go", ".rs", ".java", ".kt":
		default:
			continue
		}

		// Read file content from git
		showCmd := exec.Command("git", "-C", repo, "show", commit+":"+file)
		content, err := showCmd.Output()
		if err != nil {
			continue
		}

		// Write to temp dir
		dest := filepath.Join(tmpDir, file)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(dest, content, 0644); err != nil {
			return nil, err
		}
	}

	return runPipeline(tmpDir, configFile)
}
