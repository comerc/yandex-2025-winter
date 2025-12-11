use std::io::{self, BufRead, BufWriter, Write, Read};

const MOD: i64 = 1000000007;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем M и N
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let m: i64 = parts[0].parse().unwrap();
    let n: usize = parts[1].parse().unwrap();

    // Читаем потребности групп
    let mut w = Vec::with_capacity(n);
    let mut total_need = 0i64;

    // Читаем числа, обрабатывая случай, когда они могут быть в нескольких строках
    let mut buffer = String::new();
    reader.read_to_string(&mut buffer).unwrap();
    let tokens: Vec<&str> = buffer.split_whitespace().collect();
    
    for i in 0..n {
        let val: i64 = tokens[i].parse().unwrap();
        w.push(val);
        total_need += val;
    }

    let result = solve(m, &w, total_need);
    writeln!(writer, "{}", result).unwrap();
    writer.flush().unwrap();
}

// solve находит минимальную сумму квадратов недостачи
// Алгоритм: для минимизации суммы квадратов при фиксированной сумме
// нужно распределить недостачу максимально равномерно
fn solve(m: i64, w: &[i64], total_need: i64) -> i64 {
    let n = w.len();
    let deficit = total_need - m;

    if deficit <= 0 {
        return 0;
    }

    // Создаем пары (потребность, индекс) для сортировки
    let mut pairs: Vec<(i64, usize)> = w.iter().enumerate().map(|(i, &w)| (w, i)).collect();

    // Сортируем по убыванию потребности
    pairs.sort_by(|a, b| b.0.cmp(&a.0));

    // Распределяем недостачу оптимальным образом
    let orig_n = n;
    let mut shortfall = vec![0i64; orig_n];
    let mut remaining = deficit;
    let mut current_n = n;

    // Проходим по группам (отсортированным по убыванию потребности)
    // и распределяем недостачу
    // Ключевая идея: если avgShortfall <= minW среди оставшихся групп,
    // то можем распределить равномерно. Иначе, группа с минимальной потребностью
    // получает максимальную недостачу (равную своей потребности).
    while remaining > 0 && current_n > 0 {
        let avg_shortfall = remaining / current_n as i64;
        let min_w = pairs[current_n - 1].0; // минимальная потребность (pairs отсортирован по убыванию)

        // Если средняя недостача <= минимальной потребности среди оставшихся,
        // то все группы могут принять эту недостачу
        if avg_shortfall <= min_w {
            // Распределяем равномерно среди оставшихся групп
            let base_shortfall = remaining / current_n as i64;
            let extra = remaining % current_n as i64;

            // Распределяем базовую недостачу
            for j in 0..current_n {
                shortfall[pairs[j].1] = base_shortfall;
            }

            // Распределяем остаток по одной единице
            for j in 0..extra as usize {
                shortfall[pairs[j].1] += 1;
            }
            break;
        } else {
            // Группа с минимальной потребностью не может принять среднюю недостачу
            // Даем ей максимально возможную недостачу (равную потребности)
            shortfall[pairs[current_n - 1].1] = pairs[current_n - 1].0;
            remaining -= pairs[current_n - 1].0;
            current_n -= 1; // удаляем эту группу из рассмотрения
        }
    }

    // Вычисляем сумму квадратов недостачи по модулю
    let mut result = 0i64;
    for &s in &shortfall {
        let sq = s % MOD;
        result = (result + sq * sq % MOD) % MOD;
    }

    result
}

