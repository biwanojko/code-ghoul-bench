/// Feature-gated functionality

#[cfg(feature = "ffi")]
pub fn ffi_enabled_function() -> i32 {
    42
}

#[cfg(feature = "ffi")]
#[no_mangle]
pub extern "C" fn conditionally_exported(x: i32) -> i32 {
    x + ffi_enabled_function()
}

#[cfg(feature = "jni")]
pub fn jni_enabled_helper() -> String {
    String::from("jni-enabled")
}

/// Always available function
pub fn always_available() -> bool {
    true
}

/// Dead code regardless of features
fn unused_feature_helper() -> i32 {
    0
}

#[cfg(feature = "experimental")]
pub fn experimental_feature() -> f64 {
    3.14
}
