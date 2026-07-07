package com.example;

/**
 * Class with static initializer block.
 */
public class StaticInit {

    private static final String CONFIG;
    private static final int MAX_RETRIES;

    static {
        CONFIG = System.getenv("APP_CONFIG");
        MAX_RETRIES = 3;
        initializeSubsystems();
    }

    private static void initializeSubsystems() {
        // Initialize various subsystems
        System.out.println("Subsystems initialized with config: " + CONFIG);
    }

    public static int getMaxRetries() {
        return MAX_RETRIES;
    }

    /**
     * Unused cleanup - dead code.
     */
    public static void cleanup() {
        System.out.println("cleanup");
    }

    /**
     * Unused internal helper - dead code.
     */
    private static boolean checkConfig(String cfg) {
        return cfg != null && !cfg.isEmpty();
    }
}
