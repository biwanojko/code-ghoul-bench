package com.example.kotlin

/**
 * Singleton registry using Kotlin object declaration.
 */
object Registry {

    private val entries = mutableMapOf<String, Any>()

    fun register(key: String, value: Any) {
        entries[key] = value
    }

    fun get(key: String): Any? = entries[key]

    fun size(): Int = entries.size

    /**
     * Clear all entries - dead code.
     */
    fun clear() {
        entries.clear()
    }

    /**
     * Check if key exists - dead code.
     */
    fun contains(key: String): Boolean = entries.containsKey(key)
}

/**
 * Unused singleton - dead code.
 */
object DeadSingleton {
    fun doNothing() {}
    fun alsoDoNothing() = Unit
}
