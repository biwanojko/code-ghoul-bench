package com.example;

import java.time.Instant;

/**
 * Simple application logger.
 */
public class Logger {

    private final String prefix;
    private Level level;

    public enum Level {
        DEBUG, INFO, WARN, ERROR
    }

    public Logger(String prefix, Level level) {
        this.prefix = prefix;
        this.level = level;
    }

    public void log(Level level, String message) {
        if (level.ordinal() >= this.level.ordinal()) {
            System.out.printf("[%s] %s %s: %s%n", prefix, Instant.now(), level, message);
        }
    }

    public void info(String message) {
        log(Level.INFO, message);
    }

    public void error(String message) {
        log(Level.ERROR, message);
    }

    /**
     * Set log level - dead code.
     */
    public void setLevel(Level level) {
        this.level = level;
    }

    /**
     * Format as JSON - dead code.
     */
    private String formatJson(Level level, String message) {
        return String.format("{\"level\":\"%s\",\"msg\":\"%s\"}", level, message);
    }

    /**
     * Create a child logger - dead code.
     */
    public Logger child(String childPrefix) {
        return new Logger(prefix + "." + childPrefix, level);
    }
}
