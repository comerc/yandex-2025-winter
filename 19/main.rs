use std::env;
use std::io::{self, Read, BufWriter, Write};
use std::time::Instant;

const MOD: i64 = 1_000_000_007;

// get_memory_usage возвращает текущее использование памяти в байтах (приблизительно)
fn get_memory_usage() -> Option<u64> {
    #[cfg(target_os = "linux")]
    {
        use std::fs;
        if let Ok(content) = fs::read_to_string("/proc/self/status") {
            for line in content.lines() {
                if line.starts_with("VmRSS:") {
                    if let Some(value) = line.split_whitespace().nth(1) {
                        if let Ok(kb) = value.parse::<u64>() {
                            return Some(kb * 1024);
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
                    return Some(kb * 1024);
                }
            }
        }
    }
    
    None
}

// check_limits проверяет ограничения времени и памяти
fn check_limits(max_time_ms: u64, max_memory_mb: u64, f: impl FnOnce()) {
    if env::var("CHECK_LIMITS").is_err() {
        f();
        return;
    }

    let mem_before = get_memory_usage();
    let start = Instant::now();
    f();
    let elapsed = start.elapsed();
    let mem_after = get_memory_usage();

    let elapsed_ms = elapsed.as_millis() as u64;
    let time_ok = elapsed_ms <= max_time_ms;

    let memory_ok = if let (Some(before), Some(after)) = (mem_before, mem_after) {
        let allocated = if after > before { after - before } else { 0 };
        let max_memory_bytes = max_memory_mb * 1024 * 1024;
        allocated <= max_memory_bytes
    } else {
        true
    };

    if !time_ok || !memory_ok {
        if elapsed_ms > max_time_ms {
            eprintln!("⚠️ Превышено время: {} мс (лимит: {} мс)", elapsed_ms, max_time_ms);
        }
        if let (Some(before), Some(after)) = (mem_before, mem_after) {
            let allocated = if after > before { after - before } else { 0 };
            let memory_mb = allocated as f64 / (1024.0 * 1024.0);
            if memory_mb > max_memory_mb as f64 {
                eprintln!("⚠️ Превышена память: {:.2} МБ (лимит: {} МБ)", memory_mb, max_memory_mb);
            }
        }
    } else {
        if let (Some(before), Some(after)) = (mem_before, mem_after) {
            let allocated = if after > before { after - before } else { 0 };
            let memory_mb = allocated as f64 / (1024.0 * 1024.0);
            eprintln!("✓ Время: {} мс, Память: {:.2} МБ", elapsed_ms, memory_mb);
        } else {
            eprintln!("✓ Время: {} мс, Память: не удалось измерить", elapsed_ms);
        }
    }
}

struct Scanner<R> {
    reader: R,
    buffer: Vec<u8>,
    pos: usize,
    cap: usize,
}

impl<R: Read> Scanner<R> {
    fn new(reader: R) -> Self {
        Self {
            reader,
            buffer: vec![0; 1 << 16],
            pos: 0,
            cap: 0,
        }
    }

    fn next_i64(&mut self) -> i64 {
        let mut n = 0i64;
        loop {
            if self.pos >= self.cap {
                self.cap = self.reader.read(&mut self.buffer).unwrap_or(0);
                self.pos = 0;
                if self.cap == 0 {
                    return 0;
                }
            }
            let c = self.buffer[self.pos];
            self.pos += 1;
            if c >= b'0' && c <= b'9' {
                n = (c - b'0') as i64;
                break;
            }
        }
        loop {
            if self.pos >= self.cap {
                self.cap = self.reader.read(&mut self.buffer).unwrap_or(0);
                self.pos = 0;
                if self.cap == 0 {
                    break;
                }
            }
            let c = self.buffer[self.pos];
            if c < b'0' || c > b'9' {
                break;
            }
            self.pos += 1;
            n = n * 10 + (c - b'0') as i64;
        }
        n
    }
}

fn solve() {
    let stdin = io::stdin();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());
    let mut scanner = Scanner::new(stdin.lock());

    let m = scanner.next_i64();
    if m == 0 {
        return;
    }

    let mut n: i64 = 0;
    let mut last_count: i64 = 0;

    for _ in 0..m {
        let val = scanner.next_i64();
        n += val;
        if val > 0 {
            last_count = val;
        }
    }

    if n == 0 {
        writeln!(writer, "0").unwrap();
        return;
    }

    let s = n - last_count;
    if s == 0 {
        writeln!(writer, "0").unwrap();
        return;
    }

    // Предвычисление обратных элементов
    let mut inv = vec![0i64; (n + 1) as usize];
    inv[1] = 1;
    for i in 2..=n {
        inv[i as usize] = (MOD - (MOD / i) * inv[(MOD % i) as usize] % MOD) % MOD;
    }

    let mut sum_inv: i64 = 0;

    // Сумма 1/(k(N-k)) для k от 1 до S
    for k in 1..=s {
        // term = 1/(k*(N-k)) = inv[k] * inv[N-k]
        let term = (inv[k as usize] * inv[(n - k) as usize]) % MOD;
        sum_inv = (sum_inv + term) % MOD;
    }

    // Coeff = N*(N-1)/2
    let mut coeff = (n * (n - 1)) % MOD;
    coeff = (coeff * inv[2]) % MOD;

    let ans = (coeff * sum_inv) % MOD;
    writeln!(writer, "{}", ans).unwrap();
    writer.flush().unwrap();
}

fn main() {
    check_limits(1000, 256, solve);
}

