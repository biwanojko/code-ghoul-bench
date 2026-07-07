package com.example;

import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Simple LRU cache implementation.
 */
public class Cache<K, V> {

    private final int capacity;
    private final Map<K, V> map;

    public Cache(int capacity) {
        this.capacity = capacity;
        this.map = new LinkedHashMap<>(capacity, 0.75f, true) {
            protected boolean removeEldestEntry(Map.Entry<K, V> eldest) {
                return size() > Cache.this.capacity;
            }
        };
    }

    public synchronized V get(K key) {
        return map.get(key);
    }

    public synchronized void put(K key, V value) {
        map.put(key, value);
    }

    public synchronized int size() {
        return map.size();
    }

    /**
     * Remove a key from the cache - dead code.
     */
    public synchronized void remove(K key) {
        map.remove(key);
    }

    /**
     * Clear all entries - dead code.
     */
    public synchronized void clear() {
        map.clear();
    }

    /**
     * Check if key exists - dead code.
     */
    public synchronized boolean containsKey(K key) {
        return map.containsKey(key);
    }
}
