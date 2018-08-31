package main

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"
)

// TestSvnList requires network access to svn.apache.org
func TestSvnList(t *testing.T) {
	filename := "testdata/svn/list_apache_subversion_tags.xml"
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var l SvnList
	if err := xml.Unmarshal(f, &l); err != nil {
		t.Fatal(err)
	}

	// Check number of entries
	wantN := 235
	gotN := len(l.Entries)
	if wantN != gotN {
		t.Fatalf("want %d but got %d\n", wantN, gotN)
	}

	// Check first entry
	want := SvnListEntry{
		Kind: "dir",
		Name: "0.10.0",
		Commit: SvnCommit{
			Revision: "849186",
			Author:   "cmpilato",
			// year, month, day, hour, minute, seconds, nanoseconds
			Date: time.Date(2004, 3, 18, 17, 35, 35, 355023*1000, time.UTC),
		},
	}
	got := l.Entries[0]
	if want != got {
		t.Fatalf("want %q but got %q\n", want, got)
	}
}
