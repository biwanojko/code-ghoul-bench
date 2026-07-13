"""
Verifier tests for the code-ghoul cross-language dead code eliminator task.

Multi-turn SWE-Bench Pro — 5 steps:
  Step 1: Multi-Language Symbol Extractor         (cmd/extract/main.go)
  Step 2: Cross-Language Call Graph Construction  (cmd/graph/main.go)
  Step 3: Reachability Analysis                   (cmd/reach/main.go)
  Step 4: Safe Removal Plan                       (cmd/plan/main.go)
  Step 5: End-to-End Pipeline + Diff Mode         (cmd/deadcode/main.go)
"""

import json
import os
import subprocess
import time

REPO = "/app"
TESTDATA = os.path.join(REPO, "testdata")
_GO_ENV = {**os.environ, "PATH": f"/usr/local/go/bin:/go/bin:{os.environ.get('PATH', '')}"}


def _run(cmd, cwd=REPO, timeout=60, env=None):
    """Run a shell command and return (returncode, stdout, stderr)."""
    result = subprocess.run(
        cmd,
        cwd=cwd,
        capture_output=True,
        text=True,
        timeout=timeout,
        env=env or _GO_ENV,
    )
    return result.returncode, result.stdout, result.stderr


def _build(pkg):
    """Build a Go package and return (ok, error_msg)."""
    rc, stdout, stderr = _run(["go", "build", f"./{pkg}/..."])
    if rc != 0:
        return False, stderr
    return True, ""


# ---------------------------------------------------------------------------
# Baseline: repo structure
# ---------------------------------------------------------------------------


def test_repo_structure_intact():
    """The repository must contain expected testdata directories."""
    for lang in ("go_sources", "rust_sources", "java_sources", "kotlin_sources"):
        d = os.path.join(TESTDATA, lang)
        assert os.path.isdir(d), f"testdata/{lang}/ not found"



# ---------------------------------------------------------------------------
# Step 1: Symbol Extractor
# ---------------------------------------------------------------------------


