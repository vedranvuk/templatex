// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package templatex

import "github.com/vedranvuk/errorex"

var (
	// ErrTemplatex is the base error of templatex package.
	ErrTemplatex = errorex.New("templatex")

	// ErrParse is returned when a parse error occurs.
	ErrParse = ErrTemplatex.Wrap("parse error")
	// ErrNotFound is returned when a non-existent namespace is
	// being addressed.
	ErrNotFound = ErrTemplatex.WrapFormat("namespace '%s' not found")
)
