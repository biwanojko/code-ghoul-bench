package com.example.kotlin

/**
 * Deprecated Kotlin functions - all dead code.
 */

fun oldCompute(x: Int): Int = x + x

fun oldProcessString(s: String): String = s.uppercase()

fun oldValidate(input: String): Boolean = input.isNotEmpty()

/**
 * Old singleton - dead code.
 */
object OldManager {
    fun initialize() {}
    fun shutdown() {}
    fun getStatus(): String = "old"
}

/**
 * Old extension - dead code.
 */
fun String.oldTruncate(max: Int): String = if (length > max) substring(0, max) else this

/**
 * Old utility - dead code.
 */
private fun legacyHelper(x: Int, y: Int): Int = x * y + x
