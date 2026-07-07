package com.example.kotlin

/**
 * Native function declarations for Kotlin Multiplatform.
 */
external fun rustCompute(x: Int): Int

external fun rustChecksum(data: ByteArray): Long

/**
 * Wrapper around the native compute function.
 */
fun computeWrapper(x: Int): Int {
    return rustCompute(x)
}

/**
 * Unused native wrapper - dead code.
 */
fun unusedNativeWrapper(data: ByteArray): Long {
    return rustChecksum(data)
}

/**
 * Internal utility - dead code.
 */
private fun formatResult(result: Int): String {
    return "Result: $result"
}
