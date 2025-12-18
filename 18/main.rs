use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::collections::VecDeque;
use std::env;
use std::io::{self, BufRead, BufWriter, Write};
use std::time::Instant;

// get_memory_usage возвращает текущее использование памяти в байтах
fn get_memory_usage() -> Option<u64> {
    #[cfg(target_os = "linux")]
    {
        use std::fs;
        if let Ok(content) = fs::read_to_string("/proc/self/status") {
            for line in content.lines() {
                if line.starts_with("VmRSS:") {
                    if let Some(value) = line.split_whitespace().nth(1) {
                        if let Ok(kb) = value.parse::<u64>() {
                            return Some(kb * 1024);
                        }
                    }
                }
            }
        }
    }
    #[cfg(target_os = "macos")]
    {
        use std::process::Command;
        let pid = std::process::id().to_string();
        if let Ok(output) = Command::new("ps")
            .args(&["-o", "rss=", "-p", &pid])
            .output()
        {
            if let Ok(mem_str) = String::from_utf8(output.stdout) {
                if let Ok(kb) = mem_str.trim().parse::<u64>() {
                    return Some(kb * 1024);
                }
            }
        }
    }
    None
}

// check_limits проверяет ограничения времени и памяти
fn check_limits(max_time_ms: u64, max_memory_mb: u64, f: impl FnOnce()) {
    if env::var("CHECK_LIMITS").is_err() {
        f();
        return;
    }

    let mem_before = get_memory_usage();
    let start = Instant::now();
    f();
    let elapsed = start.elapsed();
    let mem_after = get_memory_usage();

    let elapsed_ms = elapsed.as_millis() as u64;
    let time_ok = elapsed_ms <= max_time_ms;

    let memory_ok = if let (Some(before), Some(after)) = (mem_before, mem_after) {
        let allocated = if after > before { after - before } else { 0 };
        let memory_mb = allocated as f64 / (1024.0 * 1024.0);
        let max_memory_bytes = max_memory_mb * 1024 * 1024;
        allocated <= max_memory_bytes
    } else {
        true
    };

    if !time_ok || !memory_ok {
        if elapsed_ms > max_time_ms {
            eprintln!("⚠️ Превышено время: {} мс (лимит: {} мс)", elapsed_ms, max_time_ms);
        }
        if let (Some(before), Some(after)) = (mem_before, mem_after) {
            let allocated = if after > before { after - before } else { 0 };
            let memory_mb = allocated as f64 / (1024.0 * 1024.0);
            if memory_mb > max_memory_mb as f64 {
                eprintln!("⚠️ Превышена память: {:.2} МБ (лимит: {} МБ)", memory_mb, max_memory_mb);
            }
        }
    } else {
        if let (Some(before), Some(after)) = (mem_before, mem_after) {
            let allocated = if after > before { after - before } else { 0 };
            let memory_mb = allocated as f64 / (1024.0 * 1024.0);
            eprintln!("✓ Время: {} мс, Память: {:.2} МБ", elapsed_ms, memory_mb);
        } else {
            eprintln!("✓ Время: {} мс, Память: не удалось измерить", elapsed_ms);
        }
    }
}

const INF: i64 = 1_000_000_000_000_000_000; // 1e18

#[derive(Clone)]
struct Edge {
    to: usize,
    capacity: i32,
    flow: i32,
    cost: i64,
    rev: usize,
}

struct Graph {
    adj: Vec<Vec<Edge>>,
}

impl Graph {
    fn new(n: usize) -> Self {
        Graph {
            adj: vec![Vec::new(); n],
        }
    }

    fn add_edge(&mut self, from: usize, to: usize, cap: i32, cost: i64) {
        let rev_from = self.adj[to].len();
        let rev_to = self.adj[from].len();
        self.adj[from].push(Edge {
            to,
            capacity: cap,
            flow: 0,
            cost,
            rev: rev_from,
        });
        self.adj[to].push(Edge {
            to: from,
            capacity: 0,
            flow: 0,
            cost: -cost,
            rev: rev_to,
        });
    }

