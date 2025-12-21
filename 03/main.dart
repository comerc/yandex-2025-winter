import 'dart:io';
import 'dart:typed_data';

const int MOD = 1000000007;

class Pair {
  int w;
  int index;
  Pair(this.w, this.index);
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

  int M = sc.readInt();
  int N = sc.readInt();

  final W = Int64List(N);
  int totalNeed = 0;
  for (int i = 0; i < N; i++) {
    W[i] = sc.readInt();
    totalNeed += W[i].toInt();
  }

  int result = solveTest(M, W, totalNeed);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveTest(int M, Int64List W, int totalNeed) {
  int N = W.length;
  int deficit = totalNeed - M;

  if (deficit == 0) {
    return 0;
  }

  final pairs = List<Pair>.generate(N, (i) => Pair(W[i].toInt(), i));
  pairs.sort((a, b) => b.w.compareTo(a.w)); // descending

  final shortfall = Int64List(N);
  int remaining = deficit;
  int currentN = N;

  while (remaining > 0 && currentN > 0) {
    int avgShortfall = remaining ~/ currentN;
    int minW = pairs[currentN - 1].w;

    if (avgShortfall <= minW) {
      int baseShortfall = remaining ~/ currentN;
      int extra = remaining % currentN;

      for (int j = 0; j < currentN; j++) {
        shortfall[pairs[j].index] = baseShortfall;
      }

      for (int j = 0; j < extra; j++) {
        shortfall[pairs[j].index]++;
      }
      break;
    } else {
      shortfall[pairs[currentN - 1].index] = minW;
      remaining -= minW;
      currentN--;
    }
  }

  int result = 0;
  for (int i = 0; i < N; i++) {
    int sq = shortfall[i].toInt() % MOD;
    result = (result + sq * sq % MOD) % MOD;
  }

  return result;
}
