// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

import (
	"strings"
)

// Enum represents an Enum.
type Enum struct {
	Iota    string
	Kind    Kind
	Text    string
	RawText string
	Type    string
	Value   string
}

// Format formats the constant regarding to its context (iota, position, etc.)
//
// Enum Kind = iota // e = 0
// Enum
// Enum Kind = iota + 100
// _ Kind = iota
// _
// Enum Kind = "rv"
// Non-exhaustive list.
func (e Enum) Format(pos int, useIota, commented bool) string {
	if e.Text == "" || e.Type == "" || e.Value == "" {
		return ""
	}
	p := []string{e.Text}
	if useIota {
		if e.Iota != "" {
			p = append(p, e.Type, "=", e.Iota)
		}
	} else if e.Text != unnamed || pos == 0 {
		p = append(p, e.Type, "=", e.Value)
	}
	if commented {
		p = append(p, []string{"//", "e", "=", e.Value}...)
	}
	return strings.Join(p, " ") + "\n"
}
