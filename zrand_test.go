package zrand

import "testing"

var testOp = Combine(
	Random(Uppers, 3),
	Shuffle(
		Combine(
			Random(Lowers, 3),
			Random(Numerics, 3),
			Random(Letters, 3),
		),
	),
	Random(Lowers, 3),
)

func BenchmarkZrand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildString(testOp)
	}
}

func TestZrand(t *testing.T) {
	t.Log(BuildString(testOp))
}
