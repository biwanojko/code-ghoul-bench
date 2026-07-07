/// Data parser module

/// Parse an integer from bytes
pub fn parse_int(bytes: &[u8]) -> Option<i64> {
    std::str::from_utf8(bytes).ok()?.parse().ok()
}

/// Parse a float from bytes
pub fn parse_float(bytes: &[u8]) -> Option<f64> {
    std::str::from_utf8(bytes).ok()?.parse().ok()
}

/// Parse a boolean from string
pub fn parse_bool(s: &str) -> bool {
    matches!(s.to_lowercase().as_str(), "true" | "yes" | "1")
}

/// Parse a JSON-like key=value string - dead code
fn parse_kv(s: &str) -> Option<(&str, &str)> {
    let mut parts = s.splitn(2, '=');
    Some((parts.next()?, parts.next()?))
}

/// Validate a numeric range - dead code
pub fn validate_range(val: i64, min: i64, max: i64) -> bool {
    val >= min && val <= max
}

/// Format a parsed value - dead code
pub fn format_parsed(val: i64, prefix: &str) -> String {
    format!("{}{}", prefix, val)
}
