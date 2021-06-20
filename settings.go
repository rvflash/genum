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

const goFileExt = ".go"

// Settings contains all the options exposed by Genum.
type Settings struct {
	srcFile        io.Reader
	dstDir         string
	packageName    string
	enumType       string
	enumKind       string
	comment        bool
	joinPrefix     bool
	trimPrefix     bool
	iota           bool
	textMarshaler  bool
	jsonMarshaler  bool
	xmlMarshaler   bool
	stringer       bool
	stringFormater string
}

// Commented implements the genum.Settings interface.
func (s Settings) Commented() bool {
	return s.comment
}

// DstFilename implements the genum.Settings interface.
func (s Settings) DstFilename() string {
	if s.enumKind == "" {
		return ""
	}
	if s.dstDir == "" {
		s.dstDir, _ = os.Getwd()
	}
	return filepath.Join(s.dstDir, naming.SnakeCase(s.enumType)+goFileExt)
}

// EnumTypeKind implements the genum.Settings interface.
func (s Settings) EnumTypeKind() genum.Kind {
	return genum.KindNamed(s.enumKind)
}

// EnumTypeName implements the genum.Settings interface.
func (s Settings) EnumTypeName() string {
	return naming.PascalCase(s.enumType)
}

// JoinPrefix implements the genum.Settings interface.
func (s Settings) JoinPrefix() bool {
	return s.joinPrefix
}

// Iota implements the genum.Settings interface.
func (s Settings) Iota() bool {
	return s.iota && s.EnumTypeKind().IsInteger()
}

// JSONMarshaler implements the genum.Settings interface.
func (s Settings) JSONMarshaler() bool {
	return s.jsonMarshaler
}

// PackageName implements the genum.Settings interface.
func (s Settings) PackageName() string {
	return naming.FlatCase(s.packageName)
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
	return s.stringer
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
