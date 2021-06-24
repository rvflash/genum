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

const cmdFileName = iota + 1

func main() {
	var (
		s = new(Settings)
		w = log.New(os.Stderr, genum.Command+": ", 0)
		u = "package name"
	)
	flag.StringVar(&s.packageName, "pkg", "", u)
	u = "enum type name"
	flag.StringVar(&s.enumType, "name", "Enum", u)
	u = "enum base type"
	flag.StringVar(&s.enumKind, "type", "int", u)
	u = "declare sequentially growing numeric constants"
	flag.BoolVar(&s.iota, "iota", true, u)
	u = "output file name; default dst_dir/<snake_type>.go"
	flag.StringVar(&s.dstDir, "output", "", u)
	u = "implement the fmt.Stringer interface"
	flag.BoolVar(&s.stringer, "stringer", false, u)
	u = "format used by the fmt.Stringer method:\n"
	u += "[%d] represents the enum name\n"
	u += "[%d] represents the enum value\n"
	u += "[%d] represents the enum type\n"
	u = fmt.Sprintf(u, genum.NamePos, genum.ValuePos, genum.TypePos)
	flag.StringVar(&s.stringFormater, "stringer_format", genum.NameFormat(), u)
	u = "trim the type name from the generated constant names"
	flag.BoolVar(&s.trimPrefix, "noprefix", false, u)
	u = "add the type name as prefix of each generated constant names"
	flag.BoolVar(&s.joinPrefix, "prefix", false, u)
	u = "implement the json.Marshaler and json.Unmarshaler interfaces"
	flag.BoolVar(&s.jsonMarshaler, "json", false, u)
	u = "implement the encoding.TextMarshaler and encoding.TextUnmarshaler interfaces"
	flag.BoolVar(&s.textMarshaler, "text", false, u)
	u = "implement the xml.Marshaler and xml.Unmarshaler interfaces"
	flag.BoolVar(&s.xmlMarshaler, "xml", false, u)
	u = "add in comment the values of generated constants"
	flag.BoolVar(&s.comment, "comment", false, u)
	u = "add a method \"IsValid\" to verify the set up of the constant"
	flag.BoolVar(&s.validator, "validator", false, u)
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
