use std::env;
use std::io::{self, BufRead, BufWriter, Write};
use std::time::Instant;

const MOD: i64 = 1000000007;
const INV2: i64 = 500000004; // 2^(-1) mod 10^9+7
const INV3: i64 = 333333336; // 3^(-1) mod 10^9+7

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let a: i64 = parts[0].parse().unwrap();
    let q: i64 = parts[1].parse().unwrap();
    let l: i64 = parts[2].parse().unwrap();
    let r: i64 = parts[3].parse().unwrap();

    let mut result = 0i64;
    check_limits(1000, 256, || {
        result = solve(a, q, l, r);
    });

    writeln!(writer, "{}", result).unwrap();
    writer.flush().unwrap();
}

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

// check_limits проверяет ограничения времени и памяти (работает только если установлена переменная окружения CHECK_LIMITS)
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

// solve находит количество четвёрок (n, m, k, s) таких, что (aq^n - aq^m) / (aq^k - aq^s) - целое число
fn solve(a: i64, q: i64, l: i64, r: i64) -> i64 {
    let n = r - l + 1;
    if n <= 0 {
        return 0;
    }

    let n_mod = n % MOD;

    // Специальные случаи
    if a == 0 {
        return 0;
    }

    // Случай q = 0
    if q == 0 {
        let mut zeros = 0i64;
        if l <= 0 && r >= 0 {
            zeros = 1;
        }
        if zeros == 0 {
            return 0;
        }
        let non_zeros = (n - zeros) % MOD;
        let mut valid_denoms = (2 * zeros) % MOD;
        valid_denoms = (valid_denoms * non_zeros) % MOD;
        let all_nums = (n_mod * n_mod) % MOD;
        return (valid_denoms * all_nums) % MOD;
    }

    // Случай |q| = 1
    if q == 1 {
        return 0;
    }
    if q == -1 {
        let l_is_even = l % 2 == 0;
        let (evens, odds) = if n % 2 == 0 {
            (n / 2, n / 2)
        } else {
            if l_is_even {
                (n / 2 + 1, n / 2)
            } else {
                (n / 2, n / 2 + 1)
            }
        };
        let mut valid_denoms = (2 * (evens % MOD)) % MOD;
        valid_denoms = (valid_denoms * (odds % MOD)) % MOD;
        let all_nums = (n_mod * n_mod) % MOD;
        return (valid_denoms * all_nums) % MOD;
    }

    // Общий случай: q ≠ 0, q ≠ 1, a ≠ 0
    solve_general(a, q, l, r)
}

