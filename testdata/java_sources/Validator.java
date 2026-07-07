package com.example;

import java.util.regex.Pattern;

/**
 * Input validation utilities.
 */
public class Validator {

    private static final Pattern EMAIL_PATTERN =
        Pattern.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$");

    private static final Pattern URL_PATTERN =
        Pattern.compile("^https?://[\\w.-]+(?:/[\\w./-]*)?$");

    /**
     * Validate an email address.
     */
    public static boolean isValidEmail(String email) {
        if (email == null) return false;
        return EMAIL_PATTERN.matcher(email).matches();
    }

    /**
     * Validate a URL.
     */
    public static boolean isValidUrl(String url) {
        if (url == null) return false;
        return URL_PATTERN.matcher(url).matches();
    }

    /**
     * Validate a port number - dead code.
     */
    public static boolean isValidPort(int port) {
        return port > 0 && port <= 65535;
    }

    /**
     * Validate a hostname - dead code.
     */
    public static boolean isValidHostname(String hostname) {
        if (hostname == null || hostname.isEmpty()) return false;
        return hostname.matches("^[a-zA-Z0-9.-]+$");
    }

    /**
     * Sanitize a string - dead code.
     */
    public static String sanitize(String input) {
        if (input == null) return "";
        return input.replaceAll("[<>\"'&]", "");
    }
}
