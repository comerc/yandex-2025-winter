use std::env;
use std::io::{self, BufRead, BufWriter, Write};
use std::time::Instant;

const MOD: i64 = 998244353;
const MAX_N: usize = 3005;

// get_memory_usage возвращает текущее использование памяти в байтах (приблизительно)
// Работает на Linux (читает /proc/self/status) и macOS (использует системные вызовы)
fn get_memory_usage() -> Option<u64> {
    #[cfg(target_os = "linux")]
    {
        use std::fs;
        if let Ok(content) = fs::read_to_string("/proc/self/status") {
            for line in content.lines() {
                if line.starts_with("VmRSS:") {
                    if let Some(value) = line.split_whitespace().nth(1) {
                        if let Ok(kb) = value.parse::<u64>() {
                            return Some(kb * 1024); // Конвертируем KB в байты
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
                    return Some(kb * 1024); // Конвертируем KB в байты
                }
            }
        }
    }

    None
}

// check_limits проверяет ограничения времени и памяти (работает только если установлена переменная окружения CHECK_LIMITS)
// Результаты выводятся в stderr, функция ничего не возвращает
fn check_limits(max_time_ms: u64, max_memory_mb: u64, f: impl FnOnce()) {
    // Проверяем переменную окружения
    if env::var("CHECK_LIMITS").is_err() {
        // Если переменная не установлена, просто выполняем функцию без проверок
        f();
        return;
    }

    // Измеряем память до выполнения
    let mem_before = get_memory_usage();

    // Измеряем время выполнения
    let start = Instant::now();
    f();
    let elapsed = start.elapsed();

    // Измеряем память после выполнения
    let mem_after = get_memory_usage();

    // Проверяем ограничения
    let elapsed_ms = elapsed.as_millis() as u64;
    let time_ok = elapsed_ms <= max_time_ms;

    // Вычисляем использованную память
    let memory_ok = if let (Some(before), Some(after)) = (mem_before, mem_after) {
        let allocated = if after > before { after - before } else { 0 };
        let max_memory_bytes = max_memory_mb * 1024 * 1024;
        allocated <= max_memory_bytes
    } else {
        true // Если не удалось измерить память, считаем что всё ОК
    };

    // Логируем результаты
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

fn init_stirling_sum() -> Vec<Vec<i64>> {
    let mut stirling_sum = vec![vec![0i64; MAX_N]; MAX_N];
    let mut prev_stirling = vec![0i64; MAX_N];
    let mut curr_stirling = vec![0i64; MAX_N];

    prev_stirling[0] = 1;
    stirling_sum[0][0] = 1;
    for j in 1..MAX_N {
        stirling_sum[0][j] = 1;
    }

    for i in 1..MAX_N {
        curr_stirling[0] = 0; // S(n, 0) = 0 for n >= 1
        stirling_sum[i][0] = 0;
        let mut current_sum: i64 = 0;

        for j in 1..=i {
            // S(n, k) = S(n-1, k-1) + (n-1)*S(n-1, k)
            let val = prev_stirling[j - 1] + ((i - 1) as i64 * prev_stirling[j]) % MOD;
            curr_stirling[j] = val % MOD;
            current_sum = (current_sum + curr_stirling[j]) % MOD;
            stirling_sum[i][j] = current_sum;
        }
        // Fill remaining sums with the total sum for this row
        for j in (i + 1)..MAX_N {
            stirling_sum[i][j] = current_sum;
        }

        // Update prev_stirling for next iteration
        for j in 0..=i {
            prev_stirling[j] = curr_stirling[j];
        }
    }

    stirling_sum
}

fn solve(
    n: usize,
    q: usize,
    l: usize,
    r: usize,
    b: &[usize],
    c: &[usize],
    stirling_sum: &[Vec<i64>],
) -> i64 {
    let mut in_deg = vec![0; n + 1];
    let mut out_deg = vec![0; n + 1];
    let mut adj = vec![0; n + 1];

    let mut possible = true;
    for i in 0..q {
        let u = b[i];
        let v = c[i];
        if out_deg[u] > 0 || in_deg[v] > 0 {
            possible = false;
        }
        out_deg[u] += 1;
        in_deg[v] += 1;
        adj[u] = v;
    }

    if !possible {
        return 0;
    }

    // M = n - q (number of path components)
    let m = n - q;

    // Count fixed cycles
    let mut fixed_cycles = 0;
    let mut visited = vec![false; n + 1];

    // First, traverse everything starting from nodes with in_deg == 0 (Path starts)
    for i in 1..=n {
        if in_deg[i] == 0 {
            let mut curr = i;
            while curr != 0 && !visited[curr] {
                visited[curr] = true;
                if out_deg[curr] > 0 {
                    curr = adj[curr];
                } else {
                    curr = 0;
                }
            }
        }
    }

    // Remaining unvisited nodes must be part of cycles
    for i in 1..=n {
        if !visited[i] {
            // Found a cycle
            fixed_cycles += 1;
            let mut curr = i;
            while !visited[curr] {
                visited[curr] = true;
                curr = adj[curr];
            }
        }
    }

    // Range of cycles needed from path components
    let need_l = if l >= fixed_cycles {
        l - fixed_cycles
    } else {
        0
    };
    let mut need_r = if r >= fixed_cycles {
        r - fixed_cycles
    } else {
        return 0;
    };

    if need_l > m {
        return 0;
    }
    if need_r > m {
        need_r = m;
    }

    // Answer is sum of Stirling numbers [M][k] for k in [needL, needR]
    let sub = if need_l > 0 {
        stirling_sum[m][need_l - 1]
    } else {
        0
    };
    let ans = (stirling_sum[m][need_r] - sub + MOD) % MOD;
    ans
}

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Initialize Stirling numbers table
    let stirling_sum = init_stirling_sum();

    check_limits(1000, 128, || {
        let mut line = String::new();
        reader.read_line(&mut line).unwrap();
        let t: usize = line.trim().parse().unwrap();

        for _ in 0..t {
            let mut line = String::new();
            reader.read_line(&mut line).unwrap();
            let parts: Vec<&str> = line.trim().split_whitespace().collect();
            let n: usize = parts[0].parse().unwrap();
            let q: usize = parts[1].parse().unwrap();
            let l: usize = parts[2].parse().unwrap();
            let r: usize = parts[3].parse().unwrap();

            let mut line = String::new();
            reader.read_line(&mut line).unwrap();
            let b: Vec<usize> = line
                .trim()
                .split_whitespace()
                .map(|x| x.parse().unwrap())
                .collect();

            let mut line = String::new();
            reader.read_line(&mut line).unwrap();
            let c: Vec<usize> = line
                .trim()
                .split_whitespace()
                .map(|x| x.parse().unwrap())
                .collect();

            let result = solve(n, q, l, r, &b, &c, &stirling_sum);
            writeln!(writer, "{}", result).unwrap();
        }
    });

    writer.flush().unwrap();
}