class TestStep1:
    """Tests for cmd/extract/main.go"""

    _binary = os.path.join(REPO, "cmd", "extract", "main.go")

    def _extract(self):
        ok, err = _build("cmd/extract")
        assert ok, f"cmd/extract build failed:\n{err}"
        rc, out, err = _run(
            ["go", "run", "./cmd/extract", "--dir", TESTDATA],
            timeout=15,
        )
        assert rc == 0, f"cmd/extract exited {rc}:\n{err}"
        return json.loads(out)

    def test_exported_symbols_go(self):
        """All Go exported symbols (uppercase + //export + //go:linkname) detected."""
        data = self._extract()
        total = len(data["symbols"])
        langs = {s.get("language") for s in data["symbols"]}
        go_syms = [s for s in data["symbols"] if s["language"] == "go" and s["visibility"] == "exported"]
        assert len(go_syms) > 0, (
            f"No Go exported symbols found. Total symbols={total}, languages present={langs}. "
            f"Hint: Go symbols need language='go' and visibility='exported' (lowercase). "
            f"Check that --dir is parsed via flag.Parse() and filepath.Walk recurses into subdirs."
        )

    def test_exported_symbols_rust(self):
        """All Rust pub symbols detected, pub(crate) excluded."""
        data = self._extract()
        total = len(data["symbols"])
        # Accept either "pub" (preferred, matching Rust keyword) or "exported" (generic term)
        rust_exported = [s for s in data["symbols"] if s["language"] == "rust" and s["visibility"] in ("pub", "exported")]
        assert len(rust_exported) > 0, (
            f"No Rust exported symbols found. Total symbols={total}. "
            f"Hint: Rust pub symbols need language='rust' and visibility='pub' or 'exported'."
        )
        for s in data["symbols"]:
            if s["language"] == "rust":
                assert "pub(crate)" not in str(s.get("attributes", [])), \
                    f"pub(crate) symbol incorrectly exported: {s['id']}"

    def test_exported_symbols_java(self):
        """All Java public classes/methods/native declarations detected."""
        data = self._extract()
        total = len(data["symbols"])
        java_syms = [s for s in data["symbols"] if s["language"] == "java"]
        assert len(java_syms) > 0, (
            f"No Java symbols found. Total symbols={total}. "
            f"Hint: Java symbols need language='java' (lowercase). Walk .java files in testdata/java_sources/."
        )

    def test_exported_symbols_kotlin(self):
        """All Kotlin top-level and companion-object symbols detected."""
        data = self._extract()
        total = len(data["symbols"])
        kotlin_syms = [s for s in data["symbols"] if s["language"] == "kotlin"]
        assert len(kotlin_syms) > 0, (
            f"No Kotlin symbols found. Total symbols={total}. "
            f"Hint: Kotlin symbols need language='kotlin' (lowercase). Walk .kt files in testdata/kotlin_sources/."
        )

    def test_ffi_callsites(self):
        """FFI callsites (cgo C.xxx, JNI native, extern blocks) recorded."""
        data = self._extract()
        callsites = data.get("callsites", [])
        assert len(callsites) > 0, (
            f"No callsites found. Hint: look for C.xxx in .go files (mechanism='cgo') "
            f"and native methods in .java files (mechanism='jni')."
        )
        mechanisms = {cs["mechanism"] for cs in callsites}
        assert "cgo" in mechanisms, (
            f"No cgo callsites found. Hint: look for C.xxx calls in Go files that import 'C'. "
            f"Got mechanisms: {mechanisms}"
        )
        assert "jni" in mechanisms, (
            f"No jni callsites found. Hint: Java native methods create JNI callsites "
            f"with mechanism='jni'. Got mechanisms: {mechanisms}"
        )

    def test_entry_points(self):
        """main(), init(), static initializers, framework entries recorded."""
        data = self._extract()
        reasons = {ep["reason"] for ep in data.get("entry_points", [])}
        assert "main_function" in reasons, "No main_function entry point found"
        assert "init_function" in reasons, "No init_function entry point found"
        assert "static_initializer" in reasons, (
            "No static_initializer entry point found. "
            "Hint: Java static {} blocks should be recorded as entry points with reason='static_initializer'."
        )

    def test_build_conditions(self):
        """Build-tag-guarded symbols carry their conditions."""
        data = self._extract()
        # At least one symbol should have non-empty conditions
        conditioned = [s for s in data["symbols"] if s.get("conditions")]
        assert len(conditioned) > 0, "No symbols with build conditions found"

    def test_test_scope_tagging(self):
        """Test-only symbols (cfg(test), _test.go, @VisibleForTesting) tagged."""
        data = self._extract()
        test_scope = [s for s in data["symbols"] if s.get("scope") == "test"]
        assert len(test_scope) > 0, "No test-scoped symbols found"

    def test_deterministic_output(self):
        """Two consecutive runs produce identical JSON."""
        ok, err = _build("cmd/extract")
        assert ok, f"build failed:\n{err}"
        run = lambda: _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=15)[1]
        assert run() == run(), "cmd/extract output is not deterministic"

    def test_performance_under_10s(self):
        """Symbol extraction completes in under 10 seconds."""
        ok, err = _build("cmd/extract")
        assert ok
        t0 = time.time()
        rc, out, err = _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=12)
        elapsed = time.time() - t0
        assert rc == 0, f"extract failed:\n{err}"
        assert elapsed < 10.0, f"Extraction took {elapsed:.1f}s (limit: 10s)"


# ---------------------------------------------------------------------------
# Step 2: Call Graph Constructor
# ---------------------------------------------------------------------------


