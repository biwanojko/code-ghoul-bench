/// Test utilities and test-only functions

#[cfg(test)]
mod tests {
    use super::*;

    fn test_helper_add(a: i32, b: i32) -> i32 {
        a + b
    }

    fn test_setup() -> Vec<u8> {
        vec![1, 2, 3, 4, 5]
    }
}

/// Export for benchmarking - used in tests only
#[cfg(test)]
pub fn bench_compute(iterations: u64) -> u64 {
    let mut sum = 0u64;
    for i in 0..iterations {
        sum = sum.wrapping_add(i);
    }
    sum
}

/// Test fixture builder - dead code in non-test builds
#[cfg(test)]
fn build_fixture(size: usize) -> Vec<i32> {
    (0..size as i32).collect()
}
