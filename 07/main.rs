use std::io::{self, BufRead, BufWriter, Write};

const MOD: i64 = 998244353;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Предвычисляем факториалы и обратные факториалы
    let max_n = 400001;
    let fact = precompute_factorials(max_n);
    let inv_fact = precompute_inv_factorials(&fact, max_n);

    // Читаем количество тестов
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let t: usize = line.trim().parse().unwrap();

    // Обрабатываем тесты
    for _i in 0..t {
        line.clear();
        reader.read_line(&mut line).unwrap();
        let parts: Vec<&str> = line.trim().split_whitespace().collect();
        let n: usize = parts[0].parse().unwrap();
        let s: usize = parts[1].parse().unwrap();

        let result = solve(n, s, &fact, &inv_fact);
        writeln!(writer, "{}", result).unwrap();
    }
    writer.flush().unwrap();
}

// solve считает количество замечательных массивов
// Формула: answer(n, s) = (n+1)! × C(s, n)
fn solve(n: usize, s: usize, fact: &[i64], inv_fact: &[i64]) -> i64 {
    if n > s {
        return 0;
    }

    // C(s, n) = s! / (n! * (s-n)!)
    let c = comb(s, n, fact, inv_fact);

    // (n+1)! × C(s, n)
    fact[n + 1] * c % MOD
}

// comb вычисляет C(n, k) по модулю mod
fn comb(n: usize, k: usize, fact: &[i64], inv_fact: &[i64]) -> i64 {
    if k > n {
        return 0;
    }
    fact[n] * inv_fact[k] % MOD * inv_fact[n - k] % MOD
}

// precompute_factorials предвычисляет факториалы до n
fn precompute_factorials(n: usize) -> Vec<i64> {
    let mut fact = vec![1i64; n + 1];
    for i in 1..=n {
        fact[i] = fact[i - 1] * (i as i64) % MOD;
    }
    fact
}

// precompute_inv_factorials предвычисляет обратные факториалы
fn precompute_inv_factorials(fact: &[i64], n: usize) -> Vec<i64> {
    let mut inv_fact = vec![0i64; n + 1];
    inv_fact[n] = mod_pow(fact[n], MOD - 2);
    for i in (1..=n).rev() {
        inv_fact[i - 1] = inv_fact[i] * (i as i64) % MOD;
    }
    inv_fact
}

// mod_pow вычисляет a^b mod mod
fn mod_pow(mut a: i64, mut b: i64) -> i64 {
    let mut result = 1i64;
    a %= MOD;
    while b > 0 {
        if b & 1 == 1 {
            result = result * a % MOD;
        }
        a = a * a % MOD;
        b >>= 1;
    }
    result
}







