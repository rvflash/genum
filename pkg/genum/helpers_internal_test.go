// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

import (
	"testing"

	"github.com/matryer/is"
)

func TestSumNumbers(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			sub      bool
			baseSign bool
			numSign  bool
			num      uint64
			base     uint64
			// outputs
			out  uint64
			sign bool
		}{
			"0+0":     {},
			"-1+1":    {base: 1, baseSign: true, num: 1},
			"11+-2":   {base: 11, num: 2, numSign: true, out: 9},
			"11+-20":  {base: 11, num: 20, numSign: true, out: 9, sign: true},
			"-20+21":  {base: 20, baseSign: true, num: 21, out: 1},
			"20+21":   {base: 20, num: 21, out: 41},
			"20-21":   {sub: true, base: 20, num: 21, out: 1, sign: true},
			"20-19":   {sub: true, base: 20, num: 19, out: 1},
			"-20-19":  {sub: true, base: 20, baseSign: true, num: 19, out: 39, sign: true},
			"-20--19": {sub: true, base: 20, baseSign: true, num: 19, numSign: true, out: 1, sign: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, sign := sumNumbers(tt.sub, tt.base, tt.baseSign, tt.num, tt.numSign)
			are.Equal(tt.out, out)   // mismatch out
			are.Equal(tt.sign, sign) // mismatch sign
		})
	}
}
