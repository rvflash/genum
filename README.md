# Genum

[![GoDoc](https://godoc.org/github.com/rvflash/genum?status.svg)](https://godoc.org/github.com/rvflash/genum)
[![Build Status](https://api.travis-ci.com/rvflash/genum.svg?branch=main)](https://travis-ci.com/rvflash/genum?branch=main)
[![Code Coverage](https://codecov.io/gh/rvflash/genum/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/genum)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/genum?)](https://goreportcard.com/report/github.com/rvflash/genum)

`genum` parses if provided a comma-separated values (CSV) file, or the standard input by default to automate 
the generation of a Go file, where CVS data becoming constant values of an enum type.

Given a name of a (signed or unsigned) integer, float or string type T and a package name, 
`genum` will create a new self-contained Go source file, named by default `<T>.go` with only constant values.

Each CSV line is a new constant, where the first column is the name, and/or the value, depending on the settings, 
and the second if provided, is the value.   
No check is made and records may have a variable number of fields. So we do not need to specify all values, even names.
`_` is used in the case of a missing constant name. 

Using `genum` flags, you can add implementation of the following interfaces:

* fmt.Stringer with customized output, see formatting options  to format the return based on enum type, name or value:
```go
    func (e T) String() string
```
* encoding.TextMarshaler / encoding.TextUnmarshaler
```go
    func (e T) MarshalText() (text []byte, err error)
    func (e *T) UnmarshalText(text []byte) error
```
* json.Marshaler / json.Unmarshaler
```go
    func (e T) MarshalJSON() ([]byte, error)
    func (e *T) UnmarshalJSON([]byte) error
```
* xml.Marshaler / xml.Unmarshaler
```go
    func (e T) MarshalXML(e *Encoder, start StartElement) error
    func (e *T) UnmarshalXML(d *Decoder, start StartElement) error
```

Typically, this process would be run using go generate, like this:

```go
//go:generate genum -pkg say -type hi values.csv
```


## Demo

The command `echo -e "hello\nbonjour\nguten morgen\nhola" | genum -pkg say -type hi` will generate the file `./hi.go`:

```go
// Code generated by "genum -pkg say -type hi"; DO NOT EDIT.

package say

// Hi is an enum.
type Hi int

// List of known Hi enums.
const (
	Hello Hi = iota
	Bonjour
	GutenMorgen
	Hola
)

```


## Installation

```shell
GO111MODULE=on go get github.com/rvflash/genum@vX.X.X
```


## Usage

```shell
genum [flags] [file]
```

The `genum` command generates a Go file based on the CSV values given in input (file path or standard input).
It supports the following flags:

    * `-pkg`: package name
    * `-type`: enum type name (default "Enum")
    * `-basic_type`: enum base type (default "int")
    * `-iota`: declare sequentially growing numeric constants (default true)
    * `-output`: output file name; default dst_dir/<snake_type>.go
    * `-stringer`: implement the fmt.Stringer interface
    * `-stringer_format` format used by the fmt.Stringer method (default only the enum name: "%[1]s"):
        [1] represents the enum name
        [2] represents the enum value
        [3] represents the enum type
    * `-noprefix`: trim the type name from the generated constant names
    * `-prefix`: add the type name as prefix of each generated constant names
    * `-json`: implement the json.Marshaler and json.Unmarshaler interfaces
    * `-text`: implement the encoding.TextMarshaler and encoding.TextUnmarshaler interfaces
    * `-xml`: implement the xml.Marshaler and xml.Unmarshaler interfaces
    * `-comment`: add in comment the values of generated constants