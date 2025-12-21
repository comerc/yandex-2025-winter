import 'dart:io';
import 'dart:typed_data';

const int MAX_N = 700000;
const int MAX_K = 9;

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

  // Предвычисляем префиксные суммы
  List<List<int>> prefixSums = precompute();

  int q = sc.readInt();
  for (int i = 0; i < q; i++) {
    int k = sc.readInt();
    int l = sc.readInt();
    int r = sc.readInt();
    int result = query(prefixSums, k, l, r);
    buffer.writeln(result);
  }

  stdout.write(buffer.toString());
}

List<List<int>> precompute() {
  List<List<int>> prefixSums = List.generate(MAX_K + 1, (_) => List.filled(MAX_N + 1, 0));

  // Генерируем k-интересные числа для k >= 2
  generate(1, 2, 0, MAX_N, prefixSums);

  // Все числа >= 2 являются 1-интересными
  for (int n = 2; n <= MAX_N; n++) {
    prefixSums[1][n] = 1;
  }

  // Префиксные суммы
  for (int k = 1; k <= MAX_K; k++) {
    for (int n = 1; n <= MAX_N; n++) {
      prefixSums[k][n] += prefixSums[k][n - 1];
    }
  }

  return prefixSums;
}

void generate(int product, int minFactor, int depth, int limit, List<List<int>> prefixSums) {
  for (int factor = minFactor; product * factor <= limit; factor++) {
    int newProduct = product * factor;
    int newDepth = depth + 1;

    if (newDepth >= 2 && newDepth <= MAX_K) {
      prefixSums[newDepth][newProduct] = 1;
    }

    generate(newProduct, factor + 1, newDepth, limit, prefixSums);
  }
}

int query(List<List<int>> prefixSums, int k, int l, int r) {
  if (k > MAX_K || l > MAX_N) {
    return 0;
  }
  if (r > MAX_N) {
    r = MAX_N;
  }
  if (l < 1) {
    l = 1;
  }
  return prefixSums[k][r] - (l >= 1 ? prefixSums[k][l - 1] : 0);
}
