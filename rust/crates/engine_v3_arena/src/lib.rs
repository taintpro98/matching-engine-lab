//! Engine v3: Arena allocator-based implementation (placeholder).

use std::collections::{BTreeMap, HashMap};

use engine_core::{Command, Engine, Event};

pub struct EngineV3 {
    _book: BTreeMap<(i64, i64), ()>, // placeholder: (price, timestamp) -> ()
}

impl Default for EngineV3 {
    fn default() -> Self {
        Self {
            _book: BTreeMap::new(),
        }
    }
}

impl Engine for EngineV3 {
    fn submit(&mut self, cmd: Command) -> Vec<Event> {
        let _ = cmd;
        vec![Event::Accepted]
    }

    fn reset(&mut self) {
        self._book.clear();
    }

    fn snapshot(&self) -> Vec<u8> {
        vec![]
    }

    fn load_snapshot(&mut self, _data: &[u8]) {}

    fn stats(&self) -> HashMap<String, String> {
        let mut m = HashMap::new();
        m.insert("engine".to_string(), "v3_arena".to_string());
        m
    }
}
