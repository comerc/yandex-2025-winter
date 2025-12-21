package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func checkLimits(maxTime time.Duration, maxMemoryMB int, fn func()) {
	if os.Getenv("CHECK_LIMITS") == "" {
		fn()
		return
	}
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	allocated := m2.TotalAlloc - m1.TotalAlloc
	memoryMB := float64(allocated) / (1024 * 1024)

	if elapsed > maxTime {
		fmt.Fprintf(os.Stderr, "⚠️ Превышено время: %v (лимит: %v)\n", elapsed, maxTime)
	}
	if memoryMB > float64(maxMemoryMB) {
		fmt.Fprintf(os.Stderr, "⚠️ Превышена память: %.2f МБ (лимит: %d МБ)\n", memoryMB, maxMemoryMB)
	} else {
		fmt.Fprintf(os.Stderr, "✓ Время: %v, Память: %.2f МБ\n", elapsed, memoryMB)
	}
}

type Point struct {
	x, y float64
}

var (
	points             []Point
	n                  int
	reader             *bufio.Reader
	writer             *bufio.Writer
	compPoints         []Point
	staticCands        []Point
	staticCandsIndices [][]int // indices of staticCands relevant for each point u
	solution           []Point
	eps                = 1e-13
	rnd                *rand.Rand
)

func solve() {
	// Initialize deterministic RNG
	rnd = rand.New(rand.NewSource(42))

	points = make([]Point, 10)
	compPoints = make([]Point, 10)
	staticCands = make([]Point, 0, 200)
	staticCandsIndices = make([][]int, 10)
	for i := range staticCandsIndices {
		staticCandsIndices[i] = make([]int, 0, 50)
	}
	solution = make([]Point, 0, 10)

	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	t, _ := strconv.Atoi(line)

	for i := 0; i < t; i++ {
		solveTestCase()
	}
}

func solveTestCase() {
	line, err := reader.ReadString('\n')
	for err == nil && strings.TrimSpace(line) == "" {
		line, err = reader.ReadString('\n')
	}
	if err != nil {
		return
	}
	nStr := strings.TrimSpace(line)
	n, _ = strconv.Atoi(nStr)

	// Reuse/resize points slice
	if cap(points) < n {
		points = make([]Point, n)
	}
	points = points[:n]

	for i := 0; i < n; i++ {
		line, _ = reader.ReadString('\n')
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			points[i].x, _ = strconv.ParseFloat(parts[0], 64)
			points[i].y, _ = strconv.ParseFloat(parts[1], 64)
		}
	}

	visited := 0
	totalCircles := make([]Point, 0, n)

	for i := 0; i < n; i++ {
		if (visited & (1 << i)) != 0 {
			continue
		}

		compSize := 0
		q := make([]int, 0, n)
		q = append(q, i)
		visited |= (1 << i)

		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			compPoints[compSize] = points[u]
			compSize++

			for v := 0; v < n; v++ {
				if (visited & (1 << v)) == 0 {
					if distSq(points[u], points[v]) <= 16.0+1e-7 {
						visited |= (1 << v)
						q = append(q, v)
					}
				}
			}
		}

		res := solveComponent(compSize)
		if res == nil {
			writer.WriteString("NO\n")
			return
		}
		totalCircles = append(totalCircles, res...)
	}

	writer.WriteString("YES\n")
	writer.WriteString(fmt.Sprintf("%d\n", len(totalCircles)))
	for _, c := range totalCircles {
		writer.WriteString(fmt.Sprintf("%.15f %.15f\n", c.x, c.y))
	}
}

func solveComponent(nComp int) []Point {
	// 1. Try K=1 Exact (MEC)
	if center := getMEC(nComp); center != nil {
		return []Point{*center}
	}

	// 2. Backtracking with Deterministic Candidates
	if res := runBacktrack(nComp, false); res != nil {
		return res
	}

	// 3. Backtracking with Random Candidates (Fallback)
	// Try a few batches of random candidates
	for attempt := 0; attempt < 3; attempt++ {
		if res := runBacktrack(nComp, true); res != nil {
			return res
		}
	}

	return nil
}

