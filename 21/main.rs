use std::env;
use std::f64::consts::PI;
use std::io::{self, BufRead, BufReader, BufWriter, Write};
use std::time::{Duration, Instant};

const EPS: f64 = 1e-13;

#[derive(Clone, Copy, Debug)]
struct Point {
    x: f64,
    y: f64,
}

impl Point {
    fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }
}

fn dist_sq(p1: Point, p2: Point) -> f64 {
    let dx = p1.x - p2.x;
    let dy = p1.y - p2.y;
    dx * dx + dy * dy
}

struct Rng {
    state: u64,
}

impl Rng {
    fn new(seed: u64) -> Self {
        Self { state: seed }
    }

    fn next_u64(&mut self) -> u64 {
        let mut x = self.state;
        x ^= x << 13;
        x ^= x >> 7;
        x ^= x << 17;
        self.state = x;
        x
    }

    fn next_f64(&mut self) -> f64 {
        // Generates f64 in [0, 1)
        (self.next_u64() as f64) / (u64::MAX as f64)
    }
}

fn main() {
    let check_limits_env = env::var("CHECK_LIMITS").is_ok();
    
    let start_time = Instant::now();
    let max_time = Duration::from_secs(1000); // Local limit

    let stdin = io::stdin();
    let stdout = io::stdout();
    let mut reader = BufReader::new(stdin.lock());
    let mut writer = BufWriter::new(stdout.lock());

    let mut buffer = String::new();
    
    // Read T
    if reader.read_line(&mut buffer).is_err() { return; }
    let t_str = buffer.trim();
    if t_str.is_empty() { return; }
    let t: i32 = t_str.parse().unwrap_or(0);
    buffer.clear();

    let mut solver = Solver::new();

    for _ in 0..t {
        solver.solve_test_case(&mut reader, &mut writer, &mut buffer);
    }

    writer.flush().unwrap();

    if check_limits_env {
        let elapsed = start_time.elapsed();
        eprintln!("✓ Время: {:?}", elapsed);
    }
}

struct Solver {
    points: Vec<Point>,
    comp_points: Vec<Point>,
    static_cands: Vec<Point>,
    static_cands_indices: Vec<Vec<usize>>,
    solution: Vec<Point>,
    rng: Rng,
}

impl Solver {
    fn new() -> Self {
        Self {
            points: Vec::with_capacity(10),
            comp_points: Vec::with_capacity(10),
            static_cands: Vec::with_capacity(200),
            static_cands_indices: vec![Vec::with_capacity(50); 10],
            solution: Vec::with_capacity(10),
            rng: Rng::new(42),
        }
    }

    fn solve_test_case<R: BufRead, W: Write>(&mut self, reader: &mut R, writer: &mut W, buffer: &mut String) {
        // Read N
        loop {
            buffer.clear();
            if reader.read_line(buffer).unwrap() == 0 { return; }
            if !buffer.trim().is_empty() { break; }
        }
        let n: usize = buffer.trim().parse().unwrap_or(0);
        
        self.points.clear();
        for _ in 0..n {
            buffer.clear();
            reader.read_line(buffer).unwrap();
            let parts: Vec<&str> = buffer.trim().split_whitespace().collect();
            if parts.len() >= 2 {
                let x: f64 = parts[0].parse().unwrap();
                let y: f64 = parts[1].parse().unwrap();
                self.points.push(Point::new(x, y));
            }
        }

        let mut visited = 0;
        let mut total_circles = Vec::new();

        for i in 0..n {
            if (visited & (1 << i)) != 0 { continue; }

            self.comp_points.clear();
            let mut q = Vec::new();
            q.push(i);
            visited |= 1 << i;

            // BFS for component
            let mut head = 0;
            while head < q.len() {
                let u = q[head];
                head += 1;
                self.comp_points.push(self.points[u]);
                
                for v in 0..n {
                    if (visited & (1 << v)) == 0 {
                        if dist_sq(self.points[u], self.points[v]) <= 16.0 + 1e-7 {
                            visited |= 1 << v;
                            q.push(v);
                        }
                    }
                }
            }

            let res = self.solve_component();
            match res {
                Some(circles) => total_circles.extend(circles),
                None => {
                    writeln!(writer, "NO").unwrap();
                    return;
                }
            }
        }

        writeln!(writer, "YES").unwrap();
        writeln!(writer, "{}", total_circles.len()).unwrap();
        for c in total_circles {
            writeln!(writer, "{:.15} {:.15}", c.x, c.y).unwrap();
        }
    }

