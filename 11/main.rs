use std::io::{self, BufRead, BufWriter, Write};

const MOD: i64 = 1_000_000_007;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let n: usize = line.trim().parse().unwrap();

    let result = solve(n);
    writeln!(writer, "{}", result).unwrap();
}

/// Вычисляет E(n) mod M для гиперкуба размера n
/// Решает систему линейных уравнений для ожидаемого времени в гиперкубе
/// используя модульную арифметику
fn solve(n: usize) -> i64 {
    if n == 1 {
        return 1;
    }

    let n_mod = n as i64 % MOD;

    // Метод прогонки (Thomas algorithm) в модульной арифметике
    // Прямой ход: преобразуем к виду E[k] = alpha[k]*E[k+1] + beta[k]
    let mut alpha = vec![0i64; n + 1];
    let mut beta = vec![0i64; n + 1];

    // Начальное условие: E[0] = 0
    alpha[0] = 0;
    beta[0] = 0;

    // Прямой ход: для k = 1..n-1
    for k in 1..n {
        let k_mod = k as i64 % MOD;
        let n_minus_k = (n - k) as i64 % MOD;

        // denominator = n - k*alpha[k-1]
        let denom = (n_mod - k_mod * alpha[k - 1] % MOD + MOD) % MOD;
        let denom_inv = mod_inverse(denom, MOD);

        // alpha[k] = (n-k) / denom
        alpha[k] = n_minus_k * denom_inv % MOD;

        // beta[k] = (n + k*beta[k-1]) / denom
        let numerator = (n_mod + k_mod * beta[k - 1] % MOD) % MOD;
        beta[k] = numerator * denom_inv % MOD;
    }

    // Граничное условие: E[n] = 1 + E[n-1]
    // E[n] = (1 + beta[n-1]) / (1 - alpha[n-1])
    let one_minus_alpha = (1 - alpha[n - 1] + MOD) % MOD;
    let one_minus_alpha_inv = mod_inverse(one_minus_alpha, MOD);
    let e_n = (1 + beta[n - 1]) % MOD * one_minus_alpha_inv % MOD;

    e_n
}

/// Вычисляет обратное число a^(-1) mod m используя малую теорему Ферма
fn mod_inverse(a: i64, m: i64) -> i64 {
    mod_pow(a, m - 2, m)
}

/// Вычисляет base^exp mod m с использованием быстрого возведения в степень
fn mod_pow(mut base: i64, mut exp: i64, m: i64) -> i64 {
    let mut result = 1i64;
    base %= m;
    while exp > 0 {
        if exp % 2 == 1 {
            result = result * base % m;
        }
        exp >>= 1;
        base = base * base % m;
    }
    result
}
