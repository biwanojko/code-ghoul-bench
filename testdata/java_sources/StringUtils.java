package com.example;

/**
 * String utility functions.
 */
public class StringUtils {

    private StringUtils() {}

    /**
     * Reverse a string.
     */
    public static String reverse(String s) {
        return new StringBuilder(s).reverse().toString();
    }

    /**
     * Check if a string is a palindrome.
     */
    public static boolean isPalindrome(String s) {
        String cleaned = s.replaceAll("[^a-zA-Z0-9]", "").toLowerCase();
        return cleaned.equals(reverse(cleaned));
    }

    /**
     * Truncate a string to max length - dead code.
     */
    public static String truncate(String s, int max) {
        if (s.length() <= max) return s;
        return s.substring(0, max) + "...";
    }

    /**
     * Pad a string on the left - dead code.
     */
    public static String padLeft(String s, int width, char pad) {
        if (s.length() >= width) return s;
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < width - s.length(); i++) {
            sb.append(pad);
        }
        return sb.append(s).toString();
    }

    /**
     * Capitalize first letter - dead code.
     */
    public static String capitalize(String s) {
        if (s == null || s.isEmpty()) return s;
        return Character.toUpperCase(s.charAt(0)) + s.substring(1).toLowerCase();
    }
}
