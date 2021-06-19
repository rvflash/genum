// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/genum/pkg/genum"
)

func TestKindNamed(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  string
			out genum.Kind
		}{
			"Default": {out: genum.Int},
			"FLOAT32": {in: "FLOAT32", out: genum.Float32},
			"string":  {in: "string", out: genum.String},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := genum.KindNamed(tt.in)
			are.Equal(out, tt.out)
		})
	}
}

func TestKind(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in genum.Kind
			// outputs
			isInteger bool
			isNumber  bool
			isSigned  bool
			bitSize   int
			name      string
		}{
			"Default": {in: genum.Int, bitSize: 64, isInteger: true, isNumber: true, isSigned: true, name: "int"},
			"Int":     {in: genum.Int, bitSize: 64, isInteger: true, isNumber: true, isSigned: true, name: "int"},
			"Int8":    {in: genum.Int8, bitSize: 8, isInteger: true, isNumber: true, isSigned: true, name: "int8"},
			"Int16":   {in: genum.Int16, bitSize: 16, isInteger: true, isNumber: true, isSigned: true, name: "int16"},
			"Int32":   {in: genum.Int32, bitSize: 32, isInteger: true, isNumber: true, isSigned: true, name: "int32"},
			"Int64":   {in: genum.Int64, bitSize: 64, isInteger: true, isNumber: true, isSigned: true, name: "int64"},
			"Uint":    {in: genum.Uint, bitSize: 64, isInteger: true, isNumber: true, name: "uint"},
			"Uint8":   {in: genum.Uint8, bitSize: 8, isInteger: true, isNumber: true, name: "uint8"},
			"Uint16":  {in: genum.Uint16, bitSize: 16, isInteger: true, isNumber: true, name: "uint16"},
			"Uint32":  {in: genum.Uint32, bitSize: 32, isInteger: true, isNumber: true, name: "uint32"},
			"Uint64":  {in: genum.Uint64, bitSize: 64, isInteger: true, isNumber: true, name: "uint64"},
			"Float32": {in: genum.Float32, bitSize: 32, isNumber: true, isSigned: true, name: "float32"},
			"Float64": {in: genum.Float64, bitSize: 64, isNumber: true, isSigned: true, name: "float64"},
			"String":  {in: genum.String, name: "string"},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			are.Equal(tt.bitSize, tt.in.BitSize())
			are.Equal(tt.isInteger, tt.in.IsInteger())
			are.Equal(tt.isNumber, tt.in.IsNumber())
			are.Equal(tt.isSigned, tt.in.IsSigned())
			are.Equal(tt.name, tt.in.Name())
		})
	}
}
