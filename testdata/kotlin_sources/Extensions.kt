package com.example.kotlin

/**
 * Extension functions for common types.
 */

fun String.truncate(maxLength: Int): String {
    return if (length <= maxLength) this else substring(0, maxLength) + "..."
}

fun String.isPalindrome(): Boolean {
    val cleaned = filter { it.isLetterOrDigit() }.lowercase()
    return cleaned == cleaned.reversed()
}

/**
 * Unused extension - dead code.
 */
fun String.wordCount(): Int {
    return trim().split(Regex("\\s+")).size
}

/**
 * Unused extension - dead code.
 */
fun List<Int>.average(): Double {
    if (isEmpty()) return 0.0
    return sum().toDouble() / size
}

/**
 * Unused extension - dead code.
 */
fun Int.toHex(): String = toString(16)
