package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

// checkLimits проверяет ограничения времени и памяти (работает только если установлена переменная окружения CHECK_LIMITS)
// Результаты выводятся в stderr, функция ничего не возвращает
func checkLimits(maxTime time.Duration, maxMemoryMB int, fn func()) {
	// Проверяем переменную окружения
	if os.Getenv("CHECK_LIMITS") == "" {
		// Если переменная не установлена, просто выполняем функцию без проверок
		fn()
		return
	}

	// Измеряем память до выполнения
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Измеряем время выполнения
	start := time.Now()
	fn()
	elapsed := time.Since(start)

	// Измеряем память после выполнения
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Вычисляем использованную память
	allocated := m2.TotalAlloc - m1.TotalAlloc
	memoryMB := float64(allocated) / (1024 * 1024)
	maxMemoryBytes := uint64(maxMemoryMB) * 1024 * 1024

	// Проверяем ограничения
	timeOk := elapsed <= maxTime
	memoryOk := allocated <= maxMemoryBytes

	// Логируем результаты
	if !timeOk || !memoryOk {
		if elapsed > maxTime {
			fmt.Fprintf(os.Stderr, "⚠️ Превышено время: %v (лимит: %v)\n", elapsed, maxTime)
		}
		if memoryMB > float64(maxMemoryMB) {
			fmt.Fprintf(os.Stderr, "⚠️ Превышена память: %.2f МБ (лимит: %d МБ)\n", memoryMB, maxMemoryMB)
		}
	} else {
		fmt.Fprintf(os.Stderr, "✓ Время: %v, Память: %.2f МБ\n", elapsed, memoryMB)
	}
}

// Ограничения и константы
// Сумма N и Q до 2*10^5. Глубина дерева 30.
// Максимальное количество узлов примерно 13 млн (с запасом для всех тестов).
// Используем массивы int32 для экономии памяти (256MB лимит).
const MAX_NODES = 13000000
const MAX_BITS = 29

var (
	// Глобальный пул узлов
	l    [MAX_NODES]int32 // индекс левого ребенка
	r    [MAX_NODES]int32 // индекс правого ребенка
	cnt  [MAX_NODES]int32 // количество уникальных чисел в поддереве
	memo [MAX_NODES]int32 // кэшированное значение xormex для поддерева
	ptr  int32            // указатель на свободное место
)

// Создание нового узла
func newNode() int32 {
	ptr++
	idx := ptr
	l[idx] = 0
	r[idx] = 0
	cnt[idx] = 0
	memo[idx] = 0
	return idx
}

// Пересчет значений в узле на основе детей (Push Up)
func pushUp(u int32, bit int) {
	idx0 := l[u]
	idx1 := r[u]

	c0, c1 := int32(0), int32(0)
	if idx0 != 0 {
		c0 = cnt[idx0]
	}
	if idx1 != 0 {
		c1 = cnt[idx1]
	}

	cnt[u] = c0 + c1

	m0, m1 := int32(0), int32(0)
	if idx0 != 0 {
		m0 = memo[idx0]
	}
	if idx1 != 0 {
		m1 = memo[idx1]
	}

	// Размер полного поддерева на текущем уровне
	full := int32(1) << bit

	// Логика вычисления Max MEX
	if c0 == full {
		// Левое поддерево полное. Если x-бит=0, мы закрываем диапазон [0, full-1] левым поддеревом
		// и прибавляем результат из правого.
		memo[u] = full + m1
	} else if c1 == full {
		// Правое поддерево полное. Если x-бит=1, правое становится левым (из-за XOR),
		// закрываем диапазон [0, full-1] и прибавляем результат из левого.
		memo[u] = full + m0
	} else {
		// Ни одно не полное. Мы не можем получить >= full.
		// Выбираем максимум из того, что дают дети.
		if m0 > m1 {
			memo[u] = m0
		} else {
			memo[u] = m1
		}
	}
}

// Обновление дерева (добавление/удаление числа)
// add=true (вставка), add=false (удаление)
func update(u int32, val int, bit int, add bool) {
	if bit < 0 {
		// Лист (бит -1, число полностью обработано)
		if add {
			cnt[u] = 1
			memo[u] = 1 // MEX множества {0} (пустой суффикс -> 0, следующий -> 1) равен 1
		} else {
			cnt[u] = 0
			memo[u] = 0
		}
		return
	}

	dir := (val >> bit) & 1
	var childIdx int32

	if dir == 0 {
		if l[u] == 0 {
			l[u] = newNode()
		}
		childIdx = l[u]
	} else {
		if r[u] == 0 {
			r[u] = newNode()
		}
		childIdx = r[u]
	}

	update(childIdx, val, bit-1, add)
	pushUp(u, bit)
}

// Оптимизация ввода-вывода
type FastReader struct {
	sc *bufio.Scanner
}

func NewFastReader(r *os.File) *FastReader {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	sc.Split(bufio.ScanWords)
	return &FastReader{sc: sc}
}

func (r *FastReader) ReadInt() int {
	r.sc.Scan()
	x, _ := strconv.Atoi(r.sc.Text())
	return x
}

func solve() {
	// Быстрый ввод-вывод
	reader := NewFastReader(os.Stdin)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	t := reader.ReadInt()

	// Сбрасываем указатель только один раз, так как массив рассчитан на сумму всех N и Q
	// Но для корректности нужно сбрасывать при каждом тесте, или не сбрасывать?
	// В предоставленном коде ptr = 0 было один раз перед циклом по тестам.
	// Так как MAX_NODES достаточно большой для всех тестов, это ок.
	ptr = 0

	for i := 0; i < t; i++ {
		n := reader.ReadInt()
		q := reader.ReadInt()

		a := make([]int, n)
		freq := make(map[int]int, n) // Карта частот для отслеживания дубликатов

		// Корень текущего дерева
		root := newNode()

		for j := 0; j < n; j++ {
			a[j] = reader.ReadInt()
			freq[a[j]]++
			// Вставляем в Trie только если это первое появление числа
			if freq[a[j]] == 1 {
				update(root, a[j], MAX_BITS, true)
			}
		}

		// Выводим начальный xormex
		writer.WriteString(strconv.Itoa(int(memo[root])))
		writer.WriteByte('\n')

		for k := 0; k < q; k++ {
			j := reader.ReadInt()
			v := reader.ReadInt()
			j-- // корректировка индекса к 0-based

			oldVal := a[j]
			if oldVal != v {
				// Удаляем старое значение
				freq[oldVal]--
				if freq[oldVal] == 0 {
					update(root, oldVal, MAX_BITS, false)
				}

				// Обновляем массив
				a[j] = v

				// Добавляем новое значение
				freq[v]++
				if freq[v] == 1 {
					update(root, v, MAX_BITS, true)
				}
			}

			// Выводим xormex после обновления
			writer.WriteString(strconv.Itoa(int(memo[root])))
			writer.WriteByte('\n')
		}
	}
}

func main() {
	checkLimits(4*time.Second, 256, solve)
}
