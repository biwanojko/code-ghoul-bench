/// Cryptographic utilities

/// Simple XOR cipher for demonstration
pub fn xor_cipher(data: &[u8], key: u8) -> Vec<u8> {
    data.iter().map(|b| b ^ key).collect()
}

/// Compute a simple checksum
pub fn checksum(data: &[u8]) -> u32 {
    data.iter().fold(0u32, |acc, &b| acc.wrapping_add(b as u32))
}

/// Rotate bytes left by n positions
pub fn rotate_left(data: &[u8], n: usize) -> Vec<u8> {
    if data.is_empty() {
        return Vec::new();
    }
    let n = n % data.len();
    let mut result = data[n..].to_vec();
    result.extend_from_slice(&data[..n]);
    result
}

/// FFI export: compute checksum from C
#[no_mangle]
pub extern "C" fn rust_checksum(ptr: *const u8, len: usize) -> u32 {
    if ptr.is_null() {
        return 0;
    }
    let data = unsafe { std::slice::from_raw_parts(ptr, len) };
    checksum(data)
}

/// Unused: secure erase - dead code
pub fn secure_erase(data: &mut [u8]) {
    for b in data.iter_mut() {
        *b = 0;
    }
}

/// Unused: key derivation - dead code
fn derive_key(password: &[u8], salt: &[u8]) -> Vec<u8> {
    let mut key = password.to_vec();
    key.extend_from_slice(salt);
    xor_cipher(&key, 0x5A)
}
