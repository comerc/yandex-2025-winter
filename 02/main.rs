use std::io::{self, BufRead, BufWriter, Write};

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем 10 чисел
    let mut nums = Vec::with_capacity(10);
    for _i in 0..10 {
        let mut line = String::new();
        reader.read_line(&mut line).unwrap();
        let num: i32 = line.trim().parse().unwrap();
        nums.push(num);
    }

    let result = solve(&nums);
    writeln!(writer, "{}", result).unwrap();
    writer.flush().unwrap();
}

// solve находит сумму подмножества, ближайшую к 100
// При равном расстоянии выбирает большую сумму
fn solve(nums: &[i32]) -> i32 {
    const TARGET: i32 = 100;
    let mut best_sum = 0;
    let mut best_dist = TARGET; // расстояние от 0 до 100

    // Перебираем все 2^10 = 1024 подмножества
    for mask in 0..(1 << 10) {
        let mut sum = 0;
        for i in 0..10 {
            if mask & (1 << i) != 0 {
                sum += nums[i];
            }
        }

        let dist = (sum - TARGET).abs();

        // Выбираем лучшую сумму:
        // - меньшее расстояние до 100
        // - при равном расстоянии — большую сумму
        if dist < best_dist || (dist == best_dist && sum > best_sum) {
            best_sum = sum;
            best_dist = dist;
        }
    }

    best_sum
}





