package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"runtime"
	"time"
)

// checkLimits проверяет ограничения времени и памяти
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

// Edge represents a directed edge in the graph
type Edge struct {
	to       int
	capacity int
	flow     int
	cost     int64
	rev      int // index of the reverse edge in graph[to]
}

// Graph represents the flow network
type Graph struct {
	adj [][]Edge
}

func NewGraph(n int) *Graph {
	return &Graph{
		adj: make([][]Edge, n),
	}
}

func (g *Graph) AddEdge(from, to, cap int, cost int64) {
	forward := Edge{to: to, capacity: cap, flow: 0, cost: cost, rev: len(g.adj[to])}
	backward := Edge{to: from, capacity: 0, flow: 0, cost: -cost, rev: len(g.adj[from])}
	g.adj[from] = append(g.adj[from], forward)
	g.adj[to] = append(g.adj[to], backward)
}

const INF = 1e18

// MinCostMaxFlow finds the minimum cost to send `k` units of flow from s to t
// Returns -1 if it's impossible to send `k` units
func (g *Graph) MinCostMaxFlow(s, t int, k int) int64 {
	n := len(g.adj)
	potential := make([]int64, n)

	totalFlow := 0
	minCost := int64(0)

	// Initial potentials using SPFA to handle negative costs
	dist := make([]int64, n)
	parentEdge := make([]int, n)
	parentNode := make([]int, n)

	inQueue := make([]bool, n)
	queue := make([]int, 0, n)

	// SPFA initialization
	for i := 0; i < n; i++ {
		dist[i] = INF
	}
	dist[s] = 0
	queue = append(queue, s)
	inQueue[s] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false

		for i := range g.adj[u] {
			e := &g.adj[u][i]
			if e.capacity > e.flow {
				if dist[e.to] > dist[u]+e.cost {
					dist[e.to] = dist[u] + e.cost
					if !inQueue[e.to] {
						queue = append(queue, e.to)
						inQueue[e.to] = true
					}
				}
			}
		}
	}

	// If sink is unreachable even initially (should not happen in this problem)
	if dist[t] == INF {
		return -1
	}

	// Update potentials
	for i := 0; i < n; i++ {
		if dist[i] != INF {
			potential[i] = dist[i]
		}
	}

	// Main loop with Dijkstra
	for totalFlow < k {
		// Dijkstra
		for i := 0; i < n; i++ {
			dist[i] = INF
		}
		dist[s] = 0

		pq := &PriorityQueue{}
		heap.Init(pq)
		heap.Push(pq, &Item{value: s, priority: 0})

		for pq.Len() > 0 {
			item := heap.Pop(pq).(*Item)
			u := item.value
			d := item.priority

			if d > dist[u] {
				continue
			}

			for i := range g.adj[u] {
				e := &g.adj[u][i]
				if e.capacity-e.flow > 0 {
					newDist := dist[u] + e.cost + potential[u] - potential[e.to]
					if dist[e.to] > newDist {
						dist[e.to] = newDist
						parentNode[e.to] = u
						parentEdge[e.to] = i
						heap.Push(pq, &Item{value: e.to, priority: newDist})
					}
				}
			}
		}

		if dist[t] == INF {
			return -1 // Cannot push more flow
		}

		// Update potentials
		for i := 0; i < n; i++ {
			if dist[i] != INF {
				potential[i] += dist[i]
			}
		}

		// Push flow
		push := k - totalFlow
		curr := t
		for curr != s {
			p := parentNode[curr]
			idx := parentEdge[curr]
			available := g.adj[p][idx].capacity - g.adj[p][idx].flow
			if available < push {
				push = available
			}
			curr = p
		}

		totalFlow += push
		curr = t
		for curr != s {
			p := parentNode[curr]
			idx := parentEdge[curr]
			g.adj[p][idx].flow += push
			revIdx := g.adj[p][idx].rev
			g.adj[curr][revIdx].flow -= push
			minCost += int64(push) * g.adj[p][idx].cost
			curr = p
		}
	}

	return minCost
}

// PriorityQueue implementation
type Item struct {
	value    int
	priority int64
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func solve() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var n, m int
		fmt.Fscan(reader, &n, &m)

		a := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[j])
		}

		b := make([]int64, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &b[j])
		}

		result := solveTestCase(n, m, a, b)
		fmt.Fprintln(writer, result)
	}
}

func solveTestCase(n, m int, a []int64, b []int64) int64 {
	// Source = 0, Sink = n + m + 1
	// Nodes 1..n: A
	// Nodes n+1..n+m: B
	source := 0
	sink := n + m + 1
	numNodes := n + m + 2

	g := NewGraph(numNodes)

	// Edges from Source to A
	for i := 0; i < n; i++ {
		g.AddEdge(source, i+1, 1, 0)
	}

	// Edges from A to B
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := a[i]
			div := b[j]
			rem := val % div
			cost := int64(0)
			if rem != 0 {
				cost = div - rem
			}
			g.AddEdge(i+1, n+1+j, 1, cost)
		}
	}

	// Edges from B to Sink
	// Group sizes: n = q*m + r
	q := n / m

	// Large constant to force filling the first q slots
	const M = 100000000000000 // 10^14

	for j := 0; j < m; j++ {
		if q > 0 {
			// Mandatory q slots with high priority (negative cost)
			g.AddEdge(n+1+j, sink, q, -M)
		}
		// Extra capacity for remainder
		g.AddEdge(n+1+j, sink, 1, 0)
	}

	// Calculate Min Cost for flow = n
	rawCost := g.MinCostMaxFlow(source, sink, n)

	// Adjust cost by removing the artificial negative costs
	realCost := rawCost + int64(q)*int64(m)*M

	return realCost
}

func main() {
	checkLimits(2*time.Second, 1024, solve)
}
