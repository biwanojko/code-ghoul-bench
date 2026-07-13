package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ExtractDir walks dir and extracts all symbols, callsites, and entry points.
// Used internally by the deadcode diff subcommand to build graphs per commit.
func ExtractDir(dir string) (*ExtractOutput, error) {
	var symbols []Symbol
	var callsites []Callsite
	var entryPoints []EntryPoint

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go":
			syms, calls, eps, e := extractGo(path, rel)
			if e != nil {
				return nil // skip unparseable files
			}
			symbols = append(symbols, syms...)
			callsites = append(callsites, calls...)
			entryPoints = append(entryPoints, eps...)
		case ".rs":
			syms, calls, eps := extractRust(path, rel)
			symbols = append(symbols, syms...)
			callsites = append(callsites, calls...)
			entryPoints = append(entryPoints, eps...)
		case ".java":
			syms, calls, eps := extractJava(path, rel)
			symbols = append(symbols, syms...)
			callsites = append(callsites, calls...)
			entryPoints = append(entryPoints, eps...)
		case ".kt":
			syms, calls, eps := extractKotlin(path, rel)
			symbols = append(symbols, syms...)
			callsites = append(callsites, calls...)
			entryPoints = append(entryPoints, eps...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Deduplicate and sort
	sort.Slice(symbols, func(i, j int) bool { return symbols[i].ID < symbols[j].ID })
	sort.Slice(callsites, func(i, j int) bool {
		if callsites[i].From != callsites[j].From {
			return callsites[i].From < callsites[j].From
		}
		return callsites[i].ToName < callsites[j].ToName
	})
	sort.Slice(entryPoints, func(i, j int) bool { return entryPoints[i].SymbolID < entryPoints[j].SymbolID })

	// Deduplicate entry points
	seen := map[string]bool{}
	var dedupEPs []EntryPoint
	for _, ep := range entryPoints {
		if !seen[ep.SymbolID+":"+ep.Reason] {
			seen[ep.SymbolID+":"+ep.Reason] = true
			dedupEPs = append(dedupEPs, ep)
		}
	}

	return &ExtractOutput{
		Symbols:     symbols,
		Callsites:   callsites,
		EntryPoints: dedupEPs,
	}, nil
}

// goSymID creates a symbol ID for a Go symbol
func goSymID(relFile, name string) string {
	return fmt.Sprintf("go://%s#%s", relFile, name)
}

func extractGo(path, relFile string) ([]Symbol, []Callsite, []EntryPoint, error) {
	isTest := strings.HasSuffix(relFile, "_test.go")

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}
	srcStr := string(src)

	// Detect build constraints
	var conditions []string
	for _, line := range strings.Split(srcStr, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//go:build ") {
			cond := strings.TrimPrefix(line, "//go:build ")
			conditions = append(conditions, strings.TrimSpace(cond))
		}
	}

	// Detect CGO
	hasCGO := strings.Contains(srcStr, `import "C"`) || strings.Contains(srcStr, "import \"C\"")

	// Parse file-level comments for //export and //go:linkname directives
	exportedNames := map[string]string{} // funcName -> export name
	linknameNames := map[string]bool{}
	for _, line := range strings.Split(srcStr, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//export ") {
			name := strings.TrimPrefix(line, "//export ")
			name = strings.TrimSpace(name)
			exportedNames[name] = name
		}
		if strings.HasPrefix(line, "//go:linkname ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				linknameNames[parts[1]] = true
			}
		}
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, nil, err
	}

	var symbols []Symbol
	var callsites []Callsite
	var entryPoints []EntryPoint

	pkgName := f.Name.Name

	// Walk declarations
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			name := d.Name.Name
			line := fset.Position(d.Pos()).Line
			vis := "unexported"
			if ast.IsExported(name) {
				vis = "exported"
			}
			kind := "function"
			if d.Recv != nil {
				kind = "method"
			}

			var attrs []string
			isFFIExported := false
			ffiName := ""

			// Check for //export
			if exportName, ok := exportedNames[name]; ok {
				isFFIExported = true
				ffiName = exportName
				vis = "exported"
			}
			// Check for //go:linkname
			if linknameNames[name] {
				attrs = append(attrs, "linkname")
			}

			sym := Symbol{
				ID:          goSymID(relFile, name),
				Language:    "go",
				Kind:        kind,
				Visibility:  vis,
				File:        relFile,
				Line:        line,
				FFIExported: isFFIExported,
				FFIName:     ffiName,
				Attributes:  attrs,
			}
			if isTest {
				sym.Scope = "test"
			}
			if len(conditions) > 0 {
				sym.Conditions = append(sym.Conditions, conditions...)
			}

			symbols = append(symbols, sym)
			symID := sym.ID

			// Entry points
			if name == "main" && pkgName == "main" {
				entryPoints = append(entryPoints, EntryPoint{SymbolID: symID, Reason: "main_function"})
			}
			if name == "init" {
				entryPoints = append(entryPoints, EntryPoint{SymbolID: symID, Reason: "init_function"})
				// Update the symbol kind
				symbols[len(symbols)-1].Kind = "init_function"
			}
			if isFFIExported {
				entryPoints = append(entryPoints, EntryPoint{SymbolID: symID, Reason: "ffi_export"})
			}

			// Walk function body for calls
			if d.Body != nil {
				ast.Inspect(d.Body, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					callLine := fset.Position(call.Pos()).Line
					switch fn := call.Fun.(type) {
					case *ast.SelectorExpr:
						if ident, ok := fn.X.(*ast.Ident); ok && ident.Name == "C" && hasCGO {
							callsites = append(callsites, Callsite{
								From:      symID,
								ToName:    fn.Sel.Name,
								Mechanism: "cgo",
								File:      relFile,
								Line:      callLine,
							})
						}
					case *ast.Ident:
						// Regular Go function call (unqualified)
						callsites = append(callsites, Callsite{
							From:      symID,
							ToName:    fn.Name,
							Mechanism: "go_call",
							File:      relFile,
							Line:      callLine,
						})
					}
					return true
				})
			}

		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					name := s.Name.Name
					line := fset.Position(s.Pos()).Line
					vis := "unexported"
					if ast.IsExported(name) {
						vis = "exported"
					}
					sym := Symbol{
						ID:         goSymID(relFile, name),
						Language:   "go",
						Kind:       "type",
						Visibility: vis,
						File:       relFile,
						Line:       line,
					}
					if isTest {
						sym.Scope = "test"
					}
					if len(conditions) > 0 {
						sym.Conditions = append(sym.Conditions, conditions...)
					}
					symbols = append(symbols, sym)

				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.Name == "_" {
							continue
						}
						line := fset.Position(name.Pos()).Line
						vis := "unexported"
						if ast.IsExported(name.Name) {
							vis = "exported"
						}
						kind := "constant"
						if d.Tok == token.VAR {
							kind = "variable"
						}
						sym := Symbol{
							ID:         goSymID(relFile, name.Name),
							Language:   "go",
							Kind:       kind,
							Visibility: vis,
							File:       relFile,
							Line:       line,
						}
						if isTest {
							sym.Scope = "test"
						}
						if len(conditions) > 0 {
							sym.Conditions = append(sym.Conditions, conditions...)
						}
						symbols = append(symbols, sym)
					}
				}
			}
		}
	}

	return symbols, callsites, entryPoints, nil
}

