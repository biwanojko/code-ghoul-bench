package com.example;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Simple event bus for decoupled communication.
 */
public class EventBus {

    private final Map<String, List<Runnable>> listeners = new HashMap<>();

    /**
     * Subscribe to an event type.
     */
    public void subscribe(String event, Runnable handler) {
        listeners.computeIfAbsent(event, k -> new ArrayList<>()).add(handler);
    }

    /**
     * Publish an event.
     */
    public void publish(String event) {
        List<Runnable> handlers = listeners.getOrDefault(event, List.of());
        for (Runnable h : handlers) {
            h.run();
        }
    }

    /**
     * Unsubscribe from an event - dead code.
     */
    public void unsubscribe(String event, Runnable handler) {
        List<Runnable> handlers = listeners.get(event);
        if (handlers != null) {
            handlers.remove(handler);
        }
    }

    /**
     * Get subscriber count - dead code.
     */
    public int subscriberCount(String event) {
        return listeners.getOrDefault(event, List.of()).size();
    }

    /**
     * Clear all listeners - dead code.
     */
    public void clearAll() {
        listeners.clear();
    }
}
