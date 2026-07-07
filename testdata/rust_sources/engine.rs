/// Core computation engine

/// Compute the result of an operation
pub fn compute(x: i32) -> i32 {
    x * 2
}

/// Internal helper - not exported
fn internal_helper(x: i32) -> i32 {
    x + 1
}

/// FFI-exported compute function
#[no_mangle]
pub extern "C" fn rust_compute(x: i32) -> i32 {
    compute(x)
}

/// Unused internal function - dead code
fn dead_internal(x: i32) -> i32 {
    internal_helper(x) * 3
}

/// Unused public function - dead code
pub fn unused_public() -> String {
    String::from("unused")
}
