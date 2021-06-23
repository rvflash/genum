// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	shortName      = "e"
	mixedName      = "v"
	strName        = "s"
	srcName        = "data"
	increment      = "iota"
	returnNotFound = "return \"\", false\n"
	returnErr      = "return fmt.Errorf(\"%s expects %s but got %%s\", %s)\n"
	lookupFunc     = "func lookup%[1]s(%[2]s %[1]s) (%[3]s string, ok bool) {\n"
	unnamed        = "_"
	zero           = "0"

	maxHumanScale = 10
)

// Layout returns the generation configuration based on the given settings.
func Layout(s Settings, args []string) []Configurator {
	if s == nil {
		return nil
	}
	cnf := []Configurator{
		ParseEnums(s.SrcFile(), s.EnumTypeName(), s.EnumTypeKind(), s.JoinPrefix(), s.TrimPrefix(), s.Iota()),
		PrintHeader(s.PackageName(), args, dependencies(s)),
		PrintEnums(s.EnumTypeName(), s.EnumTypeKind(), s.Iota(), s.Commented()),
	}
	if s.Stringer() || s.Validator() {
		cnf = append(cnf, PrintLookup(s.EnumTypeName(), s.EnumTypeKind()))
	}
	if s.Stringer() {
		cnf = append(cnf, PrintStringer(s.StringFormater(), s.EnumTypeName(), s.EnumTypeKind()))
	}
	if s.Validator() {
		cnf = append(cnf, PrintValidator(s.EnumTypeName()))
	}
	if s.JSONMarshaler() {
		cnf = append(cnf, PrintJSONMarshaler(s.EnumTypeName(), s.EnumTypeKind()))
	}
	if s.TextMarshaler() {
		cnf = append(cnf, PrintTextMarshaler(s.StringFormater(), s.EnumTypeName()))
	}
	if s.XMLMarshaler() {
		cnf = append(cnf, PrintXMLMarshaler(s.EnumTypeName(), s.EnumTypeKind()))
	}
	return append(cnf, WriteFile(s.DstFilename()))
}

// Generate generates the enum file based on these options.
func Generate(opts ...Configurator) (err error) {
	gen := new(Generator)
	for _, opt := range opts {
		err = opt(gen)
		if err != nil {
			return
		}
	}
	return
}

// Generator represents an enum generator.
type Generator struct {
	enums []Enum
	basic bool
	buf   bytes.Buffer
	err   error
}

func (g *Generator) advanceString(enumType string) error {
	// Variable with enums names.
	g.printf("\n")
	g.printf("var _%sNames = map[%s]string{\n", enumType, enumType)
	for _, e := range g.enums {
		g.printf("%s: %q,\n", e.Text, e.RawText)
	}
	g.printf("}\n")

	// Lookup method based on the precomputed map of Enums.
	g.printf("\n")
	g.printf(lookupFunc, enumType, shortName, strName)
	g.printf("%s, ok = _%sNames[%s]\n", strName, enumType, shortName)
	g.printf("return %s, ok\n", strName)
	g.printf("}\n")

	return nil
}

func (g *Generator) averageString(enumType string) error {
	// Lookup method (with human readable switch case)
	g.printf("\n")
	g.printf(lookupFunc, enumType, shortName, strName)
	g.printf("switch %s {\n", shortName)
	for _, e := range g.enums {
		g.printf("case %s:\n", e.Text)
		g.printf("return %q, true\n", e.RawText)
	}
	g.printf("default:\n")
	g.printf(returnNotFound)
	g.printf("}\n")
	g.printf("}\n")

	return nil
}

func (g *Generator) basicString(enumType string, enumKind Kind) error {
	var (
		buf = new(bytes.Buffer)
		pos = make([]int, len(g.enums))
	)
	for p, e := range g.enums {
		_, _ = buf.WriteString(e.RawText)
		pos[p] = buf.Len()
	}
	// Constant with all enums names concatenated together.
	g.printf("\n")
	g.printf("const _%sNames = %q\n", enumType, buf.String())
	max := buf.Len()
	buf.Reset()
	for i, v := range pos {
		if i > 0 {
			_, _ = buf.WriteString(", ")
		}
		_, _ = buf.WriteString(strconv.Itoa(v))
	}

	// Variable with positions of any enums names.
	g.printf("\n")
	g.printf("var _%sIndexes = [...]uint%d{0, %s}\n", enumType, unsignedSize(max), buf.String())

	// Lookup method based on a slice of contant names.
	g.printf("\n")
	g.printf(lookupFunc, enumType, shortName, strName)
	var guardRail string
	if enumKind.IsSigned() && g.enums[0].Value != zero {
		guardRail = shortName + " < 0 || "
		g.printf("%s -= %s\n", shortName, g.enums[0].Value)
	}
	g.printf("if %s%s >= %s(len(_%sIndexes)-1) {\n", guardRail, shortName, enumType, enumType)
	g.printf(returnNotFound)
	g.printf("}\n")
	g.printf("return _%[1]sNames[_%[1]sIndexes[%[2]s]:_%[1]sIndexes[%[2]s+1]], true\n", enumType, shortName)
	g.printf("}\n")

	return nil
}

type mode uint8

const (
	// basic is the default mode, where enums are integers.
	basic mode = iota
	// average is used when the list of enums is human scale with integers.
	average
	// advance is used to handle enums not matching constraints of other modes.
	advance
)

func (g Generator) mode() mode {
	n := len(g.enums)
	if n == 0 || g.basic {
		return basic
	}
	if n > maxHumanScale {
		return advance
	}
	return average
}

// printf formats according to a format specifier and writes to the Generator's buffer if no error has already occurred.
func (g *Generator) printf(format string, a ...interface{}) {
	if g.err == nil {
		_, g.err = fmt.Fprintf(&g.buf, format, a...)
	}
}