// solve_general обрабатывает общий случай q ≠ 0, q ≠ 1, a ≠ 0
fn solve_general(_a: i64, q: i64, l: i64, r: i64) -> i64 {
    let n = r - l + 1;
    if n <= 0 {
        return 0;
    }

    let n_mod = n % MOD;

    // 1. Вклад n = m (числитель 0). Знаменатель любой ненулевой (k != s).
    // Кол-во = N * (N * (N - 1))
    let mut ans = (n_mod * n_mod) % MOD;
    ans = (ans * ((n_mod - 1 + MOD) % MOD)) % MOD;

    // 2. Вклад n != m. Используем Block Summation (Sqrt Decomposition)
    // Формула: 2 * [ (kB)^2 - kB(2N+1) + N(N+1) ]
    // где k - множитель (A = k*B)

    // Предподсчет постоянных частей формулы
    let term_n_n1 = (n_mod * ((n_mod + 1) % MOD)) % MOD; // N(N+1)
    let term_2n_1 = ((2 * n_mod) + 1) % MOD; // 2N+1

    let limit = n - 1;
    let mut l_val = 1i64;

    // Основной цикл оптимизирован
    while l_val <= limit {
        let k_max = limit / l_val;
        let mut r_val = limit / k_max;
        if r_val > limit {
            r_val = limit;
        }

        let k_max_mod = k_max % MOD;

        // Вычисляем суммы для k от 1 до kMax
        // s1 = sum(k) = k(k+1)/2
        let mut s1 = (k_max_mod * (k_max_mod + 1)) % MOD;
        s1 = (s1 * INV2) % MOD;

        // s2 = sum(k^2) = k(k+1)(2k+1)/6
        let term_2k1 = ((2 * k_max_mod) + 1) % MOD;
        let mut s2 = (s1 * term_2k1) % MOD;
        s2 = (s2 * INV3) % MOD;

        let s0 = k_max_mod;

        // Вычисляем суммы для B на отрезке [l, r]
        let l_mod = l_val % MOD;
        let r_mod = r_val % MOD;

        // Сумма 1..r для B и B^2
        let mut sum_r1 = (r_mod * (r_mod + 1)) % MOD;
        sum_r1 = (sum_r1 * INV2) % MOD;
        let term_2r1 = ((2 * r_mod) + 1) % MOD;
        let mut sum_r2 = (sum_r1 * term_2r1) % MOD;
        sum_r2 = (sum_r2 * INV3) % MOD;

        // Сумма 1..l-1 для B и B^2
        let lm1 = (l_mod - 1 + MOD) % MOD;
        let mut sum_l1 = (lm1 * (lm1 + 1)) % MOD;
        sum_l1 = (sum_l1 * INV2) % MOD;
        let term_2l1 = ((2 * lm1) + 1) % MOD;
        let mut sum_l2 = (sum_l1 * term_2l1) % MOD;
        sum_l2 = (sum_l2 * INV3) % MOD;

        // Разность сумм (значения на отрезке)
        let ss2 = (sum_r2 - sum_l2 + MOD) % MOD;
        let ss1 = (sum_r1 - sum_l1 + MOD) % MOD;
        let ss0 = (r_mod - l_mod + 1 + MOD) % MOD;

        // Собираем итоговое выражение для блока
        // blockSum = 2 * [ s2*ss2 - s1*ss1*(2N+1) + s0*ss0*N(N+1) ]
        let p1 = (s2 * ss2) % MOD;

        let mut p2 = (term_2n_1 * s1) % MOD;
        p2 = (p2 * ss1) % MOD;

        let mut p3 = (term_n_n1 * s0) % MOD;
        p3 = (p3 * ss0) % MOD;

        let mut block_sum = (p1 - p2 + MOD) % MOD;
        block_sum = (block_sum + p3) % MOD;
        block_sum = (block_sum * 2) % MOD;

        ans = (ans + block_sum) % MOD;

        l_val = r_val + 1;
    }

    // Коррекция для q = -2
    if q == -2 && n >= 3 {
        let mut max_odd = limit;
        if max_odd % 2 == 0 {
            max_odd -= 1;
        }

        if max_odd >= 1 {
            // Кол-во нечетных чисел <= maxOdd
            let cnt = ((max_odd + 1) / 2) % MOD;

            // sumA1 = сумма нечетных A = cnt^2
            let sum_a1 = (cnt * cnt) % MOD;

            // sumA2 = сумма квадратов нечетных A = cnt(4*cnt^2 - 1)/3
            let mut term_sq = (4 * cnt * cnt) % MOD;
            term_sq = (term_sq - 1 + MOD) % MOD;
            let mut sum_a2 = (cnt * term_sq) % MOD;
            sum_a2 = (sum_a2 * INV3) % MOD;

            // Подставляем в общую формулу 2 * [ A^2 - A(2N+1) + N(N+1) ]
            let mut val = sum_a2;
            let sub_val = (term_2n_1 * sum_a1) % MOD;
            val = (val - sub_val + MOD) % MOD;
            let add_val = (term_n_n1 * cnt) % MOD;
            val = (val + add_val) % MOD;

            let mut total_odd = (val * 2) % MOD;

            // КОРРЕКЦИЯ ДЛЯ A=1 (пересечение условий)
            total_odd = (total_odd - 4 + MOD) % MOD;

            ans = (ans + total_odd) % MOD;
        }
    }

    ans
}

