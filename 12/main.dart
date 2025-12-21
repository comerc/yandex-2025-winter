import 'dart:io';
import 'dart:typed_data';

const int MOD = 998244353;

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

  int n = sc.readInt();
  int m = sc.readInt();
  Set<int> good = {};
  for (int i = 0; i < m; i++) {
    good.add(sc.readInt());
  }

  int result = solveFunc(n, good);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveFunc(int n, Set<int> good) {
  if (n < 3) {
    return 0;
  }

  // Матрица переходов 100x100
  List<List<int>> M = List.generate(100, (_) => List.filled(100, 0));

  for (int d1 = 0; d1 < 10; d1++) {
    for (int d2 = 0; d2 < 10; d2++) {
      int from = d1 * 10 + d2;
      for (int d3 = 0; d3 < 10; d3++) {
        int sum = d1 + d2 + d3;
        if (good.contains(sum)) {
          int to = d2 * 10 + d3;
          M[from][to] = 1;
        }
      }
    }
  }

  // Начальный вектор: для чисел длины 2 (d1, d2), где d1 != 0
  List<int> start = List.filled(100, 0);
  for (int d1 = 1; d1 < 10; d1++) {
    for (int d2 = 0; d2 < 10; d2++) {
      int idx = d1 * 10 + d2;
      start[idx] = 1;
    }
  }

  // Если n == 2, возвращаем количество начальных состояний
  if (n == 2) {
    int result = 0;
    for (int i = 0; i < 100; i++) {
      result = (result + start[i]) % MOD;
    }
    return result;
  }

  // Возводим матрицу в степень (n-2)
  int power = n - 2;
  List<List<int>> M_power = matrixPower(M, power);

  // Умножаем начальный вектор на матрицу
  int result = 0;
  for (int i = 0; i < 100; i++) {
    for (int j = 0; j < 100; j++) {
      result = (result + start[i] * M_power[i][j] % MOD) % MOD;
    }
  }

  return result;
}

List<List<int>> matrixPower(List<List<int>> M, int power) {
  int n = M.length;

  // Единичная матрица
  List<List<int>> result = List.generate(n, (_) => List.filled(n, 0));
  for (int i = 0; i < n; i++) {
    result[i][i] = 1;
  }

  List<List<int>> base = List.generate(n, (_) => List.filled(n, 0));
  for (int i = 0; i < n; i++) {
    for (int j = 0; j < n; j++) {
      base[i][j] = M[i][j];
    }
  }

  // Быстрое возведение в степень
  while (power > 0) {
    if (power & 1 == 1) {
      result = matrixMultiply(result, base);
    }
    base = matrixMultiply(base, base);
    power >>= 1;
  }

  return result;
}

List<List<int>> matrixMultiply(List<List<int>> A, List<List<int>> B) {
  int n = A.length;
  int m = B[0].length;
  int k = B.length;

  List<List<int>> result = List.generate(n, (_) => List.filled(m, 0));
  for (int i = 0; i < n; i++) {
    for (int j = 0; j < m; j++) {
      int sum = 0;
      for (int t = 0; t < k; t++) {
        sum = (sum + A[i][t] * B[t][j] % MOD) % MOD;
      }
      result[i][j] = sum;
    }
  }

  return result;
}