func runBacktrack(nComp int, useRandom bool) []Point {
	staticCands = staticCands[:0]
	// Type 1: Points
	for i := 0; i < nComp; i++ {
		staticCands = append(staticCands, compPoints[i])
	}
	// Type 2: Intersections
	for i := 0; i < nComp; i++ {
		for j := i + 1; j < nComp; j++ {
			p1, p2, cnt := getIntersections(compPoints[i], 1.0, compPoints[j], 1.0)
			if cnt > 0 {
				staticCands = append(staticCands, p1)
				if cnt > 1 {
					staticCands = append(staticCands, p2)
				}
			}
		}
	}

	// Random candidates
	if useRandom {
		for i := 0; i < nComp; i++ {
			for k := 0; k < 5; k++ { // 5 random points per input point
				angle := rnd.Float64() * 2 * math.Pi
				r := math.Sqrt(rnd.Float64()) * 1.0 // uniform in disk
				cx := compPoints[i].x + r*math.Cos(angle)
				cy := compPoints[i].y + r*math.Sin(angle)
				staticCands = append(staticCands, Point{cx, cy})
			}
		}
	}

	// Pre-calculate indices for each point
	for i := 0; i < nComp; i++ {
		staticCandsIndices[i] = staticCandsIndices[i][:0]
		for idx, c := range staticCands {
			if distSq(c, compPoints[i]) <= 1.0+eps {
				staticCandsIndices[i] = append(staticCandsIndices[i], idx)
			}
		}
	}

	solution = solution[:0]
	if backtrack(0, nComp) {
		res := make([]Point, len(solution))
		copy(res, solution)
		return res
	}
	return nil
}

