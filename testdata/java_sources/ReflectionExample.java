package com.example;

/**
 * Example of reflection-based class loading.
 */
public class ReflectionExample {

    /**
     * Load and use a class by name.
     */
    public static Object loadByName(String className) throws Exception {
        Class<?> cls = Class.forName("com.example.Application");
        return cls.getDeclaredConstructor().newInstance();
    }

    /**
     * Create an instance of a known class.
     */
    public static NativeBridge createBridge() {
        return new NativeBridge();
    }

    /**
     * Unused reflection helper - dead code.
     */
    private static void unusedReflectionHelper() throws Exception {
        Class.forName("com.example.Config");
    }

    /**
     * Another unused helper - dead code.
     */
    private static String getClassName(Object obj) {
        return obj.getClass().getName();
    }
}
