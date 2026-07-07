package com.example;

/**
 * Deprecated classes - entire file is dead code.
 */
public class Deprecated {

    /**
     * Old processing method - dead code.
     */
    public static String process(String input) {
        return input.toUpperCase();
    }

    /**
     * Old validation method - dead code.
     */
    public static boolean validate(String input) {
        return input != null && !input.isEmpty();
    }

    /**
     * Old factory method - dead code.
     */
    public static Deprecated create() {
        return new Deprecated();
    }

    /**
     * Old cleanup - dead code.
     */
    public void destroy() {
        // nothing
    }
}
