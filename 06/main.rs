use std::io::{self, BufRead, BufWriter, Write};

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем N и M
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let n: usize = parts[0].parse().unwrap();
    let _m: usize = parts[1].parse().unwrap(); // M - количество команд

    // Читаем координаты меток
    let mut markers_x = Vec::with_capacity(n);
    let mut markers_y = Vec::with_capacity(n);
    for _i in 0..n {
        line.clear();
        reader.read_line(&mut line).unwrap();
        let parts: Vec<&str> = line.trim().split_whitespace().collect();
        let x: i32 = parts[0].parse().unwrap();
        let y: i32 = parts[1].parse().unwrap();
        markers_x.push(x);
        markers_y.push(y);
    }

    // Читаем команды
    line.clear();
    reader.read_line(&mut line).unwrap();
    let commands = line.trim();

    let results = solve(&markers_x, &markers_y, commands);
    for result in results {
        writeln!(writer, "{}", result).unwrap();
    }
    writer.flush().unwrap();
}

// solve вычисляет сумму манхэттенских расстояний после каждой команды
fn solve(markers_x: &[i32], markers_y: &[i32], commands: &str) -> Vec<i64> {
    let n = markers_x.len();

    // Сортируем метки по x и y для быстрого вычисления суммы расстояний
    let mut sorted_x = markers_x.to_vec();
    sorted_x.sort();
    let mut sorted_y = markers_y.to_vec();
    sorted_y.sort();

    // Предвычисляем префиксные суммы для x и y
    let mut prefix_sum_x = vec![0i64; n + 1];
    let mut prefix_sum_y = vec![0i64; n + 1];
    for i in 0..n {
        prefix_sum_x[i + 1] = prefix_sum_x[i] + sorted_x[i] as i64;
        prefix_sum_y[i + 1] = prefix_sum_y[i] + sorted_y[i] as i64;
    }

    // Текущая позиция Кодеруна
    let mut cx = 0i32;
    let mut cy = 0i32;
    let mut results = Vec::with_capacity(commands.len());

    // Обрабатываем каждую команду
    for cmd in commands.chars() {
        // Обновляем позицию
        match cmd {
            'N' => cy += 1,
            'S' => cy -= 1,
            'E' => cx += 1,
            'W' => cx -= 1,
            _ => {}
        }

        // Вычисляем сумму манхэттенских расстояний
        // sum(|cx - mx| + |cy - my|) = sum(|cx - mx|) + sum(|cy - my|)

        // Для x координат
        let idx_x = sorted_x.partition_point(|&x| x <= cx);
        // sorted_x[0..idx_x] <= cx, sorted_x[idx_x..n] > cx
        let sum_x = (cx as i64) * (idx_x as i64) - prefix_sum_x[idx_x]
            + (prefix_sum_x[n] - prefix_sum_x[idx_x]) - (cx as i64) * ((n - idx_x) as i64);

        // Для y координат
        let idx_y = sorted_y.partition_point(|&y| y <= cy);
        // sorted_y[0..idx_y] <= cy, sorted_y[idx_y..n] > cy
        let sum_y = (cy as i64) * (idx_y as i64) - prefix_sum_y[idx_y]
            + (prefix_sum_y[n] - prefix_sum_y[idx_y]) - (cy as i64) * ((n - idx_y) as i64);

        let total = sum_x + sum_y;
        results.push(total);
    }

    results
}

