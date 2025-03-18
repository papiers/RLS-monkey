package benchmark

import "testing"

func Test_benchmark(t *testing.T) {
	tests := []struct {
		name   string
		engine string
	}{
		{
			name:   "vm engine",
			engine: "vm",
		},
		{
			name:   "eval engine",
			engine: "eval",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			benchmark(tt.engine)
		})
	}
}
