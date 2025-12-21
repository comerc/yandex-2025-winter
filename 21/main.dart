import 'dart:io';
import 'dart:typed_data';
import 'dart:math';

const double EPS = 1e-13;

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

  double readDouble() {
    while (_ptr < _bytes.length && _bytes[_ptr] <= 32) {
      _ptr++;
    }
    int start = _ptr;
    while (_ptr < _bytes.length && _bytes[_ptr] > 32) {
      _ptr++;
    }
    return double.parse(String.fromCharCodes(_bytes, start, _ptr));
  }
}

void solve(Scanner sc) {
  final buffer = StringBuffer();
  if (!sc.hasNext()) return;
  int t = sc.readInt();

  final solver = Solver();
  for (int i = 0; i < t; i++) {
    solver.solveTestCase(sc, buffer);
  }
  stdout.write(buffer.toString());
}

class Solver {
  // Flat arrays for points. Max N=10 -> 20 doubles.
  final Float64List points = Float64List(20);
  int n = 0;

  final Float64List compPoints = Float64List(20);
  int nComp = 0;

  final Float64List staticCands = Float64List(1000); // Max 500 candidates
  int nStaticCands = 0;

  final List<List<int>> staticCandsIndices = List.generate(10, (_) => []);

  final Float64List solution = Float64List(20);
  int nSol = 0;

  final Random _rng = Random(42);

  void solveTestCase(Scanner sc, StringBuffer buffer) {
    if (!sc.hasNext()) return;
    n = sc.readInt();

    for (int i = 0; i < n; i++) {
      points[2 * i] = sc.readDouble();
      points[2 * i + 1] = sc.readDouble();
    }

    int visited = 0;
    List<double> totalCircles = []; 

    for (int i = 0; i < n; i++) {
      if ((visited & (1 << i)) != 0) continue;

      nComp = 0;
      List<int> q = []; 
      q.add(i);
      visited |= (1 << i);

      int head = 0;
      while (head < q.length) {
        int u = q[head++];
        compPoints[2 * nComp] = points[2 * u];
        compPoints[2 * nComp + 1] = points[2 * u + 1];
        nComp++;

        for (int v = 0; v < n; v++) {
          if ((visited & (1 << v)) == 0) {
            double dx = points[2 * u] - points[2 * v];
            double dy = points[2 * u + 1] - points[2 * v + 1];
            if (dx * dx + dy * dy <= 16.0 + 1e-7) {
              visited |= (1 << v);
              q.add(v);
            }
          }
        }
      }

      if (!solveComponent(totalCircles)) {
        buffer.writeln("NO");
        return;
      }
    }

    buffer.writeln("YES");
    buffer.writeln(totalCircles.length ~/ 2);
    for (int i = 0; i < totalCircles.length; i += 2) {
      buffer.writeln("${totalCircles[i].toStringAsFixed(15)} ${totalCircles[i+1].toStringAsFixed(15)}");
    }
  }

  bool solveComponent(List<double> output) {
    // 1. K=1 Exact (MEC)
    if (getMEC(output)) return true;

    // 2. Backtracking (Deterministic)
    if (runBacktrack(false, output)) return true;

    // 3. Random fallback
    for (int i = 0; i < 3; i++) {
      if (runBacktrack(true, output)) return true;
    }

    return false;
  }

  bool runBacktrack(bool useRandom, List<double> output) {
    nStaticCands = 0;

    // Type 1: Points
    for (int i = 0; i < nComp; i++) {
      staticCands[2 * nStaticCands] = compPoints[2 * i];
      staticCands[2 * nStaticCands + 1] = compPoints[2 * i + 1];
      nStaticCands++;
    }

    // Type 2: Intersections
    for (int i = 0; i < nComp; i++) {
      for (int j = i + 1; j < nComp; j++) {
        addIntersections(
          compPoints[2 * i], compPoints[2 * i + 1], 1.0,
          compPoints[2 * j], compPoints[2 * j + 1], 1.0
        );
      }
    }

    if (useRandom) {
      for (int i = 0; i < nComp; i++) {
        for (int k = 0; k < 5; k++) {
          double angle = _rng.nextDouble() * 2 * pi;
          double r = sqrt(_rng.nextDouble()); 
          double cx = compPoints[2 * i] + r * cos(angle);
          double cy = compPoints[2 * i + 1] + r * sin(angle);
          staticCands[2 * nStaticCands] = cx;
          staticCands[2 * nStaticCands + 1] = cy;
          nStaticCands++;
        }
      }
    }

    // Pre-calculate indices
    for (int i = 0; i < nComp; i++) {
      staticCandsIndices[i].clear();
    }
    for (int idx = 0; idx < nStaticCands; idx++) {
      double cx = staticCands[2 * idx];
      double cy = staticCands[2 * idx + 1];
      for (int i = 0; i < nComp; i++) {
        double dx = cx - compPoints[2 * i];
        double dy = cy - compPoints[2 * i + 1];
        if (dx * dx + dy * dy <= 1.0 + EPS) {
          staticCandsIndices[i].add(idx);
        }
      }
    }

    nSol = 0;
    if (backtrack(0)) {
      for (int i = 0; i < nSol; i++) {
        output.add(solution[2 * i]);
        output.add(solution[2 * i + 1]);
      }
      return true;
    }
    return false;
  }
  
