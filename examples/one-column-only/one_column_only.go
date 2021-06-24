// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package one_column_only

//go:generate genum -pkg int -type int -output int hello.csv
//go:generate genum -pkg int8 -type int8 -output int8 hello.csv
//go:generate genum -pkg int16 -type int16 -output int16 hello.csv
//go:generate genum -pkg int32 -type int32 -output int32 hello.csv
//go:generate genum -pkg int64 -type int64 -output int64 hello.csv
//go:generate genum -pkg uint -type uint -output uint hello.csv
//go:generate genum -pkg uint8 -type uint8 -output uint8 hello.csv
//go:generate genum -pkg uint16 -type uint16 -output uint16 hello.csv
//go:generate genum -pkg uint32 -type uint32 -output uint32 hello.csv
//go:generate genum -pkg uint64 -type uint64 -output uint64 hello.csv
//go:generate genum -pkg float32 -type float32 -output float32 hello.csv
//go:generate genum -pkg float64 -type float64 -output float64 hello.csv
//go:generate genum -pkg string -type string -output string hello.csv
