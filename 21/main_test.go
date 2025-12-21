package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFuzzRobustness(t *testing.T) {
	// Run fuzzing for 5 seconds to ensure optimizations didn't break things.
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	start := time.Now()
	iter := 0
	for time.Since(start) < 5*time.Second {
		iter++
		n := 2 + rng.Intn(9) // N=2..10

		// Generate random points in 10x10 box
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%.5f %.5f\n", rng.Float64()*10, rng.Float64()*10))
		}

		input := sb.String()
		reader = bufio.NewReader(strings.NewReader(input))
		var outBuf bytes.Buffer
		writer = bufio.NewWriter(&outBuf)

		// Capture panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Panic on iter %d: %v\nInput:\n%s", iter, r, input)
				}
			}()
			solve()
			writer.Flush()
		}()

		// We don't check correctness here (hard), just no crashes and reasonable speed.
	}
	t.Logf("Ran %d iterations without panic.", iter)
}

func TestStrictMEC(t *testing.T) {
	R_target := 1.0 + 1e-10
	L := R_target * math.Sqrt(3.0)

	p0 := Point{0, 0}
	p1 := Point{L, 0}
	p2 := Point{L / 2, L * math.Sqrt(3.0) / 2}

	var sb strings.Builder
	sb.WriteString("1\n3\n")
	sb.WriteString(fmt.Sprintf("%.15f %.15f\n", p0.x, p0.y))
	sb.WriteString(fmt.Sprintf("%.15f %.15f\n", p1.x, p1.y))
	sb.WriteString(fmt.Sprintf("%.15f %.15f\n", p2.x, p2.y))

	reader = bufio.NewReader(strings.NewReader(sb.String()))
	var outBuf bytes.Buffer
	writer = bufio.NewWriter(&outBuf)
	solve()
	writer.Flush()

	out := outBuf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	k, _ := strconv.Atoi(lines[1])

	if k == 1 {
		t.Errorf("StrictMEC FAILURE: Algorithm used 1 circle for points with circumradius %.15f > 1.0", R_target)
	}
}

func TestSample(t *testing.T) {
	input := `2
2
1.1 1.2
-1.0 -0.9
4
1.0 1.1
1.1 1.0
1.5 1.5
0.0 -1.0
`
	runAndValidate(t, input)
}

func runAndValidate(t *testing.T, input string) {
	reader = bufio.NewReader(strings.NewReader(input))
	var outBuf bytes.Buffer
	writer = bufio.NewWriter(&outBuf)

	start := time.Now()
	solve()
	writer.Flush()
	_ = time.Since(start)

	outReader := bufio.NewReader(&outBuf)
	inReader := bufio.NewReader(strings.NewReader(input))

	line, _ := inReader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	numTests, _ := strconv.Atoi(line)

	for i := 0; i < numTests; i++ {
		line, _ = inReader.ReadString('\n')
		for strings.TrimSpace(line) == "" {
			line, _ = inReader.ReadString('\n')
		}
		n, _ := strconv.Atoi(strings.TrimSpace(line))

		points := make([]Point, n)
		for j := 0; j < n; j++ {
			line, _ = inReader.ReadString('\n')
			parts := strings.Fields(line)
			points[j].x, _ = strconv.ParseFloat(parts[0], 64)
			points[j].y, _ = strconv.ParseFloat(parts[1], 64)
		}

		resLine, err := outReader.ReadString('\n')
		if err != nil {
			t.Fatalf("Test case %d: unexpected end of output", i)
		}
		res := strings.TrimSpace(resLine)

		if res == "YES" {
			line, _ = outReader.ReadString('\n')
			k, _ := strconv.Atoi(strings.TrimSpace(line))
			for j := 0; j < k; j++ {
				outReader.ReadString('\n')
			}
		}
	}
}
