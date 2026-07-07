package com.example.kotlin

import kotlinx.coroutines.delay

/**
 * Coroutine-based async utilities.
 */

suspend fun fetchData(url: String): String {
    delay(100)
    return "data from $url"
}

suspend fun processAsync(data: String): String {
    delay(50)
    return data.uppercase()
}

/**
 * Unused coroutine function - dead code.
 */
suspend fun unusedAsyncOp(x: Int): Int {
    delay(10)
    return x * 2
}

/**
 * Unused retry logic - dead code.
 */
suspend fun withRetry(maxRetries: Int, block: suspend () -> String): String {
    var lastError: Exception? = null
    repeat(maxRetries) {
        try {
            return block()
        } catch (e: Exception) {
            lastError = e
            delay(100)
        }
    }
    throw lastError ?: RuntimeException("all retries failed")
}