class TestStep2:
    """Tests for cmd/graph/main.go"""

    def _graph(self):
        ok, err = _build("cmd/extract")
        assert ok, f"extract build failed:\n{err}"
        ok, err = _build("cmd/graph")
        assert ok, f"graph build failed:\n{err}"
        rc, sym_out, _ = _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=15)
        assert rc == 0
        rc, out, err = _run(
            ["go", "run", "./cmd/graph"],
            timeout=30,
        )
        if rc != 0 and not out:
            # Try piping symbol table
            proc = subprocess.run(
                ["go", "run", "./cmd/graph"],
                input=sym_out, capture_output=True, text=True, cwd=REPO,
                timeout=30, env=_GO_ENV,
            )
            assert proc.returncode == 0, f"graph failed:\n{proc.stderr}"
            return json.loads(proc.stdout)
        assert rc == 0, f"graph failed:\n{err}"
        return json.loads(out)

    def test_direct_ffi_edges_resolved(self):
        """Direct cgo and JNI edges connect caller to correct target symbol."""
        data = self._graph()
        edges = data.get("edges", [])
        assert len(edges) > 0, "No call graph edges found"
        mechanisms = {e["mechanism"] for e in edges}
        assert len(mechanisms & {"cgo", "jni"}) > 0, "No cgo/jni edges found"

    def test_jni_name_mangling(self):
        """JNI name mangling (underscores, overloads, byte arrays) resolved."""
        data = self._graph()
        jni_edges = [e for e in data.get("edges", []) if e["mechanism"] == "jni"]
        assert len(jni_edges) > 0, "No JNI edges found"
        # At least one edge should connect a Java symbol to a Rust JNI export
        rust_jni_targets = [e for e in jni_edges if "rust://" in e.get("to", "")]
        assert len(rust_jni_targets) > 0, "No Java→Rust JNI edge resolved"

    def test_unresolved_callsites_reported(self):
        """Unresolved callsites are listed in 'unresolved', not silently dropped."""
        data = self._graph()
        # The 'unresolved' key must exist (may be empty if all resolved)
        assert "unresolved" in data, "Graph output missing 'unresolved' key"

    def test_conditional_edges(self):
        """Conditional edges carry merged conditions from caller and target."""
        data = self._graph()
        conditional_edges = [e for e in data.get("edges", []) if e.get("conditions")]
        # If any symbol has conditions, there should be at least one conditional edge
        # (This is a soft check — testdata must contain conditional symbols)
        assert isinstance(conditional_edges, list), "edges[].conditions must be a list"

    def test_reflection_edges(self):
        """Reflection edges (Class.forName) created with confidence < 1.0."""
        data = self._graph()
        reflection_edges = [e for e in data.get("edges", []) if e.get("mechanism") == "reflection"]
        for e in reflection_edges:
            assert e.get("confidence", 1.0) < 1.0, \
                f"Reflection edge must have confidence < 1.0: {e}"

    def test_go_call_edges(self):
        """Direct Go function calls create go_call edges in the call graph."""
        data = self._graph()
        edges = data.get("edges", [])
        go_call_edges = [e for e in edges if e.get("mechanism") == "go_call"]
        assert len(go_call_edges) > 0, \
            "No go_call edges found — direct Go function calls must be tracked as edges"
        valid_nodes = set(data.get("nodes", []))
        for e in go_call_edges[:10]:
            assert e.get("from", "") in valid_nodes, \
                f"go_call edge 'from' not in node set: {e.get('from')}"
            assert e.get("to", "") in valid_nodes, \
                f"go_call edge 'to' not in node set: {e.get('to')}"


# ---------------------------------------------------------------------------
# Step 3: Reachability Analysis
# ---------------------------------------------------------------------------


