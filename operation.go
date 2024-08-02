package zrand

import (
	"errors"
)

// Op is an operation that generates bytes with given buffer and random.
type Op interface {
	// Len returns the length of the output.
	Len() int
	// BufferRequired returns the length of the buffer required.
	BufferRequired() int
	// RandomRequired returns the length of the random required.
	RandomRequired() int
	// Build generates the output with the given buffer and random.
	Build(out []byte, buf []byte, rnd []byte)
}

// Immediate create a new operation that writes the given bytes to the output.
type Immediate []byte

var (
	_ Op = Immediate{}
)

func (i Immediate) Len() int {
	return len(i)
}

func (i Immediate) BufferRequired() int {
	return 0
}

func (i Immediate) RandomRequired() int {
	return 0
}

func (i Immediate) Build(out, _, _ []byte) {
	copy(out, i)
}

// Combine create a new operation that combines multiple operations.
func Combine(ops ...Op) Op {
	return opCombine{ops: ops}
}

type opCombine struct {
	ops []Op
}

func (c opCombine) Len() (l int) {
	for _, o := range c.ops {
		l += o.Len()
	}
	return
}

func (c opCombine) BufferRequired() (b int) {
	for _, o := range c.ops {
		b += o.BufferRequired()
	}
	return
}

func (c opCombine) RandomRequired() (r int) {
	for _, o := range c.ops {
		r += o.RandomRequired()
	}
	return
}

func (c opCombine) Build(out, buf, rnd []byte) {
	for _, o := range c.ops {
		o.Build(out, buf, rnd)
		out = out[o.Len():]
		buf = buf[o.BufferRequired():]
		rnd = rnd[o.RandomRequired():]
	}
}

// Random create a new operation that generates random bytes from source.
func Random(src Op, size int) Op {
	return opRandom{src: src, size: size}
}

type opRandom struct {
	src  Op
	size int
}

func (r opRandom) Len() int {
	return r.size
}

func (r opRandom) BufferRequired() int {
	return r.src.Len() + r.src.BufferRequired()
}

func (r opRandom) RandomRequired() int {
	return r.size + r.src.RandomRequired()
}

func (r opRandom) Build(out, buf, rnd []byte) {
	src := buf[:r.src.Len()]
	buf = buf[:len(src)]
	r.src.Build(src, buf, rnd)

	if len(src) > 256 {
		panic(errors.New("random source is too large"))
	}

	rnd = rnd[r.src.RandomRequired():]

	for i := 0; i < r.size; i++ {
		out[i] = src[rnd[i]%byte(len(src))]
	}
}

// Shuffle create a new operation that shuffles the output of the source operation.
func Shuffle(o Op) Op {
	return opShuffle{src: o}
}

type opShuffle struct {
	src Op
}

func (s opShuffle) Len() int {
	return s.src.Len()
}

func (s opShuffle) BufferRequired() int {
	return s.src.BufferRequired()
}

func (s opShuffle) RandomRequired() int {
	return s.src.Len() + s.src.RandomRequired()
}

func (s opShuffle) Build(out, buf, rnd []byte) {
	s.src.Build(out, buf, rnd)

	rnd = rnd[s.src.RandomRequired():]

	for i := s.src.Len() - 1; i > 0; i-- {
		j := rnd[i] % byte(i)
		out[i], out[j] = out[j], out[i]
	}
}
