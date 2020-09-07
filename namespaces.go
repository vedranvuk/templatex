// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package templatex implements an opinionated go template parser.
package templatex

import (
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"sync"

	"github.com/vedranvuk/errorex"
)

var (
	// ErrTemplater is the base error of templater package.
	ErrTemplater = errorex.New("templater")
	// ErrParse is returned when a parse error occurs.
	ErrParse = ErrTemplater.Wrap("parse error")
	// ErrNamespaceNotFound is returned when a non-existent namespace is
	// addressed in an ExecuteNamespace call.
	ErrNamespaceNotFound = ErrTemplater.WrapFormat("template '%s' not found")
)

// Namespaces implements hierarchical template namespaces.
type Namespaces struct {
	// index is the name of template file considered as top template to
	// be executed for a given namespace.
	index string
	// ext specifies recognised extension as template files.
	ext string
	// mu protects following private fields.
	mu *sync.Mutex
	// namespaces holds the defined namespaces.
	namespaces map[string]*template.Template
}

// New returns a new, empty Namespaces instance where index specifies the
// name of a default template file to execute in a template namespace and
// ext specifies extension including dot recognized as template files.
func New(index, ext string) *Namespaces {
	return &Namespaces{
		index:      index,
		ext:        ext,
		mu:         &sync.Mutex{},
		namespaces: make(map[string]*template.Template),
	}
}

// ParseRoot parses the specified root templates directory where index
// specifies name of a default template file in a namespace to execute when
// executing a template namespace and ext specifies extension including dot
// recognized as template files.
// If an error occurs it is returned with nil Namespaces.
func ParseRoot(root, index, ext string) (*Namespaces, error) {
	t := New(index, ext)
	if err := t.ParseRoot(root); err != nil {
		return nil, err
	}
	return t, nil
}

// parseDir recursively parses a template directory and registers it as a
// namespace under specified nsname. If an error occurs it is returned.
func (ns *Namespaces) parseDir(dir, nsname string, tmpl *template.Template) error {

	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return ErrParse.WrapCause("read file infos", err)
	}

	subs := make([]string, 0, len(fis))

	for _, fi := range fis {
		if fi.IsDir() {
			subs = append(subs, fi.Name())
			continue
		}

		match, err := filepath.Match("*"+ns.ext, fi.Name())
		if err != nil {
			return ErrParse.WrapCause("file extension match", err)
		}
		if !match {
			continue
		}

		_, err = tmpl.ParseFiles(path.Join(dir, fi.Name()))
		if err != nil {
			return ErrParse.WrapCause("parse template", err)
		}
	}

	ns.namespaces[nsname] = tmpl

	for _, sub := range subs {
		fn := filepath.Join(dir, sub)
		tn := path.Join(nsname, sub)

		nt, err := tmpl.Clone()
		if err != nil {
			return ErrParse.WrapCause("clone", err)
		}
		if err := ns.parseDir(fn, tn, nt); err != nil {
			return err
		}
	}

	return nil
}

// ParseRoot recursively parses a root template directory creating a hierarchy
// of namespaces where templates in subfolders inherit templates parsed in their
// parent directories.
//
// Namespaces are registered as paths to subfolders rooted at the specified
// root directory. For example, "/", "/home", "/settings/profile".
//
// It returns an error if one occurs.
func (ns *Namespaces) ParseRoot(root string) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	return ns.parseDir(path.Clean(root), "/", template.New(""))
}

// DefinedNamespaces returns names of defined namespaces in random order.
func (ns *Namespaces) DefinedNamespaces() (result []string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	result = make([]string, 0, len(ns.namespaces))
	for tn := range ns.namespaces {
		result = append(result, tn)
	}
	return
}

// Namespace returns a namespace template by name, if found and a truth if it
// was found.
func (ns *Namespaces) Namespace(name string) (*template.Template, bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	tt, ok := ns.namespaces[name]
	return tt, ok
}

// ExecuteNamespace executes a namespace by name using specified data to w.
// Returns an error if one occurs.
func (ns *Namespaces) ExecuteNamespace(w io.Writer, name string, data interface{}) error {
	tt, found := ns.namespaces[name]
	if !found {
		return ErrNamespaceNotFound.WrapArgs(name)
	}
	return tt.ExecuteTemplate(w, ns.index+ns.ext, data)
}
