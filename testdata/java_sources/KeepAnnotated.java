package com.example;

/**
 * Classes and methods annotated with @Keep should be retained even if unreachable.
 */
public class KeepAnnotated {

    // Simulating @Keep annotation effect
    @interface Keep {}

    @Keep
    public static void keepThisMethod() {
        System.out.println("kept by annotation");
    }

    /**
     * Not annotated - dead code.
     */
    public static void notKept() {
        System.out.println("not kept");
    }

    /**
     * Not annotated - dead code.
     */
    private static void alsoNotKept() {
        System.out.println("also not kept");
    }
}
