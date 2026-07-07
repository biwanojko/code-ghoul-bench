package com.example;

import java.util.HashMap;
import java.util.Map;

/**
 * Application configuration.
 */
public class Config {

    private Map<String, String> properties = new HashMap<>();

    public Config() {
        setDefaults();
    }

    private void setDefaults() {
        properties.put("host", "localhost");
        properties.put("port", "8080");
    }

    public String get(String key) {
        return properties.getOrDefault(key, "");
    }

    public void set(String key, String value) {
        properties.put(key, value);
    }

    /**
     * Load config from file - dead code.
     */
    public static Config fromFile(String path) {
        Config cfg = new Config();
        // TODO: implement
        return cfg;
    }

    /**
     * Save config to file - dead code.
     */
    public void save(String path) {
        // TODO: implement
    }

    /**
     * Validate configuration - dead code.
     */
    public boolean validate() {
        return !get("host").isEmpty();
    }
}
