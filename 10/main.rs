use std::io::{self, BufRead, BufWriter, Write, Read};

struct Edge {
    a: usize,
    b: usize,
    w: i32,
}

fn main() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    // Читаем количество тестов
    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let t: usize = line.trim().parse().unwrap();

    // Используем scanner для чтения чисел
    let mut tokens = Vec::new();
    let mut buffer = String::new();
    reader.read_to_string(&mut buffer).unwrap();
    
    for token in buffer.split_whitespace() {
        tokens.push(token.parse::<i32>().unwrap());
    }
    
    let mut token_idx = 0;
    let mut read_int = || {
        let val = tokens[token_idx];
        token_idx += 1;
        val
    };

    for _test in 0..t {
        // Читаем n и m
        let n = read_int() as usize;
        let m = read_int() as usize;

        // Читаем массивы a, b, c
        let mut a = vec![0; m];
        let mut b = vec![0; m];
        let mut c = vec![0; m];

        for i in 0..m {
            a[i] = read_int();
        }
        for i in 0..m {
            b[i] = read_int();
        }
        for i in 0..m {
            c[i] = read_int();
        }

        let mut edges = Vec::with_capacity(m);
        for i in 0..m {
            edges.push(Edge {
                a: (a[i] - 1) as usize,
                b: (b[i] - 1) as usize,
                w: c[i],
            });
        }

        let result = solve(n, &mut edges);

        // Выводим результат
        for (i, &v) in result.iter().enumerate() {
            if i > 0 {
                write!(writer, " ").unwrap();
            }
            write!(writer, "{}", v).unwrap();
        }
        writeln!(writer).unwrap();
    }
}

// solve находит минимальный вес w для каждого размера компоненты k
fn solve(n: usize, edges: &mut [Edge]) -> Vec<i32> {
    // Сортируем рёбра по весу
    edges.sort_by_key(|e| e.w);

    // Инициализируем DSU
    let mut parent: Vec<usize> = (0..n).collect();
    let mut rank = vec![0; n];
    let mut size = vec![1; n];

    // result[k] = минимальный вес для размера k+1
    let mut result = vec![-1; n];

    // k=1: любая вершина достижима сама из себя с w=0
    result[0] = 0;

    // maxReached отслеживает максимальный размер, для которого уже найден ответ
    let mut max_reached = 1;

    // Обрабатываем рёбра в порядке возрастания веса
    for e in edges.iter() {
        // Пропускаем петли (они не меняют связность)
        if e.a == e.b {
            continue;
        }

        let root_a = find(&mut parent, e.a);
        let root_b = find(&mut parent, e.b);

        if root_a != root_b {
            // Объединяем компоненты
            let new_size = size[root_a] + size[root_b];
            union(&mut parent, &mut rank, &mut size, root_a, root_b);

            // Обновляем результат для всех новых размеров
            for k in (max_reached + 1)..=new_size {
                result[k - 1] = e.w;
            }
            if new_size > max_reached {
                max_reached = new_size;
            }
        }

        // Если достигли максимального размера, можно остановиться
        if max_reached == n {
            break;
        }
    }

    result
}

fn find(parent: &mut [usize], x: usize) -> usize {
    if parent[x] != x {
        parent[x] = find(parent, parent[x]);
    }
    parent[x]
}

fn union(parent: &mut [usize], rank: &mut [usize], size: &mut [usize], x: usize, y: usize) {
    if rank[x] < rank[y] {
        parent[x] = y;
        size[y] += size[x];
    } else if rank[x] > rank[y] {
        parent[y] = x;
        size[x] += size[y];
    } else {
        parent[y] = x;
        size[x] += size[y];
        rank[x] += 1;
    }
}

