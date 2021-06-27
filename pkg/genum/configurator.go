// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

import (
	"encoding/csv"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
)

// Configurator must be implemented by any methods acted as an enum layout generator.
type Configurator func(g *Generator) error

// ParseBitmask reads the given source as a CSV and tries to create a bitmasks list.
func ParseBitmask(data io.Reader, enumType string, joinPrefix, trimPrefix bool) Configurator {
	return func(g *Generator) error {
		r := csv.NewReader(data)
		r.FieldsPerRecord = -1 // Records may have a variable number of fields.
		g.enums = make([]Enum, 0)
		var curIota uint64
		for {
			d, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("source file: %w", err)
			}
			curIota++
			g.enums = append(g.enums, Enum{
				Iota:    bitmaskIota(curIota),
				Text:    enumName(d, enumType, joinPrefix, trimPrefix),
				RawText: enumRawName(d),
				Type:    enumType,
				Value:   bitmaskValue(curIota),
			})
		}
		enumKind := KindUnsigned(len(g.enums))
		for k := range g.enums {
			g.enums[k].Kind = enumKind
		}
		return nil
	}
}

// ParseEnums reads the given source as a CSV and tries to create a list of constants based on it.
func ParseEnums(data io.Reader, enumType string, enumKind Kind, joinPrefix, trimPrefix, useIota bool) Configurator {
	return func(g *Generator) error {
		var (
			curIota           int64 = -1
			curUint, prevUint uint64
			curSign, prevSign bool
		)
		r := csv.NewReader(data)
		r.FieldsPerRecord = -1 // Records may have a variable number of fields.
		g.enums = make([]Enum, 0)
		for {
			d, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("source file: %w", err)
			}
			e := Enum{
				Type:    enumType,
				Kind:    enumKind,
				Text:    enumName(d, enumType, joinPrefix, trimPrefix),
				RawText: enumRawName(d),
			}
			if useIota {
				curIota++
			}
			e.Value, e.Iota, curUint, curSign = enumValue(d, enumKind, curIota, prevUint, prevSign)
			g.basic = enumKind.IsInteger() && absDiff(curUint, curSign, prevUint, prevSign) == 1
			g.enums = append(g.enums, e)
			prevUint, prevSign = curUint, curSign
		}
		if len(g.enums) == 0 {
			return fmt.Errorf("source file: %w", io.ErrUnexpectedEOF)
		}
		return nil
	}
}

// PrintBitmask prints related methods to bitmask operations.
func PrintBitmask(enumType string) Configurator {
	return func(g *Generator) error {
		// Has method
		g.printf("\n")
		g.printf("// Has returns in success if this %s is set on it.\n", enumType)
		g.printf("func (%[1]s %[2]s) Has(%[1]s2 %[2]s) bool {\n", shortName, enumType)
		g.printf("return %[1]s&%[1]s2 != 0\n", shortName)
		g.printf("}\n")

		// Set method
		g.printf("\n")
		g.printf("// Set sets this %[1]s on the current %[1]s.\n", enumType)
		g.printf("func (%[1]s *%[2]s) Set(%[1]s2 %[2]s)  {\n", shortName, enumType)
		g.printf("*%[1]s |= %[1]s2\n", shortName)
		g.printf("}\n")

		// Switch method
		g.printf("\n")
		g.printf("// Switch only changes the %s value if necessary.\n", enumType)
		g.printf("// It returns true if the requested action has been done.\n")
		g.printf("func (%[1]s *%[2]s) Switch(%[1]s2 %[2]s, on bool) (done bool) {\n", shortName, enumType)
		g.printf("if %[1]s.Has(%[1]s2) == on {\n", shortName)
		g.printf("return false\n")
		g.printf("}\n")
		g.printf("if on {\n")
		g.printf("%[1]s.Set(%[1]s2)\n", shortName)
		g.printf("} else {\n")
		g.printf("%[1]s.Unset(%[1]s2)\n", shortName)
		g.printf("}\n")
		g.printf("return true\n")
		g.printf("}\n")

		// Toggle method
		g.printf("\n")
		g.printf("// Toggle toggles this %s value.\n", enumType)
		g.printf("func (%[1]s *%[2]s) Toggle(%[1]s2 %[2]s)  {\n", shortName, enumType)
		g.printf("*%[1]s ^=  %[1]s2\n", shortName)
		g.printf("}\n")

		// Unset method
		g.printf("\n")
		g.printf("// Unset clears this %s value on the current one.\n", enumType)
		g.printf("func (%[1]s *%[2]s) Unset(%[1]s2 %[2]s)  {\n", shortName, enumType)
		g.printf("*%[1]s &^=  %[1]s2\n", shortName)
		g.printf("}\n")

		return nil
	}
}

