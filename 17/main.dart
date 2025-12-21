import 'dart:io';
import 'dart:typed_data';

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
  int n = sc.readInt();

  // Find M such that 10^M > n
  int m = 0;
  int limit = 1;
  while (limit <= n) {
    limit *= 10;
    m++;
  }
  if (m == 0) {
    m = 1;
  }

  String s = '9' * m;
  buffer.writeln(s);
  buffer.writeln(s);

  stdout.write(buffer.toString());
}
