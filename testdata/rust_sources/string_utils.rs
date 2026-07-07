/// String utilities

/// Reverse a string
pub fn reverse(s: &str) -> String {
    s.chars().rev().collect()
}

/// Count occurrences of a char in a string
pub fn count_char(s: &str, c: char) -> usize {
    s.chars().filter(|&ch| ch == c).count()
}

/// Trim and lowercase a string
pub fn normalize(s: &str) -> String {
    s.trim().to_lowercase()
}

/// Check if string is a palindrome - dead code
pub fn is_palindrome(s: &str) -> bool {
    let normalized: String = s.chars().filter(|c| c.is_alphanumeric()).collect();
    let normalized = normalized.to_lowercase();
    normalized == reverse(&normalized)
}

/// Join strings with separator - dead code
pub fn join(parts: &[&str], sep: &str) -> String {
    parts.join(sep)
}

/// Split and trim - dead code
fn split_trim(s: &str, sep: char) -> Vec<String> {
    s.split(sep).map(|p| p.trim().to_string()).collect()
}

/// Capitalize first letter - dead code
pub fn capitalize(s: &str) -> String {
    let mut c = s.chars();
    match c.next() {
        None => String::new(),
        Some(f) => f.to_uppercase().collect::<String>() + c.as_str(),
    }
}