    fn solve_component(&mut self) -> Option<Vec<Point>> {
        let n_comp = self.comp_points.len();
        
        // 1. Try K=1 Exact (MEC)
        if let Some(center) = get_mec(&self.comp_points) {
            return Some(vec![center]);
        }

        // 2. Backtracking with Deterministic Candidates
        if let Some(res) = self.run_backtrack(n_comp, false) {
            return Some(res);
        }

        // 3. Backtracking with Random Candidates
        for _ in 0..3 {
            if let Some(res) = self.run_backtrack(n_comp, true) {
                return Some(res);
            }
        }

        None
    }

    fn run_backtrack(&mut self, n_comp: usize, use_random: bool) -> Option<Vec<Point>> {
        self.static_cands.clear();
        
        // Type 1: Points
        for i in 0..n_comp {
            self.static_cands.push(self.comp_points[i]);
        }

        // Type 2: Intersections
        for i in 0..n_comp {
            for j in i + 1..n_comp {
                let (p1, p2, cnt) = get_intersections(self.comp_points[i], 1.0, self.comp_points[j], 1.0);
                if cnt > 0 {
                    self.static_cands.push(p1);
                    if cnt > 1 {
                        self.static_cands.push(p2);
                    }
                }
            }
        }

        // Random candidates
        if use_random {
            for i in 0..n_comp {
                for _ in 0..5 {
                    let angle = self.rng.next_f64() * 2.0 * PI;
                    let r = self.rng.next_f64().sqrt() * 1.0;
                    let cx = self.comp_points[i].x + r * angle.cos();
                    let cy = self.comp_points[i].y + r * angle.sin();
                    self.static_cands.push(Point::new(cx, cy));
                }
            }
        }

        // Pre-calculate indices
        for i in 0..n_comp {
            self.static_cands_indices[i].clear();
        }
        
        for (idx, c) in self.static_cands.iter().enumerate() {
            for i in 0..n_comp {
                if dist_sq(*c, self.comp_points[i]) <= 1.0 + EPS {
                    self.static_cands_indices[i].push(idx);
                }
            }
        }

        self.solution.clear();
        if self.backtrack(0, n_comp) {
            return Some(self.solution.clone());
        }
        None
    }

    fn backtrack(&mut self, mask: u32, n_comp: usize) -> bool {
        if mask == (1 << n_comp) - 1 {
            return true;
        }

        let mut u = 0;
        for i in 0..n_comp {
            if (mask & (1 << i)) == 0 {
                u = i;
                break;
            }
        }

        // Try candidate closure
        let mut try_cand = |c: Point, solver: &mut Solver| -> bool {
            if dist_sq(c, solver.comp_points[u]) > 1.0 + EPS {
                return false;
            }
            for sol_c in &solver.solution {
                if dist_sq(c, *sol_c) < 4.0 - EPS {
                    return false;
                }
            }

            let mut new_mask = mask;
            for i in 0..n_comp {
                if dist_sq(c, solver.comp_points[i]) <= 1.0 + EPS {
                    new_mask |= 1 << i;
                }
            }

            solver.solution.push(c);
            if solver.backtrack(new_mask, n_comp) {
                return true;
            }
            solver.solution.pop();
            false
        };

        // 1. Static Candidates
        let indices = self.static_cands_indices[u].clone(); // Clone to avoid borrow conflict
        for &idx in &indices {
             if try_cand(self.static_cands[idx], self) {
                 return true;
             }
        }

        // 2. Dynamic Candidates
        if !self.solution.is_empty() {
             let sol_len = self.solution.len();
             // We need to collect dynamic candidates first to avoid borrowing conflicts
             // or restructure code. Collecting is safer/easier.
             let mut dyn_cands = Vec::with_capacity(50);

             // Type 3: Intersection of Boundary(P_i, 1) and Boundary(Sol_j, 2)
             for j in 0..sol_len {
                 let sol_c = self.solution[j];
                 if dist_sq(self.comp_points[u], sol_c) > 9.0 + 1e-5 {
                     continue;
                 }
                 
                 // Intersections involving u
                 let (p1, p2, cnt) = get_intersections(self.comp_points[u], 1.0, sol_c, 2.0);
                 if cnt > 0 {
                     dyn_cands.push(p1);
                     if cnt > 1 { dyn_cands.push(p2); }
                 }

                 // Other points close to u
                 for i in 0..n_comp {
                     if i == u { continue; }
                     if dist_sq(self.comp_points[u], self.comp_points[i]) > 4.0 + 1e-5 {
                         continue;
                     }
                     let (p1, p2, cnt) = get_intersections(self.comp_points[i], 1.0, sol_c, 2.0);
                     if cnt > 0 {
                         dyn_cands.push(p1);
                         if cnt > 1 { dyn_cands.push(p2); }
                     }
                 }
             }

             // Type 4: Intersection of two solution circles
             for i in 0..sol_len {
                 if dist_sq(self.comp_points[u], self.solution[i]) > 9.0 + 1e-5 { continue; }
                 for j in i+1..sol_len {
                     if dist_sq(self.comp_points[u], self.solution[j]) > 9.0 + 1e-5 { continue; }
                     let (p1, p2, cnt) = get_intersections(self.solution[i], 2.0, self.solution[j], 2.0);
                     if cnt > 0 {
                         dyn_cands.push(p1);
                         if cnt > 1 { dyn_cands.push(p2); }
                     }
                 }
             }

             for c in dyn_cands {
                 if try_cand(c, self) {
                     return true;
                 }
             }
        }

        false
    }
}

