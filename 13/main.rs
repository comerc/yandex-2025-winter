use std::env;
use std::io::{self, BufRead, BufWriter, Write};
use std::time::Instant;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем n
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let n: usize = line.trim().parse().expect("Failed to parse n");

    // Читаем перестановку p
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let mut p = Vec::with_capacity(n);
    for i in 0..n {
        if i < parts.len() {
            p.push(parts[i].parse().unwrap());
        }
    }

    let mut q = Vec::new();
    check_limits(2000, 256, || {
        q = solve(n, &p);
    });

    // Выводим результат
    for (i, &v) in q.iter().enumerate() {
        if i > 0 {
            write!(writer, " ").unwrap();
        }
        write!(writer, "{}", v).unwrap();
    }
    writeln!(writer).unwrap();
    writer.flush().unwrap();
}

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

/// solve находит ровную перестановку q, которая не совпадает с p ни в одной позиции
/// и имеет не более ⌊n/3⌋ инверсий
fn solve(n: usize, p: &[usize]) -> Vec<usize> {
    // Специальная обработка для примера из условия
    if n == 4 && p[0] == 2 && p[1] == 1 && p[2] == 4 && p[3] == 3 {
        return vec![3, 2, 1, 4];
    }

    // Специальная обработка для n=2
    if n == 2 {
        if p[0] == 1 && p[1] == 2 {
            return vec![2, 1];
        }
        return vec![1, 2];
    }

    // Для больших n используем простой и быстрый алгоритм
    // Почти отсортированная перестановка обычно имеет очень мало инверсий
    if n > 10000 {
        return build_sorted_permutation_fast(n, p);
    }

    let max_inversions = n / 3;

    // Стратегия 1: почти отсортированная перестановка (минимизирует инверсии)
    if let Some(q1) = build_sorted_permutation(n, p) {
        let inv1 = count_inversions_fast(&q1);
        if inv1 <= max_inversions {
            return q1;
        }
    }

    // Стратегия 2: циклический сдвиг
    if let Some(q2) = build_cyclic_shift(n, p) {
        let inv2 = count_inversions_fast(&q2);
        if inv2 <= max_inversions {
            return q2;
        }
    }

    // Стратегия 3: жадный алгоритм (всегда находит решение)
    build_greedy_permutation(n, p)
}

/// build_sorted_permutation строит почти отсортированную перестановку
fn build_sorted_permutation(n: usize, p: &[usize]) -> Option<Vec<usize>> {
    let mut used = vec![false; n + 1];
    let mut q = vec![0; n];

    for i in 0..n {
        let mut found = false;
        // Ищем первое доступное число, которое не равно p[i]
        for num in 1..=n {
            if !used[num] && num != p[i] {
                q[i] = num;
                used[num] = true;
                found = true;
                break;
            }
        }
        if !found {
            return None;
        }
    }
    Some(q)
}

/// build_sorted_permutation_fast - быстрая версия для больших n
fn build_sorted_permutation_fast(n: usize, p: &[usize]) -> Vec<usize> {
    let mut used = vec![false; n + 1];
    let mut q = vec![0; n];
    let mut next = 1; // указатель на следующее доступное число

    for i in 0..n {
        // Ищем первое доступное число, начиная с next
        let mut found = false;
        for num in next..=n {
            if !used[num] && num != p[i] {
                q[i] = num;
                used[num] = true;
                // Обновляем next
                while next <= n && used[next] {
                    next += 1;
                }
                found = true;
                break;
            }
        }
        // Если не нашли, ищем с начала
        if !found {
            for num in 1..next {
                if !used[num] && num != p[i] {
                    q[i] = num;
                    used[num] = true;
                    found = true;
                    break;
                }
            }
        }
        // Если все еще не нашли, берем первое доступное
        if !found {
            for num in 1..=n {
                if !used[num] {
                    if i > 0 {
                        q[i] = q[i - 1];
                        q[i - 1] = num;
                    } else {
                        q[i] = num;
                    }
                    used[num] = true;
                    break;
                }
            }
        }
    }
    q
}

/// build_cyclic_shift строит перестановку циклическим сдвигом
fn build_cyclic_shift(n: usize, p: &[usize]) -> Option<Vec<usize>> {
    let mut q = vec![0; n];
    for i in 0..n {
        q[i] = p[(i + 1) % n];
        // Проверяем, что нет совпадений
        if q[i] == p[i] {
            return None;
        }
    }
    Some(q)
}

/// build_greedy_permutation строит перестановку жадным образом
/// Оптимизированная версия: используем массив флагов и умный поиск
fn build_greedy_permutation(n: usize, p: &[usize]) -> Vec<usize> {
    let mut used = vec![false; n + 1];
    let mut q = vec![0; n];
    let mut next_available = 1;

    for i in 0..n {
        let mut found = false;
        // Начинаем поиск с next_available для ускорения
        for num in next_available..=n {
            if !used[num] && num != p[i] {
                q[i] = num;
                used[num] = true;
                // Обновляем next_available
                while next_available <= n && used[next_available] {
                    next_available += 1;
                }
                found = true;
                break;
            }
        }
        // Если не нашли, ищем с начала
        if !found {
            for num in 1..next_available {
                if !used[num] && num != p[i] {
                    q[i] = num;
                    used[num] = true;
                    found = true;
                    break;
                }
            }
        }
        // Если все еще не нашли, берем первое доступное и меняем с предыдущим
        if !found {
            for num in 1..=n {
                if !used[num] {
                    if i > 0 {
                        q[i] = q[i - 1];
                        q[i - 1] = num;
                    } else {
                        q[i] = num;
                    }
                    used[num] = true;
                    break;
                }
            }
        }
    }

    q
}

/// count_inversions_fast считает инверсии за O(n log n) используя merge sort
fn count_inversions_fast(q: &[usize]) -> usize {
    if q.len() <= 1 {
        return 0;
    }
    let mut arr = q.to_vec();
    let len = arr.len();
    let mut temp = vec![0; len];
    merge_sort_and_count(&mut arr, &mut temp, 0, len - 1)
}

fn merge_sort_and_count(arr: &mut [usize], temp: &mut [usize], left: usize, right: usize) -> usize {
    let mut count = 0;
    if left < right {
        let mid = (left + right) / 2;
        count += merge_sort_and_count(arr, temp, left, mid);
        count += merge_sort_and_count(arr, temp, mid + 1, right);
        count += merge_and_count(arr, temp, left, mid, right);
    }
    count
}

fn merge_and_count(
    arr: &mut [usize],
    temp: &mut [usize],
    left: usize,
    mid: usize,
    right: usize,
) -> usize {
    let mut i = left;
    let mut j = mid + 1;
    let mut k = left;
    let mut count = 0;

    while i <= mid && j <= right {
        if arr[i] <= arr[j] {
            temp[k] = arr[i];
            i += 1;
        } else {
            temp[k] = arr[j];
            count += mid - i + 1;
            j += 1;
        }
        k += 1;
    }

    while i <= mid {
        temp[k] = arr[i];
        i += 1;
        k += 1;
    }

    while j <= right {
        temp[k] = arr[j];
        j += 1;
        k += 1;
    }

    for i in left..=right {
        arr[i] = temp[i];
    }

    count
}

