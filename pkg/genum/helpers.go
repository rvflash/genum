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

func sumNumbers(sub bool, base uint64, baseSign bool, num uint64, numSign bool) (res uint64, resSign bool) {
	if sub || baseSign || numSign {
		b := int64(base)
		if baseSign {
			b *= -1
		}
		n := int64(num)
		if numSign {
			n *= -1
		}
		if sub {
			b -= n
		} else {
			b += n
		}
		if resSign = b < 0; resSign {
			b *= -1
		}
		return uint64(b), resSign
	}
	return base + num, baseSign
}

func calculateIota(cur uint64, curSign bool, prev uint64, prevSign bool) string {
	delta, sign := sumNumbers(true, cur, curSign, prev, prevSign)
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
	data []string, kind Kind, curIota int64, prevUint uint64, prevSign bool,
) (enumValue, enumIota string, curUint uint64, curSign bool) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("enum %q: %w", enumValue, err)
		}
	}()
	switch {
	case kind.IsInteger():
		enumValue, enumIota, curUint, curSign, err = enumIntegerValue(data, kind, curIota, prevUint, prevSign)
	case kind.IsNumber():
		enumValue, err = enumFloatValue(data)
	default:
		enumValue = enumStringValue(data)
	}
	return
}

func enumIntegerValue(
	data []string, kind Kind, curIota int64, prevUint uint64, prevSign bool,
) (enumValue, enumIota string, curUint uint64, curSign bool, err error) {
	value, ok := field(data, valuePos)
	if !ok {
		if curIota > -1 {
			if curIota == 0 {
				enumIota = increment // First value
			}
			return fmtNumber(prevUint+uint64(curIota), prevSign), enumIota, prevUint + 1, prevSign, nil
		}
		return zero, enumIota, prevUint, prevSign, nil
	}
	curUint, curSign, err = parseNumber(value, kind)
	if err != nil {
		return value, enumIota, curUint, curSign, err
	}
	if curIota > -1 {
		curUint, curSign = sumNumbers(true, curUint, curSign, uint64(curIota), false)
		enumIota = calculateIota(curUint, curSign, prevUint, prevSign)
	}
	return value, enumIota, curUint, curSign, nil
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

func parseNumber(value string, kind Kind) (abs uint64, sign bool, err error) {
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

func fmtNumber(i uint64, sign bool) string {
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
