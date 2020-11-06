// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package templatex

import (
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"sync"
)

// ParseRoot is a helper that combines New and Namespaces.ParseRoot.
// Returns a nil *Namespaces and an error if one occurs.
// For details see New and Namespaces.ParseRoot.
func ParseRoot(root, index, ext string) (*Namespaces, error) {
	t := New(index, ext)
	if err := t.ParseRoot(root); err != nil {
		return nil, err
	}
	return t, nil
}

// Namespaces implements a hierarchical template parser.
//
// It takes a directory containing templates, parses it using ParseRoot and
// creates a hierarchy of templates where templates in child directories
// contain all templates parsed along the path to that child template.
// This registers namespace paths to parsed child template directories in
// Namespaces by which they can later be addressed.
//
// For example, given the following directory structure:
//
//  /home
//    index.tmpl
//  index.tmpl
//  sidebar.tmpl
//
// Two namespaces will be defined: "/" and "/home" where "/" will contain all
// templates defined in "/index.tmpl" and "/sidebar.tmpl" and "/home" will
// contain those templates and templates defined in "/home/index.tmpl".
//
// All files with the extension "ext" specified in New inside a template
// directory are parsed into the namespace corresponding to that directory.
// Template files whose filename equals "index" specified in New are files to
// be executed as the main template file of a namespace when addressing them.
//
// User can define blocks in parent template files which execute templates in
// child template files to have child namespaces inherit content partially or
// fully. All templates inside a namespace must have a unique name in that
// namespace. Parsing a template with the same name from a child template file
// replaces content of any templates with the same name parsed earlier in any
// parent template files. This is by design of go templates.
type Namespaces struct {
	mu    *sync.Mutex // mu protects following private fields.
	index string      // index is the filename of a default namespace template file.
	ext   string      // ext specifies file extension of files recognized as templates.

	namespaces map[string]*template.Template // namespaces hold defined namespaces.
}

// New returns a new, empty *Namespaces instance where index specifies the
// name of a template file to be executed as the main template in a template
// namespace directory and ext specifies extension including dot that will be
// considered as extension of files recognized as template files.
func New(index, ext string) *Namespaces {
	return &Namespaces{
		mu:         &sync.Mutex{},
		index:      index,
		ext:        ext,
		namespaces: make(map[string]*template.Template),
	}
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

// Namespace returns a namespace template by name if found and a truth if it
// was found.
func (ns *Namespaces) Namespace(name string) (*template.Template, bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	tt, ok := ns.namespaces[name]
	return tt, ok
}

// ExecuteNamespace executes a namespace by name using specified data to w.
// Returns an error if one occurs.
func (ns *Namespaces) ExecuteNamespace(w io.Writer, path string, data interface{}) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	tt, found := ns.namespaces[path]
	if !found {
		return ErrNotFound.WrapArgs(path)
	}
	return tt.ExecuteTemplate(w, ns.index+ns.ext, data)
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

// parseDir is a recursive function that parses a template directory dir into
// tmpl and registers it as a namespace with ns under specified nsname.
// If an error occurs it is returned.
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