  void addIntersections(double x1, double y1, double r1, double x2, double y2, double r2) {
    double dx = x1 - x2;
    double dy = y1 - y2;
    double d2 = dx * dx + dy * dy;
    double d = sqrt(d2);

    if (d > r1 + r2 + EPS || d < (r1 - r2).abs() - EPS || d < 1e-9) return;

    double a = (r1 * r1 - r2 * r2 + d2) / (2 * d);
    double h = sqrt(max(0.0, r1 * r1 - a * a));

    double x0 = x1 + a * (x2 - x1) / d;
    double y0 = y1 + a * (y2 - y1) / d;

    double rx = -h * (y2 - y1) / d;
    double ry = h * (x2 - x1) / d;

    staticCands[2 * nStaticCands] = x0 - rx;
    staticCands[2 * nStaticCands + 1] = y0 - ry;
    nStaticCands++;

    staticCands[2 * nStaticCands] = x0 + rx;
    staticCands[2 * nStaticCands + 1] = y0 + ry;
    nStaticCands++;
  }

  bool backtrack(int mask) {
    if (mask == (1 << nComp) - 1) return true;

    int u = -1;
    for (int i = 0; i < nComp; i++) {
      if ((mask & (1 << i)) == 0) {
        u = i;
        break;
      }
    }

    bool tryCand(double cx, double cy) {
      double dx = cx - compPoints[2 * u];
      double dy = cy - compPoints[2 * u + 1];
      if (dx * dx + dy * dy > 1.0 + EPS) return false;

      for (int i = 0; i < nSol; i++) {
        double sx = solution[2 * i];
        double sy = solution[2 * i + 1];
        double d2 = (cx - sx) * (cx - sx) + (cy - sy) * (cy - sy);
        if (d2 < 4.0 - EPS) return false;
      }

      int newMask = mask;
      for (int i = 0; i < nComp; i++) {
        double px = compPoints[2 * i];
        double py = compPoints[2 * i + 1];
        if ((cx - px) * (cx - px) + (cy - py) * (cy - py) <= 1.0 + EPS) {
          newMask |= (1 << i);
        }
      }

      solution[2 * nSol] = cx;
      solution[2 * nSol + 1] = cy;
      nSol++;
      if (backtrack(newMask)) return true;
      nSol--;
      return false;
    }

    // 1. Static Candidates
    List<int> indices = staticCandsIndices[u];
    for (int i = 0; i < indices.length; i++) {
      int idx = indices[i];
      if (tryCand(staticCands[2 * idx], staticCands[2 * idx + 1])) return true;
    }

    // 2. Dynamic Candidates
    if (nSol > 0) {
       for (int i = 0; i < nSol; i++) {
          double sx = solution[2 * i];
          double sy = solution[2 * i + 1];

          double dux = compPoints[2 * u] - sx;
          double duy = compPoints[2 * u + 1] - sy;
          if (dux * dux + duy * duy > 9.0 + 1e-5) continue;
          
          if (tryIntersections(compPoints[2 * u], compPoints[2 * u + 1], 1.0, sx, sy, 2.0, tryCand)) return true;

          for (int k = 0; k < nComp; k++) {
             if (k == u) continue;
             double dx = compPoints[2 * k] - compPoints[2 * u];
             double dy = compPoints[2 * k + 1] - compPoints[2 * u + 1];
             if (dx * dx + dy * dy > 4.0 + 1e-5) continue;
             
             if (tryIntersections(compPoints[2 * k], compPoints[2 * k + 1], 1.0, sx, sy, 2.0, tryCand)) return true;
          }
       }
       
       for (int i = 0; i < nSol; i++) {
          double sx1 = solution[2 * i];
          double sy1 = solution[2 * i + 1];
          double dux1 = compPoints[2 * u] - sx1;
          double duy1 = compPoints[2 * u + 1] - sy1;
          if (dux1 * dux1 + duy1 * duy1 > 9.0 + 1e-5) continue;

          for (int j = i + 1; j < nSol; j++) {
             double sx2 = solution[2 * j];
             double sy2 = solution[2 * j + 1];
             double dux2 = compPoints[2 * u] - sx2;
             double duy2 = compPoints[2 * u + 1] - sy2;
             if (dux2 * dux2 + duy2 * duy2 > 9.0 + 1e-5) continue;

             if (tryIntersections(sx1, sy1, 2.0, sx2, sy2, 2.0, tryCand)) return true;
          }
       }
    }

    return false;
  }

