use std::io::{self, BufRead, BufWriter, Write};

#[derive(Clone, Copy)]
struct Point {
    x: i32,
    y: i32,
    z: i32,
    idx: usize,
}

struct Edge {
    from: usize,
    to: usize,
    weight: i32,
}

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем N
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let n: usize = line.trim().parse().unwrap();

    // Читаем точки
    let mut points = Vec::with_capacity(n);
    for i in 0..n {
        line.clear();
        reader.read_line(&mut line).unwrap();
        let parts: Vec<&str> = line.trim().split_whitespace().collect();
        let x: i32 = parts[0].parse().unwrap();
        let y: i32 = parts[1].parse().unwrap();
        let z: i32 = parts[2].parse().unwrap();
        points.push(Point { x, y, z, idx: i });
    }

    let result = solve_mst(&mut points);
    writeln!(writer, "{}", result).unwrap();
    writer.flush().unwrap();
}

fn find(parent: &mut [usize], x: usize) -> usize {
    if parent[x] != x {
        parent[x] = find(parent, parent[x]);
    }
    parent[x]
}

fn union(parent: &mut [usize], rank: &mut [usize], x: usize, y: usize) {
    if rank[x] < rank[y] {
        parent[x] = y;
    } else if rank[x] > rank[y] {
        parent[y] = x;
    } else {
        parent[y] = x;
        rank[x] += 1;
    }
}

// solve_mst находит минимальное остовное дерево для заданных точек
fn solve_mst(points: &mut [Point]) -> i64 {
    let n = points.len();

    // Строим рёбра: для каждой координаты сортируем точки и добавляем рёбра между соседями
    let mut edges = Vec::with_capacity(3 * n);

    // Сортируем по x и добавляем рёбра между соседями
    let mut sorted_by_x = points.to_vec();
    sorted_by_x.sort_by_key(|p| p.x);
    for i in 0..n - 1 {
        let cost = (sorted_by_x[i].x - sorted_by_x[i + 1].x)
            .abs()
            .min((sorted_by_x[i].y - sorted_by_x[i + 1].y).abs())
            .min((sorted_by_x[i].z - sorted_by_x[i + 1].z).abs());
        edges.push(Edge {
            from: sorted_by_x[i].idx,
            to: sorted_by_x[i + 1].idx,
            weight: cost,
        });
    }

    // Сортируем по y и добавляем рёбра между соседями
    let mut sorted_by_y = points.to_vec();
    sorted_by_y.sort_by_key(|p| p.y);
    for i in 0..n - 1 {
        let cost = (sorted_by_y[i].x - sorted_by_y[i + 1].x)
            .abs()
            .min((sorted_by_y[i].y - sorted_by_y[i + 1].y).abs())
            .min((sorted_by_y[i].z - sorted_by_y[i + 1].z).abs());
        edges.push(Edge {
            from: sorted_by_y[i].idx,
            to: sorted_by_y[i + 1].idx,
            weight: cost,
        });
    }

    // Сортируем по z и добавляем рёбра между соседями
    let mut sorted_by_z = points.to_vec();
    sorted_by_z.sort_by_key(|p| p.z);
    for i in 0..n - 1 {
        let cost = (sorted_by_z[i].x - sorted_by_z[i + 1].x)
            .abs()
            .min((sorted_by_z[i].y - sorted_by_z[i + 1].y).abs())
            .min((sorted_by_z[i].z - sorted_by_z[i + 1].z).abs());
        edges.push(Edge {
            from: sorted_by_z[i].idx,
            to: sorted_by_z[i + 1].idx,
            weight: cost,
        });
    }

    // Сортируем рёбра по весу
    edges.sort_by_key(|e| e.weight);

    // Алгоритм Крускала с DSU
    let mut parent: Vec<usize> = (0..n).collect();
    let mut rank = vec![0; n];

    let mut total_cost = 0i64;
    let mut edges_used = 0;

    for edge in edges {
        if edges_used == n - 1 {
            break;
        }
        let from_root = find(&mut parent, edge.from);
        let to_root = find(&mut parent, edge.to);
        if from_root != to_root {
            union(&mut parent, &mut rank, from_root, to_root);
            total_cost += edge.weight as i64;
            edges_used += 1;
        }
    }

    total_cost
}





