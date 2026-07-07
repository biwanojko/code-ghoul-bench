package com.example;

/**
 * Builder pattern example.
 */
public class Builder {

    private String host = "localhost";
    private int port = 8080;
    private int timeout = 30;
    private boolean debug = false;

    public Builder withHost(String host) {
        this.host = host;
        return this;
    }

    public Builder withPort(int port) {
        this.port = port;
        return this;
    }

    public Builder withTimeout(int timeout) {
        this.timeout = timeout;
        return this;
    }

    /**
     * Set debug mode - dead code (never called on builder).
     */
    public Builder withDebug(boolean debug) {
        this.debug = debug;
        return this;
    }

    public Config build() {
        Config cfg = new Config();
        cfg.set("host", host);
        cfg.set("port", String.valueOf(port));
        cfg.set("timeout", String.valueOf(timeout));
        return cfg;
    }

    /**
     * Validate builder state - dead code.
     */
    private boolean isValid() {
        return !host.isEmpty() && port > 0 && timeout > 0;
    }

    /**
     * Reset builder - dead code.
     */
    public void reset() {
        host = "localhost";
        port = 8080;
        timeout = 30;
        debug = false;
    }
}
