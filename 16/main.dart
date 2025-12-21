import 'dart:io';
import 'dart:typed_data';
import 'dart:collection';

const int MOD = 998244353;
const int MAX_N = 3005;

late List<List<int>> stirlingSum;

void init() {
  List<int> prevStirling = List.filled(MAX_N, 0);
  List<int> currStirling = List.filled(MAX_N, 0);

  prevStirling[0] = 1;
  stirlingSum[0][0] = 1;
  for (int j = 1; j < MAX_N; j++) {
    stirlingSum[0][j] = 1;
  }

  for (int i = 1; i < MAX_N; i++) {
    currStirling[0] = 0; // S(n, 0) = 0 for n >= 1
    stirlingSum[i][0] = 0;
    int currentSum = 0;

    for (int j = 1; j <= i; j++) {
      // S(n, k) = S(n-1, k-1) + (n-1)*S(n-1, k)
      int val = (prevStirling[j - 1] + (i - 1) * prevStirling[j]) % MOD;
      currStirling[j] = val;
      currentSum = (currentSum + currStirling[j]) % MOD;
      stirlingSum[i][j] = currentSum;
    }
    // Fill remaining sums with the total sum for this row
    for (int j = i + 1; j < MAX_N; j++) {
      stirlingSum[i][j] = currentSum;
    }

    // Update prevStirling for next iteration
    for (int j = 0; j <= i; j++) {
      prevStirling[j] = currStirling[j];
    }
  }
}

void main() {
  stirlingSum = List.generate(MAX_N, (_) => List.filled(MAX_N, 0));
  init();

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

int solveTestCase(int n, int q, int l, int r, List<int> b, List<int> c) {
  List<int> inDeg = List.filled(n + 1, 0);
  List<int> outDeg = List.filled(n + 1, 0);
  List<int> adj = List.filled(n + 1, 0);

  bool possible = true;
  for (int i = 0; i < q; i++) {
    int u = b[i], v = c[i];
    if (outDeg[u] > 0 || inDeg[v] > 0) {
      possible = false;
    }
    outDeg[u]++;
    inDeg[v]++;
    adj[u] = v;
  }

  if (!possible) {
    return 0;
  }

  // M = n - q (number of path components)
  int M = n - q;

  // Count fixed cycles
  int fixedCycles = 0;
  List<bool> visited = List.filled(n + 1, false);

  // First, traverse everything starting from nodes with inDeg == 0 (Path starts)
  for (int i = 1; i <= n; i++) {
    if (inDeg[i] == 0) {
      int curr = i;
      while (curr != 0 && !visited[curr]) {
        visited[curr] = true;
        if (outDeg[curr] > 0) {
          curr = adj[curr];
        } else {
          curr = 0;
        }
      }
    }
  }

  // Remaining unvisited nodes must be part of cycles
  for (int i = 1; i <= n; i++) {
    if (!visited[i]) {
      // Found a cycle
      fixedCycles++;
      int curr = i;
      while (!visited[curr]) {
        visited[curr] = true;
        curr = adj[curr];
      }
    }
  }

  // Range of cycles needed from path components
  int needL = l - fixedCycles;
  int needR = r - fixedCycles;

  if (needL < 0) {
    needL = 0;
  }
  if (needR < 0) {
    return 0;
  }
  if (needL > M) {
    return 0;
  }
  if (needR > M) {
    needR = M;
  }

  // Answer is sum of Stirling numbers [M][k] for k in [needL, needR]
  int sub = 0;
  if (needL > 0) {
    sub = stirlingSum[M][needL - 1];
  }
  int ans = (stirlingSum[M][needR] - sub + MOD) % MOD;
  return ans;
}

void solve(Scanner sc) {
  final buffer = StringBuffer();

  if (!sc.hasNext()) return;
  int t = sc.readInt();

  List<TestCase> testCases = [];
  for (int i = 0; i < t; i++) {
    if (!sc.hasNext()) return;
    int n = sc.readInt();
    if (!sc.hasNext()) return;
    int q = sc.readInt();
    if (!sc.hasNext()) return;
    int l = sc.readInt();
    if (!sc.hasNext()) return;
    int r = sc.readInt();

    List<int> b = [];
    for (int j = 0; j < q; j++) {
      if (!sc.hasNext()) return;
      b.add(sc.readInt());
    }

    List<int> c = [];
    for (int j = 0; j < q; j++) {
      if (!sc.hasNext()) return;
      c.add(sc.readInt());
    }

    testCases.add(TestCase(n, q, l, r, b, c));
  }

  for (int i = 0; i < t; i++) {
    int result = solveTestCase(testCases[i].n, testCases[i].q, testCases[i].l, testCases[i].r, testCases[i].b, testCases[i].c);
    buffer.writeln(result);
  }

  stdout.write(buffer.toString());
}

class TestCase {
  int n, q, l, r;
  List<int> b, c;
  TestCase(this.n, this.q, this.l, this.r, this.b, this.c);
}
