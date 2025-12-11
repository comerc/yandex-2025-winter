use std::io::{self, BufRead, BufWriter, Write};
use std::collections::HashMap;

const MOD: i64 = 1000000007;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем n и k
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let n: i64 = parts[0].parse().unwrap();
    let k: usize = parts[1].parse().unwrap();

    // Читаем массив a
    line.clear();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let mut a = Vec::with_capacity(k);
    for i in 0..k {
        a.push(parts[i].parse::<i32>().unwrap());
    }

    let result = solve(n, k, &a);
    writeln!(writer, "{}", result).unwrap();
    writer.flush().unwrap();
}

// solve вычисляет количество делителей числа S = n! / (A * P) по модулю 10^9 + 7
fn solve(n: i64, _k: usize, a: &[i32]) -> i64 {
    let threshold = n.min(1000000);

    // Вычисляем разложение n! на простые множители только для простых <= threshold
    let factorial_primes = factorize_factorial_small(n, threshold);

    // Вычисляем разложение A на простые множители
    let a_primes = factorize_product(a);

    // Вычисляем разложение S = n! / (A * P)
    let mut s_primes = HashMap::new();

    // Обрабатываем простые числа <= threshold
    for (prime, exp) in factorial_primes {
        // Вычитаем разложение A
        let mut exp = exp;
        if let Some(&a_exp) = a_primes.get(&(prime as i32)) {
            exp -= a_exp as i64;
        }
        if exp > 0 {
            s_primes.insert(prime, exp);
        }
    }

    // Простые числа > threshold полностью уходят в P, поэтому их не учитываем в S

    // Вычисляем количество делителей S
    let mut result = 1i64;
    for exp in s_primes.values() {
        result = (result * (exp + 1)) % MOD;
    }

    result
}

// factorize_factorial_small вычисляет разложение n! на простые множители только для простых <= threshold
// Использует формулу Лежандра: v_p(n!) = sum_{i=1}^{∞} floor(n / p^i)
fn factorize_factorial_small(n: i64, threshold: i64) -> HashMap<i64, i64> {
    let primes = sieve(threshold as i32);
    let mut result = HashMap::new();

    // Учитываем все простые числа <= threshold
    for p in primes {
        let p64 = p as i64;
        result.insert(p64, legendre(n, p64));
    }

    result
}

// process_large_primes обрабатывает простые числа в диапазоне [low, high] потоково
// Для каждого простого числа вызывается callback, не сохраняя все числа в памяти
fn process_large_primes<F>(low: i32, high: i32, mut callback: F)
where
    F: FnMut(i32),
{
    if low > high {
        return;
    }

    // Находим простые числа до sqrt(high) для фильтрации
    let sqrt_high = int_sqrt(high);
    let base_primes = sieve(sqrt_high);

    // Обрабатываем сегментами для экономии памяти
    let segment_size = 100000; // размер сегмента (100 КБ на сегмент)
    let mut segment_low = low;
    while segment_low <= high {
        let segment_high = (segment_low + segment_size - 1).min(high);

        // Создаем массив для сегмента
        let mut segment = vec![true; (segment_high - segment_low + 1) as usize];

        // Применяем решето для каждого простого числа
        for p in &base_primes {
            // Находим первое число в сегменте, кратное p
            let start = (((segment_low + p - 1) / p) * p).max(p * p);
            let mut j = start;
            while j <= segment_high {
                segment[(j - segment_low) as usize] = false;
                j += p;
            }
        }

        // Обрабатываем простые числа из сегмента
        for (i, &is_prime) in segment.iter().enumerate() {
            if is_prime {
                let num = segment_low + i as i32;
                if num >= 2 {
                    callback(num);
                }
            }
        }

        segment_low += segment_size;
    }
}

fn int_sqrt(n: i32) -> i32 {
    if n < 2 {
        return n;
    }
    let mut left = 1;
    let mut right = n;
    while left < right {
        let mid = (left + right + 1) / 2;
        if mid * mid <= n {
            left = mid;
        } else {
            right = mid - 1;
        }
    }
    left
}

// legendre вычисляет показатель простого числа p в разложении n! на простые множители
fn legendre(n: i64, p: i64) -> i64 {
    let mut result = 0i64;
    let mut power = p;
    while power <= n {
        result += n / power;
        power *= p;
    }
    result
}

// factorize_product вычисляет разложение произведения элементов массива на простые множители
fn factorize_product(arr: &[i32]) -> HashMap<i32, i32> {
    let mut result = HashMap::new();
    for &num in arr {
        let factors = factorize(num);
        for (prime, exp) in factors {
            *result.entry(prime).or_insert(0) += exp;
        }
    }
    result
}

// factorize разлагает число на простые множители
fn factorize(n: i32) -> HashMap<i32, i32> {
    let mut result = HashMap::new();
    let mut num = n;
    let mut i = 2;
    while i * i <= num {
        while num % i == 0 {
            *result.entry(i).or_insert(0) += 1;
            num /= i;
        }
        i += 1;
    }
    if num > 1 {
        *result.entry(num).or_insert(0) += 1;
    }
    result
}

// sieve возвращает список простых чисел до n (решето Эратосфена)
fn sieve(n: i32) -> Vec<i32> {
    if n < 2 {
        return Vec::new();
    }

    let mut is_prime = vec![true; (n + 1) as usize];
    is_prime[0] = false;
    is_prime[1] = false;

    let mut i = 2;
    while i * i <= n {
        if is_prime[i as usize] {
            let mut j = i * i;
            while j <= n {
                is_prime[j as usize] = false;
                j += i;
            }
        }
        i += 1;
    }

    let mut primes = Vec::new();
    for i in 2..=n {
        if is_prime[i as usize] {
            primes.push(i);
        }
    }

    primes
}

