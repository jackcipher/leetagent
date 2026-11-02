package ex01_cosine

import (
	"errors"
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestCosine(t *testing.T) {
	tests := []struct {
		name    string
		a, b    []float64
		want    float64
		wantErr error
	}{
		{"same", []float64{1, 2, 3}, []float64{2, 4, 6}, 1, nil},
		{"orthogonal", []float64{1, 0}, []float64{0, 5}, 0, nil},
		{"opposite", []float64{1, 1}, []float64{-1, -1}, -1, nil},
		{"45deg", []float64{1, 0}, []float64{1, 1}, 1 / math.Sqrt2, nil},
		{"mismatch", []float64{1}, []float64{1, 2}, 0, ErrDimMismatch},
		{"zero", []float64{0, 0}, []float64{1, 2}, 0, ErrZeroVector},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Cosine(tc.a, tc.b)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v", err)
				}
				return
			}
			if err != nil || !near(got, tc.want) {
				t.Fatalf("got %v %v want %v", got, err, tc.want)
			}
		})
	}
}

func TestDotNormNormalize(t *testing.T) {
	if d, _ := Dot([]float64{1, 2}, []float64{3, 4}); d != 11 {
		t.Errorf("Dot = %v", d)
	}
	if n := Norm([]float64{3, 4}); n != 5 {
		t.Errorf("Norm = %v", n)
	}
	in := []float64{3, 4}
	u, err := Normalize(in)
	if err != nil || !near(u[0], 0.6) || !near(u[1], 0.8) || in[0] != 3 {
		t.Errorf("Normalize = %v %v (input %v)", u, err, in)
	}
	if _, err := Normalize([]float64{0}); !errors.Is(err, ErrZeroVector) {
		t.Errorf("zero err = %v", err)
	}
}

func TestParseEmbeddingResponse(t *testing.T) {
	body := []byte(`{"model":"bge","data":[{"index":1,"embedding":[0,1]},{"index":0,"embedding":[1,0]}]}`)
	got, err := ParseEmbeddingResponse(body)
	if err != nil || len(got) != 2 || got[0][0] != 1 || got[1][1] != 1 {
		t.Fatalf("got %v %v", got, err)
	}
	for _, bad := range []string{`{"data":[]}`, `not json`, `{"data":[{"index":5,"embedding":[1]}]}`, `{"data":[{"index":0,"embedding":[1]},{"index":0,"embedding":[2]}]}`} {
		if _, err := ParseEmbeddingResponse([]byte(bad)); !errors.Is(err, ErrBadResponse) {
			t.Errorf("%s: err = %v", bad, err)
		}
	}
}
