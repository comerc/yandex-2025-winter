use std::io::{self, BufRead, BufWriter, Write};

const MAX_N: usize = 700000;

// Максимальное k, для которого существуют k-интересные числа <= maxN
// 2*3*4*5*6*7*8*9 = 362880 <= 700000
// 2*3*4*5*6*7*8*9*10 = 3628800 > 700000
// Значит maxK = 9
const MAX_K: usize = 9;

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Предвычисляем префиксные суммы k-интересных чисел
    let prefix_sums = precompute();

    // Читаем количество запросов
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let q: usize = line.trim().parse().unwrap();

    // Обрабатываем запросы
    for _i in 0..q {
        line.clear();
        reader.read_line(&mut line).unwrap();
        let parts: Vec<&str> = line.trim().split_whitespace().collect();
        if parts.len() < 3 {
            continue;
        }
        let k: usize = parts[0].parse().unwrap();
        let l: usize = parts[1].parse().unwrap();
        let r: usize = parts[2].parse().unwrap();

        let result = query(&prefix_sums, k, l, r);
        writeln!(writer, "{}", result).unwrap();
    }
    writer.flush().unwrap();
}

// precompute генерирует все k-интересные числа и строит префиксные суммы
fn precompute() -> Vec<Vec<i32>> {
    // prefix_sums[k][n] = количество k-интересных чисел от 1 до n
    let mut prefix_sums = vec![vec![0i32; MAX_N + 1]; MAX_K + 1];

    // Генерируем k-интересные числа для k >= 2 и сразу помечаем
    generate(1, 2, 0, MAX_N, &mut prefix_sums);

    // Все числа >= 2 являются 1-интересными
    for n in 2..=MAX_N {
        prefix_sums[1][n] = 1;
    }

    // Преобразуем в префиксные суммы
    for k in 1..=MAX_K {
        for n in 1..=MAX_N {
            prefix_sums[k][n] += prefix_sums[k][n - 1];
        }
    }

    prefix_sums
}

// generate рекурсивно генерирует все произведения возрастающих последовательностей множителей
fn generate(
    product: usize,
    min_factor: usize,
    depth: usize,
    limit: usize,
    prefix_sums: &mut Vec<Vec<i32>>,
) {
    let mut factor = min_factor;
    while product * factor <= limit {
        let new_product = product * factor;
        let new_depth = depth + 1;

        if new_depth >= 2 && new_depth <= MAX_K {
            prefix_sums[new_depth][new_product] = 1;
        }

        generate(new_product, factor + 1, new_depth, limit, prefix_sums);
        factor += 1;
    }
}

// query возвращает количество k-интересных чисел в диапазоне [l, r]
fn query(prefix_sums: &[Vec<i32>], k: usize, l: usize, r: usize) -> i32 {
    if k > MAX_K || l > MAX_N {
        return 0;
    }
    let r = r.min(MAX_N);
    let l = l.max(1);
    if l > r {
        return 0;
    }
    prefix_sums[k][r] - prefix_sums[k][l - 1]
}








