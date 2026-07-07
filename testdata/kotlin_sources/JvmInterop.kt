package com.example.kotlin

/**
 * Kotlin-Java interop utilities.
 */

object JvmHelper {

    @JvmStatic
    fun staticHelper(x: Int): Int = x * 2

    @JvmName("computeValue")
    fun compute(x: Int, y: Int): Int = x + y

    /**
     * Non-annotated method - regular Kotlin.
     */
    fun internalHelper(x: Int): Int = x + 1
}

/**
 * Top-level @JvmStatic equivalent.
 */
@JvmName("topLevelCompute")
fun topLevelFun(x: Int): Int = x * 3

/**
 * Unused JVM interop function - dead code.
 */
@JvmName("unusedJvmFun")
fun unusedJvmFunction(data: List<Int>): Int = data.sum()