// Rust extraction using regex
var (
	rsExternCRe = regexp.MustCompile(`(?m)extern\s+"C"\s*\{[^}]*fn\s+(\w+)`)
)

func lineOf(src string, offset int) int {
	return strings.Count(src[:offset], "\n") + 1
}

func extractRust(path, relFile string) ([]Symbol, []Callsite, []EntryPoint) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil
	}
	srcStr := string(src)

	var symbols []Symbol
	var callsites []Callsite
	var entryPoints []EntryPoint

	lines := strings.Split(srcStr, "\n")

	inCfgTest := false
	pendingNoMangle := false
	pendingExportName := ""
	pendingUsed := false
	pendingCfgFeature := ""

	type fnInfo struct {
		name        string
		line        int
		visibility  string
		ffiExported bool
		ffiName     string
		scope       string
		conditions  []string
		attributes  []string
	}

	var fns []fnInfo

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		lineNum := i + 1

		if strings.Contains(trimmed, "#[cfg(test)]") {
			inCfgTest = true
		}
		if strings.Contains(trimmed, "#[no_mangle]") {
			pendingNoMangle = true
		}
		if strings.Contains(trimmed, "#[used]") {
			pendingUsed = true
		}

		// Check export_name
		enMatch := regexp.MustCompile(`#\[export_name\s*=\s*"([^"]+)"\]`).FindStringSubmatch(trimmed)
		if len(enMatch) > 1 {
			pendingExportName = enMatch[1]
		}

		// Check cfg feature
		cfMatch := regexp.MustCompile(`#\[cfg\(feature\s*=\s*"([^"]+)"\)\]`).FindStringSubmatch(trimmed)
		if len(cfMatch) > 1 {
			pendingCfgFeature = cfMatch[1]
		}

		// Detect function declarations
		fnMatch := regexp.MustCompile(`^(pub(?:\(crate\))?\s+)?(?:unsafe\s+)?(?:async\s+)?(?:extern\s+"C"\s+)?fn\s+(\w+)`).FindStringSubmatch(trimmed)
		if fnMatch != nil {
			name := fnMatch[2]
			vis := "private"
			pubPart := strings.TrimSpace(fnMatch[1])
			if pubPart == "pub" {
				vis = "pub"
			} else if strings.HasPrefix(pubPart, "pub(crate)") {
				vis = "pub(crate)"
			}

			isFFI := pendingNoMangle && strings.Contains(trimmed, "extern")
			ffiName := name
			if pendingExportName != "" {
				ffiName = pendingExportName
				isFFI = true
			}

			var attrs []string
			if pendingUsed {
				attrs = append(attrs, "used")
			}

			scope := ""
			if inCfgTest {
				scope = "test"
			}

			var conds []string
			if pendingCfgFeature != "" {
				conds = append(conds, "feature="+pendingCfgFeature)
			}

			fi := fnInfo{
				name:        name,
				line:        lineNum,
				visibility:  vis,
				ffiExported: isFFI,
				ffiName:     ffiName,
				scope:       scope,
				conditions:  conds,
				attributes:  attrs,
			}
			fns = append(fns, fi)

			if isFFI {
				epReason := "ffi_export"
				if strings.HasPrefix(name, "Java_") {
					epReason = "jni_export"
				}
				entryPoints = append(entryPoints, EntryPoint{
					SymbolID: fmt.Sprintf("rust://%s#%s", relFile, name),
					Reason:   epReason,
				})
			}

			// Reset pending state
			pendingNoMangle = false
			pendingExportName = ""
			pendingUsed = false
			pendingCfgFeature = ""
		} else {
			// Only reset if we see a non-attribute, non-comment, non-blank line that isn't a fn
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "//") {
				if !strings.HasPrefix(trimmed, "pub") || !strings.Contains(trimmed, "fn") {
					if !strings.HasPrefix(trimmed, "unsafe") && !strings.HasPrefix(trimmed, "async") && !strings.HasPrefix(trimmed, "extern") {
						pendingNoMangle = false
						pendingExportName = ""
						pendingUsed = false
						pendingCfgFeature = ""
					}
				}
			}
		}
	}

	// Also detect extern "C" { fn ... } as callsites from this module
	externCMatches := rsExternCRe.FindAllStringSubmatch(srcStr, -1)
	for _, m := range externCMatches {
		if len(m) > 1 {
			callsites = append(callsites, Callsite{
				From:      fmt.Sprintf("rust://%s#<module>", relFile),
				ToName:    m[1],
				Mechanism: "ffi",
				File:      relFile,
				Line:      1,
			})
		}
	}

	// Convert fns to symbols
	for _, fi := range fns {
		sym := Symbol{
			ID:          fmt.Sprintf("rust://%s#%s", relFile, fi.name),
			Language:    "rust",
			Kind:        "function",
			Visibility:  fi.visibility,
			File:        relFile,
			Line:        fi.line,
			FFIExported: fi.ffiExported,
			Attributes:  fi.attributes,
			Scope:       fi.scope,
			Conditions:  fi.conditions,
		}
		if fi.ffiExported && fi.ffiName != "" && fi.ffiName != fi.name {
			sym.FFIName = fi.ffiName
		} else if fi.ffiExported {
			sym.FFIName = fi.name
		}
		symbols = append(symbols, sym)
	}

	return symbols, callsites, entryPoints
}

