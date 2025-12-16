use std::io::{self, BufRead, BufWriter, Write};

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let parts: Vec<&str> = line.trim().split_whitespace().collect();
    let r: i32 = parts[0].parse().unwrap();
    let b: i32 = parts[1].parse().unwrap();

    let (w, h) = solve(r, b);
    writeln!(writer, "{} {}", w, h).unwrap();
    writer.flush().unwrap();
}

// solve находит размеры панели W и H (W >= H) по количеству красных R и синих B плиток
// R = 2*W + 2*H - 4, B = (W-2) * (H-2)
fn solve(r: i32, b: i32) -> (i32, i32) {
    let sum = (r + 4) / 2;

    let mut d = 1;
    while d * d <= b {
        if b % d != 0 {
            d += 1;
            continue;
        }

        // Проверяем оба варианта: (W-2, H-2) = (d, B/d) и (B/d, d)
        let candidates = [b / d + 2, d + 2];
        for &w in &candidates {
            let h = sum - w;
            if (w - 2) * (h - 2) == b && w >= h {
                return (w, h);
            }
        }

        d += 1;
    }

    (0, 0)
}







