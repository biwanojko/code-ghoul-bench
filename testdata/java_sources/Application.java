package com.example;

/**
 * Main application entry point.
 */
public class Application {

    public static void main(String[] args) {
        Application app = new Application();
        app.run();
    }

    public void run() {
        NativeBridge bridge = new NativeBridge();
        int result = bridge.compute(21);
        System.out.println("Result: " + result);
    }

    /**
     * Unused helper - dead code.
     */
    private void printUsage() {
        System.out.println("Usage: application");
    }

    /**
     * Unused shutdown hook - dead code.
     */
    protected void onShutdown() {
        // cleanup
    }
}
