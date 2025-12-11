package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const mod = 1000000007

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем n и k
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(strings.TrimSpace(line))
	n, _ := strconv.Atoi(parts[0])
	k, _ := strconv.Atoi(parts[1])

	// Читаем массив a
	line, _ = reader.ReadString('\n')
	parts = strings.Fields(strings.TrimSpace(line))
	a := make([]int, k)
	for i := 0; i < k; i++ {
		a[i], _ = strconv.Atoi(parts[i])
	}

	result := solve(n, k, a)
	writer.WriteString(fmt.Sprintf("%d\n", result))
}

// solve вычисляет количество делителей числа S = n! / (A * P) по модулю 10^9 + 7
func solve(n, k int, a []int) int {
	threshold := min(n, 1000000)

	// Вычисляем разложение n! на простые множители только для простых <= threshold
	factorialPrimes := factorizeFactorialSmall(n, threshold)

	// Вычисляем разложение A на простые множители
	aPrimes := factorizeProduct(a)

	// Вычисляем разложение S = n! / (A * P)
	sPrimes := make(map[int]int)

	// Обрабатываем простые числа <= threshold
	for prime, exp := range factorialPrimes {
		// Вычитаем разложение A
		if aExp, ok := aPrimes[prime]; ok {
			exp -= aExp
		}
		if exp > 0 {
			sPrimes[prime] = exp
		}
	}

	// Простые числа > threshold полностью уходят в P, поэтому их не учитываем в S

	// Вычисляем количество делителей S
	result := 1
	for _, exp := range sPrimes {
		result = (result * (exp + 1)) % mod
	}

	return result
}

// factorizeFactorialSmall вычисляет разложение n! на простые множители только для простых <= threshold
// Использует формулу Лежандра: v_p(n!) = sum_{i=1}^{∞} floor(n / p^i)
func factorizeFactorialSmall(n, threshold int) map[int]int {
	primes := sieve(threshold)
	result := make(map[int]int)

	// Учитываем все простые числа <= threshold
	for _, p := range primes {
		result[p] = legendre(n, p)
	}

	return result
}

// processLargePrimes обрабатывает простые числа в диапазоне [low, high] потоково
// Для каждого простого числа вызывается callback, не сохраняя все числа в памяти
func processLargePrimes(low, high int, callback func(int)) {
	if low > high {
		return
	}

	// Находим простые числа до sqrt(high) для фильтрации
	sqrtHigh := intSqrt(high)
	basePrimes := sieve(sqrtHigh)

	// Обрабатываем сегментами для экономии памяти
	// Используем меньший размер сегмента для экономии памяти
	segmentSize := 100000 // размер сегмента (100 КБ на сегмент)
	for segmentLow := low; segmentLow <= high; segmentLow += segmentSize {
		segmentHigh := min(segmentLow+segmentSize-1, high)

		// Создаем массив для сегмента
		segment := make([]bool, segmentHigh-segmentLow+1)
		for i := range segment {
			segment[i] = true
		}

		// Применяем решето для каждого простого числа
		for _, p := range basePrimes {
			// Находим первое число в сегменте, кратное p
			start := max(((segmentLow+p-1)/p)*p, p*p)
			for j := start; j <= segmentHigh; j += p {
				segment[j-segmentLow] = false
			}
		}

		// Обрабатываем простые числа из сегмента
		for i, isPrime := range segment {
			if isPrime {
				num := segmentLow + i
				if num >= 2 {
					callback(num)
				}
			}
		}
	}
}

func intSqrt(n int) int {
	if n < 2 {
		return n
	}
	left, right := 1, n
	for left < right {
		mid := (left + right + 1) / 2
		if mid*mid <= n {
			left = mid
		} else {
			right = mid - 1
		}
	}
	return left
}

// legendre вычисляет показатель простого числа p в разложении n! на простые множители
func legendre(n, p int) int {
	result := 0
	power := p
	for power <= n {
		result += n / power
		power *= p
	}
	return result
}

// factorizeProduct вычисляет разложение произведения элементов массива на простые множители
func factorizeProduct(arr []int) map[int]int {
	result := make(map[int]int)
	for _, num := range arr {
		factors := factorize(num)
		for prime, exp := range factors {
			result[prime] += exp
		}
	}
	return result
}

// factorize разлагает число на простые множители
func factorize(n int) map[int]int {
	result := make(map[int]int)
	for i := 2; i*i <= n; i++ {
		for n%i == 0 {
			result[i]++
			n /= i
		}
	}
	if n > 1 {
		result[n]++
	}
	return result
}

// sieve возвращает список простых чисел до n (решето Эратосфена)
func sieve(n int) []int {
	if n < 2 {
		return []int{}
	}

	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	isPrime[0], isPrime[1] = false, false

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	primes := []int{}
	for i := range isPrime {
		if i >= 2 && isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}
