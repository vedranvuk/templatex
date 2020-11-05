// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package templatex

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDefinedTemplates(t *testing.T) {
	tt := New("index", ".html")
	if err := tt.ParseRoot("test/data"); err != nil {
		t.Fatal(err)
	}
	nss := tt.DefinedNamespaces()
	fmt.Println(nss)
}

func TestNamespaces(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	tt := New("index", ".html")
	if err := tt.ParseRoot("test/data"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(buf, "/", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(buf, "/home", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(buf, "/settings", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(buf, "/settings/preferences", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(buf, "/settings/profile", nil); err != nil {
		t.Fatal(err)
	}
	var verbose bool
	for _, v := range os.Args {
		if strings.HasPrefix(v, "-test.v") {
			verbose = true
			break
		}
	}
	if verbose {
		fmt.Println(buf.String())
	}
}
