package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const mod = 1000000007

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем M и N
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(strings.TrimSpace(line))
	M, _ := strconv.ParseInt(parts[0], 10, 64)
	N, _ := strconv.Atoi(parts[1])

	// Читаем потребности групп
	W := make([]int64, N)
	totalNeed := int64(0)

	// Читаем числа, обрабатывая случай, когда они могут быть в нескольких строках
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for i := 0; i < N && scanner.Scan(); i++ {
		val, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		W[i] = val
		totalNeed += val
	}

	result := solve(M, W, totalNeed)
	writer.WriteString(fmt.Sprintf("%d\n", result))
}

// solve находит минимальную сумму квадратов недостачи
// Алгоритм: для минимизации суммы квадратов при фиксированной сумме
// нужно распределить недостачу максимально равномерно
func solve(M int64, W []int64, totalNeed int64) int {
	N := len(W)
	deficit := totalNeed - M

	if deficit == 0 {
		return 0
	}

	// Создаем пары (потребность, индекс) для сортировки
	type pair struct {
		w     int64
		index int
	}
	pairs := make([]pair, N)
	for i := 0; i < N; i++ {
		pairs[i] = pair{w: W[i], index: i}
	}

	// Сортируем по убыванию потребности
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].w > pairs[j].w
	})

	// Распределяем недостачу оптимальным образом
	// Для минимизации суммы квадратов нужно распределить недостачу максимально равномерно
	// с учётом ограничения: недостача группы не может превышать её потребность
	origN := N
	shortfall := make([]int64, origN)
	remaining := deficit

	// Проходим по группам (отсортированным по убыванию потребности)
	// и распределяем недостачу
	// Ключевая идея: если avgShortfall <= minW среди оставшихся групп,
	// то можем распределить равномерно. Иначе, группа с минимальной потребностью
	// получает максимальную недостачу (равную своей потребности).
	for remaining > 0 && N > 0 {
		avgShortfall := remaining / int64(N)
		minW := pairs[N-1].w // минимальная потребность (pairs отсортирован по убыванию)

		// Если средняя недостача <= минимальной потребности среди оставшихся,
		// то все группы могут принять эту недостачу
		if avgShortfall <= minW {
			// Распределяем равномерно среди оставшихся групп
			baseShortfall := remaining / int64(N)
			extra := remaining % int64(N)

			// Распределяем базовую недостачу
			for j := 0; j < N; j++ {
				shortfall[pairs[j].index] = baseShortfall
			}

			// Распределяем остаток по одной единице
			for j := 0; j < int(extra); j++ {
				shortfall[pairs[j].index]++
			}
			break
		} else {
			// Группа с минимальной потребностью не может принять среднюю недостачу
			// Даем ей максимально возможную недостачу (равную потребности)
			shortfall[pairs[N-1].index] = pairs[N-1].w
			remaining -= pairs[N-1].w
			N-- // удаляем эту группу из рассмотрения
		}
	}

	N = origN // восстанавливаем для подсчёта результата

	// Вычисляем сумму квадратов недостачи по модулю
	result := int64(0)
	for i := 0; i < N; i++ {
		sq := shortfall[i] % mod
		result = (result + sq*sq%mod) % mod
	}

	return int(result)
}
