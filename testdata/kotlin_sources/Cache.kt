package com.example.kotlin

import java.util.concurrent.ConcurrentHashMap

/**
 * Thread-safe cache implementation.
 */
class KotlinCache<K, V>(private val maxSize: Int = 1000) {

    private val map = ConcurrentHashMap<K, V>()

    fun put(key: K, value: V): V? {
        if (map.size >= maxSize) {
            // Evict first entry (simplified LRU)
            map.remove(map.keys.first())
        }
        return map.put(key, value)
    }

    fun get(key: K): V? = map[key]

    fun remove(key: K): V? = map.remove(key)

    fun size(): Int = map.size

    /**
     * Clear all entries - dead code.
     */
    fun clear() = map.clear()

    /**
     * Get or compute - dead code.
     */
    fun getOrPut(key: K, compute: () -> V): V {
        return map.getOrPut(key, compute)
    }

    /**
     * Contains check - dead code.
     */
    fun contains(key: K): Boolean = map.containsKey(key)
}

/**
 * Unused global cache instance - dead code.
 */
val globalCache = KotlinCache<String, Any>()