// Java extraction
func extractJava(path, relFile string) ([]Symbol, []Callsite, []EntryPoint) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil
	}
	srcStr := string(src)
	lines := strings.Split(srcStr, "\n")

	var symbols []Symbol
	var callsites []Callsite
	var entryPoints []EntryPoint

	pkgRe := regexp.MustCompile(`^\s*package\s+([\w.]+)`)
	classRe := regexp.MustCompile(`^\s*(public|protected|private)?\s*(?:static\s+)?(?:final\s+)?(?:abstract\s+)?class\s+(\w+)`)
	methodRe := regexp.MustCompile(`^\s*(public|protected|private)?\s*(static\s+)?(native\s+)?(?:final\s+)?[\w<>\[\]]+\s+(\w+)\s*\(`)
	staticInitRe := regexp.MustCompile(`^\s*static\s*\{`)
	keepRe := regexp.MustCompile(`@Keep`)
	vftRe := regexp.MustCompile(`@VisibleForTesting`)
	loadLibRe := regexp.MustCompile(`System\.loadLibrary\("([^"]+)"\)`)
	classForNameRe := regexp.MustCompile(`Class\.forName\("([^"]+)"\)`)

	pkg := ""
	className := ""
	pendingKeep := false
	pendingVFT := false

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if m := pkgRe.FindStringSubmatch(trimmed); m != nil {
			pkg = m[1]
		}

		if strings.Contains(trimmed, "@Keep") {
			pendingKeep = true
		}
		if strings.Contains(trimmed, "@VisibleForTesting") {
			pendingVFT = true
		}

		// Class declaration
		if m := classRe.FindStringSubmatch(line); m != nil {
			className = m[2]
			id := fmt.Sprintf("java://%s.%s", pkg, className)
			vis := "package"
			if m[1] == "public" {
				vis = "public"
			} else if m[1] == "protected" {
				vis = "protected"
			} else if m[1] == "private" {
				vis = "private"
			}
			var attrs []string
			if pendingKeep {
				attrs = append(attrs, "keep")
			}
			symbols = append(symbols, Symbol{
				ID:         id,
				Language:   "java",
				Kind:       "type",
				Visibility: vis,
				File:       relFile,
				Line:       lineNum,
				Attributes: attrs,
			})
			pendingKeep = false
			pendingVFT = false
			continue
		}

		// Static initializer
		if staticInitRe.MatchString(line) {
			id := fmt.Sprintf("java://%s.%s#<static_init>", pkg, className)
			symbols = append(symbols, Symbol{
				ID:         id,
				Language:   "java",
				Kind:       "static_initializer",
				Visibility: "private",
				File:       relFile,
				Line:       lineNum,
			})
			entryPoints = append(entryPoints, EntryPoint{SymbolID: id, Reason: "static_initializer"})
			continue
		}

		// Method declaration
		if m := methodRe.FindStringSubmatch(line); m != nil {
			methodName := m[4]
			if methodName == className {
				pendingKeep = false
				pendingVFT = false
				continue
			}
			vis := "package"
			if m[1] == "public" {
				vis = "public"
			} else if m[1] == "protected" {
				vis = "protected"
			} else if m[1] == "private" {
				vis = "private"
			}
			isNative := strings.Contains(m[3], "native")

			id := fmt.Sprintf("java://%s.%s#%s", pkg, className, methodName)
			var attrs []string
			if pendingKeep {
				attrs = append(attrs, "keep")
			}
			if pendingVFT {
				attrs = append(attrs, "visible_for_testing")
			}

			sym := Symbol{
				ID:          id,
				Language:    "java",
				Kind:        "method",
				Visibility:  vis,
				File:        relFile,
				Line:        lineNum,
				FFIExported: isNative,
				Attributes:  attrs,
			}
			if isNative {
				jniName := jniMangle(pkg, className, methodName)
				sym.FFIName = jniName
				entryPoints = append(entryPoints, EntryPoint{SymbolID: id, Reason: "jni_export"})
				callsites = append(callsites, Callsite{
					From:      id,
					ToName:    jniName,
					Mechanism: "jni",
					File:      relFile,
					Line:      lineNum,
				})
			}
			symbols = append(symbols, sym)
			pendingKeep = false
			pendingVFT = false
			continue
		}

		// System.loadLibrary
		if m := loadLibRe.FindStringSubmatch(line); m != nil {
			ownerID := fmt.Sprintf("java://%s.%s", pkg, className)
			callsites = append(callsites, Callsite{
				From:      ownerID,
				ToName:    m[1],
				Mechanism: "jni_load",
				File:      relFile,
				Line:      lineNum,
			})
		}

		// Class.forName
		if m := classForNameRe.FindStringSubmatch(line); m != nil {
			ownerID := fmt.Sprintf("java://%s.%s", pkg, className)
			callsites = append(callsites, Callsite{
				From:      ownerID,
				ToName:    m[1],
				Mechanism: "reflection",
				File:      relFile,
				Line:      lineNum,
			})
		}

		// Reset pending annotations for non-annotation lines
		if trimmed != "" && !strings.HasPrefix(trimmed, "@") && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") && !strings.HasPrefix(trimmed, "*") {
			if !strings.Contains(trimmed, "class ") && !strings.Contains(trimmed, "interface ") {
				if !keepRe.MatchString(line) && !vftRe.MatchString(line) {
					pendingKeep = false
					pendingVFT = false
				}
			}
		}
	}

	return symbols, callsites, entryPoints
}

