/// FFI exports for cross-language calls

use std::ffi::CStr;
use std::os::raw::c_char;

/// Process a string from C
#[no_mangle]
pub extern "C" fn process_string(s: *const c_char) -> i32 {
    if s.is_null() {
        return -1;
    }
    let c_str = unsafe { CStr::from_ptr(s) };
    c_str.to_bytes().len() as i32
}

/// Free a string allocated by Rust
#[no_mangle]
pub extern "C" fn free_rust_string(s: *mut c_char) {
    if s.is_null() {
        return;
    }
    unsafe {
        let _ = std::ffi::CString::from_raw(s);
    }
}

/// Allocate a Rust string for C - dead code (exported but not called from Java/Go)
#[no_mangle]
pub extern "C" fn alloc_rust_string(len: usize) -> *mut c_char {
    let mut v: Vec<u8> = Vec::with_capacity(len + 1);
    v.push(0);
    let ptr = v.as_mut_ptr() as *mut c_char;
    std::mem::forget(v);
    ptr
}

/// Internal FFI helper - never called
fn ffi_internal_validate(ptr: *const c_char) -> bool {
    !ptr.is_null()
}
