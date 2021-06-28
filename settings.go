// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rvflash/genum/pkg/genum"
	"github.com/rvflash/naming"
)

// Settings contains all the options exposed by Genum.
type Settings struct {
	srcFile        io.Reader
	dstDir         string
	packageName    string
	enumType       string
	enumKind       string
	stringFormater string
	stringer       bool
	bitmask        bool
	comment        bool
	joinPrefix     bool
	trimPrefix     bool
	iota           bool
	textMarshaler  bool
	jsonMarshaler  bool
	xmlMarshaler   bool
	validator      bool
}

// Bitmask implements the genum.Settings interface.
func (s Settings) Bitmask() bool {
	return s.bitmask
}

// Commented implements the genum.Settings interface.
func (s Settings) Commented() bool {
	return s.comment
}

const goFileExt = ".go"

// DstFilename implements the genum.Settings interface.
func (s Settings) DstFilename() string {
	if s.enumType == "" {
		return ""
	}
	if s.dstDir == "" {
		s.dstDir, _ = os.Getwd()
	}
	return filepath.Join(s.dstDir, naming.SnakeCase(s.enumType)+goFileExt)
}

// TypeKind implements the genum.Settings interface.
func (s Settings) TypeKind() genum.Kind {
	return genum.KindNamed(s.enumKind)
}

// TypeName implements the genum.Settings interface.
func (s Settings) TypeName() string {
	return naming.PascalCase(s.enumType)
}

// JoinPrefix implements the genum.Settings interface.
func (s Settings) JoinPrefix() bool {
	return s.joinPrefix
}

// Iota implements the genum.Settings interface.
func (s Settings) Iota() bool {
	return s.iota && s.TypeKind().IsInteger()
}

// Validator implements the genum.Settings interface.
func (s Settings) Validator() bool {
	return s.validator
}

// JSONMarshaler implements the genum.Settings interface.
func (s Settings) JSONMarshaler() bool {
	return s.jsonMarshaler
}

// PackageName implements the genum.Settings interface.
func (s Settings) PackageName() string {
	return naming.SnakeCase(s.packageName)
}

// ReadFrom allows to read from file path given as argument or the given reader.
func (s *Settings) ReadFrom(args []string, reader io.Reader) (err error) {
	if len(args) == 0 {
		if reader == nil {
			return io.ErrClosedPipe
		}
		s.srcFile = reader
		return
	}
	s.srcFile, err = os.Open(args[0])
	return
}

// SrcFile implements the genum.Settings interface.
func (s Settings) SrcFile() io.Reader {
	return s.srcFile
}

// Stringer implements the genum.Settings interface.
func (s Settings) Stringer() bool {
	return s.stringer || s.textMarshaler
}

// StringFormater implements the genum.Settings interface.
func (s Settings) StringFormater() string {
	return s.stringFormater
}

// TextMarshaler implements the genum.Settings interface.
func (s Settings) TextMarshaler() bool {
	return s.textMarshaler
}

// TrimPrefix implements the genum.Settings interface.
func (s Settings) TrimPrefix() bool {
	return s.trimPrefix
}

// XMLMarshaler implements the genum.Settings interface.
func (s Settings) XMLMarshaler() bool {
	return s.xmlMarshaler
}