// jniMangle converts a Java native method to its JNI C name
func jniMangle(pkg, class, method string) string {
	manglePkg := strings.ReplaceAll(pkg, ".", "_")
	return fmt.Sprintf("Java_%s_%s_%s", manglePkg, class, method)
}

// Kotlin extraction
func extractKotlin(path, relFile string) ([]Symbol, []Callsite, []EntryPoint) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil
	}
	srcStr := string(src)
	lines := strings.Split(srcStr, "\n")

	var symbols []Symbol
	var entryPoints []EntryPoint

	pkgRe := regexp.MustCompile(`^\s*package\s+([\w.]+)`)
	funRe := regexp.MustCompile(`^\s*(?:private\s+|protected\s+|internal\s+|public\s+)?(?:override\s+)?(?:suspend\s+)?(?:inline\s+)?(?:external\s+)?fun\s+(?:\w+\.)?\w+\s*(?:<[^>]*>)?\s*\(`)
	funNameRe := regexp.MustCompile(`fun\s+(\w+)`)
	externalFunRe := regexp.MustCompile(`\bexternal\s+fun\b`)
	jvmStaticRe := regexp.MustCompile(`@JvmStatic`)
	jvmNameRe := regexp.MustCompile(`@JvmName\("([^"]+)"\)`)
	objectRe := regexp.MustCompile(`^\s*(?:private\s+|internal\s+|public\s+)?object\s+(\w+)`)
	sealedClassRe := regexp.MustCompile(`^\s*sealed\s+class\s+(\w+)`)
	companionRe := regexp.MustCompile(`^\s*companion\s+object`)
	extFunRe := regexp.MustCompile(`fun\s+\w+\.\w+`)

	pkg := ""
	pendingJvmStatic := false
	pendingJvmName := ""

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if m := pkgRe.FindStringSubmatch(trimmed); m != nil {
			pkg = m[1]
		}

		if jvmStaticRe.MatchString(trimmed) {
			pendingJvmStatic = true
		}
		if m := jvmNameRe.FindStringSubmatch(trimmed); m != nil {
			pendingJvmName = m[1]
		}

		// Object declaration
		if m := objectRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			id := fmt.Sprintf("kotlin://%s.%s", pkg, name)
			symbols = append(symbols, Symbol{
				ID:         id,
				Language:   "kotlin",
				Kind:       "singleton",
				Visibility: "public",
				File:       relFile,
				Line:       lineNum,
			})
			pendingJvmStatic = false
			pendingJvmName = ""
			continue
		}

		// Sealed class
		if m := sealedClassRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			id := fmt.Sprintf("kotlin://%s.%s", pkg, name)
			symbols = append(symbols, Symbol{
				ID:         id,
				Language:   "kotlin",
				Kind:       "sealed_class",
				Visibility: "public",
				File:       relFile,
				Line:       lineNum,
			})
			pendingJvmStatic = false
			pendingJvmName = ""
			continue
		}

		// Companion object
		if companionRe.MatchString(line) {
			id := fmt.Sprintf("kotlin://%s#companion", pkg)
			symbols = append(symbols, Symbol{
				ID:         id,
				Language:   "kotlin",
				Kind:       "companion",
				Visibility: "public",
				File:       relFile,
				Line:       lineNum,
			})
			pendingJvmStatic = false
			pendingJvmName = ""
			continue
		}

		// Function declaration
		if funRe.MatchString(line) {
			nameMatch := funNameRe.FindStringSubmatch(line)
			if nameMatch == nil {
				continue
			}
			name := nameMatch[1]

			vis := "public"
			if strings.Contains(line, "private ") {
				vis = "private"
			} else if strings.Contains(line, "protected ") {
				vis = "protected"
			} else if strings.Contains(line, "internal ") {
				vis = "internal"
			}

			isExternal := externalFunRe.MatchString(line)
			isExt := extFunRe.MatchString(line)

			kind := "function"
			if isExt {
				kind = "extension_function"
			}

			var attrs []string
			ffiName := ""
			if pendingJvmStatic {
				attrs = append(attrs, "jvm_static")
			}
			if pendingJvmName != "" {
				ffiName = pendingJvmName
				attrs = append(attrs, "jvm_name")
			}

			id := fmt.Sprintf("kotlin://%s#%s", pkg, name)
			sym := Symbol{
				ID:          id,
				Language:    "kotlin",
				Kind:        kind,
				Visibility:  vis,
				File:        relFile,
				Line:        lineNum,
				FFIExported: isExternal,
				FFIName:     ffiName,
				Attributes:  attrs,
			}
			if isExternal {
				entryPoints = append(entryPoints, EntryPoint{SymbolID: id, Reason: "ffi_export"})
			}
			symbols = append(symbols, sym)

			pendingJvmStatic = false
			pendingJvmName = ""
			continue
		}

		// Reset pending annotations
		if trimmed != "" && !strings.HasPrefix(trimmed, "@") && !strings.HasPrefix(trimmed, "//") {
			if !strings.Contains(trimmed, "fun ") && !strings.Contains(trimmed, "class ") && !strings.Contains(trimmed, "object ") {
				pendingJvmStatic = false
				pendingJvmName = ""
			}
		}
	}

	return symbols, nil, entryPoints
}

