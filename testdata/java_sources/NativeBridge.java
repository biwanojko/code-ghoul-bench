package com.example;

/**
 * NativeBridge provides JNI calls into the Rust native library.
 */
public class NativeBridge {

    static {
        System.loadLibrary("engine");
    }

    /**
     * Compute a value using the native Rust implementation.
     */
    public native int compute(int x);

    /**
     * Process a string natively.
     */
    public native int processString(int len);

    /**
     * Unused native method - dead code.
     */
    public native int unused();

    /**
     * Java-only helper - never called.
     */
    private int javaHelper(int x) {
        return x * 3;
    }

    /**
     * Initialize the bridge - called from static block.
     */
    private static void initialize() {
        // initialization logic
    }
}
