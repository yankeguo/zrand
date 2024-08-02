package zrand

import "crypto/rand"

var (
	Lowers   = Immediate([]byte("abcdefghijklmnopqrstuvwxyz"))
	Uppers   = Immediate([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	Letters  = Combine(Lowers, Uppers)
	Numerics = Immediate([]byte("0123456789"))
)

// Build generates bytes from the given Op.
func Build(o Op) []byte {
	out := make([]byte, o.Len())
	buf := make([]byte, o.BufferRequired())

	var rnd []byte

	if r := o.RandomRequired(); r > 0 {
		rnd = make([]byte, r)
		if _, err := rand.Read(rnd); err != nil {
			panic(err)
		}
	}

	o.Build(out, buf, rnd)

	return out
}

// BuildString generates a string from the given Op.
func BuildString(o Op) string {
	return string(Build(o))
}
