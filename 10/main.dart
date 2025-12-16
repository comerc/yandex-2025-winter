import 'dart:io';
import 'dart:typed_data';

Future<void> main() async {
  // Читаем весь ввод асинхронно через Stream с BytesBuilder
  final bytesBuilder = BytesBuilder();
  
  // Читаем все данные из stdin через Stream
  await for (final chunk in stdin) {
    bytesBuilder.add(chunk);
  }
  
  final bytes = bytesBuilder.takeBytes();
  
  // Парсим числа напрямую из байтов без декодирования в строку
  // Это намного быстрее для больших входных данных
  final tokens = <int>[];
  
  var num = 0;
  var hasNum = false;
  
  for (var i = 0; i < bytes.length; i++) {
    final byte = bytes[i];
    if (byte >= 48 && byte <= 57) { // '0'..'9'
      num = num * 10 + (byte - 48);
      hasNum = true;
    } else if (hasNum) {
      tokens.add(num);
      num = 0;
      hasNum = false;
    }
  }
  if (hasNum) {
    tokens.add(num);
  }
  
  // Конвертируем в Int32List для экономии памяти
  final tokensList = Int32List.fromList(tokens);
  var tokenIdx = 0;
  int readInt() => tokensList[tokenIdx++];

  final output = StringBuffer();
  final t = readInt();

  for (var test = 0; test < t; test++) {
    final n = readInt();
    final m = readInt();

    // Используем один буфер для всех чисел (a, b, c)
    final data = Int32List(m * 3);
    
    for (var i = 0; i < m; i++) data[i] = readInt();           // a
    for (var i = 0; i < m; i++) data[m + i] = readInt();       // b
    for (var i = 0; i < m; i++) data[m * 2 + i] = readInt();   // c

    // Сортируем индексы по весу (c хранится в data[m*2..m*3-1])
    final indices = Int32List(m);
    for (var i = 0; i < m; i++) indices[i] = i;
    indices.sort((i, j) => data[m * 2 + i].compareTo(data[m * 2 + j]));

    // DSU
    final parent = Int32List(n);
    final rank = Int32List(n);
    final size = Int32List(n);
    for (var i = 0; i < n; i++) {
      parent[i] = i;
      size[i] = 1;
    }

    final result = Int32List(n);
    result.fillRange(0, n, -1);
    result[0] = 0;

    var maxReached = 1;

    // Find с path compression (рекурсивная версия - быстрее в Dart)
    int find(int x) {
      if (parent[x] != x) {
        parent[x] = find(parent[x]);
      }
      return parent[x];
    }

    for (final ei in indices) {
      final va = data[ei] - 1;           // a[ei]
      final vb = data[m + ei] - 1;        // b[ei]
      final w = data[m * 2 + ei];         // c[ei]

      if (va == vb) continue;

      final rootA = find(va);
      final rootB = find(vb);

      if (rootA != rootB) {
        final sizeA = size[rootA];
        final sizeB = size[rootB];
        final newSize = sizeA + sizeB;

        // Union by rank
        if (rank[rootA] < rank[rootB]) {
          parent[rootA] = rootB;
          size[rootB] = newSize;
        } else {
          parent[rootB] = rootA;
          size[rootA] = newSize;
          if (rank[rootA] == rank[rootB]) {
            rank[rootA]++;
          }
        }

        // Заполняем результат только для новых размеров
        if (newSize > maxReached) {
          result.fillRange(maxReached, newSize, w);
          maxReached = newSize;
          if (maxReached == n) break;
        }
      }
    }

    // Вывод
    output.write(result[0]);
    for (var i = 1; i < n; i++) {
      output.write(' ');
      output.write(result[i]);
    }
    output.writeln();
  }

  // Выводим весь результат сразу
  stdout.write(output);
}
