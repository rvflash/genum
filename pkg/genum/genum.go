// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package genum provides methods used to generate enum golang type based on CVS data.
package genum

import (
	"fmt"
	"io"
)

// Pos represents a data position.
type Pos int

// Positions of each data in the String formater.
const (
	NamePos Pos = iota + 1
	ValuePos
	TypePos
)

// String implements the fmt.Stringer interface.
func (p Pos) String() string {
	switch p {
	case NamePos:
		return "name"
	case ValuePos:
		return "value"
	case TypePos:
		return "kind"
	default:
		return "unknown"
	}
}

// Command is the name of the command line to generate golang enums based on CSV data.
const Command = "genum"

// NameFormat returns the format used to return the enum name.
func NameFormat() string {
	return fmt.Sprintf("%%[%d]s", NamePos)
}

// TypeFormat returns the format used to return the enum type.
func TypeFormat() string {
	return fmt.Sprintf("%%[%d]s", TypePos)
}

// DefaultFormat returns the format used to return the enum string.
func DefaultFormat(valueFormat string) string {
	return TypeFormat() + "(" + valueFormat + ")"
}

// Settings must be implemented by any service wanted to share generation configuration.
type Settings interface {
	DstFilename() string
	SrcFile() io.Reader
	PackageName() string
	EnumTypeName() string
	EnumTypeKind() Kind
	Commented() bool
	JoinPrefix() bool
	TrimPrefix() bool
	Validator() bool
	Iota() bool
	JSONMarshaler() bool
	TextMarshaler() bool
	XMLMarshaler() bool
	Stringer() bool
	StringFormater() string
}

func dependencies(s Settings) map[string]struct{} {
	if s == nil {
		return nil
	}
	dep := make(map[string]struct{})
	if s.JSONMarshaler() {
		if s.EnumTypeKind().IsNumber() {
			dep["strconv"] = struct{}{}
		}
		dep["encoding/json"] = struct{}{}
	}
	if s.XMLMarshaler() {
		if s.EnumTypeKind().IsNumber() {
			dep["strconv"] = struct{}{}
		}
		dep["encoding/xml"] = struct{}{}
	}
	if s.Stringer() {
		dep["fmt"] = struct{}{}
	}
	return dep
}
