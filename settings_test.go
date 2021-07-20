// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/genum/pkg/genum"
)

const (
	pkg    = "test"
	format = "format"
)

func TestSettings_ReadFrom(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			opts   Settings
			args   []string
			reader io.Reader
			// outputs
			failed bool
		}{
			"Default": {failed: true},
			"File not found": {
				args:   []string{"testdata/oops.csv"},
				failed: true,
			},
			"Stdin": {reader: strings.NewReader("csv")},
			"File":  {args: []string{"testdata/hello.csv"}},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tt.opts.ReadFrom(tt.args, tt.reader)
			are.Equal(err != nil, tt.failed) // unexpected error
		})
	}
}

func TestSettings_SrcFile(t *testing.T) {
	t.Parallel()
	t.Run("Default", func(t *testing.T) {
		t.Parallel()
		is.NewRelaxed(t).Equal(nil, Settings{}.SrcFile())
	})
	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		is.NewRelaxed(t).True(Settings{srcFile: strings.NewReader("")}.SrcFile() != nil)
	})
}

func TestSettings(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			opts Settings
			// outputs
			dstDir         string
			packageName    string
			enumType       string
			enumKind       genum.Kind
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
		}{
			"Default": {enumKind: genum.Int},
			"Text marshal only": {
				opts:          Settings{textMarshaler: true},
				enumKind:      genum.Int,
				stringer:      true,
				textMarshaler: true,
			},
			"String only": {
				opts:     Settings{stringer: true},
				enumKind: genum.Int,
				stringer: true,
			},
			"Complete": {
				opts: Settings{
					srcFile:        nil,
					dstDir:         "",
					packageName:    pkg,
					enumType:       genum.DefaultType,
					enumKind:       genum.Uint.Name(),
					stringFormater: format,
					stringer:       true,
					bitmask:        true,
					comment:        true,
					joinPrefix:     true,
					trimPrefix:     true,
					iota:           true,
					textMarshaler:  true,
					jsonMarshaler:  true,
					xmlMarshaler:   true,
					validator:      true,
				},
				dstDir:         strings.ToLower(genum.DefaultType) + ".go",
				packageName:    pkg,
				enumKind:       genum.Uint,
				enumType:       genum.DefaultType,
				stringFormater: format,
				stringer:       true,
				bitmask:        true,
				comment:        true,
				joinPrefix:     true,
				trimPrefix:     true,
				iota:           true,
				textMarshaler:  true,
				jsonMarshaler:  true,
				xmlMarshaler:   true,
				validator:      true,
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			are.True(strings.HasSuffix(tt.opts.DstFilename(), tt.dstDir)) // mismatch dstDir
			are.Equal(tt.packageName, tt.opts.PackageName())              // mismatch packageName
			are.Equal(tt.enumType, tt.opts.TypeName())                    // mismatch enumType
			are.Equal(tt.enumKind, tt.opts.TypeKind())                    // mismatch enumKind
			are.Equal(tt.stringFormater, tt.opts.StringFormater())        // mismatch stringFormater
			are.Equal(tt.bitmask, tt.opts.Bitmask())                      // mismatch bitmask
			are.Equal(tt.comment, tt.opts.Commented())                    // mismatch comment
			are.Equal(tt.joinPrefix, tt.opts.JoinPrefix())                // mismatch joinPrefix
			are.Equal(tt.trimPrefix, tt.opts.TrimPrefix())                // mismatch trimPrefix
			are.Equal(tt.iota, tt.opts.Iota())                            // mismatch iota
			are.Equal(tt.stringer, tt.opts.Stringer())                    // mismatch stringer
			are.Equal(tt.textMarshaler, tt.opts.TextMarshaler())          // mismatch textMarshaler
			are.Equal(tt.jsonMarshaler, tt.opts.JSONMarshaler())          // mismatch jsonMarshaler
			are.Equal(tt.xmlMarshaler, tt.opts.XMLMarshaler())            // mismatch xmlMarshaler
			are.Equal(tt.validator, tt.opts.Validator())                  // mismatch validator
		})
	}
}
