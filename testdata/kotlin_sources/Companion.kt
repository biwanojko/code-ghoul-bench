package com.example.kotlin

/**
 * Companion object examples.
 */

class ConnectionPool(private val maxSize: Int) {

    companion object {
        const val DEFAULT_MAX_SIZE = 10
        const val DEFAULT_TIMEOUT = 5000L

        fun create(): ConnectionPool = ConnectionPool(DEFAULT_MAX_SIZE)

        @JvmStatic
        fun createWithSize(size: Int): ConnectionPool = ConnectionPool(size)
    }

    private val connections = mutableListOf<String>()

    fun acquire(): String? {
        return if (connections.size < maxSize) {
            val conn = "conn-${connections.size}"
            connections.add(conn)
            conn
        } else null
    }

    fun release(conn: String) {
        connections.remove(conn)
    }

    /**
     * Get pool size - dead code.
     */
    fun size(): Int = connections.size

    /**
     * Drain all connections - dead code.
     */
    fun drain() {
        connections.clear()
    }
}
