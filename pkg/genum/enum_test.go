// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"github.com/rvflash/genum/pkg/genum"
)

func TestEnum_Format(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in        genum.Enum
			pos       int
			useIota   bool
			commented bool
			// outputs
			out string
		}{
			"Default":   {},
			"Empty":     {in: genum.Enum{Text: "_", Type: "Hi"}},
			"Untyped":   {in: genum.Enum{Text: "hi", Value: "0"}},
			"Unnamed":   {in: genum.Enum{Text: "_", Type: "Hi", Value: "0"}, out: "_ Hi = 0\n"},
			"First one": {in: genum.Enum{Text: "Hello", Type: "Hi", Value: "0"}, out: "Hello Hi = 0\n"},
			"Iota": {
				in:      genum.Enum{Text: "Hello", Type: "Hi", Value: "0", Iota: "iota + 2"},
				useIota: true,
				out:     "Hello Hi = iota + 2\n",
			},
			"Commented one": {
				in:        genum.Enum{Text: "Hello", Type: "Hi", Value: "2", Iota: "iota + 2"},
				useIota:   true,
				commented: true,
				out:       "Hello Hi = iota + 2 // e = 2\n",
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := tt.in.Format(tt.pos, tt.useIota, tt.commented)
			are.Equal("", cmp.Diff(tt.out, out))
		})
	}
}
