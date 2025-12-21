import 'dart:io';
import 'dart:typed_data';
import 'dart:collection';

const int INF = 1000000000; // 1e9

class Edge {
  int to;
  int capacity;
  int flow;
  int cost;
  int rev;
  Edge(this.to, this.capacity, this.flow, this.cost, this.rev);
}

class Graph {
  List<List<Edge>> adj;

  Graph(int n) : adj = List.generate(n, (_) => []);

  void addEdge(int from, int to, int cap, int cost) {
    adj[from].add(Edge(to, cap, 0, cost, adj[to].length));
    adj[to].add(Edge(from, 0, 0, -cost, adj[from].length - 1));
  }

  int minCostMaxFlow(int s, int t, int k) {
    int n = adj.length;
    List<int> potential = List.filled(n, 0);

    int totalFlow = 0;
    int minCost = 0;


    // Main loop with Dijkstra
    while (totalFlow < k) {
      // Dijkstra
      List<int> dist = List.filled(n, INF);
      List<int> parentEdge = List.filled(n, -1);
      List<int> parentNode = List.filled(n, -1);

      dist[s] = 0;

      PriorityQueue pq = PriorityQueue();
      pq.add(Item(s, 0));

      while (pq.isNotEmpty) {
        Item item = pq.removeMin();
        int u = item.value;
        int d = item.priority;

        if (d > dist[u]) continue;

        for (int i = 0; i < adj[u].length; i++) {
          Edge e = adj[u][i];
          if (e.capacity - e.flow > 0) {
            int newDist = dist[u] + e.cost;
            if (dist[e.to] > newDist) {
              dist[e.to] = newDist;
              parentNode[e.to] = u;
              parentEdge[e.to] = i;
              pq.add(Item(e.to, newDist));
            }
          }
        }
      }

      if (dist[t] == INF) {
        return -1; // Cannot push more flow
      }


      // Push flow
      int push = k - totalFlow;
      int curr = t;
      while (curr != s) {
        int p = parentNode[curr];
        int idx = parentEdge[curr];
        int available = adj[p][idx].capacity - adj[p][idx].flow;
        if (available < push) {
          push = available;
        }
        curr = p;
      }

      totalFlow += push;
      curr = t;
      while (curr != s) {
        int p = parentNode[curr];
        int idx = parentEdge[curr];
        adj[p][idx].flow += push;
        int revIdx = adj[p][idx].rev;
        adj[curr][revIdx].flow -= push;
        minCost += push * adj[p][idx].cost;
        curr = p;
      }
    }

    return minCost;
  }
}

class Item implements Comparable<Item> {
  int value;
  int priority;
  Item(this.value, this.priority);

  @override
  int compareTo(Item other) {
    return priority.compareTo(other.priority);
  }
}

class PriorityQueue {
  List<Item> _heap = [];

  void add(Item item) {
    _heap.add(item);
    _bubbleUp(_heap.length - 1);
  }

  Item removeMin() {
    if (_heap.isEmpty) throw StateError("Heap is empty");
    Item result = _heap[0];
    Item last = _heap.removeLast();
    if (_heap.isNotEmpty) {
      _heap[0] = last;
      _bubbleDown(0);
    }
    return result;
  }

  bool get isNotEmpty => _heap.isNotEmpty;

  void _bubbleUp(int index) {
    while (index > 0) {
      int parent = (index - 1) ~/ 2;
      if (_heap[index].compareTo(_heap[parent]) < 0) {
        Item temp = _heap[index];
        _heap[index] = _heap[parent];
        _heap[parent] = temp;
        index = parent;
      } else {
        break;
      }
    }
  }

  void _bubbleDown(int index) {
    int size = _heap.length;
    while (true) {
      int left = 2 * index + 1;
      int right = 2 * index + 2;
      int smallest = index;

      if (left < size && _heap[left].compareTo(_heap[smallest]) < 0) {
        smallest = left;
      }
      if (right < size && _heap[right].compareTo(_heap[smallest]) < 0) {
        smallest = right;
      }
      if (smallest == index) break;

      Item temp = _heap[index];
      _heap[index] = _heap[smallest];
      _heap[smallest] = temp;
      index = smallest;
    }
  }
}

void main() {
  final builder = BytesBuilder(copy: false);
  stdin.listen((event) {
    builder.add(event);
  }, onDone: () {
    if (builder.isNotEmpty) {
      solve(Scanner(builder.takeBytes()));
    }
  });
}

class Scanner {
  final Uint8List _bytes;
  int _ptr = 0;

  Scanner(this._bytes);

  bool hasNext() {
    while (_ptr < _bytes.length && _bytes[_ptr] <= 32) {
      _ptr++;
    }
    return _ptr < _bytes.length;
  }

  int readInt() {
    while (_ptr < _bytes.length && _bytes[_ptr] <= 32) {
      _ptr++;
    }
    int sign = 1;
    if (_ptr < _bytes.length && _bytes[_ptr] == 45) { // '-'
      sign = -1;
      _ptr++;
    }
    int res = 0;
    while (_ptr < _bytes.length) {
      int c = _bytes[_ptr];
      if (c < 48 || c > 57) break;
      res = res * 10 + (c - 48);
      _ptr++;
    }
    return res * sign;
  }
}

void solve(Scanner sc) {
  final buffer = StringBuffer();
  
  if (!sc.hasNext()) return;
  int t = sc.readInt();

  for (int i = 0; i < t; i++) {
    if (!sc.hasNext()) return;
    int n = sc.readInt();
    if (!sc.hasNext()) return;
    int m = sc.readInt();

    List<int> a = [];
    for (int j = 0; j < n; j++) {
      if (!sc.hasNext()) return;
      a.add(sc.readInt());
    }

    List<int> b = [];
    for (int j = 0; j < m; j++) {
      if (!sc.hasNext()) return;
      b.add(sc.readInt());
    }

    int result = solveTestCase(n, m, a, b);
    buffer.writeln(result);
  }

  stdout.write(buffer.toString());
}

int solveTestCase(int n, int m, List<int> a, List<int> b) {
  int source = 0;
  int sink = n + m + 1;
  int numNodes = n + m + 2;

  Graph g = Graph(numNodes);

  // Edges from Source to A
  for (int i = 0; i < n; i++) {
    g.addEdge(source, i + 1, 1, 0);
  }

  // Edges from A to B
  for (int i = 0; i < n; i++) {
    for (int j = 0; j < m; j++) {
      int val = a[i];
      int div = b[j];
      int rem = val % div;
      int cost = 0;
      if (rem != 0) {
        cost = div - rem;
      }
      g.addEdge(i + 1, n + 1 + j, 1, cost);
    }
  }

  // Edges from B to Sink
  int q = n ~/ m;

  const int M = 1000000000; // 1e9

  for (int j = 0; j < m; j++) {
    if (q > 0) {
      g.addEdge(n + 1 + j, sink, q, -M);
    }
    g.addEdge(n + 1 + j, sink, 1, 0);
  }

  // Calculate Min Cost for flow = n
  int rawCost = g.minCostMaxFlow(source, sink, n);

  // Adjust cost by removing the artificial negative costs
  int realCost = rawCost + q * m * M;

  return realCost;
}

