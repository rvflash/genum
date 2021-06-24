// Code generated by "genum -pkg huge_iota_string_valid_text_comment -type int -validator -stringer -text -comment hello.csv"; DO NOT EDIT.

package huge_iota_string_valid_text_comment

import (
	"fmt"
)

// Enum is an enum.
type Enum int

// List of known Enum enums.
const (
	Hello   Enum = iota + -4 // e = -4
	Bonjour                  // e = -3
	Hallo                    // e = -2
	Hola                     // e = -1
	Ciao                     // e = 0
	Ossu                     // e = 1
	Yasou                    // e = 2
	Salam   Enum = iota + 5  // e = 12
	Zdravo                   // e = 13
	Salutu                   // e = 14
	SubhDin                  // e = 15
	Paka                     // e = 16
	Watdi                    // e = 17
	A                        // e = 18
)

const _EnumNames = "hellobonjourhalloholaciaoOssuyasouSalamZdravoSalutusubh dinPakaWatdiA"

var _EnumIndexes = [...]uint8{0, 5, 12, 17, 21, 25, 29, 34, 39, 45, 51, 59, 63, 68, 69}

func lookupEnum(e Enum) (s string, ok bool) {
	e -= -4
	if e < 0 || e >= Enum(len(_EnumIndexes)-1) {
		return "", false
	}
	return _EnumNames[_EnumIndexes[e]:_EnumIndexes[e+1]], true
}

// String implements the fmt.Stringer interface.
func (e Enum) String() string {
	s, ok := lookupEnum(e)
	if !ok {
		return fmt.Sprintf("%[3]s(%[2]d)", "", int(e), "Enum")
	}
	return s
}

// IsValid returns true if the Enum is a known constant.
func (e Enum) IsValid() bool {
	_, ok := lookupEnum(e)
	return ok
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e Enum) MarshalText() (text []byte, err error) {
	return []byte(e.String()), nil
}

var _EnumStrings = map[string]Enum{
	"hello":    Hello,
	"bonjour":  Bonjour,
	"hallo":    Hallo,
	"hola":     Hola,
	"ciao":     Ciao,
	"Ossu":     Ossu,
	"yasou":    Yasou,
	"Salam":    Salam,
	"Zdravo":   Zdravo,
	"Salutu":   Salutu,
	"subh din": SubhDin,
	"Paka":     Paka,
	"Watdi":    Watdi,
	"A":        A,
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Enum) UnmarshalText(text []byte) error {
	e2, ok := _EnumStrings[string(text)]
	if !ok {
		return fmt.Errorf("%q is not a known Enum", text)
	}
	*e = e2
	return nil
}
