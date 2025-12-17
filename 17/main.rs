use std::env;
use std::io::{self, BufWriter, Write};
use std::time::Instant;

// get_memory_usage returns current memory usage in bytes (approximately)
// Works on Linux (reads /proc/self/status) and macOS (uses system calls)
fn get_memory_usage() -> Option<u64> {
    #[cfg(target_os = "linux")]
    {
        use std::fs;
        if let Ok(content) = fs::read_to_string("/proc/self/status") {
            for line in content.lines() {
                if line.starts_with("VmRSS:") {
                    if let Some(value) = line.split_whitespace().nth(1) {
                        if let Ok(kb) = value.parse::<u64>() {
                            return Some(kb * 1024); // Convert KB to bytes
                        }
                    }
                }
            }
        }
    }

    #[cfg(target_os = "macos")]
    {
        use std::process::Command;
        let pid = std::process::id().to_string();
        if let Ok(output) = Command::new("ps")
            .args(&["-o", "rss=", "-p", &pid])
            .output()
        {
            if let Ok(mem_str) = String::from_utf8(output.stdout) {
                if let Ok(kb) = mem_str.trim().parse::<u64>() {
                    return Some(kb * 1024); // Convert KB to bytes
                }
            }
        }
    }

    None
}

// check_limits checks time and memory limits (only works if CHECK_LIMITS environment variable is set)
// Results are printed to stderr, function returns nothing
fn check_limits(max_time_ms: u64, max_memory_mb: u64, f: impl FnOnce()) {
    // Check environment variable
    if env::var("CHECK_LIMITS").is_err() {
        // If variable is not set, just execute the function without checks
        f();
        return;
    }

    // Measure memory before execution
    let mem_before = get_memory_usage();

    // Measure execution time
    let start = Instant::now();
    f();
    let elapsed = start.elapsed();

    // Measure memory after execution
    let mem_after = get_memory_usage();

    // Check limits
    let elapsed_ms = elapsed.as_millis() as u64;
    let time_ok = elapsed_ms <= max_time_ms;

    // Calculate used memory
    let memory_ok = if let (Some(before), Some(after)) = (mem_before, mem_after) {
        let allocated = if after > before { after - before } else { 0 };
        let memory_mb = allocated as f64 / (1024.0 * 1024.0);
        let max_memory_bytes = max_memory_mb * 1024 * 1024;
        allocated <= max_memory_bytes
    } else {
        true // If memory measurement failed, assume it's OK
    };

    // Log results
    if !time_ok || !memory_ok {
        if elapsed_ms > max_time_ms {
            eprintln!("⚠️ Превышено время: {} мс (лимит: {} мс)", elapsed_ms, max_time_ms);
        }
        if let (Some(before), Some(after)) = (mem_before, mem_after) {
            let allocated = if after > before { after - before } else { 0 };
            let memory_mb = allocated as f64 / (1024.0 * 1024.0);
            if memory_mb > max_memory_mb as f64 {
                eprintln!(
                    "⚠️ Превышена память: {:.2} МБ (лимит: {} МБ)",
                    memory_mb, max_memory_mb
                );
            }
        }
    } else {
        if let (Some(before), Some(after)) = (mem_before, mem_after) {
            let allocated = if after > before { after - before } else { 0 };
            let memory_mb = allocated as f64 / (1024.0 * 1024.0);
            eprintln!("✓ Время: {} мс, Память: {:.2} МБ", elapsed_ms, memory_mb);
        } else {
            eprintln!(
                "✓ Время: {} мс, Память: не удалось измерить (лимит: {} МБ)",
                elapsed_ms, max_memory_mb
            );
        }
    }
}

fn solve(n: i32) -> (String, String) {
    let mut m = 0;
    let mut limit = 1;
    
    while limit <= n - 1 {
        limit *= 10;
        m += 1;
    }
    
    if m == 0 {
        m = 1;
    }

    let s = "9".repeat(m);
    (s.clone(), s)
}

fn main() {
    let mut input = String::new();
    if io::stdin().read_line(&mut input).is_err() {
        return;
    }
    
    let n: i32 = match input.trim().parse() {
        Ok(num) => num,
        Err(_) => return,
    };

    check_limits(1000, 256, || {
        let (a, d) = solve(n);
        let stdout = io::stdout();
        let mut writer = BufWriter::new(stdout.lock());
        writeln!(writer, "{}", a).unwrap();
        writeln!(writer, "{}", d).unwrap();
        writer.flush().unwrap();
    });
}

