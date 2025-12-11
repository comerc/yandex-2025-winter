package main

import (
	"testing"
)

// Список артефактов для тестов
var testArtifacts = []int64{
	1, 2, 4, 6, 12, 24, 36, 48, 60, 120,
	180, 240, 360, 720, 840, 1260, 1680, 2520, 5040, 7560,
	10080, 15120, 20160, 25200, 27720, 45360, 50400, 55440, 83160, 110880,
	166320, 221760, 277200, 332640, 498960, 554400, 665280, 720720, 1081080, 1441440,
	2162160, 2882880, 3603600, 4324320, 6486480, 7207200, 8648640, 10810800, 14414400, 17297280,
	21621600, 32432400, 36756720, 43243200, 61261200, 73513440, 110270160, 122522400, 147026880, 183783600,
	245044800, 294053760, 367567200, 551350800, 698377680, 735134400, 1102701600, 1396755360, 2095133040, 2205403200,
	2327925600, 2793510720, 3491888400, 4655851200, 5587021440, 6983776800, 10475665200, 13967553600, 20951330400, 27935107200,
	41902660800, 48886437600, 64250746560, 73329656400, 80313433200, 97772875200, 128501493120, 146659312800, 160626866400, 240940299600,
	293318625600, 321253732800, 481880599200, 642507465600, 963761198400, 1124388064800, 1606268664000, 1686582097200, 1927522396800, 2248776129600,
	3212537328000, 3373164194400, 4497552259200, 6746328388800, 8995104518400, 9316358251200, 13492656777600, 18632716502400, 26985313555200, 27949074753600,
	32607253879200, 46581791256000, 48910880818800, 55898149507200, 65214507758400, 93163582512000, 97821761637600, 130429015516800, 195643523275200, 260858031033600,
	288807105787200, 391287046550400, 577614211574400, 782574093100800, 866421317361600, 1010824870255200, 1444035528936000, 1516237305382800, 1732842634723200, 2021649740510400,
	2888071057872000, 3032474610765600, 4043299481020800, 6064949221531200, 8086598962041600, 10108248702552000, 12129898443062400, 18194847664593600, 20216497405104000, 24259796886124800,
	30324746107656000, 36389695329187200, 48519593772249600, 60649492215312000, 72779390658374400, 74801040398884800, 106858629141264000, 112201560598327200, 149602080797769600, 224403121196654400,
	299204161595539200, 374005201994424000, 448806242393308800, 673209363589963200, 748010403988848000, 897612484786617600,
}

func TestExamples(t *testing.T) {
	tests := []struct {
		l, r     int64
		expected int
	}{
		{1, 1, 1},
		{5, 10, 1},
		{6, 8, 1},
	}

	for _, tt := range tests {
		result := countArtifacts(testArtifacts, tt.l, tt.r)
		if result != tt.expected {
			t.Errorf("countArtifacts(%d, %d) = %d, want %d", tt.l, tt.r, result, tt.expected)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		l, r     int64
		expected int
	}{
		{"single artifact 1", 1, 1, 1},
		{"single artifact 2", 2, 2, 1},
		{"single artifact 6", 6, 6, 1},
		{"no artifacts", 3, 3, 0},
		{"no artifacts range", 7, 11, 0},
		{"first 10 artifacts", 1, 120, 10},
		{"all artifacts", 1, 1000000000000000000, 156},
		{"empty range at start", 0, 0, 0},
		{"range with 5 (not artifact)", 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countArtifacts(testArtifacts, tt.l, tt.r)
			if result != tt.expected {
				t.Errorf("countArtifacts(%d, %d) = %d, want %d", tt.l, tt.r, result, tt.expected)
			}
		})
	}
}

func TestSpecificArtifacts(t *testing.T) {
	// Проверяем, что конкретные числа являются артефактами
	knownArtifacts := []int64{1, 2, 4, 6, 12, 24, 60, 120, 720, 5040}
	for _, a := range knownArtifacts {
		result := countArtifacts(testArtifacts, a, a)
		if result != 1 {
			t.Errorf("countArtifacts(%d, %d) = %d, want 1 (should be artifact)", a, a, result)
		}
	}

	// Проверяем, что некоторые числа НЕ являются артефактами
	notArtifacts := []int64{3, 5, 7, 8, 9, 10, 11, 100, 1000}
	for _, a := range notArtifacts {
		result := countArtifacts(testArtifacts, a, a)
		if result != 0 {
			t.Errorf("countArtifacts(%d, %d) = %d, want 0 (should NOT be artifact)", a, a, result)
		}
	}
}

func TestLargeRanges(t *testing.T) {
	tests := []struct {
		name string
		l, r int64
	}{
		{"full range", 1, 1000000000000000000},
		{"large start", 100000000000000000, 1000000000000000000},
		{"middle range", 1000000000, 10000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countArtifacts(testArtifacts, tt.l, tt.r)
			if result < 0 {
				t.Errorf("countArtifacts(%d, %d) = %d, should be non-negative", tt.l, tt.r, result)
			}
		})
	}
}

func TestArtifactsCount(t *testing.T) {
	// Проверяем общее количество артефактов
	if len(testArtifacts) != 156 {
		t.Errorf("Expected 156 artifacts, got %d", len(testArtifacts))
	}

	// Проверяем, что артефакты отсортированы
	for i := 1; i < len(testArtifacts); i++ {
		if testArtifacts[i] <= testArtifacts[i-1] {
			t.Errorf("Artifacts not sorted: artifacts[%d]=%d >= artifacts[%d]=%d",
				i-1, testArtifacts[i-1], i, testArtifacts[i])
		}
	}
}

func BenchmarkCountArtifacts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countArtifacts(testArtifacts, 1, 1000000000000000000)
	}
}

func BenchmarkCountArtifactsSmallRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countArtifacts(testArtifacts, 1, 100)
	}
}

func BenchmarkCountArtifactsManyQueries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50000; j++ {
			countArtifacts(testArtifacts, 1, 1000000000000000000)
		}
	}
}
