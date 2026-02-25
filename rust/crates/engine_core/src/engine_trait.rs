//! Engine trait for matching engine implementations.

use std::collections::HashMap;

use crate::ndjson::{Command, Event};

pub trait Engine {
    /// Submit a command and return resulting events.
    fn submit(&mut self, cmd: Command) -> Vec<Event>;

    /// Reset engine state.
    fn reset(&mut self);

    /// Serialize current state to bytes.
    fn snapshot(&self) -> Vec<u8>;

    /// Load state from bytes.
    fn load_snapshot(&mut self, data: &[u8]);

    /// Return engine statistics.
    fn stats(&self) -> HashMap<String, String>;
}
