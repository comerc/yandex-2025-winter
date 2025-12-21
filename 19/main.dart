import 'dart:io';
import 'dart:typed_data';

const int MOD = 1000000007;

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
  int m = sc.readInt();

  int N = 0;
  int lastCount = 0;

  for (int i = 0; i < m; i++) {
    if (!sc.hasNext()) return;
    int val = sc.readInt();
    N += val;
    if (val > 0) {
      lastCount = val;
    }
  }

  if (N == 0) {
    buffer.writeln(0);
  } else {
    int S = N - lastCount;
    if (S == 0) {
      buffer.writeln(0);
    } else {
      // Precompute inverses
      Int64List inv = Int64List(N + 1);
      inv[1] = 1;
      for (int i = 2; i <= N; i++) {
        inv[i] = (MOD - (MOD ~/ i) * inv[MOD % i] % MOD) % MOD;
      }

      int sumInv = 0;

      // Sum 1/(k*(N-k)) for k=1 to S
      for (int k = 1; k <= S; k++) {
        int term = (inv[k] * inv[N - k]) % MOD;
        sumInv = (sumInv + term) % MOD;
      }

      // Coeff = N*(N-1)/2
      int coeff = (N * (N - 1)) % MOD;
      coeff = (coeff * inv[2]) % MOD;

      int ans = (coeff * sumInv) % MOD;
      buffer.writeln(ans);
    }
  }

  stdout.write(buffer.toString());
}

