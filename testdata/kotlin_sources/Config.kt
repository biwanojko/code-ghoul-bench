package com.example.kotlin

/**
 * Configuration management.
 */

data class AppConfig(
    val host: String = "localhost",
    val port: Int = 8080,
    val maxConnections: Int = 100,
    val debug: Boolean = false
)

object ConfigManager {

    private var config = AppConfig()

    fun getConfig(): AppConfig = config

    fun updateConfig(newConfig: AppConfig) {
        config = newConfig
    }

    /**
     * Load from environment - dead code.
     */
    fun loadFromEnv(): AppConfig {
        return AppConfig(
            host = System.getenv("HOST") ?: "localhost",
            port = System.getenv("PORT")?.toIntOrNull() ?: 8080
        )
    }

    /**
     * Validate config - dead code.
     */
    fun validate(cfg: AppConfig): Boolean {
        return cfg.port in 1..65535 && cfg.maxConnections > 0
    }
}

/**
 * Unused config builder - dead code.
 */
fun buildConfig(host: String, port: Int): AppConfig {
    return AppConfig(host = host, port = port)
}