func backtrack(mask int, nComp int) bool {
	if mask == (1<<nComp)-1 {
		return true
	}

	// Find the first uncovered point
	u := -1
	for i := 0; i < nComp; i++ {
		if (mask & (1 << i)) == 0 {
			u = i
			break
		}
	}

	try := func(c Point) bool {
		// Quick check (redundant if pre-filtered, but needed for dynamic)
		if distSq(c, compPoints[u]) > 1.0+eps {
			return false
		}
		// Validate against existing solution
		for _, solC := range solution {
			if distSq(c, solC) < 4.0-eps {
				return false
			}
		}

		newMask := mask
		for i := 0; i < nComp; i++ {
			if distSq(c, compPoints[i]) <= 1.0+eps {
				newMask |= (1 << i)
			}
		}
		solution = append(solution, c)
		if backtrack(newMask, nComp) {
			return true
		}
		solution = solution[:len(solution)-1]
		return false
	}

	// 1. Static Candidates (Filtered)
	for _, idx := range staticCandsIndices[u] {
		if try(staticCands[idx]) {
			return true
		}
	}

	// 2. Dynamic Candidates
	if len(solution) > 0 {
		// Type 3: Intersection of Boundary(P_i, 1) and Boundary(Sol_j, 2)
		// Optimization: Only consider P_i if P_i is relevant?
		// Actually, we need to cover 'u'.
		// A circle 'c' covering 'u' must be in Disk(u, 1).
		// So 'c' is intersection of Disk(u, 1) and Boundary(Sol_j, 2).
		// Or intersection of Disk(P_i, 1) and Boundary(Sol_j, 2) THAT ALSO COVERS u.

		// Priority 1: Intersections involving Boundary(u, 1)
		for _, solC := range solution {
			// Pruning: if solC is too far from u, they can't touch and cover u
			// Max dist(u, c) = 1. Max dist(c, solC) = 2.
			// So if dist(u, solC) > 3, impossible.
			d2 := distSq(compPoints[u], solC)
			if d2 > 9.0+1e-5 {
				continue
			}

			// Intersection of Boundary(u, 1) and Boundary(SolC, 2)
			p1, p2, cnt := getIntersections(compPoints[u], 1.0, solC, 2.0)
			if cnt > 0 {
				if try(p1) {
					return true
				}
				if cnt > 1 {
					if try(p2) {
						return true
					}
				}
			}

			// Intersection of Boundary(P_i, 1) and Boundary(SolC, 2)
			// Iterate OTHER points only if they are close to u?
			// This might be overkill if Type 3a (u based) handles most.
			// But strictly, we should try all P_i.
			for i := 0; i < nComp; i++ {
				if i == u {
					continue
				}
				// Pruning: we need resulting 'c' to cover u.
				// c is on Boundary(P_i, 1).
				// So dist(c, u) <= 1.
				// Also dist(c, P_i) = 1.
				// Triangle ineq: dist(u, P_i) <= dist(u, c) + dist(c, P_i) <= 1 + 1 = 2.
				// If dist(u, P_i) > 2, then 'c' on Boundary(P_i, 1) cannot cover 'u'.
				if distSq(compPoints[u], compPoints[i]) > 4.0+1e-5 {
					continue
				}

				p1, p2, cnt := getIntersections(compPoints[i], 1.0, solC, 2.0)
				if cnt > 0 {
					if try(p1) {
						return true
					}
					if cnt > 1 {
						if try(p2) {
							return true
						}
					}
				}
			}
		}

		// Type 4: Intersection of Boundary(Sol_i, 2) and Boundary(Sol_j, 2)
		for i := 0; i < len(solution); i++ {
			// Pruning: check if Sol_i is close to u
			// dist(u, Sol_i) <= 3
			if distSq(compPoints[u], solution[i]) > 9.0+1e-5 {
				continue
			}
			for j := i + 1; j < len(solution); j++ {
				if distSq(compPoints[u], solution[j]) > 9.0+1e-5 {
					continue
				}
				p1, p2, cnt := getIntersections(solution[i], 2.0, solution[j], 2.0)
				if cnt > 0 {
					if try(p1) {
						return true
					}
					if cnt > 1 {
						if try(p2) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func getIntersections(p1 Point, r1 float64, p2 Point, r2 float64) (Point, Point, int) {
	d2 := distSq(p1, p2)
	d := math.Sqrt(d2)
	if d > r1+r2+eps || d < math.Abs(r1-r2)-eps || d < 1e-9 {
		return Point{}, Point{}, 0
	}
	a := (r1*r1 - r2*r2 + d2) / (2 * d)
	h := math.Sqrt(math.Max(0, r1*r1-a*a))
	x2 := p1.x + a*(p2.x-p1.x)/d
	y2 := p1.y + a*(p2.y-p1.y)/d
	return Point{x2 + h*(p2.y-p1.y)/d, y2 - h*(p2.x-p1.x)/d},
		Point{x2 - h*(p2.y-p1.y)/d, y2 + h*(p2.x-p1.x)/d}, 2
}

func getMEC(nComp int) *Point {
	if nComp == 0 {
		return nil
	}
	if nComp == 1 {
		return &compPoints[0]
	}
	bestR2 := 1.0 + eps
	var bestCenter *Point
	check := func(c Point, r2 float64) {
		if r2 > bestR2 {
			return
		}
		for i := 0; i < nComp; i++ {
			if distSq(c, compPoints[i]) > r2+eps {
				return
			}
		}
		bestR2 = r2
		bestCenter = &Point{c.x, c.y}
	}
	for i := 0; i < nComp; i++ {
		for j := i + 1; j < nComp; j++ {
			mid := Point{(compPoints[i].x + compPoints[j].x) / 2, (compPoints[i].y + compPoints[j].y) / 2}
			r2 := distSq(mid, compPoints[i])
			check(mid, r2)
		}
	}
	for i := 0; i < nComp; i++ {
		for j := i + 1; j < nComp; j++ {
			for k := j + 1; k < nComp; k++ {
				c, ok := getCircumcenter(compPoints[i], compPoints[j], compPoints[k])
				if ok {
					r2 := distSq(c, compPoints[i])
					check(c, r2)
				}
			}
		}
	}
	return bestCenter
}

func getCircumcenter(a, b, c Point) (Point, bool) {
	d := 2 * (a.x*(b.y-c.y) + b.x*(c.y-a.y) + c.x*(a.y-b.y))
	if math.Abs(d) < 1e-9 {
		return Point{}, false
	}
	ux := ((a.x*a.x+a.y*a.y)*(b.y-c.y) + (b.x*b.x+b.y*b.y)*(c.y-a.y) + (c.x*c.x+c.y*c.y)*(a.y-b.y)) / d
	uy := ((a.x*a.x+a.y*a.y)*(c.x-b.x) + (b.x*b.x+b.y*b.y)*(a.x-c.x) + (c.x*c.x+c.y*c.y)*(b.x-a.x)) / d
	return Point{ux, uy}, true
}

func distSq(p1, p2 Point) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return dx*dx + dy*dy
}

func main() {
	reader = bufio.NewReaderSize(os.Stdin, 4<<20)
	writer = bufio.NewWriterSize(os.Stdout, 4<<20)
	defer writer.Flush()
	checkLimits(1000*time.Second, 256, solve)
}
