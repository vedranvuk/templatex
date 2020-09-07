// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package templatex

import (
	"fmt"
	"os"
	"sort"
	"testing"
)

func TestDefinedTemplates(t *testing.T) {
	tt := New("index", ".html")
	if err := tt.ParseRoot("test/testdata"); err != nil {
		t.Fatal(err)
	}
	nss := tt.DefinedNamespaces()
	sort.Strings(nss)
	fmt.Println(nss)
}

func TestTestTemplates(t *testing.T) {
	tt := New("index", ".html")
	if err := tt.ParseRoot("test/testdata"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(os.Stdout, "/", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(os.Stdout, "/home", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(os.Stdout, "/settings", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(os.Stdout, "/settings/preferences", nil); err != nil {
		t.Fatal(err)
	}
	if err := tt.ExecuteNamespace(os.Stdout, "/settings/profile", nil); err != nil {
		t.Fatal(err)
	}
}
