package com.example.kotlin

/**
 * String utility functions.
 */

fun reverseString(s: String): String = s.reversed()

fun normalizeString(s: String): String = s.trim().lowercase()

fun countOccurrences(s: String, ch: Char): Int = s.count { it == ch }

/**
 * Unused: split and trim - dead code.
 */
fun splitAndTrim(s: String, sep: Char): List<String> {
    return s.split(sep).map { it.trim() }.filter { it.isNotEmpty() }
}

/**
 * Unused: join with separator - dead code.
 */
fun joinWith(parts: List<String>, sep: String): String = parts.joinToString(sep)

/**
 * Unused: repeat string - dead code.
 */
fun repeatStr(s: String, n: Int): String = s.repeat(n)

/**
 * Unused: check prefix/suffix - dead code.
 */
fun hasPrefixAndSuffix(s: String, prefix: String, suffix: String): Boolean {
    return s.startsWith(prefix) && s.endsWith(suffix)
}