fn get_intersections(p1: Point, r1: f64, p2: Point, r2: f64) -> (Point, Point, usize) {
    let d2 = dist_sq(p1, p2);
    let d = d2.sqrt();
    if d > r1 + r2 + EPS || d < (r1 - r2).abs() - EPS || d < 1e-9 {
        return (Point::new(0.0, 0.0), Point::new(0.0, 0.0), 0);
    }
    
    let a = (r1 * r1 - r2 * r2 + d2) / (2.0 * d);
    let h = (r1 * r1 - a * a).max(0.0).sqrt();
    
    let x2 = p1.x + a * (p2.x - p1.x) / d;
    let y2 = p1.y + a * (p2.y - p1.y) / d;
    
    let rx = -h * (p2.y - p1.y) / d;
    let ry = h * (p2.x - p1.x) / d;

    (
        Point::new(x2 - rx, y2 - ry),
        Point::new(x2 + rx, y2 + ry),
        2
    )
}

fn get_mec(points: &[Point]) -> Option<Point> {
    let n = points.len();
    if n == 0 { return None; }
    if n == 1 { return Some(points[0]); }

    let mut best_r2 = 1.0 + EPS;
    let mut best_center: Option<Point> = None;

    let check = |c: Point, r2: f64, best_r2: &mut f64, best_center: &mut Option<Point>| {
        if r2 > *best_r2 { return; }
        for p in points {
            if dist_sq(c, *p) > r2 + EPS { return; }
        }
        *best_r2 = r2;
        *best_center = Some(c);
    };

    // Case 1: Pairwise
    for i in 0..n {
        for j in i + 1..n {
            let mid = Point::new((points[i].x + points[j].x) / 2.0, (points[i].y + points[j].y) / 2.0);
            let r2 = dist_sq(mid, points[i]);
            check(mid, r2, &mut best_r2, &mut best_center);
        }
    }

    // Case 2: Triplet
    for i in 0..n {
        for j in i + 1..n {
            for k in j + 1..n {
                if let Some(c) = get_circumcenter(points[i], points[j], points[k]) {
                    let r2 = dist_sq(c, points[i]);
                    check(c, r2, &mut best_r2, &mut best_center);
                }
            }
        }
    }

    best_center
}

fn get_circumcenter(a: Point, b: Point, c: Point) -> Option<Point> {
    let d = 2.0 * (a.x * (b.y - c.y) + b.x * (c.y - a.y) + c.x * (a.y - b.y));
    if d.abs() < 1e-9 { return None; }
    
    let sa = a.x * a.x + a.y * a.y;
    let sb = b.x * b.x + b.y * b.y;
    let sc = c.x * c.x + c.y * c.y;

    let ux = (sa * (b.y - c.y) + sb * (c.y - a.y) + sc * (a.y - b.y)) / d;
    let uy = (sa * (c.x - b.x) + sb * (a.x - c.x) + sc * (b.x - a.x)) / d;
    
    Some(Point::new(ux, uy))
}

