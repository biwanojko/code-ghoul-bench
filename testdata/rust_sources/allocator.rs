/// Custom global allocator

use std::alloc::{GlobalAlloc, Layout, System};

/// CustomAllocator wraps the system allocator
pub struct CustomAllocator;

#[global_allocator]
static ALLOCATOR: CustomAllocator = CustomAllocator;

unsafe impl GlobalAlloc for CustomAllocator {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        System.alloc(layout)
    }

    unsafe fn dealloc(&self, ptr: *mut u8, layout: Layout) {
        System.dealloc(ptr, layout)
    }
}

/// Allocation stats tracking - dead code
fn track_alloc(size: usize) {
    let _ = size;
}

/// Custom realloc - dead code
fn custom_realloc(ptr: *mut u8, old_layout: Layout, new_size: usize) -> *mut u8 {
    unsafe { System.realloc(ptr, old_layout, new_size) }
}
