import 'dart:io';
import 'dart:typed_data';

const int MOD = 998244353;

void main() {
  final builder = BytesBuilder(copy: false);
  stdin.listen((event) {
    builder.add(event);
  }, onDone: () {
    final bytes = builder.takeBytes();
    final input = String.fromCharCodes(bytes);
    final lines = input.trim().split('\n');
    int idx = 0;
    int T = int.parse(lines[idx++]);
    final buffer = StringBuffer();
    for (int tIdx = 0; tIdx < T; tIdx++) {
      String s = lines[idx++];
      String t = lines[idx++];
      int ans = processTestCase(s, t);
      buffer.writeln(ans);
    }
    stdout.write(buffer.toString());
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

  String readString() {
    while (_ptr < _input.length && _input.codeUnitAt(_ptr) <= 32) {
      _ptr++;
    }
    int start = _ptr;
    while (_ptr < _input.length && _input.codeUnitAt(_ptr) > 32) {
      _ptr++;
    }
    return _input.substring(start, _ptr);
  }
}

void solve(Scanner sc) {
  final buffer = StringBuffer();

  if (!sc.hasNext()) return;
  int T = sc.readInt();
  print('T: $T');

  for (int tIdx = 0; tIdx < T; tIdx++) {
    String s = sc.readString();
    String t = sc.readString();
    print('s: "$s", t: "$t"');
    int ans = processTestCase(s, t);
    buffer.writeln(ans);
  }

  stdout.write(buffer.toString());
}

int processTestCase(String s, String t) {
  int n = s.length;
  int m = t.length;

  // 1. Предподсчет количества разбиений для строки длины k.
  // Если длина k, то есть k-1 мест для разрыва, итого 2^(k-1) способов.
  // Для k=0 (пустая строка) считаем 1 способ.
  List<int> partitions = List.filled(n + 1, 0);
  partitions[0] = 1;
  if (n > 0) {
    int p = 1;
    for (int i = 1; i <= n; i++) {
      partitions[i] = p;
      p = (p * 2) % MOD;
    }
  }

  // 2. Предподсчет LCP (Longest Common Prefix) для всех пар суффиксов s и t.
  // lcp[i][j] = длина общего префикса s[i:] и t[j:]
  List<List<int>> lcp = List.generate(n + 1, (_) => List.filled(m + 1, 0));
  for (int i = n - 1; i >= 0; i--) {
    for (int j = m - 1; j >= 0; j--) {
      if (s[i] == t[j]) {
        lcp[i][j] = 1 + lcp[i + 1][j + 1];
      } else {
        lcp[i][j] = 0;
      }
    }
  }

  // 3. Динамическое программирование
  // dp[i][u] - кол-во способов разбить суффикс s[i:], чтобы результат совпадал с t[u...]
  List<List<int>> dp = List.generate(n + 1, (_) => List.filled(m + 1, 0));

  // База: пустой суффикс s совпадает с "началом" любой подстроки t (пустой строкой)
  for (int u = 0; u <= m; u++) {
    dp[n][u] = 1;
  }

  int ans = 0;

  // Перебираем длину оставшегося суффикса s (от меньшего к большему с точки зрения "потребления" s справа налево)
  // i идет от n до 1. Мы пытаемся откусить кусок s[k...i-1].
  for (int i = n; i >= 1; i--) {
    // Длина уже сформированной части результата
    int currentResLen = n - i;

    for (int u = 0; u <= m; u++) {
      if (dp[i][u] == 0) {
        continue;
      }

      // Позиция в t, с которой мы должны сравнивать следующий кусок
      int posInT = u + currentResLen;

      // Если мы вышли за пределы t, сравнение невозможно (строка t кончилась раньше)
      if (posInT >= m) {
        continue;
      }

      // Перебираем, где отрезать следующий кусок от s (индекс k)
      // Кусок будет s[k...i-1]
      for (int k = i - 1; k >= 0; k--) {
        int chunkLen = i - k;

        // Используем LCP для быстрого сравнения s[k...] и t[posInT...]
        int val = lcp[k][posInT];

        if (val >= chunkLen) {
          // Кусок полностью совпал с частью t
          // Если мы не вышли за границы t, обновляем ДП
          if (posInT + chunkLen <= m) {
            dp[k][u] = (dp[k][u] + dp[i][u]) % MOD;
          }
        } else {
          // Куски различаются.
          // Индекс различия относительно начала куска: val
          int idxS = k + val;
          int idxT = posInT + val;

          // Проверяем, что различие произошло в пределах строки t
          if (idxT < m) {
            // Если символ в s меньше символа в t, то результат кодирования лексикографически меньше
            if (s.codeUnitAt(idxS) < t.codeUnitAt(idxT)) {
              // Мы нашли "меньший" вариант.
              // Длина совпадающего префикса = currentResLen + val.
              int matchLen = currentResLen + val;

              // Этот результат будет меньше любой подстроки t, начинающейся в u,
              // которая длиннее matchLen.
              // Максимальная длина подстроки от u: m - u.
              // Подходящие длины: matchLen+1, matchLen+2, ..., m-u.
              int count = (m - u) - matchLen;

              if (count > 0) {
                // Добавляем к ответу:
                // (способы дойти до i) * (способы разбить остаток s[0...k-1]) * (кол-во подстрок t)
                int ways = (dp[i][u] * partitions[k]) % MOD;
                int term = (ways * count) % MOD;
                ans = (ans + term) % MOD;
              }
            }
            // Если s[idxS] > t[idxT], то результат больше, ничего не делаем.
          }
        }
      }
    }
  }

  // Обработка случаев, когда вся строка s (переставленная) является строгим префиксом подстроки t.
  // Это соответствует состояниям dp[0][u].
  // Результат имеет длину n и совпадает с t[u ... u+n-1].
  // Он будет меньше любой подстроки t[u...], длина которой > n.
  for (int u = 0; u <= m; u++) {
    if (dp[0][u] > 0) {
      int count = (m - u) - n;
      if (count > 0) {
        int term = (dp[0][u] * count) % MOD;
        ans = (ans + term) % MOD;
      }
    }
  }

  return ans;
}
