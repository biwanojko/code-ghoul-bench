// Package analysis provides shared types for the code-ghoul dead code eliminator.
package analysis

// Symbol represents a code symbol extracted from source files.
// JSON field names are lowercase — use json struct tags on any embedding struct.
type Symbol struct {
	ID          string   `json:"id"`
	Language    string   `json:"language"`
	Kind        string   `json:"kind"`
	Visibility  string   `json:"visibility"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	FFIExported bool     `json:"ffi_exported"`
	FFIName     string   `json:"ffi_name,omitempty"`
	Scope       string   `json:"scope"`
	Conditions  []string `json:"conditions"`
	Attributes  []string `json:"attributes"`
}

// Callsite records a cross-language FFI call.
type Callsite struct {
	CallerID   string   `json:"caller_id"`
	TargetLang string   `json:"target_language"`
	TargetFFI  string   `json:"target_ffi_name"`
	Mechanism  string   `json:"mechanism"`
	File       string   `json:"file"`
	Line       int      `json:"line"`
	Conditions []string `json:"conditions"`
}

// EntryPoint is a symbol always considered reachable (main, init, JNI export, etc.).
type EntryPoint struct {
	SymbolID string `json:"id"`
	Reason   string `json:"reason"`
}

// ExtractOutput is the JSON output of cmd/extract.
type ExtractOutput struct {
	Symbols     []Symbol     `json:"symbols"`
	Callsites   []Callsite   `json:"callsites"`
	EntryPoints []EntryPoint `json:"entry_points"`
}
