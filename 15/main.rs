use std::env;
use std::io::{self, BufRead, BufWriter, Write};
use std::time::Instant;

const MOD: i64 = 998244353;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    check_limits(1000, 256, || {
        solve(&mut reader, &mut writer);
    });

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

fn solve<R: BufRead, W: Write>(reader: &mut R, writer: &mut W) {
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let t: usize = line.trim().parse().unwrap();

    for _t_idx in 0..t {
        let mut line = String::new();
        reader.read_line(&mut line).unwrap();
        let parts: Vec<&str> = line.trim().split_whitespace().collect();
        if parts.len() >= 2 {
            let s = parts[0];
            let t_str = parts[1];
            process_test_case(s, t_str, writer);
        } else if parts.len() == 1 {
            // Если только одно слово, читаем следующую строку
            let s = parts[0];
            let mut line2 = String::new();
            reader.read_line(&mut line2).unwrap();
            let t_str = line2.trim();
            process_test_case(s, t_str, writer);
        } else {
            // Пустая строка - читаем следующую строку как s, затем еще одну как t
            let mut line2 = String::new();
            reader.read_line(&mut line2).unwrap();
            let s = line2.trim();
            let mut line3 = String::new();
            reader.read_line(&mut line3).unwrap();
            let t_str = line3.trim();
            process_test_case(s, t_str, writer);
        }
    }
}

fn process_test_case<W: Write>(s: &str, t: &str, writer: &mut W) {
    let n = s.len();
    let m = t.len();

    // 1. Предподсчет количества разбиений для строки длины k.
    // Если длина k, то есть k-1 мест для разрыва, итого 2^(k-1) способов.
    // Для k=0 (пустая строка) считаем 1 способ.
    let mut partitions = vec![0i64; n + 1];
    partitions[0] = 1;
    if n > 0 {
        let mut p = 1i64;
        for i in 1..=n {
            partitions[i] = p;
            p = (p * 2) % MOD;
        }
    }

    // 2. Предподсчет LCP (Longest Common Prefix) для всех пар суффиксов s и t.
    // lcp[i][j] = длина общего префикса s[i:] и t[j:]
    let mut lcp = vec![vec![0usize; m + 1]; n + 1];
    for i in (0..n).rev() {
        for j in (0..m).rev() {
            if s.as_bytes()[i] == t.as_bytes()[j] {
                lcp[i][j] = 1 + lcp[i + 1][j + 1];
            } else {
                lcp[i][j] = 0;
            }
        }
    }

    // 3. Динамическое программирование
    // dp[i][u] - кол-во способов разбить суффикс s[i:], чтобы результат совпадал с t[u...]
    let mut dp = vec![vec![0i64; m + 1]; n + 1];

    // База: пустой суффикс s совпадает с "началом" любой подстроки t (пустой строкой)
    for u in 0..=m {
        dp[n][u] = 1;
    }

    let mut ans = 0i64;

    // Перебираем длину оставшегося суффикса s (от меньшего к большему с точки зрения "потребления" s справа налево)
    // i идет от n до 1. Мы пытаемся откусить кусок s[k...i-1].
    for i in (1..=n).rev() {
        // Длина уже сформированной части результата
        let current_res_len = n - i;

        for u in 0..=m {
            if dp[i][u] == 0 {
                continue;
            }

            // Позиция в t, с которой мы должны сравнивать следующий кусок
            let pos_in_t = u + current_res_len;

            // Если мы вышли за пределы t, сравнение невозможно (строка t кончилась раньше)
            if pos_in_t >= m {
                continue;
            }

            // Перебираем, где отрезать следующий кусок от s (индекс k)
            // Кусок будет s[k...i-1]
            for k in (0..i).rev() {
                let chunk_len = i - k;

                // Используем LCP для быстрого сравнения s[k...] и t[posInT...]
                let val = lcp[k][pos_in_t];

                if val >= chunk_len {
                    // Кусок полностью совпал с частью t
                    // Если мы не вышли за границы t, обновляем ДП
                    if pos_in_t + chunk_len <= m {
                        dp[k][u] = (dp[k][u] + dp[i][u]) % MOD;
                    }
                } else {
                    // Куски различаются.
                    // Индекс различия относительно начала куска: val
                    let idx_s = k + val;
                    let idx_t = pos_in_t + val;

                    // Проверяем, что различие произошло в пределах строки t
                    if idx_t < m {
                        // Если символ в s меньше символа в t, то результат кодирования лексикографически меньше
                        if s.as_bytes()[idx_s] < t.as_bytes()[idx_t] {
                            // Мы нашли "меньший" вариант.
                            // Длина совпадающего префикса = currentResLen + val.
                            let match_len = current_res_len + val;

                            // Этот результат будет меньше любой подстроки t, начинающейся в u,
                            // которая длиннее matchLen.
                            // Максимальная длина подстроки от u: m - u.
                            // Подходящие длины: matchLen+1, matchLen+2, ..., m-u.
                            let count = (m - u) - match_len;

                            if count > 0 {
                                // Добавляем к ответу:
                                // (способы дойти до i) * (способы разбить остаток s[0...k-1]) * (кол-во подстрок t)
                                let ways = (dp[i][u] * partitions[k]) % MOD;
                                let term = (ways * count as i64) % MOD;
                                ans = (ans + term) % MOD;
                            }
                        }
                        // Если s[idxS] > t[idxT], то результат больше, ничего не делаем.
                    }
                }
            }
        }
    }

    // Обработка случаев, когда вся строка s (переставленная) является строгим префиксом подстроки t.
    // Это соответствует состояниям dp[0][u].
    // Результат имеет длину n и совпадает с t[u ... u+n-1].
    // Он будет меньше любой подстроки t[u...], длина которой > n.
    for u in 0..=m {
        if dp[0][u] > 0 {
            let count = (m - u) - n;
            if count > 0 {
                let term = (dp[0][u] * count as i64) % MOD;
                ans = (ans + term) % MOD;
            }
        }
    }

    writeln!(writer, "{}", ans).unwrap();
}

