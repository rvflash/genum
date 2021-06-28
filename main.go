// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rvflash/genum/pkg/genum"
)

const (
	bitmaskUsage = `use one integer to hold multiple flags, provide bitwise operations and 
overwrite the enum base type with unsigned integer type (size in bits based on the number of values)`
	commentUsage        = "add in comment the values of generated constants"
	enumTypeUsage       = "enum type name"
	enumKindUsage       = "enum base type"
	iotaUsage           = "declare sequentially growing numeric constants"
	jsonUsage           = "implement the json.Marshaler and json.Unmarshaler interfaces"
	noPrefixUsage       = "trim the type name from the generated constant names"
	outputUsage         = "output file name; default dst_dir/<snake_type>.go"
	packageNameUsage    = "package name"
	prefixUsage         = "add the type name as prefix of each generated constant names"
	stringerUsage       = "implement the fmt.Stringer interface"
	stringFormaterUsage = `format used as returned value by the fmt.Stringer method:
[%d] represents the enum name
[%d] represents the enum value
[%d] represents the enum type`
	textUsage      = "implement the encoding.TextMarshaler and encoding.TextUnmarshaler interfaces"
	validatorUsage = `add a method "IsValid" to verify the set up of the constant`
	xmlUsage       = "implement the xml.Marshaler and xml.Unmarshaler interfaces"
)

const cmdFileName = 1

func main() {
	var (
		s = new(Settings)
		w = log.New(os.Stderr, genum.Command+": ", 0)
		u = fmt.Sprintf(stringFormaterUsage, genum.NamePos, genum.ValuePos, genum.TypePos)
	)
	flag.StringVar(&s.packageName, "pkg", "", packageNameUsage)
	flag.StringVar(&s.enumType, "name", genum.DefaultType, enumTypeUsage)
	flag.StringVar(&s.enumKind, "type", genum.DefaultKind, enumKindUsage)
	flag.StringVar(&s.dstDir, "output", "", outputUsage)
	flag.StringVar(&s.stringFormater, "stringer_format", genum.NameFormat(), u)
	flag.BoolVar(&s.stringer, "stringer", false, stringerUsage)
	flag.BoolVar(&s.trimPrefix, "noprefix", false, noPrefixUsage)
	flag.BoolVar(&s.joinPrefix, "prefix", false, prefixUsage)
	flag.BoolVar(&s.jsonMarshaler, "json", false, jsonUsage)
	flag.BoolVar(&s.textMarshaler, "text", false, textUsage)
	flag.BoolVar(&s.xmlMarshaler, "xml", false, xmlUsage)
	flag.BoolVar(&s.validator, "validator", false, validatorUsage)
	flag.BoolVar(&s.comment, "comment", false, commentUsage)
	flag.BoolVar(&s.bitmask, "bitmask", false, bitmaskUsage)
	flag.BoolVar(&s.iota, "iota", true, iotaUsage)
	flag.Parse()

	err := s.ReadFrom(flag.Args(), os.Stdin)
	if err != nil {
		w.Fatalf("source: %s", err)
	}
	err = genum.Generate(genum.Layout(s, os.Args[cmdFileName:])...)
	if err != nil {
		w.Fatal(err)
	}
}
