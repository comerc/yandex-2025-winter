package main

import (
	"testing"
)

func TestExample1(t *testing.T) {
	markersX := []int{0}
	markersY := []int{-10}
	commands := "ENE"
	expected := []int64{11, 12, 13}

	results := solve(markersX, markersY, commands)
	if len(results) != len(expected) {
		t.Fatalf("Expected %d outputs, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Output %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

func TestExample2(t *testing.T) {
	markersX := []int{0, 1, 1}
	markersY := []int{0, 1, -1}
	commands := "NESSW"
	expected := []int64{5, 4, 3, 4, 5}

	results := solve(markersX, markersY, commands)
	if len(results) != len(expected) {
		t.Fatalf("Expected %d outputs, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Output %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

func TestExample3(t *testing.T) {
	markersX := []int{-1, -1, 0, 2, 2}
	markersY := []int{-1, 0, -1, -1, 1}
	commands := "NSWNNSSENN"
	expected := []int64{13, 10, 11, 14, 19, 14, 11, 10, 13, 18}

	results := solve(markersX, markersY, commands)
	if len(results) != len(expected) {
		t.Fatalf("Expected %d outputs, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Output %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

func TestEdgeCases(t *testing.T) {
	// Одна метка, одна команда
	markersX1 := []int{0}
	markersY1 := []int{0}
	commands1 := "N"
	expected1 := []int64{1}

	results1 := solve(markersX1, markersY1, commands1)
	if len(results1) != len(expected1) {
		t.Fatalf("Expected %d outputs, got %d", len(expected1), len(results1))
	}
	if results1[0] != expected1[0] {
		t.Errorf("Expected %d, got %d", expected1[0], results1[0])
	}

	// Метка в начале координат, движение в разные стороны
	markersX2 := []int{0}
	markersY2 := []int{0}
	commands2 := "NSEW"
	expected2 := []int64{1, 0, 1, 0}

	results2 := solve(markersX2, markersY2, commands2)
	if len(results2) != len(expected2) {
		t.Fatalf("Expected %d outputs, got %d", len(expected2), len(results2))
	}

	for i, exp := range expected2 {
		if results2[i] != exp {
			t.Errorf("Output %d: expected %d, got %d", i, exp, results2[i])
		}
	}
}

func TestLargeCoordinates(t *testing.T) {
	// Тест с большими координатами
	markersX := []int{1000000, -1000000}
	markersY := []int{1000000, -1000000}
	commands := "EW"
	expected := []int64{4000000, 4000000}

	results := solve(markersX, markersY, commands)
	if len(results) != len(expected) {
		t.Fatalf("Expected %d outputs, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Output %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

// Тест производительности
func TestPerformance(t *testing.T) {
	// Симулируем N = 100000 меток, M = 300000 команд
	N := 100000
	markersX := make([]int, N)
	markersY := make([]int, N)
	for i := 0; i < N; i++ {
		markersX[i] = (i*1234567)%2000001 - 1000000
		markersY[i] = (i*7654321)%2000001 - 1000000
	}

	commands := ""
	cmds := []byte{'N', 'S', 'E', 'W'}
	for i := 0; i < 300000; i++ {
		commands += string(cmds[i%4])
	}

	results := solve(markersX, markersY, commands)
	if len(results) != 300000 {
		t.Errorf("Expected 300000 results, got %d", len(results))
	}
}

// Бенчмарк для одного запроса
func BenchmarkSolve(b *testing.B) {
	N := 100000
	markersX := make([]int, N)
	markersY := make([]int, N)
	for i := 0; i < N; i++ {
		markersX[i] = (i*1234567)%2000001 - 1000000
		markersY[i] = (i*7654321)%2000001 - 1000000
	}

	commands := "E" // Одна команда

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = solve(markersX, markersY, commands)
	}
}

// Бенчмарк для полного решения
func BenchmarkFullSolution(b *testing.B) {
	N := 100000
	markersX := make([]int, N)
	markersY := make([]int, N)
	for i := 0; i < N; i++ {
		markersX[i] = (i*1234567)%2000001 - 1000000
		markersY[i] = (i*7654321)%2000001 - 1000000
	}

	commands := ""
	cmds := []byte{'N', 'S', 'E', 'W'}
	for i := 0; i < 300000; i++ {
		commands += string(cmds[i%4])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = solve(markersX, markersY, commands)
	}
}

// Бенчмарк для множественных запросов
func BenchmarkMultipleQueries(b *testing.B) {
	N := 100000
	markersX := make([]int, N)
	markersY := make([]int, N)
	for i := 0; i < N; i++ {
		markersX[i] = (i*1234567)%2000001 - 1000000
		markersY[i] = (i*7654321)%2000001 - 1000000
	}

	commands := "N"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = solve(markersX, markersY, commands)
		}
	}
}
