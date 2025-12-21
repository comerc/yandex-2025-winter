import 'dart:io';
import 'dart:typed_data';
import 'dart:math';

const int MAX_NODES = 13000000;
const int MAX_BITS = 29;

Int32List l = Int32List(MAX_NODES);
Int32List r = Int32List(MAX_NODES);
Int32List cnt = Int32List(MAX_NODES);
Int32List memoList = Int32List(MAX_NODES);
int ptr = 0;

int newNode() {
  ptr++;
  int idx = ptr;
  l[idx] = 0;
  r[idx] = 0;
  cnt[idx] = 0;
  memoList[idx] = 0;
  return idx;
}

void pushUp(int u, int bit) {
  int idx0 = l[u];
  int idx1 = r[u];

  int c0 = idx0 != 0 ? cnt[idx0] : 0;
  int c1 = idx1 != 0 ? cnt[idx1] : 0;

  cnt[u] = c0 + c1;

  int m0 = idx0 != 0 ? memoList[idx0] : 0;
  int m1 = idx1 != 0 ? memoList[idx1] : 0;

  int full = 1 << bit;

  if (c0 == full) {
    memoList[u] = full + m1;
  } else if (c1 == full) {
    memoList[u] = full + m0;
  } else {
    memoList[u] = max(m0, m1);
  }
}

void update(int u, int val, int bit, bool add) {
  if (bit < 0) {
    if (add) {
      cnt[u] = 1;
      memoList[u] = 1;
    } else {
      cnt[u] = 0;
      memoList[u] = 0;
    }
    return;
  }

  int dir = (val >> bit) & 1;
  int childIdx;
  if (dir == 0) {
    if (l[u] == 0) l[u] = newNode();
    childIdx = l[u];
  } else {
    if (r[u] == 0) r[u] = newNode();
    childIdx = r[u];
  }

  update(childIdx, val, bit - 1, add);
  pushUp(u, bit);
}

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
  int t = sc.readInt();

  ptr = 0;

  for (int i = 0; i < t; i++) {
    if (!sc.hasNext()) return;
    int n = sc.readInt();
    if (!sc.hasNext()) return;
    int q = sc.readInt();

    List<int> a = List.filled(n, 0);
    Map<int, int> freq = {};

    int root = newNode();

    for (int j = 0; j < n; j++) {
      if (!sc.hasNext()) return;
      a[j] = sc.readInt();
      freq[a[j]] = (freq[a[j]] ?? 0) + 1;
      if (freq[a[j]] == 1) {
        update(root, a[j], MAX_BITS, true);
      }
    }

    buffer.writeln(memoList[root]);

    for (int k = 0; k < q; k++) {
      if (!sc.hasNext()) return;
      int j = sc.readInt() - 1; // 0-based
      if (!sc.hasNext()) return;
      int v = sc.readInt();

      int oldVal = a[j];
      if (oldVal != v) {
        freq[oldVal] = freq[oldVal]! - 1;
        if (freq[oldVal] == 0) {
          update(root, oldVal, MAX_BITS, false);
        }

        a[j] = v;
        freq[v] = (freq[v] ?? 0) + 1;
        if (freq[v] == 1) {
          update(root, v, MAX_BITS, true);
        }
      }

      buffer.writeln(memoList[root]);
    }
  }

  stdout.write(buffer.toString());
}

