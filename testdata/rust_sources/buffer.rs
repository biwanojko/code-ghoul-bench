/// Buffer management utilities

/// RingBuffer is a fixed-size circular buffer
pub struct RingBuffer {
    data: Vec<u8>,
    head: usize,
    tail: usize,
    size: usize,
}

impl RingBuffer {
    /// Create a new ring buffer with given capacity
    pub fn new(capacity: usize) -> Self {
        RingBuffer {
            data: vec![0u8; capacity],
            head: 0,
            tail: 0,
            size: 0,
        }
    }

    /// Push a byte into the buffer
    pub fn push(&mut self, b: u8) -> bool {
        if self.size == self.data.len() {
            return false;
        }
        self.data[self.tail] = b;
        self.tail = (self.tail + 1) % self.data.len();
        self.size += 1;
        true
    }

    /// Pop a byte from the buffer
    pub fn pop(&mut self) -> Option<u8> {
        if self.size == 0 {
            return None;
        }
        let b = self.data[self.head];
        self.head = (self.head + 1) % self.data.len();
        self.size -= 1;
        Some(b)
    }

    /// Peek at the next byte without removing it
    pub fn peek(&self) -> Option<u8> {
        if self.size == 0 {
            None
        } else {
            Some(self.data[self.head])
        }
    }
}

/// FFI: create a ring buffer - dead code (exported but not called)
#[no_mangle]
pub extern "C" fn ring_buffer_new(capacity: usize) -> *mut RingBuffer {
    let buf = Box::new(RingBuffer::new(capacity));
    Box::into_raw(buf)
}

/// FFI: free a ring buffer - dead code
#[no_mangle]
pub extern "C" fn ring_buffer_free(ptr: *mut RingBuffer) {
    if !ptr.is_null() {
        unsafe { drop(Box::from_raw(ptr)); }
    }
}

/// Utility: drain all bytes - dead code
fn drain_buffer(buf: &mut RingBuffer) -> Vec<u8> {
    let mut result = Vec::new();
    while let Some(b) = buf.pop() {
        result.push(b);
    }
    result
}
