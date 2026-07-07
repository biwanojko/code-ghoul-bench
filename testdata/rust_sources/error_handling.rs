/// Error handling utilities

use std::fmt;

/// AppError represents an application error
#[derive(Debug)]
pub struct AppError {
    pub code: i32,
    pub message: String,
}

impl fmt::Display for AppError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "error {}: {}", self.code, self.message)
    }
}

/// Create a new error
pub fn new_error(code: i32, msg: &str) -> AppError {
    AppError {
        code,
        message: msg.to_string(),
    }
}

/// Check if result is ok - dead code
pub fn is_ok(code: i32) -> bool {
    code == 0
}

/// Format error for logging - dead code
fn format_error(err: &AppError) -> String {
    format!("[ERR-{}] {}", err.code, err.message)
}

/// Wrap an error with context - dead code
pub fn wrap_error(err: AppError, context: &str) -> AppError {
    AppError {
        code: err.code,
        message: format!("{}: {}", context, err.message),
    }
}
