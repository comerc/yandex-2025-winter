package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	line, _ := reader.ReadString('\n')
	parts := strings.Fields(strings.TrimSpace(line))
	R, _ := strconv.Atoi(parts[0])
	B, _ := strconv.Atoi(parts[1])

	W, H := solve(R, B)
	writer.WriteString(fmt.Sprintf("%d %d\n", W, H))
}

// solve находит размеры панели W и H (W >= H) по количеству красных R и синих B плиток
// R = 2*W + 2*H - 4, B = (W-2) * (H-2)
func solve(R, B int) (int, int) {
	sum := (R + 4) / 2

	for d := 1; d*d <= B; d++ {
		if B%d != 0 {
			continue
		}

		// Проверяем оба варианта: (W-2, H-2) = (d, B/d) и (B/d, d)
		for _, w := range []int{B/d + 2, d + 2} {
			h := sum - w
			if (w-2)*(h-2) == B && w >= h {
				return w, h
			}
		}
	}

	return 0, 0
}
