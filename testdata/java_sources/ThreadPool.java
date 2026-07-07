package com.example;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

/**
 * Managed thread pool.
 */
public class ThreadPool {

    private final ExecutorService executor;
    private final int size;

    public ThreadPool(int size) {
        this.size = size;
        this.executor = Executors.newFixedThreadPool(size);
    }

    public void submit(Runnable task) {
        executor.submit(task);
    }

    public void shutdown() throws InterruptedException {
        executor.shutdown();
        executor.awaitTermination(30, TimeUnit.SECONDS);
    }

    public int getSize() {
        return size;
    }

    /**
     * Force shutdown - dead code.
     */
    public void shutdownNow() {
        executor.shutdownNow();
    }

    /**
     * Check if shutdown - dead code.
     */
    public boolean isShutdown() {
        return executor.isShutdown();
    }

    /**
     * Get active count - dead code.
     */
    private int getActiveCount() {
        // Not directly available from ExecutorService
        return 0;
    }
}
