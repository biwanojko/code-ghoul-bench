/// Metrics collection

use std::sync::atomic::{AtomicU64, Ordering};

static COUNTER: AtomicU64 = AtomicU64::new(0);
static ERROR_COUNT: AtomicU64 = AtomicU64::new(0);

/// Increment the request counter
pub fn increment() {
    COUNTER.fetch_add(1, Ordering::Relaxed);
}

/// Get the current counter value
pub fn get_count() -> u64 {
    COUNTER.load(Ordering::Relaxed)
}

/// Increment error counter
pub fn increment_errors() {
    ERROR_COUNT.fetch_add(1, Ordering::Relaxed);
}

/// FFI: get counter value
#[no_mangle]
pub extern "C" fn metrics_get_count() -> u64 {
    get_count()
}

/// FFI: reset counter - dead code
#[no_mangle]
pub extern "C" fn metrics_reset() {
    COUNTER.store(0, Ordering::Relaxed);
    ERROR_COUNT.store(0, Ordering::Relaxed);
}

/// Internal: compute rate - dead code
fn compute_rate(count: u64, duration_ms: u64) -> f64 {
    if duration_ms == 0 {
        return 0.0;
    }
    (count as f64) / (duration_ms as f64) * 1000.0
}