// PrintEnums prints the list of constants.
func PrintEnums(enumType string, useIota, commented bool) Configurator {
	return func(g *Generator) error {
		if enumType == "" {
			return fmt.Errorf("enum type: %w", ErrMissing)
		}
		g.printf("\n")
		g.printf("// %s is an enum.\n", enumType)
		g.printf("type %s %s\n", enumType, g.enums[0].Kind.Name())
		g.printf("\n")
		g.printf("// List of known %s enums.\n", enumType)
		g.printf("const (\n")
		for k, v := range g.enums {
			g.printf(v.Format(k, useIota, commented))
		}
		g.printf(")\n")

		return nil
	}
}

// PrintHeader prints the go file header (package, import, etc.).
func PrintHeader(pkg string, args []string, packages map[string]struct{}) Configurator {
	return func(g *Generator) error {
		if pkg == "" {
			return fmt.Errorf("package RawName: %w", ErrMissing)
		}
		g.printf("// Code generated by %q; DO NOT EDIT.\n", Command+" "+strings.Join(args, " "))
		g.printf("\n")
		g.printf("package %s\n", pkg)
		g.printf("\n")
		if len(packages) == 0 {
			return nil
		}
		g.printf("import (\n")
		for name := range packages {
			g.printf("%q\n", name)
		}
		g.printf(")\n")

		return nil
	}
}

// PrintJSONMarshaler adds methods to marshal and unmarshal the enum value as JSON data.
func PrintJSONMarshaler(enumType string, enumKind Kind) Configurator {
	return func(g *Generator) error {
		// json.Marshaler
		g.printf("\n")
		g.printf("// MarshalJSON implements the json.Marshaler interface.\n")
		g.printf("func (e %s) MarshalJSON() ([]byte, error) {\n", enumType)
		g.printf("return json.Marshal(%s)", strConvFormat(enumKind))
		g.printf("}\n")

		// json.Unmarshaler
		g.printf("// UnmarshalJSON implements the json.Unmarshaler interface.\n")
		g.printf("func (e *%s) UnmarshalJSON(data []byte) error{\n", enumType)
		g.printf("var (\n")
		g.printf("s string\n")
		g.printf("err = json.Unmarshal(data, &s)\n")
		g.printf(")\n")
		g.printf("if err != nil {\n")
		g.printf(returnErr, enumType, enumKind.Name(), srcName)
		g.printf("}\n")
		g.printf(strConvParse(enumType, enumKind))
		g.printf("return nil\n")
		g.printf("}\n")

		return nil
	}
}

// PrintValidator builds a method to check the validity of a constant.
func PrintValidator(enumType string) Configurator {
	return func(g *Generator) error {
		g.printf("\n")
		g.printf("// IsValid returns true if the %s is a known constant.\n", enumType)
		g.printf("func (%s %s) IsValid() bool {\n", shortName, enumType)
		g.printf("_, ok := lookup%s(%s)\n", enumType, shortName)
		g.printf("return ok\n")
		g.printf("}\n")

		return nil
	}
}

// PrintStringer chooses the "best" methods regarding the data to manage the String method.
func PrintStringer(format string, enumType string, enumKind Kind) Configurator {
	return func(g *Generator) error {
		g.printf("\n")
		g.printf("// String implements the fmt.Stringer interface.\n")
		g.printf("func (%s %s) String() string {\n", shortName, enumType)
		g.printf("s, ok := lookup%s(%s)\n", enumType, shortName)
		g.printf("if !ok {")
		g.printf(
			"return fmt.Sprintf(%q, %q, %v, %q)\n",
			DefaultFormat(enumKind.ValueFormat()), "", enumKind.Cast(shortName), enumType,
		)
		g.printf("}\n")
		if format != NameFormat() {
			g.printf("return fmt.Sprintf(%q, %s, %v, %q)\n", format, strName, enumKind.Cast(shortName), enumType)
		} else {
			g.printf("return %s\n", strName)
		}
		g.printf("}\n")

		return nil
	}
}

