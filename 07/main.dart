import 'dart:io';
import 'dart:typed_data';

const int MOD = 998244353;
const int MAX_N = 400001;

late List<int> fact;
late List<int> invFact;

void init() {
  fact = List.filled(MAX_N, 0);
  fact[0] = 1;
  for (int i = 1; i < MAX_N; i++) {
    fact[i] = fact[i - 1] * i % MOD;
  }

  invFact = List.filled(MAX_N, 0);
  invFact[MAX_N - 1] = modPow(fact[MAX_N - 1], MOD - 2);
  for (int i = MAX_N - 1; i > 0; i--) {
    invFact[i - 1] = invFact[i] * i % MOD;
  }
}

int modPow(int a, int b) {
  int result = 1;
  a %= MOD;
  while (b > 0) {
    if (b & 1 == 1) {
      result = result * a % MOD;
    }
    a = a * a % MOD;
    b >>= 1;
  }
  return result;
}

void main() {
  init();

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

  int T = sc.readInt();
  for (int i = 0; i < T; i++) {
    int n = sc.readInt();
    int s = sc.readInt();
    int result = solveTest(n, s);
    buffer.writeln(result);
  }

  stdout.write(buffer.toString());
}

int solveTest(int n, int s) {
  if (n > s) {
    return 0;
  }

  int c = comb(s, n);
  return fact[n + 1] * c % MOD;
}

int comb(int n, int k) {
  if (k < 0 || k > n) {
    return 0;
  }
  return fact[n] * invFact[k] % MOD * invFact[n - k] % MOD;
}
