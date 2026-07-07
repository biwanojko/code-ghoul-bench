package com.example.kotlin

/**
 * Flow-based reactive utilities.
 */

fun generateSequence(start: Int, count: Int): List<Int> {
    return (start until start + count).toList()
}

fun filterEven(numbers: List<Int>): List<Int> {
    return numbers.filter { it % 2 == 0 }
}

fun sumList(numbers: List<Int>): Int {
    return numbers.sum()
}

/**
 * Unused transform - dead code.
 */
fun doubleAll(numbers: List<Int>): List<Int> {
    return numbers.map { it * 2 }
}

/**
 * Unused chunking utility - dead code.
 */
fun chunked(numbers: List<Int>, size: Int): List<List<Int>> {
    return numbers.chunked(size)
}

/**
 * Unused reduce - dead code.
 */
fun product(numbers: List<Int>): Long {
    return numbers.fold(1L) { acc, n -> acc * n }
}