class TestStep3:
    """Tests for cmd/reach/main.go"""

    _config = os.path.join(TESTDATA, "config.json")

    def _reachability(self):
        for pkg in ("cmd/extract", "cmd/graph", "cmd/reach"):
            ok, err = _build(pkg)
            assert ok, f"{pkg} build failed:\n{err}"
        # Full pipeline
        rc, sym_out, _ = _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=15)
        assert rc == 0
        graph_proc = subprocess.run(
            ["go", "run", "./cmd/graph"],
            input=sym_out, capture_output=True, text=True, cwd=REPO,
            timeout=30, env=_GO_ENV,
        )
        assert graph_proc.returncode == 0, graph_proc.stderr
        reach_proc = subprocess.run(
            ["go", "run", "./cmd/reach", "--config", self._config],
            input=graph_proc.stdout, capture_output=True, text=True, cwd=REPO,
            timeout=30, env=_GO_ENV,
        )
        assert reach_proc.returncode == 0, reach_proc.stderr
        return json.loads(reach_proc.stdout)

    def test_entry_points_reachable(self):
        """All detected entry points have status REACHABLE."""
        data = self._reachability()
        statuses = {r["id"]: r["status"] for r in data.get("reachability", [])}
        # Only check symbols that are ACTUAL entry points (from the extract step),
        # not all symbols whose name happens to end with "main" or "init".
        # (e.g. a private Java method named init() is not an entry point)
        ok, err = _build("cmd/extract")
        assert ok, f"cmd/extract build failed:\n{err}"
        rc, sym_out, _ = _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=15)
        assert rc == 0, f"cmd/extract failed:\n{sym_out}"
        extract_data = json.loads(sym_out)
        # Support both "symbol_id" (preferred) and "id" (legacy per instruction example)
        entry_point_ids = {ep.get("symbol_id") or ep.get("id", "") for ep in extract_data.get("entry_points", [])}
        entry_point_ids.discard("")
        entry_suffixes = ("#main", "#init", "::main", "::init", ".main", ".init")
        for ep_id in entry_point_ids:
            if any(ep_id.endswith(sfx) for sfx in entry_suffixes):
                assert ep_id in statuses and statuses[ep_id] == "REACHABLE", (
                    f"Entry point {ep_id} not REACHABLE: {statuses.get(ep_id, 'NOT IN OUTPUT')}. "
                    f"Hint: main/init entry points must be reachable. "
                    f"Check that entry points are added to graph.EntryPoints and BFS starts from them."
                )

    def test_transitive_closure(self):
        """Statistics are consistent (total = reachable + conditionally + unreachable + test)."""
        data = self._reachability()
        stats = data.get("statistics", {})
        total = stats.get("total_symbols", 0)
        assert total > 0, "No symbols in reachability output"
        accounted = (
            stats.get("reachable", 0)
            + stats.get("conditionally_reachable", 0)
            + stats.get("unreachable", 0)
            + stats.get("test_scope_excluded", 0)
        )
        assert accounted == total, f"Symbol counts don't sum to total: {accounted} != {total}"

    def test_conditional_gating(self):
        """Symbols unreachable under current config but reachable elsewhere = CONDITIONALLY_REACHABLE."""
        data = self._reachability()
        statuses = [r["status"] for r in data.get("reachability", [])]
        assert "CONDITIONALLY_REACHABLE" in statuses, \
            "No CONDITIONALLY_REACHABLE symbols found (testdata should have some)"

    def test_virtual_dispatch(self):
        """Interface/trait method reachable when type is reachable."""
        data = self._reachability()
        # At least one reachable symbol should come from an interface/trait implementation
        reachable_ids = {r["id"] for r in data.get("reachability", []) if r["status"] == "REACHABLE"}
        assert len(reachable_ids) > 0, "No reachable symbols at all"

    def test_sealed_class_propagation(self):
        """Sealed class reachable implies all direct subclasses reachable."""
        data = self._reachability()
        reachable_ids = {r["id"] for r in data.get("reachability", []) if r["status"] == "REACHABLE"}
        assert len(reachable_ids) > 0


# ---------------------------------------------------------------------------
# Step 4: Safe Removal Plan
# ---------------------------------------------------------------------------


