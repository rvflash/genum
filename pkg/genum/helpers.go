// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/rvflash/naming"
)

func absDiff(cur uint64, curSign bool, prev uint64, prevSign bool) uint64 {
	if curSign == prevSign {
		return cur - prev
	}
	cInt, pInt := int64(cur), int64(prev)
	if curSign {
		cInt *= -1
	}
	if prevSign {
		pInt *= -1
	}
	cInt -= pInt
	if cInt < 0 {
		cInt *= -1
	}
	return uint64(cInt)
}

func calculateIota(cur uint64, curSign bool, prev uint64, prevSign bool) string {
	var (
		delta uint64
		sign  bool
	)
	if curSign || prevSign {
		c := int64(cur)
		if curSign {
			c *= -1
		}
		p := int64(prev)
		if prevSign {
			p *= -1
		}
		c -= p
		sign = c < 0
		if sign {
			c *= -1
		}
		delta = uint64(c)
	} else {
		delta = cur - prev
	}
	if delta <= 1 {
		if prev == 0 {
			return increment
		}
		return ""
	}
	if sign {
		return increment + " + -" + strconv.FormatUint(cur, base10)
	}
	return increment + " + " + strconv.FormatUint(cur, base10)
}

func enumName(data []string, enumType string, joinPrefix, trimPrefix bool) string {
	name, ok := field(data, namePos)
	if !ok || name == "" {
		return unnamed
	}
	name = naming.PascalCase(name)
	if trimPrefix {
		name = strings.TrimPrefix(name, enumType)
	}
	if joinPrefix {
		name = enumType + name
	}
	return name
}

func enumValue(
	data []string, kind Kind, useIota bool, prevUint uint64, prevSign bool,
) (enumValue, enumIota string, curUint uint64, curSign bool) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("enum %q: %w", enumValue, err)
		}
	}()
	switch {
	case kind.IsInteger():
		enumValue, enumIota, curUint, curSign, err = enumIntegerValue(data, kind, useIota, prevUint, prevSign)
	case kind.IsNumber():
		enumValue, err = enumFloatValue(data)
	default:
		enumValue = enumStringValue(data)
	}
	return
}

func enumIntegerValue(
	data []string, kind Kind, useIota bool, prevUint uint64, prevSign bool,
) (enumValue, enumIota string, curUint uint64, curSign bool, err error) {
	value, ok := field(data, valuePos)
	if !ok {
		if useIota {
			if prevUint == 0 {
				enumIota = increment // First value
			}
			return fmtNum(prevUint, prevSign), enumIota, prevUint + 1, curSign, nil
		}
		return zero, enumIota, prevUint, curSign, nil
	}
	curUint, curSign, err = parseNum(value, kind)
	if err != nil {
		return
	}
	if useIota {
		enumIota = calculateIota(curUint, curSign, prevUint, prevSign)
		curUint++
	}
	return
}

func enumFloatValue(data []string) (string, error) {
	value, ok := field(data, valuePos)
	if !ok {
		return zero, nil
	}
	_, err := strconv.ParseFloat(value, bits64)
	if err != nil {
		return "", err
	}
	return value, nil
}

func enumStringValue(data []string) string {
	value, ok := field(data, valuePos)
	if !ok {
		value, _ = field(data, namePos)
	}
	return fmt.Sprintf("%q", value)
}

func parseNum(value string, kind Kind) (abs uint64, sign bool, err error) {
	if !kind.IsInteger() {
		err = strconv.ErrRange
		return
	}
	if !kind.IsSigned() {
		abs, err = strconv.ParseUint(value, base10, kind.BitSize())
		return
	}
	var i int64
	i, err = strconv.ParseInt(value, base10, kind.BitSize())
	if err != nil {
		return
	}
	sign = i < 0
	if sign {
		i = -i
	}
	return uint64(i), sign, nil
}

func fmtNum(i uint64, sign bool) string {
	if sign {
		if i > math.MaxInt64 {
			return zero
		}
		return strconv.FormatInt(int64(i)*-1, base10)
	}
	return strconv.FormatUint(i, base10)
}

// unsignedSize returns the smallest possible slice of integers to use as indexes into the concatenated strings.
func unsignedSize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	case n < 1<<32:
		return 32
	default:
		return 64
	}
}

const (
	namePos int = iota
	valuePos
)

func field(data []string, column int) (s string, ok bool) {
	if len(data) <= column {
		return "", false
	}
	return data[column], true
}
