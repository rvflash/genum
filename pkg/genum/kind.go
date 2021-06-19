// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

import (
	"fmt"
	"strings"
)

const (
	base10 = 10
	bits8  = 8
	bits16 = 16
	bits32 = 32
	bits64 = 64
)

//go:generate stringer -type=Kind

// KindNamed converts s to a Kind.
func KindNamed(s string) Kind {
	s = strings.Title(strings.ToLower(s))
	for i := 0; i < len(_Kind_index)-1; i++ {
		if _Kind_name[_Kind_index[i]:_Kind_index[i+1]] == s {
			return Kind(i)
		}
	}
	return Int
}

// Kind represents a golang type.
type Kind uint8

// List of supported kinds.
const (
	Int Kind = iota
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	String
)

// BitSize returns the bit size of the type.
func (k Kind) BitSize() int {
	switch k {
	case Int8, Uint8:
		return bits8
	case Int16, Uint16:
		return bits16
	case Int32, Uint32, Float32:
		return bits32
	case Int, Int64, Uint, Uint64, Float64:
		return bits64
	default:
		return 0
	}
}

// Cast wraps the data with underlying data type.
func (k Kind) Cast(data string) string {
	return fmt.Sprintf("%s(%s)", k.Name(), data)
}

// IsNumber returns true if the type is a number (float or integer).
func (k Kind) IsNumber() bool {
	return k != String
}

// IsInteger returns true if the type is an integer.
func (k Kind) IsInteger() bool {
	switch k {
	case Float32, Float64, String:
		return false
	default:
		return true
	}
}

// IsSigned returns true if the type is a signed number.
func (k Kind) IsSigned() bool {
	switch k {
	case Int, Int8, Int16, Int32, Int64, Float32, Float64:
		return true
	default:
		return false
	}
}

// Name returns the Kind name.
func (k Kind) Name() string {
	return strings.ToLower(k.String())
}

// ValueFormat returns the formats used to output a Kind value.
func (k Kind) ValueFormat() string {
	verb := func() string {
		switch {
		case k.IsInteger():
			return "d"
		case k.IsNumber():
			return "f"
		default:
			return "q"
		}
	}
	return fmt.Sprintf("%%[%d]%s", ValuePos, verb())
}
