package zrand

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImmediate(t *testing.T) {
	i := Immediate("test")
	buf := Build(i)
	require.Equal(t, []byte("test"), buf)
	buf[0] = 'b'
	require.Equal(t, []byte("test"), Build(i))
}

func TestCombine(t *testing.T) {
	c := opCombine{
		ops: []Op{
			Immediate("test"),
			Immediate("test"),
		},
	}
	buf := Build(c)
	require.Equal(t, []byte("testtest"), buf)
	buf[0] = 'b'
	require.Equal(t, []byte("testtest"), Build(c))
}

func TestRandom(t *testing.T) {
	r := opRandom{
		src:  Immediate("test"),
		size: 3,
	}
	buf := Build(r)
	require.Len(t, buf, 3)
	for _, c := range buf {
		require.Contains(t, "test", string(c))
	}
}

func TestShuffle(t *testing.T) {
	s := opShuffle{
		src: opCombine{
			ops: []Op{
				Immediate("test"),
				Immediate("test"),
			},
		},
	}
	buf := Build(s)
	buf[0] = 'b'
	buf = Build(s)

	var strs []string
	for _, c := range buf {
		strs = append(strs, string(c))
	}

	sort.Strings(strs)

	require.Equal(t, "eesstttt", strings.Join(strs, ""))
}
