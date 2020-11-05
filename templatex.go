// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package templatex

import "github.com/vedranvuk/errorex"

var (
	// ErrTemplater is the base error of templater package.
	ErrTemplater = errorex.New("templater")
	// ErrParse is returned when a parse error occurs.
	ErrParse = ErrTemplater.Wrap("parse error")
	// ErrNamespaceNotFound is returned when a non-existent namespace is
	// addressed in an ExecuteNamespace call.
	ErrNamespaceNotFound = ErrTemplater.WrapFormat("template '%s' not found")
)
