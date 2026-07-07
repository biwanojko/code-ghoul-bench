package com.example.kotlin

/**
 * Type aliases and utility functions.
 */

typealias Handler = (String) -> String
typealias Predicate = (String) -> Boolean

fun applyHandler(input: String, handler: Handler): String = handler(input)

fun filterStrings(items: List<String>, predicate: Predicate): List<String> {
    return items.filter(predicate)
}

/**
 * Unused: compose two handlers - dead code.
 */
fun composeHandlers(first: Handler, second: Handler): Handler {
    return { input -> second(first(input)) }
}

/**
 * Unused: negate predicate - dead code.
 */
fun negatePredicate(p: Predicate): Predicate = { !p(it) }

/**
 * Unused: create length predicate - dead code.
 */
fun minLengthPredicate(minLen: Int): Predicate = { it.length >= minLen }
