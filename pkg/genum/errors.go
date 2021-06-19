// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package genum

type errGenum string

// Error implements the error interface.
func (e errGenum) Error() string {
	return string(e)
}

const (
	// ErrMissing is returned when a data is missing.
	ErrMissing = errGenum("missing data")
)
