import 'dart:io';
import 'dart:typed_data';

const int MOD = 1000000007;

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

  int result = solveFunc(n);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveFunc(int n) {
  if (n == 1) {
    return 1;
  }

  int nMod = n % MOD;

  // alpha[k] и beta[k] для прогонки
  // E[k] = alpha[k]*E[k+1] + beta[k]
  List<int> alpha = List.filled(n + 1, 0);
  List<int> beta = List.filled(n + 1, 0);

  // alpha[0] = 0, beta[0] = 0
  alpha[0] = 0;
  beta[0] = 0;

  // Прямой ход: для k = 1..n-1
  for (int k = 1; k < n; k++) {
    int kMod = k % MOD;
    int nMinusK = (n - k) % MOD;

    // denominator = n - k*alpha[k-1]
    int denom = (nMod - kMod * alpha[k - 1] % MOD + MOD) % MOD;
    int denomInv = modInverse(denom, MOD);

    // alpha[k] = (n-k) / denom
    alpha[k] = nMinusK * denomInv % MOD;

    // beta[k] = (n + k*beta[k-1]) / denom
    int numerator = (nMod + kMod * beta[k - 1] % MOD) % MOD;
    beta[k] = numerator * denomInv % MOD;
  }

  // Граничное условие: E[n] = 1 + E[n-1]
  // E[n] = (1 + beta[n-1]) / (1 - alpha[n-1])
  int oneMinusAlpha = (1 - alpha[n - 1] + MOD) % MOD;
  int oneMinusAlphaInv = modInverse(oneMinusAlpha, MOD);
  int en = (1 + beta[n - 1]) % MOD * oneMinusAlphaInv % MOD;

  return en;
}

// modInverse вычисляет обратное число a^(-1) mod m
int modInverse(int a, int m) {
  return modPow(a, m - 2, m);
}

// modPow вычисляет base^exp mod m
int modPow(int base, int exp, int m) {
  int result = 1;
  base = base % m;
  while (exp > 0) {
    if (exp % 2 == 1) {
      result = (result * base) % m;
    }
    exp = exp >> 1;
    base = (base * base) % m;
  }
  return result;
}
