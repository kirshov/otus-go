package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	type item struct {
		compare string
		offset  int64
		limit   int64
	}
	l := []item{
		{compare: "testdata/out_offset0_limit0.txt"},
		{compare: "testdata/out_offset0_limit10.txt", limit: 10},
		{compare: "testdata/out_offset0_limit1000.txt", limit: 1000},
		{compare: "testdata/out_offset0_limit10000.txt", limit: 10000},
		{compare: "testdata/out_offset100_limit1000.txt", limit: 1000, offset: 100},
		{compare: "testdata/out_offset6000_limit1000.txt", limit: 1000, offset: 6000},
	}

	for _, v := range l {
		t.Run(v.compare, func(t *testing.T) {
			err := Copy("testdata/input.txt", "testdata/out.txt", v.offset, v.limit)
			require.NoError(t, err)

			sf, err := os.ReadFile("testdata/out.txt")
			require.NoError(t, err)

			tf, err := os.ReadFile(v.compare)
			require.NoError(t, err)
			require.Equal(t, sf, tf)

			_ = os.Remove("testdata/out.txt")
		})
	}
}
