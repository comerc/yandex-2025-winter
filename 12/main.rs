use std::io::{self, BufRead, BufWriter, Write};

const MOD: i64 = 998244353;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем n и m
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    if parts.len() < 2 {
        return;
    }
    let n: i64 = parts[0].parse().unwrap();
    let m: usize = parts[1].parse().unwrap();

    // Читаем хорошие числа
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let mut good = std::collections::HashSet::new();
    for i in 0..m {
        if i < parts.len() {
            let val: i32 = parts[i].parse().unwrap();
            good.insert(val);
        }
    }

    let result = solve(n, &good);
    writeln!(writer, "{}", result).unwrap();
    writer.flush().unwrap();
}

// solve находит количество чудесных чисел длины n
// Чудесное число: без лидирующих нулей, сумма любых трех последовательных цифр - хорошее число
fn solve(n: i64, good: &std::collections::HashSet<i32>) -> i64 {
    if n < 3 {
        return 0;
    }

    // Состояние: последние две цифры (d1, d2) -> индекс = d1*10 + d2
    // Матрица переходов: M[i][j] = 1, если можно перейти от состояния i к состоянию j
    // i = d1*10 + d2, j = d2*10 + d3, переход возможен если d1+d2+d3 - хорошее число

    // Строим матрицу переходов 100x100
    let mut m = vec![vec![0i64; 100]; 100];

    for d1 in 0..10 {
        for d2 in 0..10 {
            let from = (d1 * 10 + d2) as usize;
            for d3 in 0..10 {
                let sum = d1 + d2 + d3;
                if good.contains(&sum) {
                    let to = (d2 * 10 + d3) as usize;
                    m[from][to] = 1;
                }
            }
        }
    }

    // Начальный вектор: для чисел длины 2 (d1, d2), где d1 != 0
    let mut start = vec![0i64; 100];
    for d1 in 1..10 {
        for d2 in 0..10 {
            let idx = (d1 * 10 + d2) as usize;
            start[idx] = 1;
        }
    }

    // Если n == 2, возвращаем количество начальных состояний
    if n == 2 {
        return start.iter().sum::<i64>() % MOD;
    }

    // Возводим матрицу в степень (n-2), так как у нас уже есть первые 2 цифры
    let power = n - 2;
    let m_power = matrix_power(&m, power);

    // Умножаем начальный вектор на матрицу
    let mut result = 0i64;
    for i in 0..100 {
        for j in 0..100 {
            result = (result + start[i] * m_power[i][j] % MOD) % MOD;
        }
    }

    result
}

// matrix_power возводит матрицу в степень используя быстрое возведение
fn matrix_power(m: &[Vec<i64>], power: i64) -> Vec<Vec<i64>> {
    let n = m.len();

    // Инициализируем единичную матрицу
    let mut result = vec![vec![0i64; n]; n];
    for i in 0..n {
        result[i][i] = 1;
    }

    let mut base = m.to_vec();
    let mut p = power;

    // Быстрое возведение в степень
    while p > 0 {
        if p & 1 == 1 {
            result = matrix_multiply(&result, &base);
        }
        base = matrix_multiply(&base, &base);
        p >>= 1;
    }

    result
}

// matrix_multiply умножает две матрицы по модулю
fn matrix_multiply(a: &[Vec<i64>], b: &[Vec<i64>]) -> Vec<Vec<i64>> {
    let n = a.len();
    let m = b[0].len();
    let k = b.len();

    let mut result = vec![vec![0i64; m]; n];

    for i in 0..n {
        for j in 0..m {
            let mut sum = 0i64;
            for t in 0..k {
                sum = (sum + a[i][t] * b[t][j] % MOD) % MOD;
            }
            result[i][j] = sum;
        }
    }

    result
}

