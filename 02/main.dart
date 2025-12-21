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

  final nums = Int32List(10);
  for (int i = 0; i < 10; i++) {
    nums[i] = sc.readInt();
  }

  int result = solveTest(nums);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveTest(Int32List nums) {
  const int target = 100;
  int bestSum = 0;
  int bestDist = target;

  // Перебор всех 2^10 = 1024 подмножеств
  for (int mask = 0; mask < (1 << 10); mask++) {
    int sum = 0;
    for (int i = 0; i < 10; i++) {
      if ((mask & (1 << i)) != 0) {
        sum += nums[i];
      }
    }

    int dist = (sum - target).abs();

    if (dist < bestDist || (dist == bestDist && sum > bestSum)) {
      bestSum = sum;
      bestDist = dist;
    }
  }

  return bestSum;
}