// PrintLookup adds a private method dedicated to get if exists the name of a constant with ok at true.
// Otherwise, a standard failover name with ok at false are returned.
func PrintLookup(enumType string, enumKind Kind) Configurator {
	return func(g *Generator) error {
		switch g.mode() {
		case basic:
			return g.basicString(enumType, enumKind)
		case average:
			return g.averageString(enumType)
		default:
			return g.advanceString(enumType)
		}
	}
}

// PrintTextMarshaler adds methods to marshal and unmarshal the enum String value as a text.
func PrintTextMarshaler(format string, enumType string) Configurator {
	return func(g *Generator) error {
		// encoding.TextMarshaler
		g.printf("\n")
		g.printf("// MarshalText implements the encoding.TextMarshaler interface.\n")
		g.printf("func (e %s) MarshalText() (text []byte, err error) {\n", enumType)
		g.printf("return []byte(e.String()), nil")
		g.printf("}\n")

		// encoding.TextUnmarshaler
		g.printf("var _%sStrings = map[string]%s{\n", enumType, enumType)
		var (
			v   interface{}
			err error
		)
		for k, e := range g.enums {
			if format != NameFormat() {
				v, err = e.ParseValue()
				if err != nil {
					return fmt.Errorf("enum value #%d: %w", k, err)
				}
				g.printf("%q: %s,\n", fmt.Sprintf(format, e.RawText, v, enumType), e.Text)
			} else {
				g.printf("%q: %s,\n", e.RawText, e.Text)
			}
		}
		g.printf("}\n")

		// encoding.TextUnmarshaler
		g.printf("\n")
		g.printf("// UnmarshalText implements the encoding.TextUnmarshaler interface.\n")
		g.printf("func (e *%s) UnmarshalText(text []byte) error{\n", enumType)
		g.printf("%s2, ok := _%sStrings[string(text)]\n", shortName, enumType)
		g.printf("if !ok {\n")
		g.printf("return fmt.Errorf(\"%%q is not a known %s\", text)\n", enumType)
		g.printf("}\n")
		g.printf("*%s = %s2\n", shortName, shortName)
		g.printf("return nil\n")
		g.printf("}\n")

		return nil
	}
}

// PrintXMLMarshaler adds methods to marshal and unmarshal the enum value as XML data.
func PrintXMLMarshaler(enumType string, enumKind Kind) Configurator {
	return func(g *Generator) error {
		// json.Marshaler
		g.printf("\n")
		g.printf("// MarshalXML implements the xml.Marshaler interface.\n")
		g.printf("func (e %s) MarshalXML(x *xml.Encoder, start xml.StartElement) error {\n", enumType)
		g.printf("return x.EncodeElement(%s, start)", strConvFormat(enumKind))
		g.printf("}\n")

		// json.Unmarshaler
		g.printf("// UnmarshalXML implements the xml.Unmarshaler interface.\n")
		g.printf("func (e *%s) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error{\n", enumType)
		g.printf("var (\n")
		g.printf("s string\n")
		g.printf("err = d.DecodeElement(&s, &start)\n")
		g.printf(")\n")
		g.printf("if err != nil {\n")
		g.printf(returnErr, enumType, enumKind.Name(), strName)
		g.printf("}\n")
		g.printf(strConvParse(enumType, enumKind))
		g.printf("return nil\n")
		g.printf("}\n")

		return nil
	}
}

// WriteFile tries to write the go file.
func WriteFile(filename string) Configurator {
	return func(g *Generator) error {
		if filename == "" {
			return fmt.Errorf("destination filename: %w", ErrMissing)
		}
		if g.err != nil {
			return g.err
		}
		src, err := format.Source(g.buf.Bytes())
		if err != nil {
			// Allows the user to compile the output to see the error.
			src = g.buf.Bytes()
		}
		wrr := ioutil.WriteFile(filename, src, 0600)
		if wrr != nil {
			return fmt.Errorf("destination: %w", wrr)
		}
		if err != nil {
			return fmt.Errorf("go format failed: %w", wrr)
		}
		return nil
	}
}