    fn min_cost_max_flow(&mut self, s: usize, t: usize, k: i32) -> i64 {
        let n = self.adj.len();
        let mut potential = vec![0i64; n];
        let mut dist = vec![INF; n];
        let mut parent_edge = vec![0; n];
        let mut parent_node = vec![0; n];
        
        let mut total_flow = 0;
        let mut min_cost = 0;

        // SPFA for initial potentials (to handle negative costs)
        {
            let mut in_queue = vec![false; n];
            let mut queue = VecDeque::new();
            
            dist[s] = 0;
            queue.push_back(s);
            in_queue[s] = true;

            while let Some(u) = queue.pop_front() {
                in_queue[u] = false;
                for i in 0..self.adj[u].len() {
                    let e = &self.adj[u][i];
                    if e.capacity > e.flow && dist[e.to] > dist[u] + e.cost {
                        dist[e.to] = dist[u] + e.cost;
                        if !in_queue[e.to] {
                            queue.push_back(e.to);
                            in_queue[e.to] = true;
                        }
                    }
                }
            }

            if dist[t] == INF {
                return -1;
            }

            for i in 0..n {
                if dist[i] != INF {
                    potential[i] = dist[i];
                }
            }
        }

        while total_flow < k {
            dist.fill(INF);
            dist[s] = 0;

            let mut pq = BinaryHeap::new();
            pq.push(State { cost: 0, position: s });

            while let Some(State { cost: d, position: u }) = pq.pop() {
                if d > dist[u] {
                    continue;
                }

                for i in 0..self.adj[u].len() {
                    let e = &self.adj[u][i];
                    if e.capacity > e.flow {
                        let new_dist = dist[u] + e.cost + potential[u] - potential[e.to];
                        if dist[e.to] > new_dist {
                            dist[e.to] = new_dist;
                            parent_node[e.to] = u;
                            parent_edge[e.to] = i;
                            pq.push(State { cost: new_dist, position: e.to });
                        }
                    }
                }
            }

            if dist[t] == INF {
                return -1;
            }

            for i in 0..n {
                if dist[i] != INF {
                    potential[i] += dist[i];
                }
            }

            let mut push = k - total_flow;
            let mut curr = t;
            while curr != s {
                let p = parent_node[curr];
                let idx = parent_edge[curr];
                let available = self.adj[p][idx].capacity - self.adj[p][idx].flow;
                if available < push {
                    push = available;
                }
                curr = p;
            }

            total_flow += push;
            curr = t;
            while curr != s {
                let p = parent_node[curr];
                let idx = parent_edge[curr];
                self.adj[p][idx].flow += push;
                let rev_idx = self.adj[p][idx].rev;
                self.adj[curr][rev_idx].flow -= push;
                min_cost += (push as i64) * self.adj[p][idx].cost;
                curr = p;
            }
        }

        min_cost
    }
}

#[derive(Copy, Clone, Eq, PartialEq)]
struct State {
    cost: i64,
    position: usize,
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> Ordering {
        other.cost.cmp(&self.cost)
    }
}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

fn solve() {
    let stdin = io::stdin();
    let mut reader = stdin.lock();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    let mut line = String::new();
    reader.read_line(&mut line).unwrap();
    let t: usize = line.trim().parse().unwrap();

    for _ in 0..t {
        line.clear();
        reader.read_line(&mut line).unwrap();
        let parts: Vec<usize> = line.split_whitespace()
            .map(|s| s.parse().unwrap())
            .collect();
        let n = parts[0];
        let m = parts[1];

        line.clear();
        reader.read_line(&mut line).unwrap();
        let a: Vec<i64> = line.split_whitespace()
            .map(|s| s.parse().unwrap())
            .collect();

        line.clear();
        reader.read_line(&mut line).unwrap();
        let b: Vec<i64> = line.split_whitespace()
            .map(|s| s.parse().unwrap())
            .collect();

        let result = solve_test_case(n, m, &a, &b);
        writeln!(writer, "{}", result).unwrap();
    }
    writer.flush().unwrap();
}

fn solve_test_case(n: usize, m: usize, a: &[i64], b: &[i64]) -> i64 {
    let source = 0;
    let sink = n + m + 1;
    let num_nodes = n + m + 2;

    let mut g = Graph::new(num_nodes);

    // Edges from Source to A
    for i in 0..n {
        g.add_edge(source, i + 1, 1, 0);
    }

    // Edges from A to B
    for i in 0..n {
        for j in 0..m {
            let val = a[i];
            let div = b[j];
            let rem = val % div;
            let cost = if rem != 0 { div - rem } else { 0 };
            g.add_edge(i + 1, n + 1 + j, 1, cost);
        }
    }

    // Edges from B to Sink
    let q = (n / m) as i32;
    const M: i64 = 100_000_000_000_000; // 10^14

    for j in 0..m {
        if q > 0 {
            g.add_edge(n + 1 + j, sink, q, -M);
        }
        g.add_edge(n + 1 + j, sink, 1, 0);
    }

    let raw_cost = g.min_cost_max_flow(source, sink, n as i32);
    let real_cost = raw_cost + (q as i64) * (m as i64) * M;

    real_cost
}

fn main() {
    check_limits(2000, 1024, || {
        solve();
    });
}

