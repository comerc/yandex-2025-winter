use std::env;
use std::io::{self, BufRead, BufWriter, Write, Read};
use std::time::Instant;
use std::collections::HashMap;

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
        let memory_mb = allocated as f64 / (1024.0 * 1024.0);
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
            eprintln!("✓ Время: {} мс, Память: не удалось измерить (лимит: {} МБ)", elapsed_ms, max_memory_mb);
        }
    }
}

// Константы
const MAX_NODES: usize = 13_000_000;
const MAX_BITS: usize = 29;

// Глобальные статические массивы через static mut для производительности и соответствия Go версии
// В Rust это unsafe, но в олимпиадном программировании допустимо для скорости
static mut L: [i32; MAX_NODES] = [0; MAX_NODES];
static mut R: [i32; MAX_NODES] = [0; MAX_NODES];
static mut CNT: [i32; MAX_NODES] = [0; MAX_NODES];
static mut MEMO: [i32; MAX_NODES] = [0; MAX_NODES];
static mut PTR: i32 = 0;

unsafe fn new_node() -> i32 {
    PTR += 1;
    let idx = PTR as usize;
    L[idx] = 0;
    R[idx] = 0;
    CNT[idx] = 0;
    MEMO[idx] = 0;
    PTR
}

unsafe fn push_up(u: usize, bit: usize) {
    let idx0 = L[u] as usize;
    let idx1 = R[u] as usize;

    let c0 = if idx0 != 0 { CNT[idx0] } else { 0 };
    let c1 = if idx1 != 0 { CNT[idx1] } else { 0 };

    CNT[u] = c0 + c1;

    let m0 = if idx0 != 0 { MEMO[idx0] } else { 0 };
    let m1 = if idx1 != 0 { MEMO[idx1] } else { 0 };

    // Размер полного поддерева на текущем уровне
    let full = 1 << bit;

    if c0 == full {
        MEMO[u] = full + m1;
    } else if c1 == full {
        MEMO[u] = full + m0;
    } else {
        MEMO[u] = if m0 > m1 { m0 } else { m1 };
    }
}

unsafe fn update(u: usize, val: i32, bit: i32, add: bool) {
    if bit < 0 {
        if add {
            CNT[u] = 1;
            MEMO[u] = 1;
        } else {
            CNT[u] = 0;
            MEMO[u] = 0;
        }
        return;
    }

    let dir = (val >> bit) & 1;
    let mut child_idx;

    if dir == 0 {
        if L[u] == 0 {
            L[u] = new_node();
        }
        child_idx = L[u] as usize;
    } else {
        if R[u] == 0 {
            R[u] = new_node();
        }
        child_idx = R[u] as usize;
    }

    update(child_idx, val, bit - 1, add);
    push_up(u, bit as usize);
}

fn solve() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let mut buffer = String::new();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Функция для чтения одного токена как целого числа
    let mut input_buffer = Vec::new();
    let mut cursor = 0;
    
    // Чтение всего ввода сразу (fast I/O)
    reader.read_to_end(&mut input_buffer).unwrap();
    let input_str = unsafe { std::str::from_utf8_unchecked(&input_buffer) };
    let mut tokens = input_str.split_ascii_whitespace();
    
    let mut next_int = || -> i32 {
        tokens.next().unwrap().parse().unwrap()
    };

    let t = next_int();

    unsafe {
        PTR = 0;
        
        for _ in 0..t {
            let n = next_int();
            let q = next_int();

            let mut a = Vec::with_capacity(n as usize);
            let mut freq = HashMap::with_capacity(n as usize);
            
            let root = new_node() as usize;

            for _ in 0..n {
                let val = next_int();
                a.push(val);
                *freq.entry(val).or_insert(0) += 1;
                
                if freq[&val] == 1 {
                    update(root, val, MAX_BITS as i32, true);
                }
            }

            writeln!(writer, "{}", MEMO[root]).unwrap();

            for _ in 0..q {
                let j = next_int() - 1; // 0-based
                let v = next_int();

                let old_val = a[j as usize];
                if old_val != v {
                    // Удаляем старое
                    if let Some(count) = freq.get_mut(&old_val) {
                        *count -= 1;
                        if *count == 0 {
                            update(root, old_val, MAX_BITS as i32, false);
                        }
                    }

                    // Обновляем массив
                    a[j as usize] = v;

                    // Добавляем новое
                    *freq.entry(v).or_insert(0) += 1;
                    if freq[&v] == 1 {
                        update(root, v, MAX_BITS as i32, true);
                    }
                }

                writeln!(writer, "{}", MEMO[root]).unwrap();
            }
        }
    }
    
    writer.flush().unwrap();
}

fn main() {
    check_limits(4000, 256, solve);
}

