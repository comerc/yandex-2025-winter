import 'dart:io';
import 'dart:typed_data';

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

  int R = sc.readInt();
  int B = sc.readInt();

  final result = solveTest(R, B);
  buffer.writeln('${result[0]} ${result[1]}');

  stdout.write(buffer.toString());
}

List<int> solveTest(int R, int B) {
  int sumWH = (R + 4) ~/ 2;

  for (int d = 1; d * d <= B; d++) {
    if (B % d != 0) continue;

    // Проверяем оба варианта: (W-2, H-2) = (d, B/d) и (B/d, d)
    final candidates = [B ~/ d + 2, d + 2];
    for (int w in candidates) {
      int h = sumWH - w;
      if ((w - 2) * (h - 2) == B && w >= h) {
        return [w, h];
      }
    }
  }

  return [0, 0]; // Не должно случиться
}
