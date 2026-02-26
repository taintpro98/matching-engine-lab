//! Runner: reads NDJSON, submits to engine, writes NDJSON output.

use std::io::{self, BufRead, Write};
use std::time::Instant;

use engine_core::{parse_command, serialize_event, Engine};
use engine_v1_btreemap::EngineV1;
use engine_v2_skiplist::EngineV2;
use engine_v3_arena::EngineV3;

fn percentile(sorted: &[u64], p: f64) -> u64 {
    if sorted.is_empty() {
        return 0;
    }
    let idx = ((sorted.len() as f64) * p / 100.0).min((sorted.len() - 1) as f64) as usize;
    sorted[idx]
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    let engine_name = args
        .windows(2)
        .find(|w| w[0] == "--engine")
        .map(|w| w[1].as_str())
        .unwrap_or("v1");
    let latency = args.contains(&"--latency".to_string());

    let mut engine: Box<dyn Engine> = match engine_name {
        "v1" => Box::new(EngineV1::default()),
        "v2" => Box::new(EngineV2::default()),
        "v3" => Box::new(EngineV3::default()),
        _ => {
            eprintln!("Unknown engine: {}. Use v1, v2, or v3.", engine_name);
            std::process::exit(1);
        }
    };

    let stdin = io::stdin();
    let mut stdout = io::stdout();
    let mut count = 0u64;
    let start = Instant::now();
    let mut latencies: Vec<u64> = Vec::new();

    for line in stdin.lock().lines() {
        let line = match line {
            Ok(l) => l,
            Err(e) => {
                eprintln!("Read error: {}", e);
                break;
            }
        };
        let line = line.trim();
        if line.is_empty() {
            continue;
        }

        let cmd = match parse_command(line) {
            Ok(c) => c,
            Err(e) => {
                let out = serialize_event(&engine_core::Event::Rejected {
                    reason: format!("parse error: {}", e),
                });
                writeln!(stdout, "{}", out).ok();
                continue;
            }
        };

        let cmd_start = if latency { Some(Instant::now()) } else { None };
        let events = engine.submit(cmd);
        if let Some(s) = cmd_start {
            latencies.push(s.elapsed().as_nanos() as u64);
        }
        for event in &events {
            writeln!(stdout, "{}", serialize_event(event)).ok();
        }
        count += 1;
    }

    let elapsed = start.elapsed();
    let secs = elapsed.as_secs_f64();
    let ops_per_sec = if secs > 0.0 { count as f64 / secs } else { 0.0 };
    eprintln!("Processed {} commands in {:?} ({:.0} ops/sec)", count, elapsed, ops_per_sec);

    if latency && !latencies.is_empty() {
        latencies.sort_unstable();
        eprintln!(
            "Latency (ns): p50={} p99={} p999={}",
            percentile(&latencies, 50.0),
            percentile(&latencies, 99.0),
            percentile(&latencies, 99.9)
        );
    }
}