  bool tryIntersections(double x1, double y1, double r1, double x2, double y2, double r2, bool Function(double, double) consumer) {
    double dx = x1 - x2;
    double dy = y1 - y2;
    double d2 = dx * dx + dy * dy;
    double d = sqrt(d2);

    if (d > r1 + r2 + EPS || d < (r1 - r2).abs() - EPS || d < 1e-9) return false;

    double a = (r1 * r1 - r2 * r2 + d2) / (2 * d);
    double h = sqrt(max(0.0, r1 * r1 - a * a));

    double x0 = x1 + a * (x2 - x1) / d;
    double y0 = y1 + a * (y2 - y1) / d;

    double rx = -h * (y2 - y1) / d;
    double ry = h * (x2 - x1) / d;

    if (consumer(x0 - rx, y0 - ry)) return true;
    if (consumer(x0 + rx, y0 + ry)) return true;
    return false;
  }
  
  bool getMEC(List<double> output) {
      if (nComp == 0) return false;
      if (nComp == 1) {
          output.add(compPoints[0]);
          output.add(compPoints[1]);
          return true;
      }
      
      double bestR2 = 1.0 + EPS;
      double bestX = 0, bestY = 0;
      bool found = false;
      
      void check(double cx, double cy, double r2) {
          if (r2 > bestR2) return;
          for (int i = 0; i < nComp; i++) {
              double dx = cx - compPoints[2 * i];
              double dy = cy - compPoints[2 * i + 1];
              if (dx * dx + dy * dy > r2 + EPS) return;
          }
          bestR2 = r2;
          bestX = cx;
          bestY = cy;
          found = true;
      }
      
      // Pairwise
      for (int i = 0; i < nComp; i++) {
          for (int j = i + 1; j < nComp; j++) {
              double cx = (compPoints[2 * i] + compPoints[2 * j]) / 2.0;
              double cy = (compPoints[2 * i + 1] + compPoints[2 * j + 1]) / 2.0;
              double dx = cx - compPoints[2 * i];
              double dy = cy - compPoints[2 * i + 1];
              check(cx, cy, dx * dx + dy * dy);
          }
      }
      
      // Triplet
      for (int i = 0; i < nComp; i++) {
          for (int j = i + 1; j < nComp; j++) {
              for (int k = j + 1; k < nComp; k++) {
                  double x1 = compPoints[2 * i], y1 = compPoints[2 * i + 1];
                  double x2 = compPoints[2 * j], y2 = compPoints[2 * j + 1];
                  double x3 = compPoints[2 * k], y3 = compPoints[2 * k + 1];
                  
                  double D = 2 * (x1 * (y2 - y3) + x2 * (y3 - y1) + x3 * (y1 - y2));
                  if (D.abs() < 1e-9) continue;
                  
                  double s1 = x1 * x1 + y1 * y1;
                  double s2 = x2 * x2 + y2 * y2;
                  double s3 = x3 * x3 + y3 * y3;
                  
                  double ux = (s1 * (y2 - y3) + s2 * (y3 - y1) + s3 * (y1 - y2)) / D;
                  double uy = (s1 * (x2 - x3) + s2 * (x3 - x1) + s3 * (x1 - x2)) / D;
                  
                  double dx = ux - x1;
                  double dy = uy - y1;
                  check(ux, uy, dx * dx + dy * dy);
              }
          }
      }
      
      if (found) {
          output.add(bestX);
          output.add(bestY);
          return true;
      }
      return false;
  }
}
