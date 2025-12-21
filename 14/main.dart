import 'dart:io';
import 'dart:typed_data';

const int MOD = 1000000007;
const int INV2 = 500000004; // 2^(-1) mod 10^9+7
const int INV3 = 333333336; // 3^(-1) mod 10^9+7

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

  int a = sc.readInt();
  int q = sc.readInt();
  int L = sc.readInt();
  int R = sc.readInt();

  int result = solveFunc(a, q, L, R);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveFunc(int a, int q, int L, int R) {
  int N = R - L + 1;
  if (N <= 0) {
    return 0;
  }

  int nMod = N % MOD;

  // Специальные случаи
  if (a == 0) {
    return 0;
  }

  // Случай q = 0
  if (q == 0) {
    int zeros = 0;
    if (L <= 0 && R >= 0) {
      zeros = 1;
    }
    if (zeros == 0) {
      return 0;
    }
    int nonZeros = (N - zeros) % MOD;
    int validDenoms = (2 * zeros) % MOD;
    validDenoms = (validDenoms * nonZeros) % MOD;
    int allNums = (nMod * nMod) % MOD;
    return (validDenoms * allNums) % MOD;
  }

  // Случай q = 1
  if (q == 1) {
    return 0;
  }

  // Случай q = -1
  if (q == -1) {
    int evens = 0, odds = 0;
    bool lIsEven = (L % 2 == 0);
    if (N % 2 == 0) {
      evens = N ~/ 2;
      odds = N ~/ 2;
    } else {
      if (lIsEven) {
        evens = N ~/ 2 + 1;
        odds = N ~/ 2;
      } else {
        evens = N ~/ 2;
        odds = N ~/ 2 + 1;
      }
    }
    int validDenoms = (2 * (evens % MOD)) % MOD;
    validDenoms = (validDenoms * (odds % MOD)) % MOD;
    int allNums = (nMod * nMod) % MOD;
    return (validDenoms * allNums) % MOD;
  }

  // Общий случай
  return solveGeneral(a, q, L, R);
}

int solveGeneral(int a, int q, int L, int R) {
  int N = R - L + 1;
  if (N <= 0) {
    return 0;
  }

  int nMod = N % MOD;

  // 1. Вклад n = m (числитель 0). Знаменатель любой ненулевой (k != s).
  // Кол-во = N * (N * (N - 1))
  int ans = (nMod * nMod) % MOD;
  ans = (ans * ((nMod - 1 + MOD) % MOD)) % MOD;

  // 2. Вклад n != m. Используем Block Summation
  int termN_N1 = (nMod * ((nMod + 1) % MOD)) % MOD; // N(N+1)
  int term2N_1 = ((2 * nMod) + 1) % MOD; // 2N+1

  int limit = N - 1;
  int l = 1;

  while (l <= limit) {
    int kMax = limit ~/ l;
    int r = limit ~/ kMax;
    if (r > limit) {
      r = limit;
    }

    int kMaxMod = kMax % MOD;

    // s1 = sum(k) = k(k+1)/2
    int s1 = (kMaxMod * (kMaxMod + 1)) % MOD;
    s1 = (s1 * INV2) % MOD;

    // s2 = sum(k^2) = k(k+1)(2k+1)/6
    int term2k1 = ((2 * kMaxMod) + 1) % MOD;
    int s2 = (s1 * term2k1) % MOD;
    s2 = (s2 * INV3) % MOD;

    int s0 = kMaxMod;

    // Вычисляем суммы для B на отрезке [l, r]
    int lMod = l % MOD;
    int rMod = r % MOD;

    // Сумма 1..r для B и B^2
    int sumR1 = (rMod * (rMod + 1)) % MOD;
    sumR1 = (sumR1 * INV2) % MOD;
    int term2r1 = ((2 * rMod) + 1) % MOD;
    int sumR2 = (sumR1 * term2r1) % MOD;
    sumR2 = (sumR2 * INV3) % MOD;

    // Сумма 1..l-1 для B и B^2
    int lm1 = (lMod - 1 + MOD) % MOD;
    int sumL1 = (lm1 * (lm1 + 1)) % MOD;
    sumL1 = (sumL1 * INV2) % MOD;
    int term2l1 = ((2 * lm1) + 1) % MOD;
    int sumL2 = (sumL1 * term2l1) % MOD;
    sumL2 = (sumL2 * INV3) % MOD;

    // Разность сумм
    int ss2 = (sumR2 - sumL2 + MOD) % MOD;
    int ss1 = (sumR1 - sumL1 + MOD) % MOD;
    int ss0 = (rMod - lMod + 1 + MOD) % MOD;

    // blockSum = 2 * [ s2*ss2 - s1*ss1*(2N+1) + s0*ss0*N(N+1) ]
    int p1 = (s2 * ss2) % MOD;

    int p2 = (term2N_1 * s1) % MOD;
    p2 = (p2 * ss1) % MOD;

    int p3 = (termN_N1 * s0) % MOD;
    p3 = (p3 * ss0) % MOD;

    int blockSum = (p1 - p2 + MOD) % MOD;
    blockSum = (blockSum + p3) % MOD;
    blockSum = (blockSum * 2) % MOD;

    ans = (ans + blockSum) % MOD;

    l = r + 1;
  }

  // Коррекция для q = -2
  if (q == -2 && N >= 3) {
    int maxOdd = limit;
    if (maxOdd % 2 == 0) {
      maxOdd--;
    }

    if (maxOdd >= 1) {
      int cnt = ((maxOdd + 1) ~/ 2) % MOD;

      int sumA1 = (cnt * cnt) % MOD;

      int termSq = (4 * cnt * cnt) % MOD;
      termSq = (termSq - 1 + MOD) % MOD;
      int sumA2 = (cnt * termSq) % MOD;
      sumA2 = (sumA2 * INV3) % MOD;

      int val = sumA2;
      int subVal = (term2N_1 * sumA1) % MOD;
      val = (val - subVal + MOD) % MOD;
      int addVal = (termN_N1 * cnt) % MOD;
      val = (val + addVal) % MOD;

      int totalOdd = (val * 2) % MOD;

      totalOdd = (totalOdd - 4 + MOD) % MOD;

      ans = (ans + totalOdd) % MOD;
    }
  }

  return ans;
}