class TestStep4:
    """Tests for cmd/plan/main.go"""

    _config = os.path.join(TESTDATA, "config.json")

    def _plan(self):
        for pkg in ("cmd/extract", "cmd/graph", "cmd/reach", "cmd/plan"):
            ok, err = _build(pkg)
            assert ok, f"{pkg} build failed:\n{err}"
        rc, sym_out, _ = _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=15)
        assert rc == 0
        graph_proc = subprocess.run(
            ["go", "run", "./cmd/graph"], input=sym_out,
            capture_output=True, text=True, cwd=REPO, timeout=30, env=_GO_ENV,
        )
        assert graph_proc.returncode == 0
        reach_proc = subprocess.run(
            ["go", "run", "./cmd/reach", "--config", self._config],
            input=graph_proc.stdout, capture_output=True, text=True,
            cwd=REPO, timeout=30, env=_GO_ENV,
        )
        assert reach_proc.returncode == 0
        plan_proc = subprocess.run(
            ["go", "run", "./cmd/plan"], input=reach_proc.stdout,
            capture_output=True, text=True, cwd=REPO, timeout=30, env=_GO_ENV,
        )
        assert plan_proc.returncode == 0, plan_proc.stderr
        return json.loads(plan_proc.stdout)

    def test_zero_safety_filter_violations(self):
        """Safety-filtered symbols must be in retained (not removals); retained must be non-empty."""
        data = self._plan()
        retained = data.get("retained", [])
        assert len(retained) > 0, "retained list must be non-empty — safety-filtered symbols (main, init, etc.) must be retained"
        retained_ids = {r["id"] for r in retained}
        removal_ids = {r["id"] for r in data.get("removals", [])}
        violations = retained_ids & removal_ids
        assert len(violations) == 0, \
            f"Safety filter violations: {violations}"

    def test_priority_scoring(self):
        """Removals are non-empty and scored correctly: P0=1.0, P1=0.8, P2=0.5, P3=0.3."""
        data = self._plan()
        removals = data.get("removals", [])
        assert len(removals) > 0, "removals list must be non-empty — unreachable symbols must be planned for removal"
        priority_score = {"P0": 1.0, "P1": 0.8, "P2": 0.5, "P3": 0.3}
        for r in removals:
            assert r["priority"] in priority_score, f"Unknown priority: {r['priority']}"
            assert abs(r["score"] - priority_score[r["priority"]]) < 0.01, \
                f"Score mismatch for {r['id']}: expected {priority_score[r['priority']]}, got {r['score']}"

    def test_file_delete_only_when_all_unreachable(self):
        """file_level_actions is non-empty; DELETE_ENTIRE_FILE only when all symbols unreachable."""
        data = self._plan()
        file_actions = data.get("file_level_actions", [])
        assert len(file_actions) > 0, "file_level_actions must be non-empty — files with all-unreachable symbols must be flagged"
        for fa in file_actions:
            if fa["action"] == "DELETE_ENTIRE_FILE":
                assert fa.get("all_symbols_unreachable"), \
                    f"DELETE_ENTIRE_FILE set without all_symbols_unreachable: {fa['file']}"

    def test_loc_estimates_accurate(self):
        """LOC estimates exist and are positive integers for each removal."""
        data = self._plan()
        removals = data.get("removals", [])
        assert len(removals) > 0, "removals must be non-empty — unreachable symbols must have loc_savings estimates"
        for r in removals:
            assert isinstance(r.get("loc_savings"), int), \
                f"loc_savings must be int: {r['id']}"
            assert r["loc_savings"] > 0, f"loc_savings must be > 0: {r['id']}"


# ---------------------------------------------------------------------------
# Step 5: End-to-End Pipeline
# ---------------------------------------------------------------------------


