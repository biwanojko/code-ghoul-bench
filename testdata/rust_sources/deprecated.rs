/// Deprecated functions - all dead code

/// Old computation function - replaced by engine::compute
pub fn old_compute(x: i32) -> i32 {
    x + x
}

/// Old string processor - replaced by string_utils
pub fn old_process_string(s: &str) -> String {
    s.to_uppercase()
}

/// Old FFI export - replaced by engine::rust_compute
#[no_mangle]
pub extern "C" fn old_rust_compute(x: i32) -> i32 {
    old_compute(x)
}

/// Legacy initialization - dead code
fn legacy_init() -> bool {
    true
}

/// Legacy cleanup - dead code
fn legacy_cleanup() {
    // nothing to do
}

/// Old error handler - dead code
pub fn old_handle_error(code: i32) -> String {
    format!("old error: {}", code)
}
