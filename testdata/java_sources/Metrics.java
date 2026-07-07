package com.example;

import java.util.concurrent.atomic.AtomicLong;

/**
 * Application metrics collector.
 */
public class Metrics {

    private final AtomicLong requestCount = new AtomicLong(0);
    private final AtomicLong errorCount = new AtomicLong(0);
    private final AtomicLong totalLatencyMs = new AtomicLong(0);

    public void recordRequest(long latencyMs) {
        requestCount.incrementAndGet();
        totalLatencyMs.addAndGet(latencyMs);
    }

    public void recordError() {
        errorCount.incrementAndGet();
    }

    public long getRequestCount() {
        return requestCount.get();
    }

    public long getErrorCount() {
        return errorCount.get();
    }

    /**
     * Get average latency - dead code.
     */
    public double getAverageLatency() {
        long count = requestCount.get();
        if (count == 0) return 0.0;
        return (double) totalLatencyMs.get() / count;
    }

    /**
     * Reset all metrics - dead code.
     */
    public void reset() {
        requestCount.set(0);
        errorCount.set(0);
        totalLatencyMs.set(0);
    }

    /**
     * Export metrics as string - dead code.
     */
    private String export() {
        return String.format("requests=%d errors=%d", requestCount.get(), errorCount.get());
    }
}
