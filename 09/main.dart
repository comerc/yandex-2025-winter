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
  int k = sc.readInt();
  List<int> a = [];
  for (int i = 0; i < k; i++) {
    a.add(sc.readInt());
  }

  int result = solveFunc(n, k, a);
  buffer.writeln(result);

  stdout.write(buffer.toString());
}

int solveFunc(int n, int k, List<int> a) {
  int threshold = n < 1000000 ? n : 1000000;

  // Разложение n! для простых <= threshold
  Map<int, int> factorialPrimes = factorizeFactorialSmall(n, threshold);

  // Разложение A
  Map<int, int> aPrimes = factorizeProduct(a);

  // Разложение S = n! / A (простые > threshold игнорируем)
  Map<int, int> sPrimes = {};

  for (var entry in factorialPrimes.entries) {
    int prime = entry.key;
    int exp = entry.value;
    if (aPrimes.containsKey(prime)) {
      exp -= aPrimes[prime]!;
    }
    if (exp > 0) {
      sPrimes[prime] = exp;
    }
  }

  // Количество делителей
  int result = 1;
  for (int exp in sPrimes.values) {
    result = (result * (exp + 1)) % MOD;
  }

  return result;
}

Map<int, int> factorizeFactorialSmall(int n, int threshold) {
  List<int> primes = sieve(threshold);
  Map<int, int> result = {};
  for (int p in primes) {
    result[p] = legendre(n, p);
  }
  return result;
}

int legendre(int n, int p) {
  int result = 0;
  int power = p;
  while (power <= n) {
    result += n ~/ power;
    if (power > n ~/ p) break; // Предотвращение переполнения
    power *= p;
  }
  return result;
}

Map<int, int> factorizeProduct(List<int> arr) {
  Map<int, int> result = {};
  for (int num in arr) {
    Map<int, int> factors = factorize(num);
    for (var entry in factors.entries) {
      result[entry.key] = (result[entry.key] ?? 0) + entry.value;
    }
  }
  return result;
}

Map<int, int> factorize(int n) {
  Map<int, int> result = {};
  for (int i = 2; i * i <= n; i++) {
    while (n % i == 0) {
      result[i] = (result[i] ?? 0) + 1;
      n ~/= i;
    }
  }
  if (n > 1) {
    result[n] = (result[n] ?? 0) + 1;
  }
  return result;
}

List<int> sieve(int n) {
  if (n < 2) return [];
  List<bool> isPrime = List.filled(n + 1, true);
  isPrime[0] = false;
  isPrime[1] = false;
  for (int i = 2; i * i <= n; i++) {
    if (isPrime[i]) {
      for (int j = i * i; j <= n; j += i) {
        isPrime[j] = false;
      }
    }
  }
  List<int> primes = [];
  for (int i = 2; i <= n; i++) {
    if (isPrime[i]) {
      primes.add(i);
    }
  }
  return primes;
}
