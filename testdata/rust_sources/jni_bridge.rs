/// JNI bridge for Android

/// JNI implementation for com.example.NativeBridge.compute
#[no_mangle]
pub extern "C" fn Java_com_example_NativeBridge_compute(
    _env: *mut std::ffi::c_void,
    _obj: *mut std::ffi::c_void,
    x: i32,
) -> i32 {
    x * 2 + 1
}

/// JNI implementation for NativeBridge.processString
#[no_mangle]
pub extern "C" fn Java_com_example_NativeBridge_processString(
    _env: *mut std::ffi::c_void,
    _obj: *mut std::ffi::c_void,
    len: i32,
) -> i32 {
    len * 3
}

/// Unused JNI function - dead code
#[no_mangle]
pub extern "C" fn Java_com_example_NativeBridge_unused(
    _env: *mut std::ffi::c_void,
    _obj: *mut std::ffi::c_void,
) -> i32 {
    0
}

/// Private helper for JNI functions - dead code
fn jni_helper(x: i32) -> i32 {
    x + 100
}
