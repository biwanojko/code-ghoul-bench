package com.example.kotlin

/**
 * Kotlin application entry point.
 */
fun main() {
    val app = Application()
    app.run()
}

class Application {
    fun run() {
        val result = computeWrapper(21)
        println("Result: $result")
    }

    /**
     * Unused shutdown method - dead code.
     */
    fun shutdown() {
        println("shutting down")
    }

    /**
     * Unused configuration loader - dead code.
     */
    private fun loadConfig(): Map<String, String> {
        return mapOf("host" to "localhost", "port" to "8080")
    }
}