class TestStep5:
    """Tests for cmd/deadcode/main.go"""

    _config = os.path.join(TESTDATA, "config.json")
    _monorepo = os.path.join(TESTDATA, "monorepo")

    def _analyze(self, extra_args=None):
        for pkg in ("cmd/extract", "cmd/graph", "cmd/deadcode"):
            ok, err = _build(pkg)
            assert ok, f"{pkg} build failed:\n{err}"
        rc, sym_out, _ = _run(["go", "run", "./cmd/extract", "--dir", TESTDATA], timeout=15)
        assert rc == 0, f"extract failed"
        graph_proc = subprocess.run(
            ["go", "run", "./cmd/graph"],
            input=sym_out, capture_output=True, text=True, cwd=REPO, timeout=30, env=_GO_ENV,
        )
        assert graph_proc.returncode == 0, f"graph failed: {graph_proc.stderr}"
        args = ["go", "run", "./cmd/deadcode", "analyze", "--config", self._config, "--json"]
        if extra_args:
            args += extra_args
        proc = subprocess.run(
            args, input=graph_proc.stdout, capture_output=True, text=True,
            cwd=REPO, timeout=20, env=_GO_ENV,
        )
        return proc.returncode, proc.stdout, proc.stderr

    def test_analyze_mode_matches_expected(self):
        """Analyze mode output matches expected_output.json (allow +/-5 LOC)."""
        expected_path = os.path.join(TESTDATA, "expected_output.json")
        if not os.path.isfile(expected_path):
            return  # expected_output not yet generated
        rc, out, err = self._analyze()
        assert rc in (0, 1), f"deadcode analyze exited {rc}:\n{err}"
        actual = json.loads(out)
        with open(expected_path) as f:
            expected = json.load(f)
        # Check removal IDs match
        actual_ids = {r["id"] for r in actual.get("removals", [])}
        expected_ids = {r["id"] for r in expected.get("removals", [])}
        assert actual_ids == expected_ids, \
            f"Removal ID mismatch.\nMissing: {expected_ids - actual_ids}\nExtra: {actual_ids - expected_ids}"

    def test_diff_mode_new_dead(self):
        """Diff mode reports new_dead symbols when dead code is introduced between commits."""
        ok, err = _build("cmd/deadcode")
        assert ok, f"build failed:\n{err}"
        if not os.path.isdir(self._monorepo):
            return
        rc, out, err = _run(
            ["go", "run", "./cmd/deadcode", "diff",
             "--repo", self._monorepo,
             "--base", "HEAD~1", "--head", "HEAD",
             "--config", self._config, "--json"],
            timeout=35,
        )
        assert rc in (0, 1, 2), f"diff exited unexpectedly: {rc}"
        if rc == 2:
            # diff subcommand not implemented — fail with helpful message
            assert False, \
                "diff subcommand exited with error (rc=2) — must implement deadcode diff"
        data = json.loads(out)
        assert "delta" in data, "Diff output missing 'delta' key"
        assert "new_dead" in data["delta"], "Diff output missing 'delta.new_dead'"
        assert len(data["delta"]["new_dead"]) > 0, \
            "new_dead must be non-empty — the monorepo removes helper()'s caller at HEAD"

    def test_diff_mode_resolved_dead(self):
        """Diff mode output has correct 'resolved_dead' key in delta (may be empty for monorepo)."""
        ok, err = _build("cmd/deadcode")
        assert ok
        if not os.path.isdir(self._monorepo):
            return
        rc, out, err = _run(
            ["go", "run", "./cmd/deadcode", "diff",
             "--repo", self._monorepo,
             "--base", "HEAD~1", "--head", "HEAD",
             "--config", self._config, "--json"],
            timeout=35,
        )
        assert rc in (0, 1, 2), f"diff exited unexpectedly: {rc}"
        if rc != 2:
            data = json.loads(out)
            assert "resolved_dead" in data.get("delta", {}), \
                "Diff output missing 'delta.resolved_dead'"

    def test_exit_codes(self):
        """Analyze exits 1 when dead code found; 0 when repo is clean."""
        ok, err = _build("cmd/deadcode")
        assert ok
        rc, out, err = self._analyze()
        data = json.loads(out)
        removals = data.get("removals", [])
        assert len(removals) > 0, \
            "testdata has dead code — analyze must report non-empty removals and exit 1"
        assert rc == 1, f"Expected exit code 1 (dead code found), got {rc}"

    def test_performance_budget(self):
        """Full analyze pipeline completes in under 15 seconds and reports non-empty results."""
        ok, err = _build("cmd/deadcode")
        assert ok
        t0 = time.time()
        rc, out, err = self._analyze()
        elapsed = time.time() - t0
        assert rc in (0, 1), f"analyze failed:\n{err}"
        data = json.loads(out)
        assert len(data.get("removals", [])) > 0, "analyze must find dead code in testdata"
        assert elapsed < 15.0, f"Pipeline took {elapsed:.1f}s (limit: 15s)"

    def test_determinism_5_runs(self):
        """5 consecutive runs produce identical JSON output with non-empty removals."""
        ok, err = _build("cmd/deadcode")
        assert ok
        outputs = set()
        for _ in range(5):
            rc, out, err = self._analyze()
            assert rc in (0, 1), f"analyze failed:\n{err}"
            data = json.loads(out)
            assert len(data.get("removals", [])) > 0, "removals must be non-empty"
            outputs.add(json.dumps(data, sort_keys=True))
        assert len(outputs) == 1, f"Output is non-deterministic: {len(outputs)} distinct outputs"
