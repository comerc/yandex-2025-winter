const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

let inputLines = [];
rl.on('line', (input) => {
  inputLines.push(input);
});
rl.on('close', () => {
  solve();
});

function solve() {
  let line = 0;
  const T = parseInt(inputLines[line++]);
  let results = [];
  for (let test = 0; test < T; test++) {
    const [n, m] = inputLines[line++].split(' ').map(Number);
    const A = inputLines[line++].split(' ').map(Number);
    const B = inputLines[line++].split(' ').map(Number);
    const answer = minOperationsToMakeSuccessful(A, B, n, m);
    console.log(answer);
  }
}

class PriorityQueue {
  constructor(comparator) {
    this.heap = [];
    this.comparator = comparator || ((a, b) => a < b);
  }
  push(item) {
    this.heap.push(item);
    this._siftUp(this.heap.length - 1);
  }
  pop() {
    if (this.heap.length === 0) return null;
    let top = this.heap[0];
    let bottom = this.heap.pop();
    if (this.heap.length > 0) {
      this.heap[0] = bottom;
      this._siftDown(0);
    }
    return top;
  }
  isEmpty() {
    return this.heap.length === 0;
  }
  _siftUp(index) {
    let item = this.heap[index];
    while (index > 0) {
      let parent = (index - 1) >> 1;
      if (this.comparator(item, this.heap[parent])) {
        this.heap[index] = this.heap[parent];
        index = parent;
      } else {
        break;
      }
    }
    this.heap[index] = item;
  }
  _siftDown(index) {
    let item = this.heap[index];
    let n = this.heap.length;
    while (index * 2 + 1 < n) {
      let child = index * 2 + 1;
      if (
        child + 1 < n &&
        this.comparator(this.heap[child + 1], this.heap[child])
      ) {
        child++;
      }
      if (this.comparator(this.heap[child], item)) {
        this.heap[index] = this.heap[child];
        index = child;
      } else {
        break;
      }
    }
    this.heap[index] = item;
  }
}

function addEdge(graph, from, to, capacity, cost) {
  graph[from].push({ to, capacity, cost, reverse: graph[to].length });
  graph[to].push({
    to: from,
    capacity: 0,
    cost: -cost,
    reverse: graph[from].length - 1,
  });
}

function minCostFlow(N, graph, s, t, flowLimit) {
  let phi = new Array(N).fill(0);
  let prevv = new Array(N);
  let preve = new Array(N);
  let totalCost = 0;
  let flow = 0;
  while (flow < flowLimit) {
    let dist = new Array(N).fill(Infinity);
    dist[s] = 0;
    let queue = new PriorityQueue((a, b) => a[0] < b[0]);
    queue.push([0, s]);
    while (!queue.isEmpty()) {
      let [d, v] = queue.pop();
      if (dist[v] < d) continue;
      for (let i = 0; i < graph[v].length; i++) {
        let e = graph[v][i];
        if (e.capacity > 0) {
          let newDist = dist[v] + e.cost + phi[v] - phi[e.to];
          if (dist[e.to] > newDist) {
            dist[e.to] = newDist;
            prevv[e.to] = v;
            preve[e.to] = i;
            queue.push([newDist, e.to]);
          }
        }
      }
    }
    if (dist[t] === Infinity) {
      break;
    }
    for (let i = 0; i < N; i++) {
      if (dist[i] < Infinity) {
        phi[i] += dist[i];
      }
    }
    let d = flowLimit - flow;
    for (let v = t; v !== s; v = prevv[v]) {
      d = Math.min(d, graph[prevv[v]][preve[v]].capacity);
    }
    flow += d;
    totalCost += d * (phi[t] - phi[s]);
    for (let v = t; v !== s; v = prevv[v]) {
      let e = graph[prevv[v]][preve[v]];
      e.capacity -= d;
      graph[v][e.reverse].capacity += d;
    }
  }
  if (flow < flowLimit) {
    // В данной задаче всегда должен быть поток, но на всякий случай
    return Infinity;
  }
  return totalCost;
}

function minOperationsToMakeSuccessful(A, B, n, m) {
  const q = Math.floor(n / m);
  const r = n % m;
  const N = n + m + 3;
  const S = 0;
  const T = n + m + 2;
  const U = n + m + 1;
  let graph = new Array(N);
  for (let i = 0; i < N; i++) graph[i] = [];
  // S -> элементы
  for (let j = 1; j <= n; j++) {
    addEdge(graph, S, j, 1, 0);
  }
  // элементы -> типы
  for (let j = 1; j <= n; j++) {
    let a = A[j - 1];
    for (let i = 1; i <= m; i++) {
      let b = B[i - 1];
      let cost = (b - (a % b)) % b;
      addEdge(graph, j, n + i, 1, cost);
    }
  }
  // типы -> T и типы -> U
  for (let i = 1; i <= m; i++) {
    addEdge(graph, n + i, T, q, 0);
    addEdge(graph, n + i, U, 1, 0);
  }
  // U -> T
  addEdge(graph, U, T, r, 0);
  let cost = minCostFlow(N, graph, S, T, n);
  return cost;
}
