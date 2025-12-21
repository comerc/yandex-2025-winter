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

  String readLine() {
    while (_ptr < _input.length && _input.codeUnitAt(_ptr) <= 32) {
      _ptr++;
    }
    int start = _ptr;
    while (_ptr < _input.length && _input.codeUnitAt(_ptr) != 10) {
      _ptr++;
    }
    if (_ptr < _input.length) _ptr++;
    return _input.substring(start, _ptr - 1);
  }
}

void solve(Scanner sc) {
  final buffer = StringBuffer();

  int N = sc.readInt();
  int M = sc.readInt();

  final markersX = Int32List(N);
  final markersY = Int32List(N);
  for (int i = 0; i < N; i++) {
    markersX[i] = sc.readInt();
    markersY[i] = sc.readInt();
  }

  String commands = sc.readLine();

  final results = solveTest(markersX, markersY, commands);
  for (int i = 0; i < results.length; i++) {
    buffer.writeln(results[i]);
  }

  stdout.write(buffer.toString());
}

List<int> solveTest(Int32List markersX, Int32List markersY, String commands) {
  int N = markersX.length;

  final sortedX = Int32List.fromList(markersX);
  sortedX.sort();
  final sortedY = Int32List.fromList(markersY);
  sortedY.sort();

  final prefixSumX = Int64List(N + 1);
  final prefixSumY = Int64List(N + 1);
  for (int i = 0; i < N; i++) {
    prefixSumX[i + 1] = prefixSumX[i] + sortedX[i];
    prefixSumY[i + 1] = prefixSumY[i] + sortedY[i];
  }

  int cx = 0, cy = 0;
  final results = <int>[];

  for (int i = 0; i < commands.length; i++) {
    final cmd = commands[i];
    switch (cmd) {
      case 'N':
        cy++;
        break;
      case 'S':
        cy--;
        break;
      case 'E':
        cx++;
        break;
      case 'W':
        cx--;
        break;
    }

    // sum |cx - mx|
    int idxX = _upperBound(sortedX, cx);
    int leftCount = idxX;
    int rightCount = N - idxX;
    int sumX = cx * leftCount - prefixSumX[idxX].toInt() +
               (prefixSumX[N] - prefixSumX[idxX]).toInt() - cx * rightCount;

    int idxY = _upperBound(sortedY, cy);
    int leftCountY = idxY;
    int rightCountY = N - idxY;
    int sumY = cy * leftCountY - prefixSumY[idxY].toInt() +
               (prefixSumY[N] - prefixSumY[idxY]).toInt() - cy * rightCountY;

    int total = sumX + sumY;
    results.add(total);
  }

  return results;
}

int _upperBound(Int32List list, int value) {
  int low = 0;
  int high = list.length;
  while (low < high) {
    int mid = (low + high) ~/ 2;
    if (list[mid] <= value) {
      low = mid + 1;
    } else {
      high = mid;
    }
  }
  return low;
}
