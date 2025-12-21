import 'dart:io';
import 'dart:typed_data';

class Point {
  int x, y, z, idx;
  Point(this.x, this.y, this.z, this.idx);
}

class Edge {
  int from, to, weight;
  Edge(this.from, this.to, this.weight);
}

void main() {
  final builder = BytesBuilder(copy: false);
  stdin.listen((event) {
    builder.add(event);
  }, onDone: () {
    if (builder.isNotEmpty) {
      final input = String.fromCharCodes(builder.takeBytes());
      final sc = Scanner(input);
      solve(sc);
    }
  });
}

class Scanner {
  final String _input;
  int _ptr = 0;

  Scanner(this._input);

  bool hasNext() {
    return _ptr < _input.length;
  }

  int readInt() {
    while (_ptr < _input.length && _input.codeUnitAt(_ptr) <= 32) {
      _ptr++;
    }
    int sign = 1;
    if (_ptr < _input.length && _input[_ptr] == '-') {
      sign = -1;
      _ptr++;
    }
    int res = 0;
    while (_ptr < _input.length) {
      int c = _input.codeUnitAt(_ptr);
      if (c < 48 || c > 57) break;
      res = res * 10 + (c - 48);
      _ptr++;
    }
    return res * sign;
  }
}

void solve(Scanner sc) {
  final buffer = StringBuffer();

  int N = sc.readInt();

  final points = List<Point>.generate(N, (i) => Point(0, 0, 0, i));
  for (int i = 0; i < N; i++) {
    points[i].x = sc.readInt();
    points[i].y = sc.readInt();
    points[i].z = sc.readInt();
  }

  int result = solveMST(points);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveMST(List<Point> points) {
  int N = points.length;

  final edges = <Edge>[];

  // Сортировка по x
  final sortedByX = List<Point>.from(points);
  sortedByX.sort((a, b) => a.x.compareTo(b.x));
  for (int i = 0; i < N - 1; i++) {
    int cost = min3(
      (sortedByX[i].x - sortedByX[i + 1].x).abs(),
      (sortedByX[i].y - sortedByX[i + 1].y).abs(),
      (sortedByX[i].z - sortedByX[i + 1].z).abs(),
    );
    edges.add(Edge(sortedByX[i].idx, sortedByX[i + 1].idx, cost));
  }

  // Сортировка по y
  final sortedByY = List<Point>.from(points);
  sortedByY.sort((a, b) => a.y.compareTo(b.y));
  for (int i = 0; i < N - 1; i++) {
    int cost = min3(
      (sortedByY[i].x - sortedByY[i + 1].x).abs(),
      (sortedByY[i].y - sortedByY[i + 1].y).abs(),
      (sortedByY[i].z - sortedByY[i + 1].z).abs(),
    );
    edges.add(Edge(sortedByY[i].idx, sortedByY[i + 1].idx, cost));
  }

  // Сортировка по z
  final sortedByZ = List<Point>.from(points);
  sortedByZ.sort((a, b) => a.z.compareTo(b.z));
  for (int i = 0; i < N - 1; i++) {
    int cost = min3(
      (sortedByZ[i].x - sortedByZ[i + 1].x).abs(),
      (sortedByZ[i].y - sortedByZ[i + 1].y).abs(),
      (sortedByZ[i].z - sortedByZ[i + 1].z).abs(),
    );
    edges.add(Edge(sortedByZ[i].idx, sortedByZ[i + 1].idx, cost));
  }

  // Сортировка рёбер по весу
  edges.sort((a, b) => a.weight.compareTo(b.weight));

  // DSU
  final parent = Int32List(N);
  final rank = Int32List(N);
  for (int i = 0; i < N; i++) {
    parent[i] = i;
    rank[i] = 0;
  }

  int totalCost = 0;
  int edgesUsed = 0;

  for (final edge in edges) {
    if (edgesUsed == N - 1) break;
    int fromRoot = find(parent, edge.from);
    int toRoot = find(parent, edge.to);
    if (fromRoot != toRoot) {
      union(parent, rank, fromRoot, toRoot);
      totalCost += edge.weight;
      edgesUsed++;
    }
  }

  return totalCost;
}

int find(Int32List parent, int x) {
  if (parent[x] != x) {
    parent[x] = find(parent, parent[x]);
  }
  return parent[x];
}

void union(Int32List parent, Int32List rank, int x, int y) {
  if (rank[x] < rank[y]) {
    parent[x] = y;
  } else if (rank[x] > rank[y]) {
    parent[y] = x;
  } else {
    parent[y] = x;
    rank[x]++;
  }
}

int min3(int a, int b, int c) {
  return a < b ? (a < c ? a : c) : (b < c ? b : c);
}
