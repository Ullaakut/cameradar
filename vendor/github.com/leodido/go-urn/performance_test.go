package urn

import (
	"fmt"
	"testing"
)

var benchs = []testCase{
	tests[14],
	tests[2],
	tests[6],
	tests[10],
	tests[11],
	tests[13],
	tests[20],
	tests[23],
	tests[33],
	tests[45],
	tests[47],
	tests[48],
	tests[50],
	tests[52],
	tests[53],
	tests[57],
	tests[62],
	tests[63],
	tests[67],
	tests[60],
}

// This is here to avoid compiler optimizations that
// could remove the actual call we are benchmarking
// during benchmarks
var benchParseResult *URN

func BenchmarkParse(b *testing.B) {
	for ii, tt := range benchs {
		tt := tt
		outcome := (map[bool]string{true: "ok", false: "no"})[tt.ok]
		b.Run(
			fmt.Sprintf("%s/%02d/%s/", outcome, ii, rxpad(string(tt.in), 45)),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					benchParseResult, _ = Parse(tt.in)
				}
			},
		)
	}
}
